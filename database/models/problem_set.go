package models

import (
	"encoding/json"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type ProblemSet struct {
	ID uint `gorm:"primaryKey" json:"id"`

	ClassID     uint   `sql:"index" json:"class_id" gorm:"not null"`
	Name        string `json:"name" gorm:"not null;size:255"`
	Description string `json:"description"`

	Problems []Problem `json:"problems" gorm:"many2many:problems_in_problem_sets"`
	Scores   []Grade   `json:"scores"`

	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type Grade struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID       uint `json:"user_id"`
	ProblemSetID uint `json:"problem_set_id"`

	ScoreDetail datatypes.JSON `json:"score_detail"`
	TotalScore  uint           `json:"total_score"`
}

func (p *ProblemSet) AddProblems(ids []uint) error {
	existingIds := make([]uint, len(p.Problems))
	for i, problem := range p.Problems {
		existingIds[i] = problem.ID
	}
	var problems []Problem
	if err := base.DB.Not("id in ?", existingIds).Preload("TestCases").Find(&problems, ids).Error; err != nil {
		return err
	}
	return base.DB.Model(p).Association("Problems").Append(&problems)
}

func (p *ProblemSet) DeleteProblems(ids []uint) error {
	var problems []Problem
	if err := base.DB.Find(&problems, ids).Error; err != nil {
		return err
	}
	return base.DB.Model(p).Association("Problems").Delete(&problems)
}

// TODO: register this function for event submission_judge_finished\
func UpdateGrade(submission Submission) error {
	problemSet := ProblemSet{}
	if err := base.DB.First(&problemSet, submission.ProblemSetId).Error; err != nil {
		return errors.Wrap(err, "could not find problem set for updating grade")
	}
	return problemSet.UpdateGrade(submission)
}

func (p *ProblemSet) UpdateGrade(submission Submission) error {
	var grades []Grade
	scoresDetail := make(map[uint]uint)
	var originalTotalScore uint = 0
	if err := base.DB.Model(p).Association("Scores").Find(&grades, "user_id", submission.UserID); err != nil {
		return errors.Wrap(err, "could not find grade for updating grade")
	}
	if len(grades) == 1 {
		b, err := grades[0].ScoreDetail.MarshalJSON()
		if err != nil {
			return errors.Wrap(err, "could not marshal json for original score detail while updating grade")
		}
		err = json.Unmarshal(b, &scoresDetail)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal json for original score detail while updating grade")
		}
		originalTotalScore = grades[0].TotalScore
		if err = base.DB.Model(p).Association("Scores").Delete(grades); err != nil {
			return errors.Wrap(err, "could not delete grade for updating grade")
		}
	} else if len(grades) > 1 {
		return errors.New("there are two score records for the same user in the same question group")
	}
	originalScore := scoresDetail[submission.ProblemID]
	scoresDetail[submission.ProblemID] = submission.Score
	updatedGrade := Grade{
		UserID:     submission.UserID,
		TotalScore: originalTotalScore - originalScore + submission.Score,
	}
	b, err := json.Marshal(scoresDetail)
	if err != nil {
		return errors.Wrap(err, "could not marshal json for updated score detail while updating grade")
	}
	if err = updatedGrade.ScoreDetail.UnmarshalJSON(b); err != nil {
		return errors.Wrap(err, "could not unmarshal json for updated score detail while updating grade")
	}
	if err = base.DB.Model(p).Association("Scores").Append(&updatedGrade); err != nil {
		return errors.Wrap(err, "could not replace grade for updating grade")
	}
	return nil
}
