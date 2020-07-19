package models

import "time"

type User struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5"`
	Nickname string `json:"nickname" validate:"required,max=30,min=5"`
	Email    string `gorm:"unique_index" json:"email" validate:"required,email,max=30,min=5"`
	Password string `json:"-" validate:"required,max=30,min=5"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
