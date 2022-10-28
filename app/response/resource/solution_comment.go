package resource

import "github.com/EduOJ/backend/database/models"

type SolutionComment struct {
	ID uint `gorm:"primaryKey" json:"id"`

	SolutionID  uint   `sql:"index" json:"solution_id" gorm:"not null"`
	FatherNode  uint   `json:"father_node" gorm:"not null"`
	Description string `json:"description"`
	Speaker     string `json:"speaker"`
}

func (sc *SolutionComment) convert(solutionComment *models.SolutionComment) {
	sc.ID = solutionComment.ID

	sc.SolutionID = solutionComment.SolutionID
	sc.FatherNode = solutionComment.FatherNode
	sc.Description = solutionComment.Description
	sc.Speaker = solutionComment.Speaker
}

func GetSolutionComment(solutionComment *models.SolutionComment) *SolutionComment {
	sc := SolutionComment{}
	sc.convert(solutionComment)
	return &sc
}
