package request

type LoginRequest struct {
	UsernameOrEmail string `json:"username" form:"username" query:"username" validate:"required,min=5,max=30,username|email"`
	Password        string `json:"password" form:"password" query:"password" validate:"required,min=5,max=30"`
	RememberMe      bool   `json:"remember_me"`
}

type RegisterRequest struct {
	Username string `json:"username" form:"username" query:"username" validate:"required,max=30,min=5,username"`
	Nickname string `json:"nickname" form:"nickname" query:"nickname" validate:"required,max=30,min=1"`
	Email    string `json:"email" form:"email" query:"email" validate:"required,email,max=320,min=5"`
	Password string `json:"password" form:"password" query:"password" validate:"required,max=30,min=5"`
}

type EmailRegistered struct {
	Email string `json:"email" form:"email" query:"email" validate:"required,email,max=320,min=5"`
}
