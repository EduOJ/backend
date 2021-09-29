package resource

import (
	"database/sql"
	"encoding/json"
	"github.com/EduOJ/backend/database/models"
)

type TestCaseForAdmin struct {
	ID uint `json:"id"`

	ProblemID uint `sql:"index" json:"problem_id"`
	Score     uint `json:"score"` // 0 for 平均分配
	Sample    bool `json:"sample"`

	InputFileName  string `json:"input_file_name"`
	OutputFileName string `json:"output_file_name"`
}

type TestCase struct {
	ID uint `json:"id"`

	ProblemID uint `sql:"index" json:"problem_id"`
	Score     uint `json:"score"` // 0 for 平均分配
	Sample    bool `json:"sample"`
}

type ProblemForAdmin struct {
	ID                 uint   `json:"id"`
	Name               string `sql:"index" json:"name"`
	Description        string `json:"description"`
	AttachmentFileName string `json:"attachment_file_name"`
	Public             bool   `json:"public"`
	Privacy            bool   `json:"privacy"`

	MemoryLimit       uint64   `json:"memory_limit"` // Byte
	TimeLimit         uint     `json:"time_limit"`   // ms
	LanguageAllowed   []string `json:"language_allowed"`
	BuildArg          string   `json:"build_arg"` // E.g.  O2=false
	CompareScriptName string   `json:"compare_script_name"`

	TestCases []TestCaseForAdmin `json:"test_cases"`
	Tags      []Tag              `json:"tags"`
}

type Problem struct {
	ID                 uint   `json:"id"`
	Name               string `sql:"index" json:"name"`
	Description        string `json:"description"`
	AttachmentFileName string `json:"attachment_file_name"`

	MemoryLimit       uint64   `json:"memory_limit"` // Byte
	TimeLimit         uint     `json:"time_limit"`   // ms
	LanguageAllowed   []string `json:"language_allowed"`
	CompareScriptName string   `json:"compare_script_name"`

	TestCases []TestCase `json:"test_cases"`
	Tags      []Tag      `json:"tags"`
}

type Tag struct {
	Name string
}

func (t *Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Name)
}

func (t *Tag) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &t.Name)
}

func (t *Tag) Convert(tt *models.Tag) {
	t.Name = tt.Name
}

type ProblemSummary struct {
	ID                 uint   `json:"id"`
	Name               string `sql:"index" json:"name"`
	AttachmentFileName string `json:"attachment_file_name"`
	Passed             bool   `json:"passed"`

	MemoryLimit       uint64   `json:"memory_limit"` // Byte
	TimeLimit         uint     `json:"time_limit"`   // ms
	LanguageAllowed   []string `json:"language_allowed"`
	CompareScriptName string   `json:"compare_script_name"`
	Tags              []Tag    `json:"tags"`
}

type ProblemSummaryForAdmin struct {
	ID                 uint   `json:"id"`
	Name               string `sql:"index" json:"name"`
	AttachmentFileName string `json:"attachment_file_name"`
	Public             bool   `json:"public"`
	Privacy            bool   `json:"privacy"`
	Passed             bool   `json:"passed"`

	MemoryLimit       uint64   `json:"memory_limit"` // Byte
	TimeLimit         uint     `json:"time_limit"`   // ms
	LanguageAllowed   []string `json:"language_allowed"`
	BuildArg          string   `json:"build_arg"` // E.g.  O2=false
	CompareScriptName string   `json:"compare_script_name"`

	Tags []Tag `json:"tags"`
}

func (t *TestCaseForAdmin) convert(testCase *models.TestCase) {
	t.ID = testCase.ID
	t.ProblemID = testCase.ProblemID
	t.Score = testCase.Score
	t.Sample = testCase.Sample
	t.InputFileName = testCase.InputFileName
	t.OutputFileName = testCase.OutputFileName
}

func (t *TestCase) convert(testCase *models.TestCase) {
	t.ID = testCase.ID
	t.ProblemID = testCase.ProblemID
	t.Score = testCase.Score
	t.Sample = testCase.Sample
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
	p.LanguageAllowed = problem.LanguageAllowed
	p.CompareScriptName = problem.CompareScriptName

	p.Public = problem.Public
	p.Privacy = problem.Privacy
	p.BuildArg = problem.BuildArg

	p.Tags = make([]Tag, len(problem.Tags))

	for i, t := range problem.Tags {
		p.Tags[i].Convert(&t)
	}

	p.TestCases = make([]TestCaseForAdmin, len(problem.TestCases))
	for i, testCase := range problem.TestCases {
		p.TestCases[i].convert(&testCase)
	}

}

