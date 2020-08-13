package request

type GetUserRequest struct {
}

type GetUsersRequest struct {
	Search string `json:"search" form:"search" query:"search"`

	Limit  int `json:"limit" form:"limit" query:"limit" validate:"max=100,min=0"`
	Offset int `json:"offset" form:"offset" query:"offset" validate:"min=0"`

	// OrderBy example: username.DESC
	OrderBy string `json:"order_by" form:"order_by" query:"order_by"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" form:"old_password" query:"old_password" validate:"required,max=30,min=5"`
	NewPassword string `json:"new_password" form:"new_password" query:"new_password" validate:"required,max=30,min=5"`
}
