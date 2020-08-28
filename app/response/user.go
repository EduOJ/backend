package response

type GetUserResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*UserProfile `json:"user"`
	} `json:"data"`
}

type GetUsersResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Users  []UserProfile `json:"users"`
		Total  int           `json:"total"`
		Count  int           `json:"count"`
		Offset int           `json:"offset"`
		Prev   *string       `json:"prev"`
		Next   *string       `json:"next"`
	} `json:"data"`
}

type GetMeResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*UserProfileForMe `json:"user"`
	} `json:"data"`
}

type UpdateMeResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*UserProfileForMe `json:"user"`
	} `json:"data"`
}
