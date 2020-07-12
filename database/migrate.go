package database

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
)

func Migrate() {
	base.DB.AutoMigrate(&models.Log{})
}
