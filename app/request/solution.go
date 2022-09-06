package request

type CreateSolutionRequest struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required,max=255"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`
	// attachment_file(optional)
	Public *bool `json:"public" form:"public" query:"public" validate:"required"`
}

type GetSolutionRequest struct {
}

type GetSolutionsRequest struct {
	Search string `json:"search" form:"search" query:"search"`
	UserID uint   `json:"user_id" form:"user_id" query:"user_id" validate:"min=0,required_with=Tried Passed"`

	Limit  int `json:"limit" form:"limit" query:"limit" validate:"max=100,min=0"`
	Offset int `json:"offset" form:"offset" query:"offset" validate:"min=0"`

	Tags string `json:"tags" form:"tags" query:"tags"`
}
