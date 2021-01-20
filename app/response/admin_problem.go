package response

import (
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
)

type CreateProblemResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemForAdmin `json:"problem"`
	} `json:"data"`
}

type AdminGetProblemResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemForAdmin `json:"problem"`
	} `json:"data"`
}

type AdminGetProblemsResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Problems []resource.ProblemForAdmin `json:"problems"`
		Total    int                        `json:"total"`
		Count    int                        `json:"count"`
		Offset   int                        `json:"offset"`
		Prev     *string                    `json:"prev"`
		Next     *string                    `json:"next"`
	} `json:"data"`
}

type UpdateProblemResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemForAdmin `json:"problem"`
	} `json:"data"`
}

type CreateTestCaseResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.TestCaseForAdmin `json:"test_case"`
	} `json:"data"`
}

type UpdateTestCaseResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.TestCaseForAdmin `json:"test_case"`
	} `json:"data"`
}
