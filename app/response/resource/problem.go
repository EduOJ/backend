package resource

import (
	"github.com/leoleoasd/EduOJBackend/database/models"
	"strings"
)

type TestCaseForAdmin struct {
	ID uint `json:"id"`

	ProblemID uint `sql:"index" json:"problem_id"`
	Score     uint `json:"score"` // 0 for 平均分配

	InputFileName  string `json:"input_file_name"`
	OutputFileName string `json:"output_file_name"`
}

type TestCase struct {
	ID uint `json:"id"`

	ProblemID uint `sql:"index" json:"problem_id"`
	Score     uint `json:"score"` // 0 for 平均分配
}

type ProblemForAdmin struct {
	ID                 uint   `json:"id"`
	Name               string `sql:"index" json:"name"`
	Description        string `json:"description"`
	AttachmentFileName string `json:"attachment_file_name"`
	Public             bool   `json:"public"`
	Privacy            bool   `json:"privacy"`

	MemoryLimit        uint64   `json:"memory_limit"` // Byte
	TimeLimit          uint     `json:"time_limit"`   // ms
	LanguageAllowed    []string `json:"language_allowed"`
	CompileEnvironment string   `json:"compile_environment"` // E.g.  O2=false
	CompareScriptID    uint     `json:"compare_script_id"`

	TestCases []TestCaseForAdmin `json:"test_cases"`
}

type Problem struct {
	ID                 uint   `json:"id"`
	Name               string `sql:"index" json:"name"`
	Description        string `json:"description"`
	AttachmentFileName string `json:"attachment_file_name"`

	MemoryLimit     uint64   `json:"memory_limit"` // Byte
	TimeLimit       uint     `json:"time_limit"`   // ms
	LanguageAllowed []string `json:"language_allowed"`
	CompareScriptID uint     `json:"compare_script_id"`

	TestCases []TestCase `json:"test_cases"`
}

func (t *TestCaseForAdmin) convert(testCase *models.TestCase) {
	t.ID = testCase.ID
	t.ProblemID = testCase.ProblemID
	t.Score = testCase.Score
	t.InputFileName = testCase.InputFileName
	t.OutputFileName = testCase.OutputFileName
}

func (t *TestCase) convert(testCase *models.TestCase) {
	t.ID = testCase.ID
	t.ProblemID = testCase.ProblemID
	t.Score = testCase.Score
}

func GetTestCaseForAdmin(testCase *models.TestCase) *TestCaseForAdmin {
	t := TestCaseForAdmin{}
	t.convert(testCase)
	return &t
}

func GetTestCase(testCase *models.TestCase) *TestCase {
	t := TestCase{}
	t.convert(testCase)
	return &t
}

func (p *ProblemForAdmin) convert(problem *models.Problem) {
	p.ID = problem.ID
	p.Name = problem.Name
	p.Description = problem.Description
	p.AttachmentFileName = problem.AttachmentFileName
	p.MemoryLimit = problem.MemoryLimit
	p.TimeLimit = problem.TimeLimit
	p.LanguageAllowed = strings.Split(problem.LanguageAllowed, ",")
	p.CompareScriptID = problem.CompareScriptID

	p.Public = problem.Public
	p.Privacy = problem.Privacy
	p.CompileEnvironment = problem.CompileEnvironment

	p.TestCases = make([]TestCaseForAdmin, len(problem.TestCases))
	for i, testCase := range problem.TestCases {
		p.TestCases[i].convert(&testCase)
	}
}

func (p *Problem) convert(problem *models.Problem) {
	p.ID = problem.ID
	p.Name = problem.Name
	p.Description = problem.Description
	p.AttachmentFileName = problem.AttachmentFileName
	p.MemoryLimit = problem.MemoryLimit
	p.TimeLimit = problem.TimeLimit
	p.LanguageAllowed = strings.Split(problem.LanguageAllowed, ",")
	p.CompareScriptID = problem.CompareScriptID

	p.TestCases = make([]TestCase, len(problem.TestCases))
	for i, testCase := range problem.TestCases {
		p.TestCases[i].convert(&testCase)
	}
}

func GetProblemForAdmin(problem *models.Problem) *ProblemForAdmin {
	p := ProblemForAdmin{}
	p.convert(problem)
	return &p
}

func GetProblemForAdminSlice(problems []models.Problem) (profiles []ProblemForAdmin) {
	profiles = make([]ProblemForAdmin, len(problems))
	for i, problem := range problems {
		profiles[i].convert(&problem)
	}
	return
}

func GetProblem(problem *models.Problem) *Problem {
	p := Problem{}
	p.convert(problem)
	return &p
}

func GetProblemSlice(problems []models.Problem) (profiles []Problem) {
	profiles = make([]Problem, len(problems))
	for i, problem := range problems {
		profiles[i].convert(&problem)
	}
	return
}
