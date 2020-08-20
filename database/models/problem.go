package models

import "time"

type TestCase struct {
	ID uint `gorm:"primary_key" json:"id"`

	ProblemID uint `sql:"index" json:"problem_id"`
	Score     uint `json:"score"` // 0 for 平均分配

	InputFileName  string `json:"input_file_name"`
	OutputFileName string `json:"output_file_name"`

	CreatedAt time.Time  `sql:"index" json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

type Problem struct {
	ID                 uint   `gorm:"primary_key" json:"id"`
	Name               string `sql:"index" json:"name"`
	Description        string `gorm:"type:'TEXT'" json:"description"`
	AttachmentFileName string `json:"attachment_file_name"`
	Public             bool   `json:"public" gorm:"default:true"`
	Privacy            bool   `json:"privacy" gorm:"default:true"`

	MemoryLimit        uint   `json:"memory_limit"`        // Byte
	TimeLimit          uint   `json:"time_limt"`           // ms
	LanguageAllowed    string `json:"language_allowed"`    // E.g.    cpp,c,java,python
	CompileEnvironment string `json:"compile_environment"` // E.g.  O2=false
	CompareScriptID    uint   `json:"compare_script_id"`

	TestCases []TestCase `json:"test_cases"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
