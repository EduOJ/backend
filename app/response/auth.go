package response

type RegisterResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		User  UserProfileForMe `json:"user"`
		Token string           `json:"token"`
	} `json:"data"`
}

type LoginResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		User  UserProfileForMe `json:"user"`
		Token string           `json:"token"`
	} `json:"data"`
}
