package response

import "github.com/leoleoasd/EduOJBackend/database/models"

type GetUserResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*models.User `json:"user"`
	} `json:"data"`
}

type GetMeResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*models.User `json:"user"`
	} `json:"data"`
}

type GetUsersResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Users  []models.User `json:"users"` // TODO:modify models.users
		Total  int           `json:"total"`
		Count  int           `json:"count"`
		Offset int           `json:"offset"`
		Prev   *string       `json:"prev"`
		Next   *string       `json:"next"`
	} `json:"data"`
}

type UpdateMeResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*models.User `json:"user"`
	} `json:"data"`
}
