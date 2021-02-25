package models

import (
	"gorm.io/gorm"
	"time"
)

type ProblemSet struct {
	ID uint `gorm:"primaryKey" json:"id"`

	ClassID     uint   `sql:"index" json:"class_id" gorm:"not null"`
	Name        string `json:"name" gorm:"not null;size:255"`
	Description string `json:"description"`

	Problems []Problem `json:"problems" gorm:"many2many:problems_in_problem_sets"`
	Scores   []Grade   `json:"scores"`

	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type Grade struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID       uint `json:"user_id"`
	ProblemSetID uint `json:"problem_set_id"`

	ScoreDetail string `json:"score_detail"`
	TotalScore  int    `json:"total_score"`
}
