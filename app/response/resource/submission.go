package resource

import (
	"database/sql"
	"time"
)

type Submission struct {
	ID uint `json:"id"`

	UserID       uint          `sql:"index" json:"user_id"`
	ProblemID    uint          `sql:"index" json:"problem_id"`
	ProblemSetId sql.NullInt32 `sql:"index" json:"problem_set_id"`
	Language     string        `json:"language"`

	Judged bool   `json:"judged"`
	Score  uint   `json:"score"`
	Status string `json:"status"`

	CreatedAt time.Time `sql:"index" json:"created_at"`
}

type SubmissionForAdmin struct {
	ID uint `json:"id"`

	UserID       uint          `sql:"index" json:"user_id"`
	ProblemID    uint          `sql:"index" json:"problem_id"`
	ProblemSetId sql.NullInt32 `sql:"index" json:"problem_set_id"`
	Language     string        `json:"language"`
	// TODO: remove file name?
	FileName string `json:"file_name"`
	Priority uint8  `json:"priority"`

	Judged bool   `json:"judged"`
	Score  uint   `json:"score"`
	Status string `json:"status"`

	CodeUrl string `json:"code_url"`

	Runs []Run `json:"runs"`

	CreatedAt time.Time `sql:"index" json:"created_at"`
}

type Run struct {
	ID uint `json:"id"`

	UserID       uint          `sql:"index" json:"user_id"`
	ProblemID    uint          `sql:"index" json:"problem_id"`
	ProblemSetId sql.NullInt32 `sql:"index" json:"problem_set_id"`
	TestCaseID   uint          `json:"test_case_id"`
	Sample       bool          `json:"sample"`
	SubmissionID uint          `json:"submission_id"`

	Judged     bool   `json:"judged"`
	Status     string `json:"status"`      // AC WA TLE MLE OLE
	MemoryUsed uint   `json:"memory_used"` // Byte
	TimeUsed   uint   `json:"time_used"`   // ms

	CompilerOutputUrl string `json:"compiler_output_url"`

	CreatedAt time.Time `sql:"index" json:"created_at"`
}

type RunForAdmin struct {
	ID uint `json:"id"`

	UserID       uint          `sql:"index" json:"user_id"`
	ProblemID    uint          `sql:"index" json:"problem_id"`
	ProblemSetId sql.NullInt32 `sql:"index" json:"problem_set_id"`
	TestCaseID   uint          `json:"test_case_id"`
	Sample       bool          `json:"sample"`
	SubmissionID uint          `json:"submission_id"`
	Priority     uint8         `json:"priority"`

	Judged             bool   `json:"judged"`
	Status             string `json:"status"`      // AC WA TLE MLE OLE
	MemoryUsed         uint   `json:"memory_used"` // Byte
	TimeUsed           uint   `json:"time_used"`   // ms
	OutputStrippedHash string `json:"output_stripped_hash"`

	CompilerOutputUrl string `json:"compiler_output_url"`
	InputUrl          string `json:"input_url"`
	OutputUrl         string `json:"output_url"`
	ComparerOutputUrl string `json:"comparer_output_url"`

	CreatedAt time.Time `sql:"index" json:"created_at"`
}
