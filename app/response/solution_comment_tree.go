package response

import "github.com/EduOJ/backend/app/response/resource"

type GetCommentTreeResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.SolutionCommentTree `json:"solution_comment_tree"`
	} `json:"data"`
}
