package notification

import (
	"fmt"
	"github.com/EduOJ/backend/base/event"
	"github.com/EduOJ/backend/database/models"
	"github.com/pkg/errors"
)

//func SendMessage is used to use FireEvent to launch a listener already registered like "email_send_message"
//if err != nil it means the listener is bad:not registered?error name?...
//mErr != nil means sendmessage failed but not because of the Notification Channel developer,you may check the sender's account
func SendMessage(receiver *models.User, title string, message string) error {
	method := receiver.PreferredNoticeMethod
	result, err := event.FireEvent(fmt.Sprintf("%s_send_message", method), receiver, title, message)
	flag := false
	for _, m := range RegistedPreferredNoticedMethod {
		if m == method {
			flag = true
			break
		}
	}
	if !flag {
		return errors.New("notice method not registered")
	}
	if err != nil {
		panic(err)
	}
	if mErr := result[0][0].(error); mErr != nil {
		return errors.Wrap(mErr, "failed to send message")
	}
	return nil
}

