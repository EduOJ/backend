package response

import "github.com/leoleoasd/EduOJBackend/database/models"

type GetUserResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		models.User `json:"user"`
	} `json:"data"`
}

type GetUsersResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Users  []models.User `json:"users"` // TODO:modify models.users
		Limit  int           `json:"limit"`
		Offset int           `json:"offset"`
		Prev   string        `json:"prev"`
		Next   string        `json:"next"`
	} `json:"data"`
}
