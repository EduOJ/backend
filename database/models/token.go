package models

import "time"

type Token struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Token      string `gorm:"unique_index" json:"token"`
	UserID     uint
	User       User
	RememberMe bool      `json:"remember_me"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
