package models

import "time"

type Image struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Filename  string `gorm:"filename"`
	FilePath  string `gorm:"filepath"`
	UserID    uint
	User      User
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
