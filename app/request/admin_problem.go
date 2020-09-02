package request

type AdminCreateProblemRequest struct {
	Name               string `sql:"index" json:"name" validate:"required,max=255"`
	Description        string `json:"description"`
	AttachmentFileName string `json:"attachment_file_name" validate:"required,max=255"`
	Public             *bool  `json:"public"`
	Privacy            *bool  `json:"privacy"`

	MemoryLimit        uint   `json:"memory_limit"`                                 // Byte
	TimeLimit          uint   `json:"time_limit"`                                   // ms
	LanguageAllowed    string `json:"language_allowed" validate:"required,max=255"` // E.g.    cpp,c,java,python
	CompileEnvironment string `json:"compile_environment" validate:"max=255"`       // E.g.  O2=false
	CompareScriptID    uint   `json:"compare_script_id"`
}

type AdminGetProblemRequest struct {
}

type AdminGetProblemsRequest struct {
	Search string `json:"search" form:"search" query:"search"`

	Limit  int `json:"limit" form:"limit" query:"limit" validate:"max=100,min=0"`
	Offset int `json:"offset" form:"offset" query:"offset" validate:"min=0"`

	// OrderBy example: name.DESC
	OrderBy string `json:"order_by" form:"order_by" query:"order_by"`
}

type AdminUpdateProblemRequest struct {
	Name               string `sql:"index" json:"name" validate:"required,max=255"`
	Description        string `json:"description"`
	AttachmentFileName string `json:"attachment_file_name" validate:"required,max=255"`
	Public             *bool  `json:"public"`
	Privacy            *bool  `json:"privacy"`

	MemoryLimit        uint   `json:"memory_limit"`                                 // Byte
	TimeLimit          uint   `json:"time_limit"`                                   // ms
	LanguageAllowed    string `json:"language_allowed" validate:"required,max=255"` // E.g.    cpp,c,java,python
	CompileEnvironment string `json:"compile_environment" validate:"max=255"`       // E.g.  O2=false
	CompareScriptID    uint   `json:"compare_script_id"`
}

type AdminDeleteProblemRequest struct {
}

type AdminCreateTestCase struct {
	ProblemID uint `sql:"index" json:"problem_id" validate:"required"`
	Score     uint `json:"score"` // 0 for 平均分配

	InputFileName  string `json:"input_file_name" validate:"required,max=255"`
	OutputFileName string `json:"output_file_name" validate:"required,max=255"`
}

type AdminGetTestCase struct {
}

type AdminGetTestCases struct {
	Search string `json:"search" form:"search" query:"search"`

	Limit  int `json:"limit" form:"limit" query:"limit" validate:"max=100,min=0"`
	Offset int `json:"offset" form:"offset" query:"offset" validate:"min=0"`

	// OrderBy example: score.DESC
	OrderBy string `json:"order_by" form:"order_by" query:"order_by"`
}

type AdminUpdateTestCase struct {
	ProblemID uint `sql:"index" json:"problem_id" validate:"required"`
	Score     uint `json:"score"` // 0 for 平均分配

	InputFileName  string `json:"input_file_name" validate:"required,max=255"`
	OutputFileName string `json:"output_file_name" validate:"required,max=255"`
}

type AdminDeleteTestCase struct {
}
