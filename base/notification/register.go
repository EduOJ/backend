package notification

import (
	"fmt"
	"github.com/EduOJ/backend/base/event"
	"github.com/EduOJ/backend/database/models"
)

var RegistedPreferedNoticedMethod []string

func register(name string) {
	RegistedPreferedNoticedMethod = append(RegistedPreferedNoticedMethod, name)
	//..
}
//记得写对应的迁移

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
//func init() {
//	event.RegisterListener("register_rouer", func(e *echo.Echo) {
//		e.POST("/bind_sms", ...)
//	})
//}
