package response

import "github.com/leoleoasd/EduOJBackend/database/models"

type RegisterResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		models.User `json:"user"`
		Token       string `json:"token"`
	} `json:"data"`
}

type LoginResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		models.User `json:"user"`
		Token       string `json:"token"`
	} `json:"data"`
}
