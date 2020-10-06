package response

import (
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
)

type GetProblemResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.Problem `json:"problem"`
	} `json:"data"`
}

type GetProblemsResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Problems []resource.Problem `json:"problems"`
		Total    int                `json:"total"`
		Count    int                `json:"count"`
		Offset   int                `json:"offset"`
		Prev     *string            `json:"prev"`
		Next     *string            `json:"next"`
	} `json:"data"`
}
