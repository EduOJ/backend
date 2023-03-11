package models

import (
	"fmt"
	"hash/fnv"
	"testing"
	"time"

	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/database"
	"github.com/stretchr/testify/assert"
)

func hashStringToTime(s string) time.Time {
	h := fnv.New32()
	if _, err := h.Write([]byte(s)); err != nil {
		panic(err)
	}
	return time.Unix(int64(h.Sum32()), 0).UTC()
}

func createProblemForTest(t *testing.T, name string, id uint) *Problem {
	problem := Problem{
		Name:        fmt.Sprintf("%s_%d_name", name, id),
		Description: fmt.Sprintf("%s_%d_description", name, id),
		TestCases: []TestCase{
			{
				Score:          10,
				Sample:         true,
				InputFileName:  fmt.Sprintf("%s_%d_1.in", name, id),
				OutputFileName: fmt.Sprintf("%s_%d_1.out", name, id),
			},
			{
				Score:          20,
				Sample:         false,
				InputFileName:  fmt.Sprintf("%s_%d_2.in", name, id),
				OutputFileName: fmt.Sprintf("%s_%d_2.out", name, id),
			},
		},
		LanguageAllowed: database.StringArray([]string{""}),
	}
	assert.NoError(t, base.DB.Create(&problem).Error)
	return &problem
}

