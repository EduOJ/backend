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

type GetMeRequest struct {
}

type UpdateMeRequest struct {
	Username string `json:"username" form:"username" query:"username" validate:"required,max=30,min=5,username"`
	Nickname string `json:"nickname" form:"nickname" query:"nickname" validate:"required,max=30,min=1"`
	Email    string `json:"email" form:"email" query:"email" validate:"required,email,max=320,min=5"`
}
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" form:"old_password" query:"old_password" validate:"required,max=30,min=5"`
	NewPassword string `json:"new_password" form:"new_password" query:"new_password" validate:"required,max=30,min=5"`
}

type GetClassesIManageRequest struct {
}

type GetClassesITakeRequest struct {
}

type GetUserProblemInfoRequest struct {
}

type VerifyEmailRequest struct {
	Token string `json:"token" form:"token" query:"token" validate:"required,max=5,min=5"`
}
