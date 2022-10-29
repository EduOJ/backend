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

func GetSolutionComments(solutionComments []models.SolutionComment) []SolutionComment {
	scs := make([]SolutionComment, 0)
	for _, sc := range solutionComments {
		scs = append(scs, *GetSolutionComment(&sc))
	}
	return scs
}

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

type SolutionCommentTree struct {
	Roots []SolutionCommentNode `json:"roots"`
}

func (t *SolutionCommentTree) BuildTree(commentNodes []SolutionCommentNode) {
	for _, commentNode := range commentNodes {
		if commentNode.FatherNode == 0 {
			commentNode.GetKids(commentNodes)
			t.Roots = append(t.Roots, commentNode)
		}
	}
}

func GetSolutionCommentTree(commentNodes []SolutionCommentNode) *SolutionCommentTree {
	sct := SolutionCommentTree{}
	sct.BuildTree(commentNodes)
	return &sct
}
