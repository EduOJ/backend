package response

import "github.com/leoleoasd/EduOJBackend/app/response/resource"

type CreateProblemSetResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemSetDetail `json:"problem_set"`
	} `json:"data"`
}

type CloneProblemSetResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemSetDetail `json:"problem_set"`
	} `json:"data"`
}

type GetProblemSetResponseForAdmin struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemSetDetail `json:"problem_set"`
	} `json:"data"`
}

type GetProblemSetResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemSet `json:"problem_set"`
	} `json:"data"`
}

type UpdateProblemSetResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemSetDetail `json:"problem_set"`
	} `json:"data"`
}

type AddProblemInSetResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemSetDetail `json:"problem_set"`
	} `json:"data"`
}

type DeleteProblemInSetResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemSetDetail `json:"problem_set"`
	} `json:"data"`
}
