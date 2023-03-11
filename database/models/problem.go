package models

import (
	"context"
	"fmt"
	"time"

	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/log"
	"github.com/EduOJ/backend/database"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type TestCase struct {
	ID uint `gorm:"primaryKey" json:"id"`

	ProblemID uint `sql:"index" json:"problem_id" gorm:"not null"`
	Score     uint `json:"score" gorm:"default:0;not null"` // 0 for 平均分配
	Sample    bool `json:"sample" gorm:"default:false;not null"`

	InputFileName  string `json:"input_file_name" gorm:"size:255;default:'';not null"`
	OutputFileName string `json:"output_file_name" gorm:"size:255;default:'';not null"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type ProblemTag struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	ProblemID uint `gorm:"index"`
	Name      string
	CreatedAt time.Time `json:"created_at"`
}

type Problem struct {
	ID                 uint   `gorm:"primaryKey" json:"id"`
	Name               string `sql:"index" json:"name" gorm:"size:255;default:'';not null"`
	Description        string `json:"description"`
	AttachmentFileName string `json:"attachment_file_name" gorm:"size:255;default:'';not null"`
	Public             bool   `json:"public" gorm:"default:false;not null"`
	Privacy            bool   `json:"privacy" gorm:"default:false;not null"`

	MemoryLimit       uint64               `json:"memory_limit" gorm:"default:0;not null;type:bigint"`               // Byte
	TimeLimit         uint                 `json:"time_limit" gorm:"default:0;not null"`                             // ms
	LanguageAllowed   database.StringArray `json:"language_allowed" gorm:"size:255;default:'';not null;type:string"` // E.g.    cpp,c,java,python
	BuildArg          string               `json:"build_arg" gorm:"size:2047;default:'';not null"`                   // E.g.  O2=false
	CompareScriptName string               `json:"compare_script_name" gorm:"default:0;not null"`
	CompareScript     Script               `json:"compare_script"`

	TestCases []TestCase `json:"test_cases"`
	Tags      []Tag      `json:"tags" gorm:"OnDelete:CASCADE"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type Tag struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	ProblemID uint
	Name      string
	CreatedAt time.Time `json:"created_at"`
}

func (p Problem) GetID() uint {
	return p.ID
}

func (p Problem) TypeName() string {
	return "problem"
}

func (p *Problem) LoadTestCases() {
	err := base.DB.Model(p).Association("TestCases").Find(&p.TestCases)
	if err != nil {
		panic(err)
	}
}

func (p *Problem) LoadTags() {
	err := base.DB.Model(p).Association("Tags").Find(&p.Tags)
	if err != nil {
		panic(err)
	}
}

func (p *Problem) AfterDelete(tx *gorm.DB) (err error) {
	if err := tx.Where("problem_id = ?", p.ID).Delete(&Submission{}).Error; err != nil {
		return err
	}
	return tx.Where("problem_id = ?", p.ID).Delete(&TestCase{}).Error
}

func (t *TestCase) AfterDelete(tx *gorm.DB) (err error) {
	err = base.Storage.RemoveObject(context.Background(), "problems", fmt.Sprintf("%d/input/%d.in", t.ProblemID, t.ID), minio.RemoveObjectOptions{})
	if err != nil {
		log.Errorf("Error occurred in TestCase afterDelete, %+v\n", err)
	}
	err = base.Storage.RemoveObject(context.Background(), "problems", fmt.Sprintf("%d/output/%d.out", t.ProblemID, t.ID), minio.RemoveObjectOptions{})
	if err != nil {
		log.Errorf("Error occurred in TestCase afterDelete, %+v\n", err)
	}
	return tx.Where("test_case_id = ?", t.ID).Delete(&Run{}).Error
}
