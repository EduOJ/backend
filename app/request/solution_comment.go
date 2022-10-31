package request

type CreateSolutionCommentRequest struct {
	SolutionID  uint   `json:"solutionId" form:"solutionId" query:"solutionId" validate:"required"`
	FatherNode  uint   `json:"fatherNode" form:"fatherNode" query:"fatherNode" validate:"required"`
	Description string `json:"reply" form:"reply" query:"reply" validate:"required"`
	Speaker     string `json:"speaker" form:"speaker" query:"speaker" validate:"required"`
	IsRoot      bool   `json:"isRoot" form:"isRoot" query:"isRoot" validate:"required"`
}

type UpdateSolutionCommentRequest struct {
	SolutionID  uint   `json:"solutionId" form:"solutionId" query:"solutionId" validate:"required"`
	FatherNode  uint   `json:"fatherNode" form:"fatherNode" query:"fatherNode" validate:"required"`
	Description string `json:"reply" form:"reply" query:"reply" validate:"required"`
	Speaker     string `json:"speaker" form:"speaker" query:"speaker" validate:"required"`
}
