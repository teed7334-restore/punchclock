package base

import (
	"net/smtp"
	"os"
)

//SendMail 寄發通知郵件
func SendMail(to []string, subject string, content string) {
	host := os.Getenv("mail.host") + ":" + os.Getenv("mail.port")
	user := os.Getenv("mail.user")
	passwd := os.Getenv("mail.password")
	auth := smtp.PlainAuth("", user, passwd, host)
	message := []byte(
		"Subject: " + subject + "\r\n" +
			"To: " + to[0] + "\r\n" +
			"From: " + user + "\r\n" +
			"Content-Type: text/plain; charset=UTF-8" + "\r\n" +
			"\r\n" +
			content + "\r\n" +
			"\r\n",
	)
	smtp.SendMail(host, auth, user, to, message)
}
