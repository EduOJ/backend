package database

import (
	"fmt"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
)

func Migrate() {
	err := base.DB.AutoMigrate(&models.Log{}, &models.User{}, &models.Token{}, &models.Config{}).Error
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
}
