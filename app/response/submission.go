package response

import "github.com/leoleoasd/EduOJBackend/app/response/resource"

type CreateSubmissionResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.Submission `json:"submission"`
	} `json:"data"`
}

type CreateSubmissionResponseForAdmin struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.SubmissionForAdmin `json:"submission"`
	} `json:"data"`
}

type GetSubmissionResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.Submission `json:"submission"`
	} `json:"data"`
}

type GetSubmissionResponseForAdmin struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.SubmissionForAdmin `json:"submission"`
	} `json:"data"`
}

type GetSubmissionsResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Submissions []resource.Submission `json:"submissions"`
		Total       int                   `json:"total"`
		Count       int                   `json:"count"`
		Offset      int                   `json:"offset"`
		Prev        *string               `json:"prev"`
		Next        *string               `json:"next"`
	} `json:"data"`
}

type GetRunResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.Run `json:"run"`
	} `json:"data"`
}

type GetRunResponseForAdmin struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.RunForAdmin `json:"run"`
	} `json:"data"`
}
