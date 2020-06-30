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
		return ut.Add("checkUsername", "{0}必须是由字母开头的4-16位字母和数字组成的字符串", true)
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
	return VerifyUsernameFormat(fl.Field().String())
}

func VerifyUsernameFormat(username string) bool {
	if ok, _ := regexp.MatchString(`^[a-zA-Z]{1}[a-zA-Z0-9]{3,15}$`, username); !ok {
		return false
	}
	return true
}

func VerifyEmailFormat(email string) bool {
	ok := regexp.MustCompile("^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$").
		MatchString(email)
	return ok
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
