package models

import (
	"time"
)

type WebauthnCredential struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	UserID    uint
	Content   string
	CreatedAt time.Time `json:"created_at"`
}
