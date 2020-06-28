package user_service

import (
	"userSystem/models"
	"userSystem/pkg/errmsg"
	"userSystem/pkg/gredis"
	"userSystem/pkg/mail"
	"userSystem/pkg/util"
)

type ActionType interface {
	GetEmail() string
	GetCodeStr() string
	GetPriKeyStr() string
}

type Registered struct {
	Email string
}

func (r *Registered) GetEmail() string {
	return r.Email
}
func (r *Registered) GetCodeStr() string {
	return "registeredCode-" + r.Email
}
func (r *Registered) GetPriKeyStr() string {
	return "registeredPriKey-" + r.Email
}

type RecoverPassword struct {
	Email string
}

func (r *RecoverPassword) GetEmail() string {
	return r.Email
}
func (r *RecoverPassword) GetCodeStr() string {
	return "recoverPasswordCode-" + r.Email
}
func (r *RecoverPassword) GetPriKeyStr() string {
	return "recoverPasswordPriKey-" + r.Email
}

type Login struct {
	Email string
}

func (l *Login) GetEmail() string {
	return l.Email
}
func (l *Login) GetCodeStr() string {
	return ""
}
func (l *Login) GetPriKeyStr() string {
	return "loginPriKey-" + l.Email
}

//判断用户是否已存在
func CheckUser(username, email string) (bool, error) {
	var errStr string
	if username != "" {
		isExistUsername, err := models.CheckUser(username, "")
		if err != nil {
			return false, err
		}
		if isExistUsername {
			errStr += "用户名已注册"
		}
	}
	if email != "" {
		isExistEmail, err := models.CheckUser("", email)
		if err != nil {
			return false, err
		}
		if isExistEmail {
			if errStr != "" {
				errStr += " | "
			}
			errStr += "邮箱已注册"
		}
	}
	if errStr != "" {
		return true, errmsg.NewBadMsg(errStr)
	}
	return false, nil
}

//发送验证码
func SendCode(actionType ActionType) (string, error) {
	codeStr := actionType.GetCodeStr()
	priKeyStr := actionType.GetPriKeyStr()
	//距离上一次发送要超过2分钟的时间才能重新发送
	times, err := gredis.GetTTL(codeStr)
	if err != nil {
		return "", err
	}
	if times > 13*60 {
		return "", errmsg.NewBadMsg("距离您上次发送验证码不足两分钟，请检查收件箱或垃圾箱，或两分钟后再尝试重新获取")
	}
	//生成RSA密钥对
	pubKey, priKey, err := util.GenRsaKey(2048)
	if err != nil {
		return "", err
	}
	//生成6位数字验证码
	code := util.GetRandomCode(6)
	//将验证码保存到Redis缓存，15分钟有效期
	if err := gredis.Set(codeStr, code, 15*60); err != nil {
		return "", err
	}
	//将私钥保存到Redis缓存，15分钟有效期
	if err := gredis.Set(priKeyStr, priKey, 15*60); err != nil {
		//发生错误，立即删除验证码缓存
		gredis.Delete(codeStr)
		return "", err
	}
	//发送验证码
	if err := mail.SendMail(actionType.GetEmail(), code); err != nil {
		//发生错误，立即删除验证码缓存和私钥缓存
		gredis.Delete(codeStr)
		gredis.Delete(priKeyStr)
		return "", err
	}
	return pubKey, nil
}

//使用RSA私钥解密密码
func DecryptPassword(actionType ActionType, password string) (string, error) {
	priKey, err := gredis.Get(actionType.GetPriKeyStr())
	if err != nil {
		return "", err
	}
	if priKey == "" {
		return "", errmsg.NewBadMsg("密钥获取失败，请重新获取验证码后再注册")
	}
	pw, err := util.PrivateDecrypt(priKey, password)
	if err != nil {
		return "", errmsg.NewBadMsg("密钥验证失败，请重新获取验证码后再注册")
	}
	return pw, nil
}

//注册
func UserRegistered(actionType ActionType, username, password, code string) error {
	//从缓存获取验证码
	gcode, err := gredis.Get(actionType.GetCodeStr())
	if err != nil {
		return err
	}
	if gcode == "" {
		return errmsg.NewBadMsg("当前验证码已失效，请重新获取")
	}
	if gcode != code {
		return errmsg.NewBadMsg("您输入的验证码错误，请重新输入")
	}
	//生成salt随机盐值
	salt := util.GetRandomString(8)
	//md5 密码+salt
	md5PW := util.MD5(password + salt)
	//插入到数据库
	err = models.User{
		UserID:   util.GetUUIDString(false),
		UserName: username,
		Email:    actionType.GetEmail(),
		Password: md5PW,
		Salt:     salt,
	}.InsertUser()
	if err != nil {
		return err
	}
	//删除缓存验证码
	gredis.Delete(actionType.GetCodeStr())
	//删除缓存私钥
	gredis.Delete(actionType.GetPriKeyStr())
	return nil
}
