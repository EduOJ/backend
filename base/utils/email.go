package utils

import (
	"github.com/EduOJ/backend/base"
	"github.com/go-mail/mail"
	"github.com/spf13/viper"
)

type DialSender interface {
	DialAndSend(m ...*mail.Message) error
}

var sender = (DialSender)(&base.Mail)

func SendMail(address string, subject string, message string) error {

	m := mail.NewMessage()
	m.SetHeader("From", viper.GetString("email.from"))
	m.SetHeader("To", address)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	return sender.DialAndSend(m)
}
