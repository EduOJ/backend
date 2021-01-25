package resource_test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"testing"
	"time"
)

func createTestCaseForTest(name string, problemId uint, id uint) (testCase models.TestCase) {
	testCase = models.TestCase{
		ID:             id,
		ProblemID:      problemId,
		Score:          id,
		InputFileName:  fmt.Sprintf("test_%s_problem_%d_test_case_%d_input", name, problemId, id),
		OutputFileName: fmt.Sprintf("test_%s_problem_%d_test_case_%d_output", name, problemId, id),
		CreatedAt:      time.Date(int(problemId), 1, 1, 1, 1, 1, 1, time.FixedZone("test_zone", 0)),
		UpdatedAt:      time.Date(int(problemId), 2, 2, 2, 2, 2, 2, time.FixedZone("test_zone", 0)),
		DeletedAt:      nil,
	}
	return
}

func createProblemForTest(name string, id uint, testCaseCount uint) (problem models.Problem) {
	problem = models.Problem{
		ID:                 id,
		Name:               fmt.Sprintf("test_%s_problem_%d", name, id),
		Description:        fmt.Sprintf("test_%s_problem_%d_desc", name, id),
		AttachmentFileName: fmt.Sprintf("test_%s_problem_%d_attachment", name, id),
		Public:             true,
		Privacy:            false,
		MemoryLimit:        1024,
		TimeLimit:          1000,
		LanguageAllowed:    fmt.Sprintf("test_%s_language_allowed_%d,test_language", name, id),
		CompileEnvironment: fmt.Sprintf("test_%s_compile_environment_%d", name, id),
		CompareScriptID:    1,
		TestCases:          make([]models.TestCase, testCaseCount),
		CreatedAt:          time.Date(int(id), 1, 1, 1, 1, 1, 1, time.FixedZone("test_zone", 0)),
		UpdatedAt:          time.Date(int(id), 2, 2, 2, 2, 2, 2, time.FixedZone("test_zone", 0)),
		DeletedAt:          nil,
	}
	for i := range problem.TestCases {
		problem.TestCases[i] = createTestCaseForTest(name, id, uint(i))
	}
	return
}

func TestGetTestCaseAndGetTestCaseForAdmin(t *testing.T) {
	testCase := createTestCaseForTest("get_test_case", 0, 0)
	t.Run("testGetTestCase", func(t *testing.T) {
		actualT := resource.GetTestCase(&testCase)
		expectedT := resource.TestCase{
			ID:        0,
			ProblemID: 0,
			Score:     0,
		}
		assert.Equal(t, expectedT, actualT)
	})
	t.Run("testGetTestCaseForAdmin", func(t *testing.T) {
		actualT := resource.GetTestCaseForAdmin(&testCase)
		expectedT := resource.TestCaseForAdmin{
			ID:             0,
			ProblemID:      0,
			Score:          0,
			InputFileName:  "test_get_test_case_problem_0_test_case_0_input",
			OutputFileName: "test_get_test_case_problem_0_test_case_0_output",
		}
		assert.Equal(t, expectedT, actualT)
	})
}

func TestGetProblemAndGetProblemForAdmin(t *testing.T) {
	problem := createProblemForTest("get_problem", 0, 2)
	t.Run("testGetProblem", func(t *testing.T) {
		actualP := resource.GetProblem(&problem)
		expectedP := resource.Problem{
			ID:                 0,
			Name:               "test_get_problem_problem_0",
			Description:        "test_get_problem_problem_0_desc",
			AttachmentFileName: "test_get_problem_problem_0_attachment",
			MemoryLimit:        1024,
			TimeLimit:          1000,
			LanguageAllowed:    []string{"test_get_problem_language_allowed_0", "test_language"},
			CompareScriptID:    1,
			TestCases: []resource.TestCase{
				{
					ID:        0,
					ProblemID: 0,
					Score:     0,
				}, {
					ID:        1,
					ProblemID: 0,
					Score:     1,
				},
			},
		}
		assert.Equal(t, &expectedP, actualP)
	})
	t.Run("testGetProblemForAdmin", func(t *testing.T) {
		actualP := resource.GetProblemForAdmin(&problem)
		expectedP := resource.ProblemForAdmin{
			ID:                 0,
			Name:               "test_get_problem_problem_0",
			Description:        "test_get_problem_problem_0_desc",
			AttachmentFileName: "test_get_problem_problem_0_attachment",
			Public:             true,
			Privacy:            false,
			MemoryLimit:        1024,
			TimeLimit:          1000,
			LanguageAllowed:    []string{"test_get_problem_language_allowed_0", "test_language"},
			CompileEnvironment: "test_get_problem_compile_environment_0",
			CompareScriptID:    1,
			TestCases: []resource.TestCaseForAdmin{
				{
					ID:             0,
					ProblemID:      0,
					Score:          0,
					InputFileName:  "test_get_problem_problem_0_test_case_0_input",
					OutputFileName: "test_get_problem_problem_0_test_case_0_output",
				}, {
					ID:             1,
					ProblemID:      0,
					Score:          1,
					InputFileName:  "test_get_problem_problem_0_test_case_1_input",
					OutputFileName: "test_get_problem_problem_0_test_case_1_output",
				},
			},
		}
		assert.Equal(t, &expectedP, actualP)
	})
}

