package response

import (
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
)

type GetUserResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.UserProfile `json:"user"`
	} `json:"data"`
}

type GetUsersResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Users  []resource.UserProfile `json:"users"`
		Total  int                    `json:"total"`
		Count  int                    `json:"count"`
		Offset int                    `json:"offset"`
		Prev   *string                `json:"prev"`
		Next   *string                `json:"next"`
	} `json:"data"`
}

type GetMeResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.UserProfileForAdmin `json:"user"`
	} `json:"data"`
}

type UpdateMeResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.UserProfileForAdmin `json:"user"`
	} `json:"data"`
}
