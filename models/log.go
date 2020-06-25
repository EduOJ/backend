package models

import (
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/base/logging"
)

type Log struct {
	gorm.Model
	Level  logging.Level
	Text   string
	Caller string
}
