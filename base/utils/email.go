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

type fakeSender struct {
	messages []*mail.Message
}

func (f *fakeSender) DialAndSend(m ...*mail.Message) error {
	f.messages = append(f.messages, m...)
	return nil
}

func SetTest() {
	sender = &fakeSender{}
}

func GetTestMessages() []*mail.Message {
	return sender.(*fakeSender).messages
}

func SendMail(address string, subject string, message string) error {

	m := mail.NewMessage()
	m.SetHeader("From", viper.GetString("email.from"))
	m.SetHeader("To", address)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	return sender.DialAndSend(m)
}
