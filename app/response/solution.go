package response

import (
	"github.com/EduOJ/backend/app/response/resource"
)

type GetSolutionResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.Problem `json:"problem"`
	} `json:"data"`
}

type GetSolutionsResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Problems []resource.ProblemSummary `json:"problems"`		// 这里是否需要把problem都改成solution? --Noah
		Total    int                       `json:"total"`
		Count    int                       `json:"count"`
		Offset   int                       `json:"offset"`
		Prev     *string                   `json:"prev"`
		Next     *string                   `json:"next"`
	} `json:"data"`
}

type CreateSolutionResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemForAdmin `json:"problem"`
	} `json:"data"`
}
