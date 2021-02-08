package models

import "time"

type Script struct {
	Name      string    `gorm:"primaryKey" json:"name"`
	Filename  string    `json:"file_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
