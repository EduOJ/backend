package resource

import "github.com/EduOJ/backend/database/models"

type Solution struct {
	ID          uint   `json:"id"`
	Name        string `sql:"index" json:"name"`
	Description string `json:"description"`
}

func GetSolution(problem *models.Problem) *Problem {
	p := Problem{}
	p.convert(problem)
	return &p
}
