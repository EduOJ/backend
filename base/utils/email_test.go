package utils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	_, err = message.WriteTo(&x)
	assert.NoError(t, err)
	messageRead := strings.SplitN(x.String(), "\r\n\r\n", 2)[1]
	assert.Equal(t, "123", messageRead)
}
