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
}
