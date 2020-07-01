package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"userSystem/pkg/app"
	"userSystem/pkg/gredis"
	"userSystem/pkg/util"
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
		userId, err := user_service.CheckUser(body.Email)
		if userId != "" && appG.HasError(err) {
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
		userId, err := user_service.CheckUser(body.Email)
		if userId == "" && appG.HasError(err) {
			return
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
	Email    string `json:"email" validate:"required,checkEmail"`
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
	userId, err := user_service.CheckUser(body.Email)
	if userId != "" && appG.HasError(err) {
		return
	}
	//判断用户名是否已注册(不可以存在)
	userId, err = user_service.CheckUser(body.UserName)
	if userId != "" && appG.HasError(err) {
		return
	}
	//解密密码
	pwVal, err := user_service.DecryptPassword(
		&user_service.Registered{Email: body.Email},
		body.Password)
	if appG.HasError(err) {
		return
	}
	//验证密码强度
	if appG.HasError(validator.VerifyPasswordFormat(pwVal)) {
		return
	}
	//注册
	err = user_service.UserRegistered(
		body.Email,
		body.UserName,
		pwVal,
		body.Code)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse("注册成功")
}

// @Summary 登录时，当用户输入用户名或邮箱(二选一)后，就调用该接口判断当前用户是否注册
// @Tags 用户
// @Produce json
// @Param query query string false "用户名或邮箱(二选一)"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/find [get]
func Find(c *gin.Context) {
	appG := app.Gin{C: c}
	query := c.Query("query")
	if query == "" {
		appG.BadResponse("缺少必要参数")
		return
	}
	if validator.VerifyEmailFormat(query) || validator.VerifyUsernameFormat(query) {
		//判断邮箱或用户名是否已注册(必须存在)
		userId, err := user_service.CheckUser(query)
		if userId == "" && appG.HasError(err) {
			return
		}
		//生成用于登录的密钥对
		pubKey, err := user_service.GenerateLoginKey(userId)
		if appG.HasError(err) {
			return
		}
		appG.SuccessResponse(pubKey)
	} else {
		appG.BadResponse("参数不合法")
	}
}

type LoginBody struct {
	User     string `json:"user" validate:"required,checkUsername|checkEmail"`
	Password string `json:"password" validate:"required,base64"`
}

// @Summary 用户登录
// @Tags 用户
// @Produce json
// @Param data body LoginBody true "登录信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/login [post]
func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	var body LoginBody
	if !appG.ParseRequest(&body) {
		return
	}
	//判断邮箱是否已注册(必须存在)
	userId, err := user_service.CheckUser(body.User)
	if userId == "" && appG.HasError(err) {
		return
	}
	//判断密码输错次数，防止暴力破解
	numStr, err := gredis.Get(user_service.LoginErrNum(userId))
	if appG.HasError(err) {
		return
	}
	if numStr != "" {
		num, err := strconv.Atoi(numStr)
		if appG.HasError(err) {
			return
		}
		if num > 2 {
			appG.BadResponse("当前账号今日登录失败次数超过3次，为保证您的账号安全，系统已锁定当前账号，您可明天再登录或立即重置密码后使用新密码登录！")
			return
		}
	}
	//解密验证
	pwVal, err := user_service.DecryptPassword(
		&user_service.Login{UserId: userId},
		body.Password)
	if appG.HasError(err) {
		return
	}
	//登录
	err = user_service.UserLogin(userId, pwVal)
	if appG.HasError(err) {
		return
	}
	//生成jwt-token
	token, err := util.GenerateToken(userId)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(token)
}

type RecoverPasswordBody struct {
	Email       string `json:"email" validate:"required,checkEmail"`
	NewPassword string `json:"newPassword" validate:"required,base64"`
	Code        string `json:"code" validate:"required,number,len=6"`
}

// @Summary 忘记密码找回
// @Tags 用户
// @Produce json
// @Param data body RecoverPasswordBody true "密码修改信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/forget [post]
func Forget(c *gin.Context) {
	appG := app.Gin{C: c}
	var body RecoverPasswordBody
	if !appG.ParseRequest(&body) {
		return
	}
	//判断邮箱是否已注册(必须存在)
	userId, err := user_service.CheckUser(body.Email)
	if userId == "" && appG.HasError(err) {
		return
	}
	//解密验证
	pwVal, err := user_service.DecryptPassword(
		&user_service.RecoverPassword{Email: body.Email},
		body.NewPassword)
	if appG.HasError(err) {
		return
	}
	//验证密码强度
	if appG.HasError(validator.VerifyPasswordFormat(pwVal)) {
		return
	}
	//重置密码
	err = user_service.ResetPassword(userId, body.Email, pwVal, body.Code)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse("修改密码成功")
}
