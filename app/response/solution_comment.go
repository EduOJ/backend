package response

import "github.com/EduOJ/backend/app/response/resource"

type GetSolutionCommentResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.SolutionComment `json:"solution_comment"`
	} `json:"data"`
}

type CreateSolutionCommentResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.SolutionComment `json:"solution_comment_create"`
	} `json:"data"`
}

type UpdateSolutionCommentResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.SolutionComment `json:"solution_comment"`
	} `json:"data"`
}

type GetSolutionCommentsResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		SolutionComments []resource.SolutionComment `json:"solution_comments"`
	} `json:"data"`
}

type GetCommentTreeResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.SolutionCommentTree `json:"solution_comment_tree"`
	} `json:"data"`
}
