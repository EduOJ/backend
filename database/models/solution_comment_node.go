package models

import (
	"time"

	"gorm.io/gorm"
)

type SolutionCommentNode struct {
	ID uint `gorm:"primaryKey" json:"id"`

	SolutionID  uint   `sql:"index" json:"solution_id" gorm:"not null"`
	FatherNode  uint   `json:"father_node" gorm:"not null"`
	Description string `json:"description"`
	Speaker     string `json:"speaker"`

	Kids []SolutionCommentNode `json:"kids"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}