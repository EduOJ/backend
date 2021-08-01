package controller

import (
	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/base/utils"
	"github.com/coreos/etcd/client"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/smtp"
)

func SendMessage (senderId uint, recerverId uint, message string,c echo.Context) (err error){
	sender := utils.GetSender(senderId)
	receiverPreferedNoticeMethod := utils.GetReceiverPreferedNoticeMethod(recerverId)
	receiverNoticeAddress := utils.GetReceiverNoticeAddress(recerverId)
	if receiverPreferedNoticeMethod == "email" {
		//use mail module to send message
		//not valid now
		smtp.SendMail(receiverNoticeAddress,nil,sender,nil,nil)
		return c.JSON(http.StatusOK,response.Response{"SENDMESSAGE_SUCCESSFUL",nil,client.User{}})
	}
	//add more if
	return c.JSON(http.StatusBadRequest, response.ErrorResp("SENDMESSAGE_ERROR",nil))
}
