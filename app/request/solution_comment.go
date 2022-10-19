package request

type CreateSolutionCommentRequest struct {
	SolutionID  uint   `json:"solution_id" from:"solution_id" query:"solution_id" validate:"required"`
	FatherNode  uint   `json:"father_node" from:"father_node" query:"father_node" validate:"father_node"`
	Description string `json:"description" from:"description" query:"description" validate:"description"`
	Speaker     string `json:"speaker" from:"speaker" query:"speaker" validate:"speaker"`
}

type UpdateSolutionCommentRequest struct {
	SolutionID  uint   `json:"solution_id" from:"solution_id" query:"solution_id" validate:"required"`
	FatherNode  uint   `json:"father_node" from:"father_node" query:"father_node" validate:"father_node"`
	Description string `json:"description" from:"description" query:"description" validate:"description"`
	Speaker     string `json:"speaker" from:"speaker" query:"speaker" validate:"speaker"`
}
