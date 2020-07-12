package models

import (
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/base/log"
)

type Log struct {
	gorm.Model
	Level   log.Level
	Message string
	Caller  string
}
