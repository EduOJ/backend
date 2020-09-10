package response

import (
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
)

type RegisterResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		User  resource.UserProfileForAdmin `json:"user"`
		Token string                       `json:"token"`
	} `json:"data"`
}

type LoginResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		User  resource.UserProfileForAdmin `json:"user"`
		Token string                       `json:"token"`
	} `json:"data"`
}
