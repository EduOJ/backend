package response

import "github.com/EduOJ/backend/app/response/resource"

type ReadCauseResponseForAdmin struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.CauseForAdmin `json:"cause"`
	} `json:"data"`
}

type ReadCauseResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.Cause `json:"cause"`
	} `json:"data"`
}

type ReadCausesResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Causes []resource.CauseForAdmin `json:"causes"`
	} `json:"data"`
}

type UpdateCauseResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.CauseForAdmin `json:"cause"`
	} `json:"data"`
}
