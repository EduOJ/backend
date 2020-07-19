package models

import "time"

type Config struct {
	ID        uint `gorm:"primary_key"`
	Key       string
	Value     *string `gorm:"default:''"` // 可能是空字符串, 因此得是指针
	CreatedAt time.Time
	UpdatedAt time.Time
}
