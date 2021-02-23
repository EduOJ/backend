package models

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"gorm.io/gorm"
	"time"
)

type Class struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name" gorm:"size:255;default:'';not null"`
	CourseName  string `json:"course_name" gorm:"size:255;default:'';not null"`
	Description string `json:"description"`
	InviteCode  string `json:"invite_code" gorm:"size:255;default:'';not null"`
	Managers    []User `json:"managers" gorm:"many2many:user_manage_classes"`
	Students    []User `json:"students" gorm:"many2many:user_in_classes"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (c Class) GetID() uint {
	return c.ID
}

func (c Class) TypeName() string {
	return "class"
}

func (c *Class) AddStudents(ids []uint) error {
	var users []User
	if err := base.DB.Find(&users, ids).Error; err != nil {
		return err
	}
	return base.DB.Model(c).Association("Students").Append(&users)
}

func (c *Class) DeleteStudents(ids []uint) error {
	var users []User
	if err := base.DB.Find(&users, ids).Error; err != nil {
		return err
	}
	return base.DB.Model(c).Association("Students").Delete(&users)
}
