package request

type GetProblemRequest struct {
}

type GetProblemsRequest struct {
	Search string `json:"search" form:"search" query:"search"`

	Limit  int `json:"limit" form:"limit" query:"limit" validate:"max=100,min=0"`
	Offset int `json:"offset" form:"offset" query:"offset" validate:"min=0"`
}
