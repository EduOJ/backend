package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestUsernameValidator(t *testing.T) {
	v := New()
	t.Run("testUsernameValidatorSuccess", func(t *testing.T) {
		assert.NoError(t, v.Var("abcdefghijklmnopqrstuvwxyz", "username"))
		assert.NoError(t, v.Var("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "username"))
		assert.NoError(t, v.Var("1234567890", "username"))
		assert.NoError(t, v.Var("_____", "username"))
		assert.NoError(t, v.Var("abcABC123_", "username"))
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
