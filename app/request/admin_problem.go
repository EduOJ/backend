package request

type CreateProblemRequest struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required,max=255"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`
	// attachment_file(optional)
	Public  *bool `json:"public" form:"public" query:"public" validate:"required"`
	Privacy *bool `json:"privacy" form:"privacy" query:"privacy" validate:"required"`

	MemoryLimit        uint64 `json:"memory_limit" form:"memory_limit" query:"memory_limit" validate:"required"`                     // Byte
	TimeLimit          uint   `json:"time_limit" form:"time_limit" query:"time_limit" validate:"required"`                           // ms
	LanguageAllowed    string `json:"language_allowed" form:"language_allowed" query:"language_allowed" validate:"required,max=255"` // E.g.    cpp,c,java,python
	CompileEnvironment string `json:"compile_environment" form:"compile_environment" query:"compile_environment" validate:"max=255"` // E.g.  O2=false
	CompareScriptID    uint   `json:"compare_script_id" form:"compare_script_id" query:"compare_script_id" validate:"required"`
}

type UpdateProblemRequest struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required,max=255"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`
	// attachment_file(optional)
	Public  *bool `json:"public" form:"public" query:"public" validate:"required"`
	Privacy *bool `json:"privacy" form:"privacy" query:"privacy" validate:"required"`

	MemoryLimit        uint64 `json:"memory_limit" form:"memory_limit" query:"memory_limit" validate:"required"`                     // Byte
	TimeLimit          uint   `json:"time_limit" form:"time_limit" query:"time_limit" validate:"required"`                           // ms
	LanguageAllowed    string `json:"language_allowed" form:"language_allowed" query:"language_allowed" validate:"required,max=255"` // E.g.    cpp,c,java,python
	CompileEnvironment string `json:"compile_environment" form:"compile_environment" query:"compile_environment" validate:"max=255"` // E.g.  O2=false
	CompareScriptID    uint   `json:"compare_script_id" form:"compare_script_id" query:"compare_script_id" validate:"required"`
}

type DeleteProblemRequest struct {
}

type CreateTestCaseRequest struct {
	Score uint `json:"score" form:"score" query:"score"` // 0 for 平均分配

	// input_file(required)
	// output_file(required)
}

type GetTestCaseInputFileRequest struct {
}

type GetTestCaseOutputFileRequest struct {
}

type UpdateTestCaseRequest struct {
	Score uint `json:"score" form:"score" query:"score"` // 0 for 平均分配

	// input_file(optional)
	// output_file(optional)
}

type DeleteTestCaseRequest struct {
}

type DeleteTestCasesRequest struct {
}
