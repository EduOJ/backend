package resource

import (
	"github.com/EduOJ/backend/database/models"
)

type SolutionCommentNode struct {
	ID uint `gorm:"primaryKey" json:"id"`

	SolutionID  uint   `sql:"index" json:"solution_id" gorm:"not null"`
	FatherNode  uint   `json:"father_node" gorm:"not null"`
	Description string `json:"description"`
	Speaker     string `json:"speaker"`

	Kids []SolutionCommentNode `json:"kids"`
}

func (cn *SolutionCommentNode) ConvertCommentToNode(solutionComment *models.SolutionComment) {
	cn.ID = solutionComment.ID
	cn.SolutionID = solutionComment.SolutionID
	cn.FatherNode = solutionComment.FatherNode
	cn.Description = solutionComment.Description
	cn.Speaker = solutionComment.Speaker
	cn.Kids = make([]SolutionCommentNode, 0)
}

func (commentNode *SolutionCommentNode) GetKids(commentNodes []SolutionCommentNode) {
	for _, cn := range commentNodes {
		if cn.FatherNode == commentNode.FatherNode {
			cn.GetKids(commentNodes)
			commentNode.Kids = append(commentNode.Kids, cn)
		}
	}
}
