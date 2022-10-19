package resource

import "github.com/EduOJ/backend/database/models"

type Solution struct {
	ID uint `json:"id"`

	ProblemID   uint   `json:"problem_id"`
	Name        string `sql:"index" json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Likes       uint   `json:"likes"`
}

func (s *Solution) convert(solution *models.Solution) {
	s.ID = solution.ID

	s.ProblemID = solution.ProblemID
	s.Name = solution.Name
	s.Description = solution.Description
	s.Likes = solution.Likes
}

func GetSolution(solution *models.Solution) *Solution {
	s := Solution{}
	s.convert(solution)
	return &s
}
