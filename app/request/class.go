package request

type CreateClassRequest struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required, max=255"`
	CourseName  string `json:"course_name" form:"course_name" query:"course_name" validate:"required, max=255"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`
}

type GetClassRequest struct {
}

type GetClassesRequest struct {
}

type GetMyClassesRequest struct {
}

type UpdateClassRequest struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required, max=255"`
	CourseName  string `json:"course_name" form:"course_name" query:"course_name" validate:"required, max=255"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`
}

type AddStudentsRequest struct {
	UserIds []uint
}

type RemoveStudentsRequest struct {
	UserIds []uint
}

type DeleteClassRequest struct {
}
