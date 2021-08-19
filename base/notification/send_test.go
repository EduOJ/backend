package notification_test

import (
	"errors"
	"fmt"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/event"
	"github.com/EduOJ/backend/base/notification"
	"github.com/EduOJ/backend/base/utils"
	"github.com/EduOJ/backend/database"
	"github.com/EduOJ/backend/database/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendMessage(t *testing.T) {
	defer database.SetupDatabaseForTest()()
	t.Parallel()

	receiver1 := models.User{
		Username:             "test_send_message_1",
		Nickname:             "test_send_message_1_nick",
		Email:                "test_send_message_1@mail.com",
		Password:             utils.HashPassword("test_send_message_1_password"),
		PreferredNoticeMethod: "test_send_message_registered_method",
		NoticeAccount:        `{"test_send_message_registered_method": "test_send_message_registered_account"}`,
	}
	receiver2 := models.User{
		Username:             "test_send_message_2",
		Nickname:             "test_send_message_2nick",
		Email:                "test_send_message_2@mail.com",
		Password:             utils.HashPassword("test_send_message_2_password"),
		PreferredNoticeMethod: "test_send_message_unregistered_method",
		NoticeAccount:        `{"test_send_message_unregistered_method": "test_send_message_unregistered_account"}`,
	}
	assert.NoError(t, base.DB.Create(&receiver1).Error)
	assert.NoError(t, base.DB.Create(&receiver2).Error)

	notification.RegistedPreferredNoticedMethod = append(notification.RegistedPreferredNoticedMethod, "test_send_message_registered_method")
	event.RegisterListener("test_send_message_registered_method_send_message", func(receiver *models.User, title, message string) error {
		if title == "send_message_success_title" {
			return nil
		} else {
			return errors.New(fmt.Sprintf("test send message %s, %s, %s", title, message, receiver.NoticeAccount))
		}
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, notification.SendMessage(&receiver1,"send_message_success_title","send_message_success_message"))
	})

	t.Run("NotRegistered", func(t *testing.T) {
		t.Parallel()
		assert.ErrorIs(t, notification.SendMessage(&receiver2,"send_message_not_registered_title","send_message_not_registered_message"), notification.ErrNoticeMethodNotRigisted)
	})

	t.Run("SendFailed", func(t *testing.T) {
		t.Parallel()
		err := notification.SendMessage(&receiver1,"send_message_fail_title","send_message_fail_message")
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Sprintf("failed to send message: test send message send_message_fail_title, send_message_fail_message, %s", receiver1.NoticeAccount), err.Error())
	})

}

