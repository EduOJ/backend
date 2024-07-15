package utils

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/log"
	"github.com/EduOJ/backend/database/models"
	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var gradeLock sync.Mutex

func UpdateGrade(submission *models.Submission) error {
	gradeLock.Lock()
	defer gradeLock.Unlock()
	if submission.ProblemSetID == 0 {
		return nil
	}
	if submission.ProblemSet == nil {
		problemSet := models.ProblemSet{}
		if err := base.DB.First(&problemSet, submission.ProblemSetID).Error; err != nil {
			return errors.Wrap(err, "could not get problem set for updating grade")
		}
		submission.ProblemSet = &problemSet
	}
	if time.Now().After(submission.ProblemSet.EndTime) {
		return nil
	}
	grade := models.Grade{}
	detail := make(map[uint]uint)
	var err error
	err = base.DB.First(&grade, "problem_set_id = ? and user_id = ?", submission.ProblemSetID, submission.UserID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			problemSet := models.ProblemSet{}
			if err := base.DB.First(&problemSet, submission.ProblemSetID).Error; err != nil {
				return err
			}
			grade = models.Grade{
				UserID:       submission.UserID,
				ProblemSetID: submission.ProblemSetID,
				ClassID:      problemSet.ClassID,
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
	if detail[submission.ProblemID] < submission.Score {
		detail[submission.ProblemID] = submission.Score
	}
	grade.Detail, err = json.Marshal(detail)
	if err != nil {
		return err
	}
	return base.DB.Save(&grade).Error
}

func RefreshGrades(problemSet *models.ProblemSet) error {
	gradeLock.Lock()
	defer gradeLock.Unlock()
	if err := base.DB.Delete(&models.Grade{}, "problem_set_id = ?", problemSet.ID).Error; err != nil {
		return err
	}
	var grades []*models.Grade
	for _, u := range problemSet.Class.Students {
		grade := models.Grade{
			UserID:       u.ID,
			ProblemSetID: problemSet.ID,
			ClassID:      problemSet.ClassID,
			Detail:       nil,
			Total:        0,
		}
		detail := make(map[uint]uint)
		for _, p := range problemSet.Problems {
			var score uint = 0
			submission := models.Submission{}
			err := base.DB.
				Where("user_id = ?", u.ID).
				Where("problem_id = ?", p.ID).
				Where("problem_set_id = ?", problemSet.ID).
				Where("created_at < ?", problemSet.EndTime).
				Order("score desc").
				Order("created_at desc").
				First(&submission).Error
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.Wrap(err, "could not get submission when refreshing grades")
				}
			} else {
				score = submission.Score
			}
			detail[p.ID] = score
			grade.Total += score
		}
		var err error
		grade.Detail, err = json.Marshal(detail)
		if err != nil {
			return errors.Wrap(err, "could not marshal grade detail when refreshing grades")
		}
		grades = append(grades, &grade)
	}
	if err := base.DB.Create(&grades).Error; err != nil {
		return errors.Wrap(err, "could not create grades when refreshing grades")
	}
	problemSet.Grades = grades
	return nil
}

// CreateEmptyGrades Creates empty grades(score 0 for all the problems)
//
//	for users who don't have a grade for this problem set.
func CreateEmptyGrades(problemSet *models.ProblemSet) error {
	gradeLock.Lock()
	defer gradeLock.Unlock()

	// Create empty grade JSON object
	detail := make(map[uint]uint)
	for _, p := range problemSet.Problems {
		detail[p.ID] = 0
	}
	emptyDetail, err := json.Marshal(detail) // turn map to json
	if err != nil {
		log.Errorf("Error marshalling grade detail for empty grade: %v", err)
		return errors.Wrap(err, "could not marshal grade detail for empty grade")
	}

	// json log
	log.Debugf("Empty detail JSON: %s", emptyDetail)

	// Record students who have a grade
	gradeSet := make(map[uint]bool)
	for _, g := range problemSet.Grades {
		gradeSet[g.UserID] = true
	}

	// Generate empty grade slice
	grades := make([]*models.Grade, 0, len(problemSet.Class.Students)-len(problemSet.Grades))
	for _, u := range problemSet.Class.Students {
		if gradeSet[u.ID] {
			continue
		}
		newGrade := models.Grade{
			UserID:       u.ID,
			ProblemSetID: problemSet.ID,
			ClassID:      problemSet.ClassID,
			Detail:       emptyDetail,
			Total:        0,
		}
		grades = append(grades, &newGrade)
	}

	// Store empty grades into DB
	if len(grades) > 0 {
		if err = base.DB.Create(&grades).Error; err != nil {
			return errors.Wrap(err, "could not create empty grades")
		}
	}

	// Update problem set
	problemSet.Grades = append(problemSet.Grades, grades...)
	return nil
}
