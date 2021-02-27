package request

import "time"

type CreateProblemSetRequest struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required,max=255"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`

	StartAt time.Time `json:"start_at" form:"start_at" query:"start_at" validate:"required"`
	EndAt   time.Time `json:"end_at" form:"end_at" query:"end_at" validate:"required,gtefield=StartAt"`
}

type CloneProblemSetRequest struct {
	SourceClassID      uint `json:"source_class_id" form:"source_class_id" query:"source_class_id" validate:"required"`
	SourceProblemSetID uint `json:"source_problem_set_id" form:"source_problem_set_id" query:"source_problem_set_id" validate:"required"`
}

type GetProblemSetRequest struct {
}

type UpdateProblemSetRequest struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required,max=255"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`

	StartAt time.Time `json:"start_at" form:"start_at" query:"start_at" validate:"required"`
	EndAt   time.Time `json:"end_at" form:"end_at" query:"end_at" validate:"required,gtefield=StartAt"`
}

type AddProblemsInSetRequest struct {
	ProblemIDs []uint `json:"problem_ids" form:"problem_ids" query:"problem_ids" validate:"required"`
}

type DeleteProblemsInSetRequest struct {
	ProblemIDs []uint `json:"problem_ids" form:"problem_ids" query:"problem_ids" validate:"required"`
}

type DeleteProblemSetRequest struct {
}
