package models

import "time"

type Script struct {
	Name      string `gorm:"primaryKey"`
	Filename  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
