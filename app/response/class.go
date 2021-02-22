package response

import "github.com/leoleoasd/EduOJBackend/app/response/resource"

type CreateClassResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ClassDetail `json:"class"`
	} `json:"data"`
}

type GetClassResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.Class `json:"class"`
	} `json:"data"`
}

type GetClassResponseForAdmin struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ClassDetail `json:"class"`
	} `json:"data"`
}

type GetClassesResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Classes []resource.Class `json:"classes"`
		Total   int              `json:"total"`
		Count   int              `json:"count"`
		Offset  int              `json:"offset"`
		Prev    *string          `json:"prev"`
		Next    *string          `json:"next"`
	} `json:"data"`
}

type UpdateClassResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ClassDetail `json:"class"`
	} `json:"data"`
}

type AddStudentsResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ClassDetail `json:"class"`
	} `json:"data"`
}

type RemoveStudentsResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ClassDetail `json:"class"`
	} `json:"data"`
}
