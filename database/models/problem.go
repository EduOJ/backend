package models

import "time"

type TestCase struct {
	ID uint `gorm:"primary_key" json:"id"`

	ProblemID uint `sql:"index" json:"problem_id" gorm:"not null"`
	Score     uint `json:"score" gorm:"default:0;not null"` // 0 for 平均分配

	InputFileName  string `json:"input_file_name" gorm:"size:255;default:'';not null"`
	OutputFileName string `json:"output_file_name" gorm:"size:255;default:'';not null"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type Problem struct {
	ID                 uint   `gorm:"primary_key" json:"id"`
	Name               string `sql:"index" json:"name" gorm:"size:255;default:'';not null"`
	Description        string `json:"description;default:'';not null"`
	AttachmentFileName string `json:"attachment_file_name" gorm:"size:255;default:'';not null"`
	Public             bool   `json:"public" gorm:"default:false;not null"`
	Privacy            bool   `json:"privacy" gorm:"default:false;not null"`

	MemoryLimit        uint   `json:"memory_limit" gorm:"default:0;not null"`                   // Byte
	TimeLimit          uint   `json:"time_limit" gorm:"default:0;not null"`                     // ms
	LanguageAllowed    string `json:"language_allowed" gorm:"size:255;default:'';not null"`     // E.g.    cpp,c,java,python
	CompileEnvironment string `json:"compile_environment" gorm:"size:2047;default:'';not null"` // E.g.  O2=false
	CompareScriptID    uint   `json:"compare_script_id" gorm:"default:0;not null"`

	TestCases []TestCase `json:"test_cases"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"deleted_at"`
}
