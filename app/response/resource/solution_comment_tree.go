package resource

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