func TestAddProblemsAndDeleteProblemsByID(t *testing.T) {
	t.Parallel()

	t.Run("AddSuccess", func(t *testing.T) {
		problem1 := createProblemForTest(t, "add_success", 1)
		problem2 := createProblemForTest(t, "add_success", 2)
		problem3 := createProblemForTest(t, "add_success", 3)
		problem4 := createProblemForTest(t, "add_success", 4)
		t.Parallel()
		problemSet := ProblemSet{
			Name:        "test_add_students_success_problem_set_name",
			Description: "test_add_students_success_problem_set_description",
			Problems: []*Problem{
				problem1,
				problem2,
			},
			StartTime: hashStringToTime("test_add_students_success_problem_set_start_time"),
			EndTime:   hashStringToTime("test_add_students_success_problem_set_end_time"),
		}
		assert.NoError(t, base.DB.Create(&problemSet).Error)
		assert.NoError(t, problemSet.AddProblems([]uint{
			problem3.ID,
			problem4.ID,
		}))
		assert.Equal(t, ProblemSet{
			ID:          problemSet.ID,
			Name:        "test_add_students_success_problem_set_name",
			Description: "test_add_students_success_problem_set_description",
			Problems: []*Problem{
				problem1,
				problem2,
				problem3,
				problem4,
			},
			StartTime: hashStringToTime("test_add_students_success_problem_set_start_time"),
			EndTime:   hashStringToTime("test_add_students_success_problem_set_end_time"),
			CreatedAt: problemSet.CreatedAt,
			UpdatedAt: problemSet.UpdatedAt,
		}, problemSet)
	})
	t.Run("AddInEmptySet", func(t *testing.T) {
		problem1 := createProblemForTest(t, "add_in_empty_set", 1)
		problem2 := createProblemForTest(t, "add_in_empty_set", 2)
		t.Parallel()
		problemSet := ProblemSet{
			Name:        "test_add_students_in_empty_set_problem_set_name",
			Description: "test_add_students_in_empty_set_problem_set_description",
			Problems:    []*Problem{},
			StartTime:   hashStringToTime("test_add_students_in_empty_set_problem_set_start_time"),
			EndTime:     hashStringToTime("test_add_students_in_empty_set_problem_set_end_time"),
		}
		assert.NoError(t, base.DB.Create(&problemSet).Error)
		assert.NoError(t, problemSet.AddProblems([]uint{
			problem1.ID,
			problem2.ID,
		}))
		assert.Equal(t, ProblemSet{
			ID:          problemSet.ID,
			Name:        "test_add_students_in_empty_set_problem_set_name",
			Description: "test_add_students_in_empty_set_problem_set_description",
			Problems: []*Problem{
				problem1,
				problem2,
			},
			StartTime: hashStringToTime("test_add_students_in_empty_set_problem_set_start_time"),
			EndTime:   hashStringToTime("test_add_students_in_empty_set_problem_set_end_time"),
			CreatedAt: problemSet.CreatedAt,
			UpdatedAt: problemSet.UpdatedAt,
		}, problemSet)
	})
	t.Run("AddNothing", func(t *testing.T) {
		problem1 := createProblemForTest(t, "add_nothing", 1)
		problem2 := createProblemForTest(t, "add_nothing", 2)
		t.Parallel()
		problemSet := ProblemSet{
			Name:        "test_add_students_nothing_problem_set_name",
			Description: "test_add_students_nothing_problem_set_description",
			Problems:    []*Problem{},
			StartTime:   hashStringToTime("test_add_students_nothing_problem_set_start_time"),
			EndTime:     hashStringToTime("test_add_students_nothing_problem_set_end_time"),
		}
		assert.NoError(t, base.DB.Create(&problemSet).Error)
		assert.NoError(t, problemSet.AddProblems([]uint{
			problem1.ID,
			problem2.ID,
		}))
		assert.Equal(t, ProblemSet{
			ID:          problemSet.ID,
			Name:        "test_add_students_nothing_problem_set_name",
			Description: "test_add_students_nothing_problem_set_description",
			Problems: []*Problem{
				problem1,
				problem2,
			},
			StartTime: hashStringToTime("test_add_students_nothing_problem_set_start_time"),
			EndTime:   hashStringToTime("test_add_students_nothing_problem_set_end_time"),
			CreatedAt: problemSet.CreatedAt,
			UpdatedAt: problemSet.UpdatedAt,
		}, problemSet)
	})
	t.Run("AddExistingInSet", func(t *testing.T) {
		problem1 := createProblemForTest(t, "add_existing_in_set", 1)
		problem2 := createProblemForTest(t, "add_existing_in_set", 2)
		problem3 := createProblemForTest(t, "add_existing_in_set", 3)
		t.Parallel()
		problemSet := ProblemSet{
			Name:        "test_add_students_existing_in_set_problem_set_name",
			Description: "test_add_students_existing_in_set_problem_set_description",
			Problems: []*Problem{
				problem1,
				problem2,
			},
			StartTime: hashStringToTime("test_add_students_existing_in_set_problem_set_start_time"),
			EndTime:   hashStringToTime("test_add_students_existing_in_set_problem_set_end_time"),
		}
		assert.NoError(t, base.DB.Create(&problemSet).Error)
		assert.NoError(t, problemSet.AddProblems([]uint{
			problem2.ID,
			problem3.ID,
		}))
		assert.Equal(t, ProblemSet{
			ID:          problemSet.ID,
			Name:        "test_add_students_existing_in_set_problem_set_name",
			Description: "test_add_students_existing_in_set_problem_set_description",
			Problems: []*Problem{
				problem1,
				problem2,
				problem3,
			},
			StartTime: hashStringToTime("test_add_students_existing_in_set_problem_set_start_time"),
			EndTime:   hashStringToTime("test_add_students_existing_in_set_problem_set_end_time"),
			CreatedAt: problemSet.CreatedAt,
			UpdatedAt: problemSet.UpdatedAt,
		}, problemSet)
	})
	t.Run("AddNonExisting", func(t *testing.T) {
		problem1 := createProblemForTest(t, "add_non_exist", 1)
		problem2 := createProblemForTest(t, "add_non_exist", 2)
		problem3 := createProblemForTest(t, "add_non_exist", 3)
		t.Parallel()
		problemSet := ProblemSet{
			Name:        "test_add_students_non_existing_problem_set_name",
			Description: "test_add_students_non_existing_problem_set_description",
			Problems: []*Problem{
				problem1,
				problem2,
			},
			StartTime: hashStringToTime("test_add_students_non_existing_problem_set_start_time"),
			EndTime:   hashStringToTime("test_add_students_non_existing_problem_set_end_time"),
		}
		assert.NoError(t, base.DB.Create(&problemSet).Error)
		assert.NoError(t, problemSet.AddProblems([]uint{
			0,
			problem3.ID,
		}))
		assert.Equal(t, ProblemSet{
			ID:          problemSet.ID,
			Name:        "test_add_students_non_existing_problem_set_name",
			Description: "test_add_students_non_existing_problem_set_description",
			Problems: []*Problem{
				problem1,
				problem2,
				problem3,
			},
			StartTime: hashStringToTime("test_add_students_non_existing_problem_set_start_time"),
			EndTime:   hashStringToTime("test_add_students_non_existing_problem_set_end_time"),
			CreatedAt: problemSet.CreatedAt,
			UpdatedAt: problemSet.UpdatedAt,
		}, problemSet)
	})
	t.Run("DeleteSuccess", func(t *testing.T) {
		problem1 := createProblemForTest(t, "delete_success", 1)
		problem2 := createProblemForTest(t, "delete_success", 2)
		problem3 := createProblemForTest(t, "delete_success", 3)
		t.Parallel()
		problemSet := ProblemSet{
			Name:        "test_delete_students_success_problem_set_name",
			Description: "test_delete_students_success_problem_set_description",
			Problems: []*Problem{
				problem1,
				problem2,
				problem3,
			},
			StartTime: hashStringToTime("test_delete_students_success_problem_set_start_time"),
			EndTime:   hashStringToTime("test_delete_students_success_problem_set_end_time"),
		}
		assert.NoError(t, base.DB.Create(&problemSet).Error)
		assert.NoError(t, problemSet.DeleteProblems([]uint{
			problem2.ID,
			problem3.ID,
		}))
		assert.Equal(t, ProblemSet{
			ID:          problemSet.ID,
			Name:        "test_delete_students_success_problem_set_name",
			Description: "test_delete_students_success_problem_set_description",
			Problems: []*Problem{
				problem1,
			},
			StartTime: hashStringToTime("test_delete_students_success_problem_set_start_time"),
			EndTime:   hashStringToTime("test_delete_students_success_problem_set_end_time"),
			CreatedAt: problemSet.CreatedAt,
			UpdatedAt: problemSet.UpdatedAt,
		}, problemSet)
	})
	t.Run("DeleteNothing", func(t *testing.T) {
		problem1 := createProblemForTest(t, "delete_nothing", 1)
		problem2 := createProblemForTest(t, "delete_nothing", 2)
		problem3 := createProblemForTest(t, "delete_nothing", 3)
		t.Parallel()
		problemSet := ProblemSet{
			Name:        "test_delete_students_nothing_problem_set_name",
			Description: "test_delete_students_nothing_problem_set_description",
			Problems: []*Problem{
				problem1,
				problem2,
				problem3,
			},
			StartTime: hashStringToTime("test_delete_students_nothing_problem_set_start_time"),
			EndTime:   hashStringToTime("test_delete_students_nothing_problem_set_end_time"),
		}
		assert.NoError(t, base.DB.Create(&problemSet).Error)
		assert.NoError(t, problemSet.DeleteProblems([]uint{}))
		assert.Equal(t, ProblemSet{
			ID:          problemSet.ID,
			Name:        "test_delete_students_nothing_problem_set_name",
			Description: "test_delete_students_nothing_problem_set_description",
			Problems: []*Problem{
				problem1,
				problem2,
				problem3,
			},
			StartTime: hashStringToTime("test_delete_students_nothing_problem_set_start_time"),
			EndTime:   hashStringToTime("test_delete_students_nothing_problem_set_end_time"),
			CreatedAt: problemSet.CreatedAt,
			UpdatedAt: problemSet.UpdatedAt,
		}, problemSet)
	})
	t.Run("DeleteInEmptySet", func(t *testing.T) {
		problem1 := createProblemForTest(t, "delete_in_empty_set", 1)
		problem2 := createProblemForTest(t, "delete_in_empty_set", 2)
		t.Parallel()
		problemSet := ProblemSet{
			Name:        "test_delete_students_in_empty_set_problem_set_name",
			Description: "test_delete_students_in_empty_set_problem_set_description",
			Problems:    []*Problem{},
			StartTime:   hashStringToTime("test_delete_students_in_empty_set_problem_set_start_time"),
			EndTime:     hashStringToTime("test_delete_students_in_empty_set_problem_set_end_time"),
		}
		assert.NoError(t, base.DB.Create(&problemSet).Error)
		assert.NoError(t, problemSet.DeleteProblems([]uint{
			problem1.ID,
			problem2.ID,
		}))
		assert.Equal(t, ProblemSet{
			ID:          problemSet.ID,
			Name:        "test_delete_students_in_empty_set_problem_set_name",
			Description: "test_delete_students_in_empty_set_problem_set_description",
			Problems:    nil,
			StartTime:   hashStringToTime("test_delete_students_in_empty_set_problem_set_start_time"),
			EndTime:     hashStringToTime("test_delete_students_in_empty_set_problem_set_end_time"),
			CreatedAt:   problemSet.CreatedAt,
			UpdatedAt:   problemSet.UpdatedAt,
		}, problemSet)
	})
	t.Run("DeleteNotBelongTo", func(t *testing.T) {
		problem1 := createProblemForTest(t, "delete_not_belong_to", 1)
		problem2 := createProblemForTest(t, "delete_not_belong_to", 2)
		problem3 := createProblemForTest(t, "delete_not_belong_to", 3)
		t.Parallel()
		problemSet := ProblemSet{
			Name:        "test_delete_students_not_belong_to_problem_set_name",
			Description: "test_delete_students_not_belong_to_problem_set_description",
			Problems: []*Problem{
				problem1,
				problem2,
			},
			StartTime: hashStringToTime("test_delete_students_not_belong_to_problem_set_start_time"),
			EndTime:   hashStringToTime("test_delete_students_not_belong_to_problem_set_end_time"),
		}
		assert.NoError(t, base.DB.Create(&problemSet).Error)
		assert.NoError(t, problemSet.DeleteProblems([]uint{
			problem2.ID,
			problem3.ID,
		}))
		assert.Equal(t, ProblemSet{
			ID:          problemSet.ID,
			Name:        "test_delete_students_not_belong_to_problem_set_name",
			Description: "test_delete_students_not_belong_to_problem_set_description",
			Problems: []*Problem{
				problem1,
			},
			StartTime: hashStringToTime("test_delete_students_not_belong_to_problem_set_start_time"),
			EndTime:   hashStringToTime("test_delete_students_not_belong_to_problem_set_end_time"),
			CreatedAt: problemSet.CreatedAt,
			UpdatedAt: problemSet.UpdatedAt,
		}, problemSet)
	})
}
