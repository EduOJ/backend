package register

import (
	"bytes"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/log"
	"github.com/EduOJ/backend/base/utils"
	"github.com/EduOJ/backend/database/models"
)

func SendVerificationEmail(user *models.User) {
	go func() {
		verification := models.EmailVerificationToken{
			User:  user,
			Email: user.Email,
			Token: utils.RandStr(5),
			Used:  false,
		}
		if err := base.DB.Save(&verification).Error; err != nil {
			log.Error("Error saving email verification code:", err)
			return
		}
		buf := new(bytes.Buffer)
		if err := base.Template.Execute(buf, map[string]string{
			"Code":     verification.Token,
			"Nickname": user.Nickname,
		}); err != nil {
			log.Errorf("%+v\n", err)
			return
		}
		if err := utils.SendMail(user.Email, "Your email verification code", buf.String()); err != nil {
			log.Errorf("%+v\n", err)
			return
		}
		return
	}()
}
