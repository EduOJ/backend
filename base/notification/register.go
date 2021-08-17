package notification

import (
	"github.com/EduOJ/backend/base/event"
	"github.com/pkg/errors"
)

var RegistedPreferredNoticedMethod []string

//func Register is used to add a new method in RegistedNoticeMethod
//todo registe the method as a new event
func Register(name string,sendmessage func(string)error) error{
	for _,m := range RegistedPreferredNoticedMethod {
		if m == name {
			return errors.New("notice method already registered")
		}
	}
	RegistedPreferredNoticedMethod = append(RegistedPreferredNoticedMethod, name)
	eventname := name + "_send_message"
	event.RegisterListener(eventname,sendmessage)
	//..
	return nil
}
//记得写对应的迁移
