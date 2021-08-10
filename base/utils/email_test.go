package utils

import (
	"bytes"
	"github.com/go-mail/mail"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type fakeSender struct {
	messages []*mail.Message
}

func (f *fakeSender) DialAndSend(m ...*mail.Message) error {
	f.messages = append(f.messages, m...)
	return nil
}

func TestSendMail(t *testing.T) {
	fakeSend := &fakeSender{}
	sender = fakeSend
	err := SendMail("a.com", "123", "123")
	assert.NoError(t, err)
	message := fakeSend.messages[0]
	assert.Equal(t, []string{""}, message.GetHeader("From"))
	assert.Equal(t, []string{"a.com"}, message.GetHeader("To"))
	x := bytes.Buffer{}
	message.WriteTo(&x)
	messageRead := strings.SplitN(x.String(), "\r\n\r\n", 2)[1]
	assert.Equal(t, "123", messageRead)
}
