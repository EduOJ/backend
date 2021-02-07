package response

import (
	"github.com/leoleoasd/EduOJBackend/database/models"
	"time"
)

type GetTaskResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		RunID              uint
		Language           models.Language `json:"language"`
		TestCaseID         uint            `json:"test_case_id"`
		InputFile          string          `json:"input_file"`  // pre-signed url
		OutputFile         string          `json:"output_file"` // same as above
		CodeFile           string          `json:"code_file"`
		TestCaseUpdatedAt  time.Time       `json:"test_case_updated_at"`
		MemoryLimit        uint64          `json:"memory_limit"`        // Byte
		TimeLimit          uint            `json:"time_limit"`          // ms
		CompileEnvironment string          `json:"compile_environment"` // E.g.  O2=false
		CompareScriptID    uint            `json:"compare_script_id"`
	} `json:"data"`
}
