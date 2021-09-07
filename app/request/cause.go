package request

type ReadCauseRequest struct {
}

type ReadCausesRequest struct {
}

type UpdateCauseRequest struct {
	Description string `json:"description" form:"description" query:"description"`
	Point       string `json:"point" form:"point" query:"point" validate:"max=100,min=0"`
}

type DeleteCauseRequest struct {
}
