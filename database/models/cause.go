package models

import (
	"gorm.io/gorm"
	"time"
)

type Cause struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProblemID  uint      `sql:"index" json:"problem_id"`
	Problem    *Problem  `json:"problem"`
	TestCaseID uint      `sql:"index" json:"test_case_id"`
	TestCase   *TestCase `json:"test_case"`

	Hash        string `json:"output_stripped_hash" gorm:"index;not null;size:255;default:''"`
	Description string `json:"description"`

	// Point: Points to be subtracted for this cause
	Point  uint `json:"point" gorm:"default:0;not null"`
	Marked bool `json:"marked" gorm:"default:false;not null"`
	Count  uint `json:"count" gorm:"default:0;not null"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_at"`
}
