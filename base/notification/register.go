package notification

import (
	"fmt"
	"github.com/EduOJ/backend/base/event"
	"github.com/EduOJ/backend/database/models"
)

var RegistedPreferedNoticedMethod []string

//func Register is used to add a new method in RegistedNoticeMethod
//todo check if the new method has been alreaded registed
func Register(name string) {
	RegistedPreferedNoticedMethod = append(RegistedPreferedNoticedMethod, name)
	//..
}
//记得写对应的迁移

//func SendMessage is used to use FireEvent to launch a listener already regiusted like "email_send_message"
//if err != nil it means the listener is bad:not registed?error name?...
//mErr != nil means sendmessage faild but not because of the Notification Channel developer,you may check the sender's account
func SendMessage(receiver *models.User, title string, message string) {
	method := receiver.PreferedNoticeMethod
	result, err := event.FireEvent(fmt.Sprintf("%s_send_message", method), receiver, title, message)
	if err != nil {
		//panic
		//事件不存在？
	}
	if mErr := result[0][0].(error); mErr != nil {
		//手机号不存在？
	}

}

//func:show how registed method be used by users
//realize:check every user in db and analysis user.PreferedNoticeMethod and print the result
func ShowUsedMethod() {

}

//func:delete the mrthod in RegistedNoticeMethod by name
//realize:traverse the slice RegistedNoticeMethod and check if the name is in
// 			if so :remove ;
//			if not :do nothing and print log
func DeleteRegistedMethod(name string) {

}


//func init() {
//	event.RegisterListener("register_rouer", func(e *echo.Echo) {
//		e.POST("/email_send_message", ...)
//	})
//}
