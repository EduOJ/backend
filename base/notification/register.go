package notification

import (
	"errors"
	"github.com/EduOJ/backend/base/event"
)

var RegistedPreferredNoticedMethod []string
var ErrMethodAlreadyExist = errors.New("notice method already registered")

//func Register is used to add a new method in RegistedNoticeMethod
//todo registe the method as a new event
func Register(name string,sendmessage func(string)error) error{
	for _,m := range RegistedPreferredNoticedMethod {
		if m == name {
			return ErrMethodAlreadyExist
		}
	}
	RegistedPreferredNoticedMethod = append(RegistedPreferredNoticedMethod, name)
	eventname := name + "_send_message"
	event.RegisterListener(eventname,sendmessage)
	//..
	return nil
}

