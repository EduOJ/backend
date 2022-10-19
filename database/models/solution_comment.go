package models

import (
	"time"

	"gorm.io/gorm"
)

type SolutionComment struct {
	ID uint `gorm:"primaryKey" json:"id"`

	SolutionID  uint   `sql:"index" json:"solution_id" gorm:"not null"`
	FatherNode  uint   `json:"father_node" gorm:"not null"`
	Description string `json:"description"`
	Speaker     string `json:"speaker"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (sc SolutionComment) GetID() uint {
	return sc.ID
}

func (sc SolutionComment) GetSolutionID() uint {
	return sc.SolutionID
}

func (sc SolutionComment) GetFatherNode() uint {
	return sc.FatherNode
}
