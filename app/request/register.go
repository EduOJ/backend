package request

type UserRequest struct {
	Username string `validate:"required,max=30,min=5"`
	Nickname string `validate:"required,max=30,min=5"`
	Email    string `validate:"required,email,max=30,min=5"`
	Password string `validate:"required,max=30,min=5"`
}