func TestGetProblemSliceAndGetProblemForAdminSlice(t *testing.T) {
	problem1 := createProblemForTest("get_problem_slice", 1, 1)
	problem2 := createProblemForTest("get_problem_slice", 2, 2)
	t.Run("testGetProblemSlice", func(t *testing.T) {
		actualPS := resource.GetProblemSlice([]models.Problem{problem1, problem2})
		expectedPS := []resource.Problem{
			{
				ID:                 1,
				Name:               "test_get_problem_slice_problem_1",
				Description:        "test_get_problem_slice_problem_1_desc",
				AttachmentFileName: "test_get_problem_slice_problem_1_attachment",
				MemoryLimit:        1024,
				TimeLimit:          1000,
				LanguageAllowed:    []string{"test_get_problem_slice_language_allowed_1", "test_language"},
				CompareScriptID:    1,
				TestCases: []resource.TestCase{
					{
						ID:        0,
						ProblemID: 1,
						Score:     0,
					},
				},
			}, {
				ID:                 2,
				Name:               "test_get_problem_slice_problem_2",
				Description:        "test_get_problem_slice_problem_2_desc",
				AttachmentFileName: "test_get_problem_slice_problem_2_attachment",
				MemoryLimit:        1024,
				TimeLimit:          1000,
				LanguageAllowed:    []string{"test_get_problem_slice_language_allowed_2", "test_language"},
				CompareScriptID:    1,
				TestCases: []resource.TestCase{
					{
						ID:        0,
						ProblemID: 2,
						Score:     0,
					}, {
						ID:        1,
						ProblemID: 2,
						Score:     1,
					},
				},
			},
		}
		assert.Equal(t, expectedPS, actualPS)
	})
	t.Run("testGetProblemForAdminSlice", func(t *testing.T) {
		actualPS := resource.GetProblemForAdminSlice([]models.Problem{problem1, problem2})
		expectedPS := []resource.ProblemForAdmin{
			{
				ID:                 1,
				Name:               "test_get_problem_slice_problem_1",
				Description:        "test_get_problem_slice_problem_1_desc",
				AttachmentFileName: "test_get_problem_slice_problem_1_attachment",
				Public:             true,
				Privacy:            false,
				MemoryLimit:        1024,
				TimeLimit:          1000,
				LanguageAllowed:    []string{"test_get_problem_slice_language_allowed_1", "test_language"},
				CompileEnvironment: "test_get_problem_slice_compile_environment_1",
				CompareScriptID:    1,
				TestCases: []resource.TestCaseForAdmin{
					{
						ID:             0,
						ProblemID:      1,
						Score:          0,
						InputFileName:  "test_get_problem_slice_problem_1_test_case_0_input",
						OutputFileName: "test_get_problem_slice_problem_1_test_case_0_output",
					},
				},
			}, {
				ID:                 2,
				Name:               "test_get_problem_slice_problem_2",
				Description:        "test_get_problem_slice_problem_2_desc",
				AttachmentFileName: "test_get_problem_slice_problem_2_attachment",
				Public:             true,
				Privacy:            false,
				MemoryLimit:        1024,
				TimeLimit:          1000,
				LanguageAllowed:    []string{"test_get_problem_slice_language_allowed_2", "test_language"},
				CompileEnvironment: "test_get_problem_slice_compile_environment_2",
				CompareScriptID:    1,
				TestCases: []resource.TestCaseForAdmin{
					{
						ID:             0,
						ProblemID:      2,
						Score:          0,
						InputFileName:  "test_get_problem_slice_problem_2_test_case_0_input",
						OutputFileName: "test_get_problem_slice_problem_2_test_case_0_output",
					}, {
						ID:             1,
						ProblemID:      2,
						Score:          1,
						InputFileName:  "test_get_problem_slice_problem_2_test_case_1_input",
						OutputFileName: "test_get_problem_slice_problem_2_test_case_1_output",
					},
				},
			},
		}
		assert.Equal(t, expectedPS, actualPS)
	})
}
