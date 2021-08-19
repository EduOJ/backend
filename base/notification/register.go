package notification

import (
	"errors"
	"github.com/EduOJ/backend/base/event"
)

var RegisteredPreferredNoticedMethod []string
var ErrMethodAlreadyExist = errors.New("notice method already registered")

//func Register is used to add a new method in RegisteredNoticeMethod

func Register(name string,sendmessage func(string)error) error{
	for _,m := range RegisteredPreferredNoticedMethod {
		if m == name {
			return ErrMethodAlreadyExist
		}
	}
	RegisteredPreferredNoticedMethod = append(RegisteredPreferredNoticedMethod, name)
	eventname := name + "_send_message"
	event.RegisterListener(eventname,sendmessage)
	//..
	return nil
}

