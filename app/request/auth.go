package request

type LoginRequest struct {
	// The username or email of the user logging in.
	UsernameOrEmail string `json:"username" form:"username" query:"username" validate:"required,min=5,max=320,username|email" example:"username"`
	// The password of the user logging in.
	Password string `json:"password" form:"password" query:"password" validate:"required,min=5,max=30" example:"password"`
	// If true, the created token will last longer.
	RememberMe bool `json:"remember_me" example:"false"`
}

type RegisterRequest struct {
	Username string `json:"username" form:"username" query:"username" validate:"required,max=30,min=5,username"`
	Nickname string `json:"nickname" form:"nickname" query:"nickname" validate:"required,max=30,min=1"`
	Email    string `json:"email" form:"email" query:"email" validate:"required,email,max=320,min=5"`
	Password string `json:"password" form:"password" query:"password" validate:"required,max=30,min=5"`
}

type UpdateEmailRequest struct {
	Email string `json:"email" form:"email" query:"email" validate:"required,email,max=320,min=5"`
}

type EmailRegisteredRequest struct {
	Email string `json:"email" form:"email" query:"email" validate:"required,email,max=320,min=5"`
}

type RequestResetPasswordRequest struct {
	UsernameOrEmail string `json:"username" form:"username" query:"username" validate:"required,min=5,username|email"`
}

type DoResetPasswordRequest struct {
	UsernameOrEmail string `json:"username" form:"username" query:"username" validate:"required,min=5,username|email"`
	Token           string `json:"token" form:"token" query:"token" validate:"required,max=5,min=5"`
	Password        string `json:"password" form:"password" query:"password" validate:"required,max=30,min=5"`
}
