package notification_test

import (
	"errors"
	"fmt"
	"github.com/EduOJ/backend/base/event"
	"github.com/EduOJ/backend/base/notification"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testSendingFunc(message string) error {
	return errors.New(fmt.Sprintf("notification testing %s", message))
}

func TestRegister(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, notification.Register("register_success_method",testSendingFunc))
		found := false
		for _,m := range notification.RegisteredPreferredNoticedMethod {
			if m == "register_success_method" {
				found = true
				break
			}
		}
		assert.True(t, found)
		result, err := event.FireEvent("register_success_method_send_message", "register_success_message")
		assert.Equal(t, len(result), 1)
		assert.Equal(t, len(result[0]), 1)
		assert.Equal(t, "notification testing register_success_message", result[0][0].(error).Error())
		assert.NoError(t, err)
	})

	t.Run("AlreadyExist", func(t *testing.T) {
		t.Parallel()
		event.RegisterListener("register_existing_method_send_message",testSendingFunc)
		assert.ErrorIs(t, notification.Register("register_existing_method", func(s string) error {
			return errors.New("another sending message method")
		}), notification.ErrMethodAlreadyExist)
		result, err := event.FireEvent("register_existing_method_send_message", "register_existing_message")
		assert.Equal(t, len(result), 1)
		assert.Equal(t, len(result[0]), 1)
		assert.Equal(t, "notification testing register_existing_message", result[0][0].(error).Error())
		assert.NoError(t, err)
	})

}


