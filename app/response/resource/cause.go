package resource

import "github.com/EduOJ/backend/database/models"

type CauseForAdmin struct {
	ID         uint `son:"id"`
	ProblemID  uint `json:"problem_id"`
	TestCaseID uint `json:"test_case_id"`

	Hash        string `json:"output_stripped_hash"`
	Description string `json:"description"`

	// Point: Points to be subtracted for this cause
	Point  uint `json:"point"`
	Marked bool `json:"marked"`
	Count  uint `json:"count"`
}

type Cause struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProblemID  uint      `sql:"index" json:"problem_id"`
	Problem    *Problem  `json:"problem"`
	TestCaseID uint      `sql:"index" json:"test_case_id"`
	TestCase   *TestCase `json:"test_case"`

	Hash        string `json:"output_stripped_hash" gorm:"index;not null;size:255;default:''"`
	Description string `json:"description"`
	Marked      bool   `json:"marked" gorm:"default:false;not null"`
}

func (c *CauseForAdmin) convert(cause *models.Cause) {
	c.ID = cause.ID
	c.ProblemID = cause.ProblemID
	c.TestCaseID = cause.TestCaseID
	c.Hash = cause.Hash
	c.Description = cause.Description
	c.Point = cause.Point
	c.Marked = cause.Marked
	c.Count = cause.Count
}

func (c *Cause) convert(cause *models.Cause) {
	c.ID = cause.ID
	c.ProblemID = cause.ProblemID
	c.TestCaseID = cause.TestCaseID
	c.Hash = cause.Hash
	c.Description = cause.Description
	c.Marked = cause.Marked
}

func GetCauseForAdmin(cause *models.Cause) *CauseForAdmin {
	c := CauseForAdmin{}
	c.convert(cause)
	return &c
}

func GetCause(cause *models.Cause) *Cause {
	c := Cause{}
	c.convert(cause)
	return &c
}

func GetCauseForAdminSlice(causes []*models.Cause) []CauseForAdmin {
	causeSlice := make([]CauseForAdmin, len(causes))
	for i, cause := range causes {
		causeSlice[i].convert(cause)
	}
	return causeSlice
}
