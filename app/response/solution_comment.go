package response

import "github.com/EduOJ/backend/app/response/resource"

type GetSolutionCommentResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.SolutionComment `json:"solution"`
	} `json:"data"`
}

type CreateSolutionCommentResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.SolutionComment `json:"solution"`
	} `json:"data"`
}

type UpdateSolutionCommentResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.SolutionComment `json:"solutiong_comment"`
	} `json:"data"`
}
