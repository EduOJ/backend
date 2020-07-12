package models

import (
	"github.com/jinzhu/gorm"
)

type Log struct {
	gorm.Model
	Level   int
	Message string
	Caller  string
}
