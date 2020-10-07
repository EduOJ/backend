package request

type AdminGetLogsRequest struct {
	Levels string `json:"levels" form:"levels" query:"levels"`

	Limit  int `json:"limit" form:"limit" query:"limit" validate:"max=100,min=0"`
	Offset int `json:"offset" form:"offset" query:"offset" validate:"min=0"`
}
