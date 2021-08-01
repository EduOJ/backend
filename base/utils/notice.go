package utils

import (
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/database/models"
)

func GetSender (id uint) string{
	user := models.User{}
	PanicIfDBError(base.DB.Model(&user).Where("id = ?",id),"could not get receiver")
	return user.Nickname
}
func GetReceiverPreferedNoticeMethod (id uint) string{
	user := models.User{}
	PanicIfDBError(base.DB.Model(&user).Where("id = ?",id),"could not get receiver")
	return user.PreferedNoticeMethod
}

func GetReceiverNoticeAddress (id uint) string{
	user := models.User{}
	PanicIfDBError(base.DB.Model(&user).Where("id = ?",id),"could not get receiver")
	return user.NoticeAddress
}