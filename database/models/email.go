package models

import (
	"gorm.io/gorm"
	"time"
)

type EmailVerificationToken struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint
	User   *User
	Email  string
	Token  string

	Used bool

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
