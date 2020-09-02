package response

import (
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
)

type GetProblemResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemProfile `json:"problem"`
	} `json:"data"`
}

type GetProblemsResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Problems []resource.ProblemProfile `json:"problems"`
		Total    int                       `json:"total"`
		Count    int                       `json:"count"`
		Offset   int                       `json:"offset"`
		Prev     *string                   `json:"prev"`
		Next     *string                   `json:"next"`
	} `json:"data"`
}

type GetTestCase struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.TestCaseProfile `json:"test_case"`
	} `json:"data"`
}

type GetTestCasesResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		TestCases []resource.TestCaseProfile `json:"test_cases"`
		Total     int                        `json:"total"`
		Count     int                        `json:"count"`
		Offset    int                        `json:"offset"`
		Prev      *string                    `json:"prev"`
		Next      *string                    `json:"next"`
	} `json:"data"`
}
