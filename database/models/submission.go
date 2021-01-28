package models

import (
	"time"
)

const DEFAULT_PRIORITY = uint8(127)

type Submission struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID       uint     `sql:"index" json:"user_id"`
	User         *User    `json:"user"`
	ProblemID    uint     `sql:"index" json:"problem_id"`
	Problem      *Problem `json:"problem"`
	ProblemSetId uint     `sql:"index" json:"problem_set_id"`
	Language     string   `json:"language"`
	FileName     string   `json:"file_name"`
	Priority     uint8    `json:"priority"`

	Judged bool   `json:"judged"`
	Score  uint   `json:"score"`
	Status string `json:"status"`

	Runs []Run `json:"runs"`

	CreatedAt time.Time `sql:"index" json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

type Run struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID       uint        `sql:"index" json:"user_id"`
	User         *User       `json:"user"`
	ProblemID    uint        `sql:"index" json:"problem_id"`
	Problem      *Problem    `json:"problem"`
	ProblemSetId uint        `sql:"index" json:"problem_set_id"`
	TestCaseID   uint        `json:"test_case_id"`
	TestCase     *TestCase   `json:"test_case"`
	Sample       bool        `json:"sample" gorm:"not null"`
	SubmissionID uint        `json:"submission_id"`
	Submission   *Submission `json:"submission"`
	Priority     uint8       `json:"priority"`

	Judged             bool   `json:"judged"`
	Status             string `json:"status"`      // AC WA TLE MLE OLE
	MemoryUsed         uint   `json:"memory_used"` // Byte
	TimeUsed           uint   `json:"time_used"`   // ms
	OutputStrippedHash string `json:"output_stripped_hash"`

	CreatedAt time.Time `sql:"index" json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
