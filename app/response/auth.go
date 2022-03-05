package response

import (
	"github.com/EduOJ/backend/app/response/resource"
)

type RegisterResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		User  resource.UserForAdmin `json:"user"`
		Token string                `json:"token"`
	} `json:"data"`
}

type LoginResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		User  resource.UserForAdmin `json:"user"`
		Token string                `json:"token"`
	} `json:"data"`
}

type UpdateEmailResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.UserForAdmin `json:"user"`
	} `json:"data"`
}

type RequestResetPasswordResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    interface{} `json:"data"`
}

type ResendEmailVerificationResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    interface{} `json:"data"`
}

type EmailVerificationResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    interface{} `json:"data"`
}
