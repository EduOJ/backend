package request

type CreateClassRequest struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required,max=255"`
	CourseName  string `json:"course_name" form:"course_name" query:"course_name" validate:"required,max=255"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`
}

type GetClassRequest struct {
}

type UpdateClassRequest struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required,max=255"`
	CourseName  string `json:"course_name" form:"course_name" query:"course_name" validate:"required,max=255"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`
}

type AddStudentsRequest struct {
	UserIds []uint `json:"user_ids" form:"user_ids" query:"user_ids" validate:"required"`
}

type DeleteStudentsRequest struct {
	UserIds []uint `json:"user_ids" form:"user_ids" query:"user_ids" validate:"required"`
}

type RefreshInviteCodeRequest struct {
}

type JoinClassRequest struct {
	InviteCode string `json:"invite_code" form:"invite_code" query:"invite_code" validate:"required,max=255"`
}

type DeleteClassRequest struct {
}
