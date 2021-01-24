package models

import "time"

type Submission struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID       uint   `sql:"index" json:"user_id"`
	ProblemID    uint   `sql:"index" json:"problem_id"`
	ProblemSetId *uint  `sql:"index" gorm:"nullable" json:"problem_set_id"`
	Language     string `json:"language"`
	FileName     string `json:"file_name"`

	Judged bool `json:"judged"`
	Score  uint `json:"score"`

	Runs []Run `json:"runs"`

	CreatedAt time.Time `sql:"index" json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

type Run struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID       uint  `sql:"index" json:"user_id"`
	ProblemID    uint  `sql:"index" json:"problem_id"`
	ProblemSetId *uint `sql:"index" gorm:"nullable" json:"problem_set_id"`
	SubmissionID uint  `json:"submission_id"`

	Judged                 bool   `json:"judged"`
	Status                 string `json:"status"`      // AC WA TLE MLE OLE
	MemoryUsed             uint   `json:"memory_used"` // Byte
	TimeUsed               uint   `json:"time_used"`   // ms
	FileName               string `json:"file_name"`
	OutputFileName         string `json:"output_file_name"`
	OutputStrippedHash     string `json:"output_stripped_hash"`
	CompilerOutputFileName string `json:"compiler_output_file_name"`
	ComparerOutputFileName string `json:"comparer_output_file_name"`

	CreatedAt time.Time `sql:"index" json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
