package request

import (
	"time"
)

type CreateSubmissionRequest struct {
	ProblemID uint `sql:"index" json:"problem_id"`
	// TODO: problem set?
	ProblemSetId *int   `sql:"index" gorm:"nullable" json:"problem_set_id"`
	Language     string `json:"language"`
}

type GetSubmissionRequest struct {
}

type GetSubmissionsRequest struct {
	// TODO: are these needed?
	Problem   string    `json:"problem" form:"problem" query:"problem"`
	User      string    `json:"user" form:"user" query:"user"`
	Language  string    `json:"language" form:"language" query:"language"`
	Status    string    `json:"status" form:"status" query:"status"`
	MinScore  uint      `json:"min_score" form:"min_score" query:"min_score" validate:"max=100,min=0"`
	MaxScore  uint      `json:"max_score" form:"max_score" query:"max_score" validate:"max=100,min=0"`
	StartTime time.Time `json:"start_time" form:"start_time" query:"start_time"`
	EndTime   time.Time `json:"end_time" form:"end_time" query:"end_time"`

	Limit  int `json:"limit" form:"limit" query:"limit" validate:"max=100,min=0"`
	Offset int `json:"offset" form:"offset" query:"offset" validate:"min=0"`

	OrderBy string `json:"order_by" form:"order_by" query:"order_by"`
}

type GetRunRequest struct {
}
