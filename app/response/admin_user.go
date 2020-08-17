package response

import "github.com/leoleoasd/EduOJBackend/database/models"

type PostUserResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*models.User `json:"user"`
	} `json:"data"`
}

type PutUserResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*models.User `json:"user"`
	} `json:"data"`
}

type AdminGetUserResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*models.User `json:"user"`
	} `json:"data"`
}

type AdminGetUsersResponse struct {
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
