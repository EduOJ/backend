package controller

import (
	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/base/utils"
	"github.com/coreos/etcd/client"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/smtp"
)

func SendMessage (recerverName string, message string,c echo.Context) (err error){
	receiver,_ := utils.FindUser(recerverName)
	if receiver.PreferedNoticeMethod == "email" {
		//use mail module to send message
		//not valid now
		smtp.SendMail("nil",nil,"EduOJ",nil,nil)
		return c.JSON(http.StatusOK,response.Response{"SENDMESSAGE_SUCCESSFUL",nil,client.User{}})
	}
	//add more if
	return c.JSON(http.StatusBadRequest, response.ErrorResp("SENDMESSAGE_ERROR",nil))
}
