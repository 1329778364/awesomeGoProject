package mail

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"strconv"
	"userSystem/pkg/setting"
)

var (
	codeFormat = func(code string) string {
		return fmt.Sprintf("【awesomeGoProject】验证码：%s，15分钟内有效。如非本人操作，请忽略。", code)
	}
)

func SendMail(mailTo string, code string) error {
	port, err := strconv.Atoi(setting.MailSetting.Port)
	if err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "awesomeGoProject"+"<"+setting.MailSetting.User+">")
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", "awesomeGoProject")
	m.SetBody("text/html", codeFormat(code))
	d := gomail.NewDialer(setting.MailSetting.Host, port, setting.MailSetting.User, setting.MailSetting.Pass)
	return d.DialAndSend(m)
}
