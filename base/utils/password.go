package utils

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (hashed string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	hashed = string(hash)
	return
}

func VerifyPassword(password, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

func GetUserFromToken(token string) (user models.User, err error) {
	t := models.Token{}
	err = base.DB.Preload("User").Where("token = ?", token).First(&t).Error
	if err != nil {
		return
	}
	return t.User, nil
}
