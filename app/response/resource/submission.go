package resource

import (
	"github.com/leoleoasd/EduOJBackend/database/models"
	"time"
)

type Submission struct {
	ID uint `json:"id"`

	UserID       uint   `json:"user_id"`
	ProblemID    uint   `json:"problem_id"`
	ProblemSetId uint   `json:"problem_set_id"` // 0 means not in problem set
	Language     string `json:"language"`

	Judged bool   `json:"judged"`
	Score  uint   `json:"score"`
	Status string `json:"status"`

	CreatedAt time.Time `json:"created_at"`
}

func (s *Submission) convert(submission *models.Submission) {
	s.ID = submission.ID
	s.UserID = submission.UserID
	s.ProblemID = submission.ProblemID
	s.ProblemSetId = submission.ProblemSetId
	s.Language = submission.Language
	s.Judged = submission.Judged
	s.Score = submission.Score
	s.Status = submission.Status
	s.CreatedAt = submission.CreatedAt
}

func GetSubmission(submission *models.Submission) *Submission {
	s := Submission{}
	s.convert(submission)
	return &s
}

func GetSubmissionSlice(submissions []models.Submission) []Submission {
	s := make([]Submission, len(submissions))
	for i, submission := range submissions {
		s[i].convert(&submission)
	}
	return s
}

type SubmissionDetail struct {
	ID uint `json:"id"`

	UserID       uint   `json:"user_id"`
	ProblemID    uint   `json:"problem_id"`
	ProblemSetId uint   `json:"problem_set_id"`
	Language     string `json:"language"`
	FileName     string `json:"file_name"`
	Priority     uint8  `json:"priority"`

	Judged bool   `json:"judged"`
	Score  uint   `json:"score"`
	Status string `json:"status"`

	Runs []Run `json:"runs"`

	CreatedAt time.Time `json:"created_at"`
}

func (s *SubmissionDetail) convert(submission *models.Submission) {
	s.ID = submission.ID
	s.UserID = submission.UserID
	s.ProblemID = submission.ProblemID
	s.ProblemSetId = submission.ProblemSetId
	s.Language = submission.Language
	s.FileName = submission.FileName
	s.Priority = submission.Priority
	s.Judged = submission.Judged
	s.Score = submission.Score
	s.Status = submission.Status
	s.Runs = GetRunSlice(submission.Runs)
	s.CreatedAt = submission.CreatedAt
}

func GetSubmissionDetail(submission *models.Submission) *SubmissionDetail {
	s := SubmissionDetail{}
	s.convert(submission)
	return &s
}

func GetSubmissionDetailSlice(submissions []models.Submission) []SubmissionDetail {
	s := make([]SubmissionDetail, len(submissions))
	for i, submission := range submissions {
		s[i].convert(&submission)
	}
	return s
}

type Run struct {
	ID uint `json:"id"`

	UserID       uint  `json:"user_id"`
	ProblemID    uint  `json:"problem_id"`
	ProblemSetId uint  `json:"problem_set_id"`
	TestCaseID   uint  `json:"test_case_id"`
	Sample       bool  `json:"sample"`
	SubmissionID uint  `json:"submission_id"`
	Priority     uint8 `json:"priority"`

	Judged     bool   `json:"judged"`
	Status     string `json:"status"`      // AC WA TLE MLE OLE
	MemoryUsed uint   `json:"memory_used"` // Byte
	TimeUsed   uint   `json:"time_used"`   // ms

	CreatedAt time.Time `json:"created_at"`
}

func (r *Run) convert(run *models.Run) {
	r.ID = run.ID
	r.UserID = run.UserID
	r.ProblemID = run.ProblemID
	r.ProblemSetId = run.ProblemSetId
	r.TestCaseID = run.TestCaseID
	r.Sample = run.Sample
	r.SubmissionID = run.SubmissionID
	r.Priority = run.Priority
	r.Judged = run.Judged
	r.Status = run.Status
	r.MemoryUsed = run.MemoryUsed
	r.TimeUsed = run.TimeUsed
	r.CreatedAt = run.CreatedAt
}

func GetRun(run *models.Run) *Run {
	r := Run{}
	r.convert(run)
	return &r
}

func GetRunSlice(runs []models.Run) []Run {
	r := make([]Run, len(runs))
	for i, run := range runs {
		r[i].convert(&run)
	}
	return r
}
