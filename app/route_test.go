package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsernameValidator(t *testing.T) {
	v := validator.New()
	assert.Nil(t, v.RegisterValidation("username", utils.ValidateUsername))
	t.Run("testUsernameValidatorSuccess", func(t *testing.T) {
		assert.Nil(t, v.Var("abcdefghijklmnopqrstuvwxyz", "username"))
		assert.Nil(t, v.Var("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "username"))
		assert.Nil(t, v.Var("1234567890", "username"))
		assert.Nil(t, v.Var("_____", "username"))
		assert.Nil(t, v.Var("abcABC123_", "username"))
	})
	tests := []struct {
		field string
		err   string
	}{
		{
			field: "test_username_with_@",
			err:   "Key: '' Error:Field validation for '' failed on the 'username' tag",
		},
		{
			field: "test_username_with_non_ascii_char_中文",
			err:   "Key: '' Error:Field validation for '' failed on the 'username' tag",
		},
	}
	t.Run("testUsernameValidatorFail", func(t *testing.T) {
		for _, test := range tests {
			err := v.Var(test.field, "username")
			e, ok := err.(validator.ValidationErrors)
			assert.True(t, ok)
			assert.NotNil(t, e)
			assert.Equal(t, test.err, e.Error())
		}
	})
}
