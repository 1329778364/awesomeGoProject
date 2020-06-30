package v1

import (
	"github.com/gin-gonic/gin"
	"userSystem/pkg/app"
	"userSystem/pkg/validator"
	"userSystem/service/user_service"
)

type MailType int

const (
	RegisteredType      MailType = iota + 1 //注册
	RecoverPasswordType                     //忘记密码
)

type MailCodeBody struct {
	Email string   `json:"email" validate:"required,email"`
	Type  MailType `json:"type" validate:"required,gt=0"`
}

// @Summary 用户注册前和忘记密码的时候请求发送验证码
// @Tags 用户
// @Produce json
// @Param data body MailCodeBody true "发送验证码"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/sendMailCode [post]
func SendMailCode(c *gin.Context) {
	appG := app.Gin{C: c}
	var body MailCodeBody
	if !appG.ParseRequest(&body) {
		return
	}
	switch body.Type {
	case RegisteredType:
		//判断邮箱是否已注册(不可以存在)
		_, err := user_service.CheckUser("", body.Email)
		if appG.HasError(err) {
			return
		}
		//发送验证码
		pubKey, err := user_service.SendCode(&user_service.Registered{Email: body.Email})
		if appG.HasError(err) {
			return
		}
		appG.SuccessResponse(pubKey)
		return
	case RecoverPasswordType:
		//判断邮箱是否已注册(必须存在)
		if ok, err := user_service.CheckUser("", body.Email); !ok {
			if appG.HasError(err) {
				return
			} else {
				appG.BadResponse("用户不存在")
				return
			}
		}
		//发送验证码
		pubKey, err := user_service.SendCode(&user_service.RecoverPassword{Email: body.Email})
		if appG.HasError(err) {
			return
		}
		appG.SuccessResponse(pubKey)
		return
	default:
		appG.BadResponse("UnKnow Type")
		return
	}
}

type RegisteredBody struct {
	UserName string `json:"username" validate:"required,checkUsername"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,base64"`
	Code     string `json:"code" validate:"required,number,len=6"`
}

// @Summary 注册用户
// @Tags 用户
// @Produce json
// @Param data body RegisteredBody true "注册信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/registered [post]
func Registered(c *gin.Context) {
	appG := app.Gin{C: c}
	var body RegisteredBody
	if !appG.ParseRequest(&body) {
		return
	}
	//判断邮箱是否已注册(不可以存在)
	_, err := user_service.CheckUser(body.UserName, body.Email)
	if appG.HasError(err) {
		return
	}
	//解密验证
	pwVal, err := user_service.DecryptPassword(
		&user_service.Registered{Email: body.Email},
		body.Password)
	if appG.HasError(err) {
		return
	}
	if !validator.VerifyPasswordFormat(pwVal) {
		appG.BadResponse("密码必须包含数字、英文大小写字母、特殊符号（特殊符号包括: !@#~$%^&*()+|_），长度必须大于等于8位且小于等于16位")
		return
	}
	//注册
	err = user_service.UserRegistered(
		&user_service.Registered{Email: body.Email},
		body.UserName,
		pwVal,
		body.Code)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse("注册成功")
}

// @Summary 登录时，当用户输入用户名或邮箱后，就调用该接口判断当前手机号是否注册
// @Tags 用户
// @Produce json
// @Param username query string false "用户名"
// @Param email query string false "邮箱"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/find [post]
func Find(c *gin.Context) {
	appG := app.Gin{C: c}
	username := c.Query("username")
	email := c.Query("email")
	if username != "" && email == "" {
		if !validator.VerifyUsernameFormat(username) {
			appG.BadResponse("用户名不合法")
			return
		}
		appG.SuccessResponse(username)
	} else if username == "" && email != "" {
		if !validator.VerifyEmailFormat(email) {
			appG.BadResponse("邮箱不合法")
			return
		}
		appG.SuccessResponse(email)
	} else {
		appG.BadResponse("请输入用户名或邮箱")
	}
}

//TODO: 用户登录
//TODO: 找回密码