func (p *ProblemSummaryForAdmin) convert(problem *models.Problem, passed sql.NullBool) {
	p.ID = problem.ID
	p.Name = problem.Name
	p.AttachmentFileName = problem.AttachmentFileName
	p.MemoryLimit = problem.MemoryLimit
	p.TimeLimit = problem.TimeLimit
	p.LanguageAllowed = problem.LanguageAllowed
	p.CompareScriptName = problem.CompareScriptName

	p.Public = problem.Public
	p.Privacy = problem.Privacy
	p.BuildArg = problem.BuildArg
	p.Passed = passed.Bool

	p.Tags = make([]Tag, len(problem.Tags))

	for i, t := range problem.Tags {
		p.Tags[i].Convert(&t)
	}
}

func (p *Problem) convert(problem *models.Problem) {
	p.ID = problem.ID
	p.Name = problem.Name
	p.Description = problem.Description
	p.AttachmentFileName = problem.AttachmentFileName
	p.MemoryLimit = problem.MemoryLimit
	p.TimeLimit = problem.TimeLimit
	p.LanguageAllowed = problem.LanguageAllowed
	p.CompareScriptName = problem.CompareScriptName

	p.Tags = make([]Tag, len(problem.Tags))

	for i, t := range problem.Tags {
		p.Tags[i].Convert(&t)
	}

	p.TestCases = make([]TestCase, len(problem.TestCases))
	for i, testCase := range problem.TestCases {
		p.TestCases[i].convert(&testCase)
	}
}

func (p *ProblemSummary) convert(problem *models.Problem, passed sql.NullBool) {
	p.ID = problem.ID
	p.Name = problem.Name
	p.AttachmentFileName = problem.AttachmentFileName
	p.MemoryLimit = problem.MemoryLimit
	p.TimeLimit = problem.TimeLimit
	p.LanguageAllowed = problem.LanguageAllowed
	p.CompareScriptName = problem.CompareScriptName
	p.Passed = passed.Bool

	p.Tags = make([]Tag, len(problem.Tags))

	for i, t := range problem.Tags {
		p.Tags[i].Convert(&t)
	}
}

func GetProblemForAdmin(problem *models.Problem) *ProblemForAdmin {
	p := ProblemForAdmin{}
	p.convert(problem)
	return &p
}

func GetProblemForAdminSlice(problems []*models.Problem) (profiles []ProblemForAdmin) {
	profiles = make([]ProblemForAdmin, len(problems))
	for i, problem := range problems {
		profiles[i].convert(problem)
	}
	return
}

func GetProblemSummaryForAdmin(problem *models.Problem, passed sql.NullBool) *ProblemSummaryForAdmin {
	p := ProblemSummaryForAdmin{}
	p.convert(problem, passed)
	return &p
}

func GetProblemSummaryForAdminSlice(problems []*models.Problem, passed []sql.NullBool) (summaries []ProblemSummaryForAdmin) {
	summaries = make([]ProblemSummaryForAdmin, len(problems))
	for i, problem := range problems {
		summaries[i].convert(problem, passed[i])
	}
	return
}

func GetProblem(problem *models.Problem) *Problem {
	p := Problem{}
	p.convert(problem)
	return &p
}

func GetProblemSlice(problems []*models.Problem) (profiles []Problem) {
	profiles = make([]Problem, len(problems))
	for i, problem := range problems {
		profiles[i].convert(problem)
	}
	return
}

func GetProblemSummary(problem *models.Problem, passed sql.NullBool) *ProblemSummary {
	p := ProblemSummary{}
	p.convert(problem, passed)
	return &p
}

func GetProblemSummarySlice(problems []*models.Problem, passed []sql.NullBool) (summaries []ProblemSummary) {
	summaries = make([]ProblemSummary, len(problems))
	for i, problem := range problems {
		summaries[i].convert(problem, passed[i])
	}
	return
}
