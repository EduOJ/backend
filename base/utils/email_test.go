package utils

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSendMail(t *testing.T) {
	t.Parallel()
	SetTest()
	err := SendMail("a.com", "123", "123")
	assert.NoError(t, err)
	message := GetTestMessages()[0]
	assert.Equal(t, []string{""}, message.GetHeader("From"))
	assert.Equal(t, []string{"a.com"}, message.GetHeader("To"))
	x := bytes.Buffer{}
	message.WriteTo(&x)
	messageRead := strings.SplitN(x.String(), "\r\n\r\n", 2)[1]
	assert.Equal(t, "123", messageRead)
}
