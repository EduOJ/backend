package models

import (
	"gorm.io/gorm"
	"time"
)

type Reaction struct {
	ID uint `gorm:"primaryKey"`
	BelongType string
	BelongID uint

	Details string

	LastDealType string

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`

	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_at"`

}
