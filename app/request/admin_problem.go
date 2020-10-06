package request

type AdminCreateProblemRequest struct {
	Name               string `json:"name" form:"name" query:"name" validate:"required,max=255"`
	Description        string `json:"description" form:"description" query:"description"`
	AttachmentFileName string `json:"attachment_file_name" form:"attachment_file_name" query:"attachment_file_name" validate:"required,max=255"`
	Public             *bool  `json:"public" form:"public" query:"public"`
	Privacy            *bool  `json:"privacy" form:"privacy" query:"privacy"`

	MemoryLimit        uint   `json:"memory_limit" form:"memory_limit" query:"memory_limit"`                                         // Byte
	TimeLimit          uint   `json:"time_limit" form:"time_limit" query:"time_limit"`                                               // ms
	LanguageAllowed    string `json:"language_allowed" form:"language_allowed" query:"language_allowed" validate:"required,max=255"` // E.g.    cpp,c,java,python
	CompileEnvironment string `json:"compile_environment" form:"compile_environment" query:"compile_environment" validate:"max=255"` // E.g.  O2=false
	CompareScriptID    uint   `json:"compare_script_id" form:"compare_script_id" query:"compare_script_id"`
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
	Name               string `json:"name" form:"name" query:"name" validate:"required,max=255"`
	Description        string `json:"description" form:"description" query:"description"`
	AttachmentFileName string `json:"attachment_file_name" form:"attachment_file_name" query:"attachment_file_name" validate:"required,max=255"`
	Public             *bool  `json:"public" form:"public" query:"public"`
	Privacy            *bool  `json:"privacy" form:"privacy" query:"privacy"`

	MemoryLimit        uint   `json:"memory_limit" form:"memory_limit" query:"memory_limit"`                                         // Byte
	TimeLimit          uint   `json:"time_limit" form:"time_limit" query:"time_limit"`                                               // ms
	LanguageAllowed    string `json:"language_allowed" form:"language_allowed" query:"language_allowed" validate:"required,max=255"` // E.g.    cpp,c,java,python
	CompileEnvironment string `json:"compile_environment" form:"compile_environment" query:"compile_environment" validate:"max=255"` // E.g.  O2=false
	CompareScriptID    uint   `json:"compare_script_id" form:"compare_script_id" query:"compare_script_id"`
}

type AdminDeleteProblemRequest struct {
}

type AdminCreateTestCase struct {
	Score uint `json:"score" form:"score" query:"score"` // 0 for 平均分配

	// TODO: replace to file blob.
	InputFileName  string `json:"input_file_name" form:"input_file_name" query:"input_file_name" validate:"required,max=255"`
	OutputFileName string `json:"output_file_name" form:"output_file_name" query:"output_file_name" validate:"required,max=255"`
}

type AdminUpdateTestCase struct {
	Score uint `json:"score" form:"score" query:"score"` // 0 for 平均分配
}

type AdminDeleteTestCase struct {
}

type AdminDeleteTestCases struct {
}
