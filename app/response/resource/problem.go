package resource

import "github.com/leoleoasd/EduOJBackend/database/models"

type TestCaseProfileForAdmin struct {
	ID uint `gorm:"primary_key" json:"id"`

	ProblemID uint `sql:"index" json:"problem_id" gorm:"not null"`
	Score     uint `json:"score" gorm:"default:0;not null"` // 0 for 平均分配

	InputFileName  string `json:"input_file_name" gorm:"size:255;default:'';not null"`
	OutputFileName string `json:"output_file_name" gorm:"size:255;default:'';not null"`
}

type TestCaseProfile struct {
	ID uint `gorm:"primary_key" json:"id"`

	ProblemID uint `sql:"index" json:"problem_id" gorm:"not null"`
	Score     uint `json:"score" gorm:"default:0;not null"` // 0 for 平均分配
}

type ProblemProfileForAdmin struct {
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

	TestCases []TestCaseProfileForAdmin `json:"test_cases"`
}

type ProblemProfile struct {
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

	TestCases []TestCaseProfile `json:"test_cases"`
}

func (t *TestCaseProfileForAdmin) Convert(testCase *models.TestCase) {
	t.ID = testCase.ID
	t.ProblemID = testCase.ProblemID
	t.Score = testCase.Score
	t.InputFileName = testCase.InputFileName
	t.OutputFileName = testCase.OutputFileName
}

func (t *TestCaseProfile) Convert(testCase *models.TestCase) {
	t.ID = testCase.ID
	t.ProblemID = testCase.ProblemID
	t.Score = testCase.Score
}

func GetTestCaseProfileForAdmin(testCase *models.TestCase) *TestCaseProfileForAdmin {
	t := TestCaseProfileForAdmin{}
	t.Convert(testCase)
	return &t
}

func GetTestCaseProfile(testCase *models.TestCase) *TestCaseProfile {
	t := TestCaseProfile{}
	t.Convert(testCase)
	return &t
}

func (p *ProblemProfileForAdmin) Convert(problem *models.Problem) {
	p.ID = problem.ID
	p.Name = problem.Name
	p.Description = problem.Description
	p.AttachmentFileName = problem.AttachmentFileName
	p.Public = problem.Public
	p.Privacy = problem.Privacy
	p.MemoryLimit = problem.MemoryLimit
	p.TimeLimit = problem.TimeLimit
	p.LanguageAllowed = problem.LanguageAllowed
	p.CompileEnvironment = problem.CompileEnvironment
	p.CompareScriptID = problem.CompareScriptID

	p.TestCases = make([]TestCaseProfileForAdmin, len(problem.TestCases))
	for i, testCase := range problem.TestCases {
		p.TestCases[i].Convert(&testCase)
	}
}

func (p *ProblemProfile) Convert(problem *models.Problem) {
	p.ID = problem.ID
	p.Name = problem.Name
	p.Description = problem.Description
	p.AttachmentFileName = problem.AttachmentFileName
	p.Public = problem.Public
	p.Privacy = problem.Privacy
	p.MemoryLimit = problem.MemoryLimit
	p.TimeLimit = problem.TimeLimit
	p.LanguageAllowed = problem.LanguageAllowed
	p.CompileEnvironment = problem.CompileEnvironment
	p.CompareScriptID = problem.CompareScriptID

	p.TestCases = make([]TestCaseProfile, len(problem.TestCases))
	for i, testCase := range problem.TestCases {
		p.TestCases[i].Convert(&testCase)
	}
}

func GetProblemProfileForAdmin(problem *models.Problem) *ProblemProfileForAdmin {
	p := ProblemProfileForAdmin{}
	p.Convert(problem)
	return &p
}

func GetProblemProfileForAdminSlice(problems []models.Problem) (profiles []ProblemProfileForAdmin) {
	profiles = make([]ProblemProfileForAdmin, len(problems))
	for i, problem := range problems {
		profiles[i].Convert(&problem)
	}
	return
}

func GetProblemProfile(problem *models.Problem) *ProblemProfile {
	p := ProblemProfile{}
	p.Convert(problem)
	return &p
}
