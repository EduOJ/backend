package resource

import (
	"github.com/EduOJ/backend/database/models"
)

type Solution struct {
	ID uint `json:"id"`

	ProblemID   uint   `json:"problem_id"`
	Name        string `sql:"index" json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Likes       string `json:"likes"`
}

type Likes struct {
	Count  int  `json:"count"`
	IsLike bool `json:"isLike"`
}

func (s *Solution) convert(solution *models.Solution) {
	s.ID = solution.ID

	s.ProblemID = solution.ProblemID
	s.Name = solution.Name
	s.Author = solution.Author
	s.Description = solution.Description
	s.Likes = solution.Likes
}

func (l *Likes) convert(likes *models.Likes) {
	l.Count = likes.Count
	l.IsLike = likes.IsLike
}

func GetSolution(solution *models.Solution) *Solution {
	s := Solution{}
	s.convert(solution)
	return &s
}

func GetSolutions(solutions []*models.Solution) (profiles []Solution) {
	profiles = make([]Solution, len(solutions))
	for i, solution := range solutions {
		profiles[i].convert(solution)
	}
	return
}

func GetLikes(likes *models.Likes) *Likes {
	l := Likes{}
	l.convert(likes)
	return &l
}
