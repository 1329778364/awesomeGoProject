package validator

import (
	chinese "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"regexp"
	"strings"
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
	// 根据自定义的标记注册翻译
	Validate.RegisterTranslation("checkUsername", trans, func(ut ut.Translator) error {
		return ut.Add("checkUsername", "{0}必须以英文字母开头，由字母、数字、下划线组成的4-16位字符，其中下划线不可以出现在开头或结尾", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("checkUsername", fe.Field())
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
	username := fl.Field().String()
	if ok, _ := regexp.MatchString(`^[a-zA-Z][a-z0-9A-Z]*(_[a-z0-9A-Z]+)*$`, username); !ok {
		return false
	}
	return true
}

func VerifyPasswordFormat(password string) bool {
	if len(password) < 8 || len(password) > 16 {
		return false
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	if ok, _ := regexp.MatchString(num, password); !ok {
		return false
	}
	if ok, _ := regexp.MatchString(a_z, password); !ok {
		return false
	}
	if ok, _ := regexp.MatchString(A_Z, password); !ok {
		return false
	}
	if ok, _ := regexp.MatchString(symbol, password); !ok {
		return false
	}
	return true
}
