package response

import "github.com/leoleoasd/EduOJBackend/database/models"

type RegisterResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		models.User `json:"user"`
		Token       string `json:"token"`
	} `json:"data"`
}
