package log

import (
	"time"
)

type Log struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Level     *int      `json:"level"`
	Message   string    `json:"message"`
	Caller    string    `json:"caller"`
	CreatedAt time.Time `json:"created_at"`
}
