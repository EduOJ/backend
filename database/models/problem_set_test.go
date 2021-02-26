package models

import (
	"encoding/json"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"hash/fnv"
	"testing"
	"time"
)

func hashStringToTime(s string) time.Time {
	h := fnv.New32()
	if _, err := h.Write([]byte(s)); err != nil {
		panic(err)
	}
	return time.Unix(int64(h.Sum32()), 0).UTC()
}

func createJSONForTest(t *testing.T, in interface{}) datatypes.JSON {
	j := datatypes.JSON{}
	b, err := json.Marshal(in)
	assert.NoError(t, err)
	assert.NoError(t, j.UnmarshalJSON(b))
	return j
}

func TestUpdateGrade(t *testing.T) {
	t.Parallel()

	user1 := User{
		Username: "test_update_grade_1_username",
		Nickname: "test_update_grade_1_nickname",
		Email:    "test_update_grade_1@email.com",
		Password: "test_update_grade_1_password",
	}
	user2 := User{
		Username: "test_update_grade_2_username",
		Nickname: "test_update_grade_2_nickname",
		Email:    "test_update_grade_2@email.com",
		Password: "test_update_grade_2_password",
	}
	assert.NoError(t, base.DB.Create(&user1).Error)
	assert.NoError(t, base.DB.Create(&user2).Error)
	problem1 := Problem{
		Name:        "test_update_grade_1_name",
		Description: "test_update_grade_1_description",
		MemoryLimit: 1024,
		TimeLimit:   1000,
	}
	problem2 := Problem{
		Name:        "test_update_grade_2_name",
		Description: "test_update_grade_2_description",
		MemoryLimit: 2048,
		TimeLimit:   2000,
	}
	assert.NoError(t, base.DB.Create(&problem1).Error)
	assert.NoError(t, base.DB.Create(&problem2).Error)
	problemSet := ProblemSet{
		Name:        "test_update_grade_name",
		Description: "test_update_grade_description",
		Problems: []Problem{
			problem1,
			problem2,
		},
		Scores: []Grade{
			{
				UserID: user1.ID,
				ScoreDetail: createJSONForTest(t, map[uint]uint{
					problem1.ID: 11,
					problem2.ID: 12,
				}),
				TotalScore: 23,
			},
		},
		StartAt: hashStringToTime("test_update_grade_start_at"),
		EndAt:   hashStringToTime("test_update_grade_end_at"),
	}
	assert.NoError(t, base.DB.Create(&problemSet).Error)
	assert.NoError(t, problemSet.UpdateGrade(Submission{
		UserID:    user1.ID,
		ProblemID: problem2.ID,
		Score:     120,
	}))
	assert.Equal(t, []Grade{
		{
			ID:           problemSet.Scores[0].ID,
			UserID:       user1.ID,
			ProblemSetID: problemSet.ID,
			ScoreDetail: createJSONForTest(t, map[uint]uint{
				problem1.ID: 11,
				problem2.ID: 120,
			}),
			TotalScore: 131,
		},
	}, problemSet.Scores)
	assert.NoError(t, problemSet.UpdateGrade(Submission{
		UserID:    user2.ID,
		ProblemID: problem1.ID,
		Score:     21,
	}))
	assert.Equal(t, []Grade{
		{
			ID:           problemSet.Scores[0].ID,
			UserID:       user1.ID,
			ProblemSetID: problemSet.ID,
			ScoreDetail: createJSONForTest(t, map[uint]uint{
				problem1.ID: 11,
				problem2.ID: 120,
			}),
			TotalScore: 131,
		},
		{
			ID:           problemSet.Scores[1].ID,
			UserID:       user2.ID,
			ProblemSetID: problemSet.ID,
			ScoreDetail: createJSONForTest(t, map[uint]uint{
				problem1.ID: 21,
			}),
			TotalScore: 21,
		},
	}, problemSet.Scores)
	assert.NoError(t, problemSet.UpdateGrade(Submission{
		UserID:    user2.ID,
		ProblemID: problem2.ID,
		Score:     22,
	}))
	assert.Equal(t, []Grade{
		{
			ID:           problemSet.Scores[0].ID,
			UserID:       user1.ID,
			ProblemSetID: problemSet.ID,
			ScoreDetail: createJSONForTest(t, map[uint]uint{
				problem1.ID: 11,
				problem2.ID: 120,
			}),
			TotalScore: 131,
		},
		{
			ID:           problemSet.Scores[1].ID,
			UserID:       user2.ID,
			ProblemSetID: problemSet.ID,
			ScoreDetail: createJSONForTest(t, map[uint]uint{
				problem1.ID: 21,
				problem2.ID: 22,
			}),
			TotalScore: 43,
		},
	}, problemSet.Scores)
	assert.NoError(t, problemSet.UpdateGrade(Submission{
		UserID:    user2.ID,
		ProblemID: problem2.ID,
		Score:     220,
	}))
	assert.Equal(t, []Grade{
		{
			ID:           problemSet.Scores[0].ID,
			UserID:       user1.ID,
			ProblemSetID: problemSet.ID,
			ScoreDetail: createJSONForTest(t, map[uint]uint{
				problem1.ID: 11,
				problem2.ID: 120,
			}),
			TotalScore: 131,
		},
		{
			ID:           problemSet.Scores[1].ID,
			UserID:       user2.ID,
			ProblemSetID: problemSet.ID,
			ScoreDetail: createJSONForTest(t, map[uint]uint{
				problem1.ID: 21,
				problem2.ID: 220,
			}),
			TotalScore: 241,
		},
	}, problemSet.Scores)
}
