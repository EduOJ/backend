package response

import (
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
)

type AdminCreateProblemResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemProfileForAdmin `json:"problem"`
	} `json:"data"`
}

type AdminGetProblemResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemProfileForAdmin `json:"problem"`
	} `json:"data"`
}

type AdminGetProblemsResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Problems []resource.ProblemProfileForAdmin `json:"problems"`
		Total    int                               `json:"total"`
		Count    int                               `json:"count"`
		Offset   int                               `json:"offset"`
		Prev     *string                           `json:"prev"`
		Next     *string                           `json:"next"`
	} `json:"data"`
}

type AdminUpdateProblemResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemProfileForAdmin `json:"problem"`
	} `json:"data"`
}

type AdminCreateTestCase struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.TestCaseProfileForAdmin `json:"test_case"`
	} `json:"data"`
}

type AdminUpdateTestCase struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.TestCaseProfileForAdmin `json:"test_case"`
	} `json:"data"`
}
