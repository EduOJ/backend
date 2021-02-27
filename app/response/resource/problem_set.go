package resource

import (
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"time"
)

type ProblemSetDetail struct {
	ID uint `json:"id"`

	ClassID     uint   `json:"class_id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Problems []Problem `json:"problems"`
	Scores   []Grade   `json:"scores"`

	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

type ProblemSet struct {
	ID uint `json:"id"`

	ClassID     uint   `json:"class_id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Problems []Problem `json:"problems"`

	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

type Grade struct {
	ID uint `json:"id"`

	UserID       uint `json:"user_id"`
	ProblemSetID uint `json:"problem_set_id"`

	ScoreDetail string `json:"score_detail"`
	TotalScore  uint   `json:"total_score"`
}

func (p *ProblemSetDetail) convert(problemSet *models.ProblemSet) {
	p.ID = problemSet.ID
	p.ClassID = problemSet.ClassID
	p.Name = problemSet.Name
	p.Description = problemSet.Description
	p.Problems = GetProblemSlice(problemSet.Problems)
	p.Scores = GetGradeSlice(problemSet.Scores)
	p.StartAt = problemSet.StartAt
	p.EndAt = problemSet.EndAt
}

func (p *ProblemSet) convert(problemSet *models.ProblemSet) {
	p.ID = problemSet.ID
	p.ClassID = problemSet.ClassID
	p.Name = problemSet.Name
	p.Description = problemSet.Description
	p.Problems = GetProblemSlice(problemSet.Problems)
	p.StartAt = problemSet.StartAt
	p.EndAt = problemSet.EndAt
}

func GetProblemSet(problemSet *models.ProblemSet) *ProblemSet {
	p := ProblemSet{}
	p.convert(problemSet)
	return &p
}

func GetProblemSetDetail(problemSet *models.ProblemSet) *ProblemSetDetail {
	p := ProblemSetDetail{}
	p.convert(problemSet)
	return &p
}

func GetProblemSetSlice(problemSetSlice []*models.ProblemSet) (ps []ProblemSet) {
	ps = make([]ProblemSet, len(problemSetSlice))
	for i, problemSet := range problemSetSlice {
		ps[i].convert(problemSet)
	}
	return
}

func (g *Grade) convert(grade *models.Grade) {
	g.ID = grade.ID
	g.UserID = grade.UserID
	g.ProblemSetID = grade.ProblemSetID
	b, err := grade.ScoreDetail.MarshalJSON()
	if err != nil {
		panic(errors.Wrap(err, "could not marshal json for converting grade"))
	}
	g.ScoreDetail = string(b)
	g.TotalScore = grade.TotalScore
}

func GetGrade(grade *models.Grade) *Grade {
	g := Grade{}
	g.convert(grade)
	return &g
}

func GetGradeSlice(grades []*models.Grade) (g []Grade) {
	g = make([]Grade, len(grades))
	for i, grade := range grades {
		g[i].convert(grade)
	}
	return
}
