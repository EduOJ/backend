package notification

import (
	"fmt"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/event"
	"github.com/EduOJ/backend/base/utils"
	"github.com/EduOJ/backend/database"
	"github.com/EduOJ/backend/database/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegister(t *testing.T) {
	t.Parallel()
	assert.NoError(t, Register("register_success_method"))
	_, found := registeredPreferredNoticedMethod["register_success_method"]
	assert.True(t, found)
	assert.ErrorIs(t, Register("register_success_method"), ErrMethodAlreadyExist)
}

func TestSendMessage(t *testing.T) {
	t.Parallel()

	receiver1 := models.User{
		Username:              "test_send_message_1",
		Nickname:              "test_send_message_1_nick",
		Email:                 "test_send_message_1@mail.com",
		Password:              utils.HashPassword("test_send_message_1_password"),
		PreferredNoticeMethod: "test_send_message_registered_method",
		NoticeAccount:         `{"test_send_message_registered_method": "test_send_message_registered_account"}`,
	}
	receiver2 := models.User{
		Username:              "test_send_message_2",
		Nickname:              "test_send_message_2nick",
		Email:                 "test_send_message_2@mail.com",
		Password:              utils.HashPassword("test_send_message_2_password"),
		PreferredNoticeMethod: "test_send_message_unregistered_method",
		NoticeAccount:         `{"test_send_message_unregistered_method": "test_send_message_unregistered_account"}`,
	}
	receiver3 := models.User{
		Username:              "test_send_message_3",
		Nickname:              "test_send_message_3nick",
		Email:                 "test_send_message_3@mail.com",
		Password:              utils.HashPassword("test_send_message_3_password"),
		PreferredNoticeMethod: "test_send_message_registered_method",
		NoticeAccount:         ``,
	}

	event.RegisterListener("test_send_message_registered_method_send_message", func(account interface{}, title, message string, extras map[string]interface{}) error {
		if title == "send_message_success_title" {
			assert.Equal(t, "test_send_message_registered_account", account)
			return nil
		} else {
			return errors.New(fmt.Sprintf("test send message %s, %s, %s", title, message, account))
		}
	})
	assert.NoError(t, Register("test_send_message_registered_method"))

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, SendMessage(&receiver1, "send_message_success_title", "send_message_success_message", map[string]interface{}{}))
	})

	t.Run("NotRegistered", func(t *testing.T) {
		t.Parallel()
		assert.ErrorIs(t, SendMessage(&receiver2, "send_message_not_registered_title", "send_message_not_registered_message", map[string]interface{}{}), ErrNoticeMethodNotRigisted)
	})
	t.Run("NoNoticeAccount", func(t *testing.T) {
		t.Parallel()
		err := SendMessage(&receiver3, "send_message_no_account_title", "send_message_no_account_message", map[string]interface{}{})
		assert.NotNil(t, err)
		assert.Equal(t, "receiver's test_send_message_registered_method account not found!", err.(error).Error())
	})

	t.Run("SendFailed", func(t *testing.T) {
		t.Parallel()
		err := SendMessage(&receiver1, "send_message_fail_title", "send_message_fail_message", map[string]interface{}{})
		assert.NotNil(t, err)
		assert.Equal(t, "failed to send message: test send message send_message_fail_title, send_message_fail_message, test_send_message_registered_account", err.Error())
	})
}

func TestSetAccount(t *testing.T) {
	t.Parallel()
	user := models.User{
		Username:              "test_set_account_username",
		Nickname:              "test_set_account_nickname",
		Email:                 "test_set_account@mail.com",
		Password:              "test_set_account_pwd",
		PreferredNoticeMethod: "test_set_account",
		NoticeAccount:         "",
	}
	assert.NoError(t, Register("test_account_method"))
	assert.NoError(t, base.DB.Create(&user).Error)
	assert.NoError(t, SetAccount("test_account_method", &user, "test_account_content"))
	databaseUser := models.User{}
	assert.NoError(t, base.DB.First(&databaseUser, user.ID).Error)
	assert.Equal(t, user.NoticeAccount, databaseUser.NoticeAccount)
}

func TestMain(m *testing.M) {
	defer database.SetupDatabaseForTest()()
	m.Run()
}
