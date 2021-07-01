package models

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct{

	ID       uint `gorm:"primaryKey"`

	UserID uint
	Writer User `gorm:"foreignKey:UserID"`

	Content string


	IfDeleted bool

	FirstID uint
	FirstType string

	FatherID uint
	FatherType string


	Detail string


	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_at"`

}






