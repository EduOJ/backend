package models

import (
	"time"

	"gorm.io/gorm"
)

type SolutionCommentTree struct {
	Roots []SolutionCommentNode `json:"roots"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
