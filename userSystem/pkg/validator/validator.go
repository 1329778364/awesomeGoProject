package validator

import (
	chinese "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"regexp"
	"strings"
	"userSystem/pkg/errmsg"
)

var (
	Validate *validator.Validate
	trans    ut.Translator
)

func InitValidator() {
	// 中文翻译
	zh := chinese.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")
	Validate = validator.New()
	// 验证器注册翻译器
	zhTranslations.RegisterDefaultTranslations(Validate, trans)
	// 自定义验证方法
	Validate.RegisterValidation("checkUsername", checkUsername)
	Validate.RegisterValidation("checkEmail", checkEmail)
	// 根据自定义的标记注册翻译
	Validate.RegisterTranslation("checkUsername", trans, func(ut ut.Translator) error {
		return ut.Add("checkUsername", "{0}必须是由字母开头的4-16位字母和数字组成的字符串", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("checkUsername", fe.Field())
		return t
	})
	Validate.RegisterTranslation("checkEmail", trans, func(ut ut.Translator) error {
		return ut.Add("checkEmail", "{0}不合法", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("checkEmail", fe.Field())
		return t
	})
}

func Translate(errs validator.ValidationErrors) string {
	var errList []string
	for _, e := range errs {
		errList = append(errList, e.Translate(trans))
	}
	return strings.Join(errList, "|")
}

func checkUsername(fl validator.FieldLevel) bool {
	return VerifyUsernameFormat(fl.Field().String())
}

func VerifyUsernameFormat(username string) bool {
	if ok, _ := regexp.MatchString(`^[a-zA-Z]{1}[a-zA-Z0-9]{3,15}$`, username); !ok {
		return false
	}
	return true
}

func checkEmail(fl validator.FieldLevel) bool {
	return VerifyEmailFormat(fl.Field().String())
}

func VerifyEmailFormat(email string) bool {
	return regexp.MustCompile(`^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`).
		MatchString(email)
}

func VerifyPasswordFormat(password string) error {
	err := errmsg.NewBadMsg("密码必须包含数字、英文大小写字母、特殊符号（特殊符号包括: !@#~$%^&*()+|_），长度必须大于等于8位且小于等于16位")
	if len(password) < 8 || len(password) > 16 {
		return err
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	if ok, _ := regexp.MatchString(num, password); !ok {
		return err
	}
	if ok, _ := regexp.MatchString(a_z, password); !ok {
		return err
	}
	if ok, _ := regexp.MatchString(A_Z, password); !ok {
		return err
	}
	if ok, _ := regexp.MatchString(symbol, password); !ok {
		return err
	}
	return nil
}
