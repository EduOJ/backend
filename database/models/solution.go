package models

import (
	"time"

	"gorm.io/gorm"
)

type Solution struct {
	ID uint `gorm:"primaryKey" json:"id"`

	ProblemID   uint   `sql:"index" json:"problem_id"`
	Name        string `json:"name"`
	Author      string `json:"auther"`
	Description string `json:"description"`
	Likes       string `json:"likes"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type Likes struct {
	Count  int  `json:"count"`
	IsLike bool `json:"isLike"`
}

func (s Solution) GetID() uint {
	return s.ID
}

func (s Solution) GetProblemID() uint {
	return s.ProblemID
}
