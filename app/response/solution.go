package response

import (
	"github.com/EduOJ/backend/app/response/resource"
)

type CreateSolutionResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.Solution `json:"solution"`
	} `json:"data"`
}

type GetSolutionsResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Solutions []resource.Solution `json:"solutions"`
	} `json:"data"`
}

type GetLikesResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Likes resource.Likes `json:"likes"`
	} `json:"data"`
}
