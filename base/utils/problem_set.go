package utils

import (
	"encoding/json"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/database/models"
	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func UpdateGrade(submission *models.Submission) error {
	// TODO: check problem set end time, don't update if problem set ends
	if submission.ProblemSetID == 0 {
		return nil
	}
	grade := models.Grade{}
	detail := make(map[uint]uint)
	var err error
	err = base.DB.First(&grade, "problem_set_id = ? and user_id = ?", submission.ProblemSetID, submission.UserID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			grade = models.Grade{
				UserID:       submission.UserID,
				ProblemSetID: submission.ProblemSetID,
				Detail:       datatypes.JSON("{}"),
				Total:        0,
			}
		} else {
			return err
		}
	}
	err = json.Unmarshal(grade.Detail, &detail)
	if err != nil {
		return err
	}
	detail[submission.ProblemID] = submission.Score
	grade.Detail, err = json.Marshal(detail)
	return base.DB.Save(&grade).Error
}
