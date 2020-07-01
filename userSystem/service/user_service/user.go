package user_service

import (
	"userSystem/models"
	"userSystem/pkg/errmsg"
	"userSystem/pkg/gredis"
	"userSystem/pkg/mail"
	"userSystem/pkg/util"
	"userSystem/pkg/validator"
)

type ActionType interface {
	ID() string
	CodeKey() string
	PubKey() string
	PriKey() string
}

type Registered struct {
	Email string
}

func (r *Registered) ID() string {
	return r.Email
}
func (r *Registered) CodeKey() string {
	return "registeredCode-" + r.Email
}
func (r *Registered) PubKey() string {
	return ""
}
func (r *Registered) PriKey() string {
	return "registeredPriKey-" + r.Email
}

type RecoverPassword struct {
	Email string
}

func (r *RecoverPassword) ID() string {
	return r.Email
}
func (r *RecoverPassword) CodeKey() string {
	return "recoverPasswordCode-" + r.Email
}
func (r *RecoverPassword) PubKey() string {
	return ""
}
func (r *RecoverPassword) PriKey() string {
	return "recoverPasswordPriKey-" + r.Email
}

type Login struct {
	UserId string
}

func (l *Login) ID() string {
	return l.UserId
}
func (l *Login) CodeKey() string {
	return ""
}
func (l *Login) PubKey() string {
	return "loginPubKey-" + l.UserId
}
func (l *Login) PriKey() string {
	return "loginPriKey-" + l.UserId
}

var (
	LoginErrNum = func(userId string) string {
		return "loginErrNum-" + userId
	}
)

//判断用户是否已存在
func CheckUser(query string) (userId string, err error) {
	if validator.VerifyEmailFormat(query) {
		userId, err = models.CheckUser("", query)
		if err != nil {
			return "", err
		}
		if userId != "" {
			return userId, errmsg.NewBadMsg("邮箱已注册")
		}
	} else if validator.VerifyUsernameFormat(query) {
		userId, err = models.CheckUser(query, "")
		if err != nil {
			return "", err
		}
		if userId != "" {
			return userId, errmsg.NewBadMsg("用户名已注册")
		}
	}
	return "", errmsg.NewBadMsg("用户不存在")
}

//发送验证码(注册或忘记密码)
func SendCode(actionType ActionType) (string, error) {
	//距离上一次发送要超过1分钟的时间才能重新发送
	times, err := gredis.GetTTL(actionType.CodeKey())
	if err != nil {
		return "", err
	}
	if times > 14*60 {
		return "", errmsg.NewBadMsg("距离您上次发送验证码不足一分钟，请检查收件箱或垃圾箱，或两分钟后再尝试重新获取")
	}
	//生成6位数字验证码
	code := util.GetRandomCode(6)
	//发送验证码
	if err := mail.SendMail(actionType.ID(), code); err != nil {
		return "", err
	}
	//生成RSA密钥对
	pubKey, priKey, err := util.GenRsaKey(2048)
	if err != nil {
		return "", err
	}
	//将验证码和私钥保存到Redis缓存，15分钟有效期
	if err := gredis.Set(map[string]string{actionType.CodeKey(): code, actionType.PriKey(): priKey},
		15*60); err != nil {
		//发生错误，立即删除验证码缓存
		gredis.Delete(actionType.CodeKey())
		gredis.Delete(actionType.PriKey())
		return "", err
	}
	return pubKey, nil
}

//生成用于登录的密钥对
func GenerateLoginKey(userId string) (string, error) {
	actionType := &Login{UserId: userId}
	//原来的公钥有效时间大于15分钟，继续使用
	times, err := gredis.GetTTL(actionType.PubKey())
	if err != nil {
		return "", err
	}
	if times > 15*60 {
		pubKey, err := gredis.Get(actionType.PubKey())
		if err != nil {
			return "", err
		}
		return pubKey, nil
	}
	//生成RSA密钥对
	pubKey, priKey, err := util.GenRsaKey(2048)
	if err != nil {
		return "", err
	}
	//将密钥对保存到Redis缓存，30分钟有效期
	if err := gredis.Set(map[string]string{
		actionType.PubKey(): pubKey,
		actionType.PriKey(): priKey},
		30*60); err != nil {
		//发生错误，立即删除密钥对缓存
		gredis.Delete(actionType.PubKey())
		gredis.Delete(actionType.PriKey())
		return "", err
	}
	return pubKey, nil
}

//使用RSA私钥解密密码(注册或登录或忘记密码)
func DecryptPassword(actionType ActionType, password string) (string, error) {
	priKey, err := gredis.Get(actionType.PriKey())
	if err != nil {
		return "", err
	}
	if priKey == "" {
		return "", errmsg.NewBadMsg("密钥失效，请重新获取验证码后再操作")
	}
	pw, err := util.PrivateDecrypt(priKey, password)
	if err != nil {
		return "", errmsg.NewBadMsg("密钥验证失败，请检查公钥")
	}
	return pw, nil
}

//注册
func UserRegistered(email, username, password, code string) error {
	actionType := &Registered{Email: email}
	//从缓存获取验证码
	gcode, err := gredis.Get(actionType.CodeKey())
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
		Email:    actionType.ID(),
		Password: md5PW,
		Salt:     salt,
	}.InsertUser()
	if err != nil {
		return err
	}
	//删除缓存验证码
	gredis.Delete(actionType.CodeKey())
	//删除缓存私钥
	gredis.Delete(actionType.PriKey())
	return nil
}

//登录
func UserLogin(userId, password string) error {
	actionType := &Login{UserId: userId}
	//从数据库获取用户
	user, err := models.GetUserById(userId)
	if err != nil {
		return err
	}
	//判断密码是否正确
	if util.MD5(password+user.Salt) != user.Password {
		//记录密码输出次数
		if err := gredis.Incr(LoginErrNum(userId), util.GetRemainSecondsOneDay()); err != nil {
			return err
		}
		return errmsg.NewBadMsg("密码错误")
	}
	//删除缓存密钥对
	gredis.Delete(actionType.PubKey())
	gredis.Delete(actionType.PriKey())
	return nil
}

//重置密码
func ResetPassword(userId, email, password, code string) error {
	actionType := &RecoverPassword{Email: email}
	//从数据库获取用户
	user, err := models.GetUserById(userId)
	if err != nil {
		return err
	}
	//从缓存获取验证码
	gcode, err := gredis.Get(actionType.CodeKey())
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
	//修改密码
	err = user.UpdatePassword(md5PW, salt)
	if err != nil {
		return err
	}
	//删除缓存验证码
	gredis.Delete(actionType.CodeKey())
	//删除缓存私钥
	gredis.Delete(actionType.PriKey())
	//删除密码错误记录
	gredis.Delete(LoginErrNum(userId))
	return nil
}
