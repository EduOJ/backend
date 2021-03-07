package utils

import (
	"encoding/json"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/database/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"testing"
	"time"
)

func createJSONForTest(t *testing.T, in interface{}) datatypes.JSON {
	j := datatypes.JSON{}
	b, err := json.Marshal(in)
	assert.NoError(t, err)
	assert.NoError(t, j.UnmarshalJSON(b))
	return j
}

func TestUpdateGrade(t *testing.T) {
	t.Parallel()

	user1 := models.User{
		Username: "test_update_grade_1_username",
		Nickname: "test_update_grade_1_nickname",
		Email:    "test_update_grade_1@email.com",
		Password: "test_update_grade_1_password",
	}
	user2 := models.User{
		Username: "test_update_grade_2_username",
		Nickname: "test_update_grade_2_nickname",
		Email:    "test_update_grade_2@email.com",
		Password: "test_update_grade_2_password",
	}
	assert.NoError(t, base.DB.Create(&user1).Error)
	assert.NoError(t, base.DB.Create(&user2).Error)
	problem1 := models.Problem{
		Name:        "test_update_grade_1_name",
		Description: "test_update_grade_1_description",
		MemoryLimit: 1024,
		TimeLimit:   1000,
	}
	problem2 := models.Problem{
		Name:        "test_update_grade_2_name",
		Description: "test_update_grade_2_description",
		MemoryLimit: 2048,
		TimeLimit:   2000,
	}
	assert.NoError(t, base.DB.Create(&problem1).Error)
	assert.NoError(t, base.DB.Create(&problem2).Error)

	t.Run("SubmissionNotInProblemSet", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, UpdateGrade(&models.Submission{
			ProblemSetID: 0,
			UserID:       user1.ID,
			ProblemID:    problem2.ID,
			Score:        100,
		}))
	})
	t.Run("EndedProblemSet", func(t *testing.T) {
		t.Parallel()
		problemSet := models.ProblemSet{
			Name:        "test_update_grade_ended_name",
			Description: "test_update_grade_ended_description",
			Problems: []*models.Problem{
				&problem1,
				&problem2,
			},
			Grades: []*models.Grade{
				{
					UserID: user1.ID,
					Detail: createJSONForTest(t, map[uint]uint{
						problem1.ID: 33,
						problem2.ID: 44,
					}),
					Total: 77,
				},
			},
			StartTime: time.Now().Add(-2 * time.Hour),
			EndTime:   time.Now().Add(-1 * time.Hour),
		}
		assert.NoError(t, base.DB.Create(&problemSet).Error)
		assert.NoError(t, UpdateGrade(&models.Submission{
			ProblemSetID: problemSet.ID,
			UserID:       user1.ID,
			ProblemID:    problem2.ID,
			Score:        100,
		}))
		assert.NoError(t, base.DB.Preload("Grades").First(&problemSet, problemSet.ID).Error)
		assert.Equal(t, []*models.Grade{
			{
				ID:           problemSet.Grades[0].ID,
				UserID:       user1.ID,
				ProblemSetID: problemSet.ID,
				Detail: createJSONForTest(t, map[uint]uint{
					problem1.ID: 33,
					problem2.ID: 44,
				}),
				Total:     77,
				CreatedAt: problemSet.Grades[0].CreatedAt,
				UpdatedAt: problemSet.Grades[0].UpdatedAt,
			},
		}, problemSet.Grades)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		problemSet := models.ProblemSet{
			Name:        "test_update_grade_name",
			Description: "test_update_grade_description",
			Problems: []*models.Problem{
				&problem1,
				&problem2,
			},
			Grades: []*models.Grade{
				{
					UserID: user1.ID,
					Detail: createJSONForTest(t, map[uint]uint{
						problem1.ID: 11,
						problem2.ID: 12,
					}),
					Total: 23,
				},
			},
			StartTime: time.Now().Add(-1 * time.Hour),
			EndTime:   time.Now().Add(time.Hour),
		}
		assert.NoError(t, base.DB.Create(&problemSet).Error)
		assert.NoError(t, UpdateGrade(&models.Submission{
			ProblemSetID: problemSet.ID,
			UserID:       user1.ID,
			ProblemID:    problem2.ID,
			Score:        120,
		}))
		assert.NoError(t, base.DB.Preload("Grades").First(&problemSet, problemSet.ID).Error)
		assert.Equal(t, []*models.Grade{
			{
				ID:           problemSet.Grades[0].ID,
				UserID:       user1.ID,
				ProblemSetID: problemSet.ID,
				Detail: createJSONForTest(t, map[uint]uint{
					problem1.ID: 11,
					problem2.ID: 120,
				}),
				Total:     131,
				CreatedAt: problemSet.Grades[0].CreatedAt,
				UpdatedAt: problemSet.Grades[0].UpdatedAt,
			},
		}, problemSet.Grades)
		assert.NoError(t, UpdateGrade(&models.Submission{
			ProblemSetID: problemSet.ID,
			UserID:       user2.ID,
			ProblemID:    problem1.ID,
			Score:        21,
		}))
		assert.NoError(t, base.DB.Preload("Grades").First(&problemSet, problemSet.ID).Error)
		assert.Equal(t, []*models.Grade{
			{
				ID:           problemSet.Grades[0].ID,
				UserID:       user1.ID,
				ProblemSetID: problemSet.ID,
				Detail: createJSONForTest(t, map[uint]uint{
					problem1.ID: 11,
					problem2.ID: 120,
				}),
				Total:     131,
				CreatedAt: problemSet.Grades[0].CreatedAt,
				UpdatedAt: problemSet.Grades[0].UpdatedAt,
			},
			{
				ID:           problemSet.Grades[1].ID,
				UserID:       user2.ID,
				ProblemSetID: problemSet.ID,
				Detail: createJSONForTest(t, map[uint]uint{
					problem1.ID: 21,
				}),
				Total:     21,
				CreatedAt: problemSet.Grades[1].CreatedAt,
				UpdatedAt: problemSet.Grades[1].UpdatedAt,
			},
		}, problemSet.Grades)
		assert.NoError(t, UpdateGrade(&models.Submission{
			ProblemSetID: problemSet.ID,
			UserID:       user2.ID,
			ProblemID:    problem2.ID,
			Score:        22,
		}))
		assert.NoError(t, base.DB.Preload("Grades").First(&problemSet, problemSet.ID).Error)
		assert.Equal(t, []*models.Grade{
			{
				ID:           problemSet.Grades[0].ID,
				UserID:       user1.ID,
				ProblemSetID: problemSet.ID,
				Detail: createJSONForTest(t, map[uint]uint{
					problem1.ID: 11,
					problem2.ID: 120,
				}),
				Total:     131,
				CreatedAt: problemSet.Grades[0].CreatedAt,
				UpdatedAt: problemSet.Grades[0].UpdatedAt,
			},
			{
				ID:           problemSet.Grades[1].ID,
				UserID:       user2.ID,
				ProblemSetID: problemSet.ID,
				Detail: createJSONForTest(t, map[uint]uint{
					problem1.ID: 21,
					problem2.ID: 22,
				}),
				Total:     43,
				CreatedAt: problemSet.Grades[1].CreatedAt,
				UpdatedAt: problemSet.Grades[1].UpdatedAt,
			},
		}, problemSet.Grades)
		assert.NoError(t, UpdateGrade(&models.Submission{
			ProblemSetID: problemSet.ID,
			UserID:       user2.ID,
			ProblemID:    problem2.ID,
			Score:        220,
		}))
		assert.NoError(t, base.DB.Preload("Grades").First(&problemSet, problemSet.ID).Error)
		assert.Equal(t, []*models.Grade{
			{
				ID:           problemSet.Grades[0].ID,
				UserID:       user1.ID,
				ProblemSetID: problemSet.ID,
				Detail: createJSONForTest(t, map[uint]uint{
					problem1.ID: 11,
					problem2.ID: 120,
				}),
				Total:     131,
				CreatedAt: problemSet.Grades[0].CreatedAt,
				UpdatedAt: problemSet.Grades[0].UpdatedAt,
			},
			{
				ID:           problemSet.Grades[1].ID,
				UserID:       user2.ID,
				ProblemSetID: problemSet.ID,
				Detail: createJSONForTest(t, map[uint]uint{
					problem1.ID: 21,
					problem2.ID: 220,
				}),
				Total:     241,
				CreatedAt: problemSet.Grades[1].CreatedAt,
				UpdatedAt: problemSet.Grades[1].UpdatedAt,
			},
		}, problemSet.Grades)
	})
}
