package response

import (
	"github.com/EduOJ/backend/app/response/resource"
)

type GetUserResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.User `json:"user"`
	} `json:"data"`
}

type GetUsersResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Users  []resource.User `json:"users"`
		Total  int             `json:"total"`
		Count  int             `json:"count"`
		Offset int             `json:"offset"`
		Prev   *string         `json:"prev"`
		Next   *string         `json:"next"`
	} `json:"data"`
}

type GetMeResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.UserForAdmin `json:"user"`
	} `json:"data"`
}

type UpdateMeResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.UserForAdmin `json:"user"`
	} `json:"data"`
}

type EditPreferedNoticeWay struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.UserForAdmin `json:"user"`
	} `json:"data"`
}

type GetClassesIManageResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Classes []resource.Class `json:"classes"`
	} `json:"data"`
}

type GetClassesITakeResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Classes []resource.Class `json:"classes"`
	} `json:"data"`
}

type GetUserProblemInfoResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		TriedCount  int `json:"tried_count"`
		PassedCount int `json:"passed_count"`
		Rank        int `json:"rank"` // TODO: develop this
	} `json:"data"`
}
