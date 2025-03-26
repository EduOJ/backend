package request

type CreateSolutionRequest struct {
	ProblemID   uint   `json:"problemID" form:"problemID" query:"problemID" validate:"required"`
	Name        string `json:"name" form:"name" query:"name" validate:"required,max=255"`
	Author      string `json:"author" form:"author" query:"author" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`
}

type GetSolutionsRequest struct {
	ProblemID string `json:"problemID" form:"problemID" query:"problemID" validate:"required"`
}

type UpdateSolutionRequest struct {
	ProblemID   uint   `json:"problemID" form:"problemID" query:"problemID" validate:"required"`
	Name        string `json:"name" form:"name" query:"name" validate:"required,max=255"`
	Author      string `json:"author" form:"author" query:"author" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`
	Likes       string `json:"likes" form:"likes" query:"likes" validate:"required"`
}

type LikesRequest struct {
	SolutionId uint `json:"solutionId" form:"solutionId" query:"solutionId" validate:"required"`
	UserId     uint `json:"userId" form:"userId" query:"userId" validate:"required"`
	IsLike     int  `json:"isLike" form:"isLike" query:"isLike"`
}

// type GetSolutionsRequest struct {
// 	Search string `json:"search" form:"search" query:"search"`
// 	UserID uint   `json:"user_id" form:"user_id" query:"user_id" validate:"min=0,required_with=Tried Passed"`

// 	Limit  int `json:"limit" form:"limit" query:"limit" validate:"max=100,min=0"`
// 	Offset int `json:"offset" form:"offset" query:"offset" validate:"min=0"`

// 	Tags string `json:"tags" form:"tags" query:"tags"`
// }
