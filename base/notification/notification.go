// This package is to send external message to users in their preferred way
// Work for developing notification channels:
// 	* Registering the channel name in init using function Register.
// 	* Operating account data using custom API via function SetAccount.
// 		* Register event listeners in group "register_route".
// 	* Registering a event listener for sending message in group "{channel name}_send_message".
// To send a message:
// 	* Use function SendMessage.

package notification

import (
	"encoding/json"
	"fmt"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/event"
	"github.com/EduOJ/backend/database/models"
	"github.com/pkg/errors"
	"sync"
)

var registeredPreferredNoticedMethod = make(map[string]struct{})
var ErrMethodAlreadyExist = errors.New("notice method already registered")
var ErrNoticeMethodNotRigisted = errors.New("notice methord not registered")
var methodLock sync.RWMutex

// func Register is used to add a new method in RegisteredNoticeMethod
func Register(name string) error {
	methodLock.Lock()
	defer methodLock.Unlock()
	if _, found := registeredPreferredNoticedMethod[name]; found {
		return ErrMethodAlreadyExist
	}
	registeredPreferredNoticedMethod[name] = struct{}{}

	return nil
}

func CheckNoticeMethod(name string) bool {
	methodLock.RLock()
	defer methodLock.RUnlock()
	_, found := registeredPreferredNoticedMethod[name]
	return found
}

// func SendMessage is used to use FireEvent to launch a listener already registered like "email_send_message"
// if err != nil it means the listener is bad:not registered?error name?...
// mErr != nil means sendmessage failed but not because of the Notification Channel developer,you may check the sender's account
func SendMessage(receiver *models.User, title string, message string, extras map[string]interface{}) error {
	method := receiver.PreferredNoticeMethod
	if !CheckNoticeMethod(method) {
		return ErrNoticeMethodNotRigisted
	}
	var data map[string]interface{}
	if receiver.NoticeAccount == "" {
		receiver.NoticeAccount = "{}"
	}
	if err := json.Unmarshal([]byte(receiver.NoticeAccount), &data); err != nil {
		return errors.Wrap(err, "could not unmarshal notice account")
	}
	var account interface{}
	if method == "email" {
		account = receiver.Email
	} else {
		var found bool
		account, found = data[method]
		if !found {
			return errors.New(fmt.Sprintf("receiver's %s account not found!", method))
		}
	}
	result, err := event.FireEvent(fmt.Sprintf("%s_send_message", method), account, title, message, extras)
	if err != nil {
		panic(err)
	}
	if result[0][0] == nil {
		return nil
	}
	return errors.Wrap(result[0][0].(error), "failed to send message")
	return nil
}

func SetAccount(method string, user *models.User, account interface{}) error {
	if !CheckNoticeMethod(method) {
		return ErrNoticeMethodNotRigisted
	}
	var data map[string]interface{}
	if user.NoticeAccount == "" {
		user.NoticeAccount = "{}"
	}
	if err := json.Unmarshal([]byte(user.NoticeAccount), &data); err != nil {
		return errors.Wrap(err, "could not unmarshal notice account")
	}
	data[method] = account
	b, err := json.Marshal(data)
	user.NoticeAccount = string(b)
	if err != nil {
		return errors.Wrap(err, "could not marshal notice account")
	}
	if err = base.DB.Save(user).Error; err != nil {
		return errors.Wrap(err, "could not save notice account")
	}
	return nil
}
