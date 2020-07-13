package models

import (
	"time"
)

type Log struct {
	ID        uint `gorm:"primary_key"`
	Level     *int
	Message   string
	Caller    string
	CreatedAt time.Time
}
