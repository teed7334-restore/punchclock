package base

import (
	"net/smtp"
)

//SendMail 寄發通知郵件
func SendMail(to []string, subject string, content string) {
	host := cfg.Mail.Host + ":" + cfg.Mail.Port
	auth := smtp.PlainAuth("", cfg.Mail.User, cfg.Mail.Password, cfg.Mail.Host)
	message := []byte(
		"Subject: " + subject + "\r\n" +
			"To: " + to[0] + "\r\n" +
			"From: " + cfg.Mail.User + "\r\n" +
			"Content-Type: text/plain; charset=UTF-8" + "\r\n" +
			"\r\n" +
			content + "\r\n" +
			"\r\n",
	)
	smtp.SendMail(host, auth, cfg.Mail.User, to, message)
}
