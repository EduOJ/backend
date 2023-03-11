package models

import (
	"time"

	"github.com/EduOJ/backend/base"
	"gorm.io/gorm"
)

type Class struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `json:"name" gorm:"size:255;default:'';not null"`
	CourseName  string  `json:"course_name" gorm:"size:255;default:'';not null"`
	Description string  `json:"description"`
	InviteCode  string  `json:"invite_code" gorm:"size:255;default:'';not null"`
	Managers    []*User `json:"managers" gorm:"many2many:user_manage_classes"`
	Students    []*User `json:"students" gorm:"many2many:user_in_classes"`

	ProblemSets []*ProblemSet `json:"problem_sets"`

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
	if len(ids) == 0 {
		return nil
	}
	existingIds := make([]uint, len(c.Students))
	for i, s := range c.Students {
		existingIds[i] = s.ID
	}
	var users []User
	query := base.DB
	if len(existingIds) != 0 {
		query = base.DB.Where("id not in (?)", existingIds)
	}
	if err := query.Find(&users, ids).Error; err != nil {
		return err
	}
	return base.DB.Model(c).Association("Students").Append(&users)
}

func (c *Class) DeleteStudents(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	var users []User
	if err := base.DB.Find(&users, ids).Error; err != nil {
		return err
	}
	return base.DB.Model(c).Association("Students").Delete(&users)
}

func (c *Class) AfterDelete(tx *gorm.DB) error {
	var problemSets []ProblemSet
	err := tx.Find(&problemSets, "class_id = ?", c.ID).Error
	if err != nil {
		return err
	}
	if len(problemSets) != 0 {
		return tx.Delete(&problemSets).Error
	}
	return nil
}
