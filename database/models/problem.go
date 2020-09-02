package models

import "time"

type TestCase struct {
	ID uint `gorm:"primary_key" json:"id"`

	ProblemID uint `sql:"index" json:"problem_id"`
	Score     uint `json:"score"` // 0 for 平均分配

	InputFileName  string `json:"input_file_name" gorm:"size:255"`
	OutputFileName string `json:"output_file_name" gorm:"size:255"`

	CreatedAt time.Time  `sql:"index" json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

type Problem struct {
	ID                 uint   `gorm:"primary_key" json:"id"`
	Name               string `sql:"index" json:"name" gorm:"size:255"`
	Description        string `json:"description"`
	AttachmentFileName string `json:"attachment_file_name" gorm:"size:255"`
	Public             bool   `json:"public"`
	Privacy            bool   `json:"privacy"`

	MemoryLimit        uint   `json:"memory_limit"`                         // Byte
	TimeLimit          uint   `json:"time_limit"`                           // ms
	LanguageAllowed    string `json:"language_allowed" gorm:"size:255"`     // E.g.    cpp,c,java,python
	CompileEnvironment string `json:"compile_environment" gorm:"size:2047"` // E.g.  O2=false
	CompareScriptID    uint   `json:"compare_script_id"`

	TestCases []TestCase `json:"test_cases"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
