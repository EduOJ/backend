package request

type LoginRequest struct {
	UsernameOrEmail string `json:"username" form:"username" query:"username" validate:"required,min=5,max=30"`
	Password        string `json:"password" form:"password" query:"password" validate:"required,min=5,max=30"`
}

type RegisterRequest struct {
	Username string `json:"username" form:"username" query:"username" validate:"required,max=30,min=5,printascii"`
	Nickname string `json:"nickname" form:"nickname" query:"nickname" validate:"required,max=30,min=5"`
	Email    string `json:"email" form:"email" query:"email" validate:"required,email,max=50,min=5"`
	Password string `json:"password" form:"password" query:"password" validate:"required,max=30,min=5"`
}

type LoggedRequest struct {
	Token string `validate:"required,alpha,len=32"`
}