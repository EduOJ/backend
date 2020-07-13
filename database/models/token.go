package models

import "time"

type Token struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	Token     string `gorm:"unique_index" json:"token"`
	UserID    uint
	User      User
	CreatedAt time.Time `json:"created_at"`
}
