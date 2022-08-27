package models

import (
	"time"

	"gorm.io/gorm"
)

type Solution struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `sql:"index" json:"name" gorm:"size:255;default:'';not null"`
	Description string `json:"description"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
