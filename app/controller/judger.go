package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"sync"
	"time"
)

var taskLock sync.Mutex

func getRun() *models.Run {
	//taskLock.Lock()
	//defer taskLock.Unlock()
	run := models.Run{}
	err := base.DB.Order("priority desc").
		Order("id asc").
		Preload("Problem").
		Preload("TestCase").
		Preload("Submission.Language").
		First(&run, "status = 'PENDING'").Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			panic(errors.Wrap(err, "could not query run"))
		}
	}
	return &run
}

func generateResponse(run *models.Run) response.GetTaskResponse {
	inputUrl, err := utils.GetPresignedURL("problems", fmt.Sprintf("%d/input/%d.in", run.Problem.ID, run.TestCase.ID), run.TestCase.InputFileName)
	if err != nil {
		panic(errors.Wrap(err, "could not get problem input file"))
	}
	outputUrl, err := utils.GetPresignedURL("problems", fmt.Sprintf("%d/output/%d.out", run.Problem.ID, run.TestCase.ID), run.TestCase.OutputFileName)
	if err != nil {
		panic(errors.Wrap(err, "could not get problem output file"))
	}
	codeUrl, err := utils.GetPresignedURL("submissions", fmt.Sprintf("%d/code", run.Submission.ID), run.Submission.FileName)
	if err != nil {
		panic(errors.Wrap(err, "could not get problem code file"))
	}
	return response.GetTaskResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			RunID              uint            `json:"run_id"`
			Language           models.Language `json:"language"`
			TestCaseID         uint            `json:"test_case_id"`
			InputFile          string          `json:"input_file"`
			OutputFile         string          `json:"output_file"`
			CodeFile           string          `json:"code_file"`
			TestCaseUpdatedAt  time.Time       `json:"test_case_updated_at"`
			MemoryLimit        uint64          `json:"memory_limit"`
			TimeLimit          uint            `json:"time_limit"`
			CompileEnvironment string          `json:"compile_environment"`
			CompareScriptName  string          `json:"compare_script_name"`
		}{
			RunID:              run.ID,
			Language:           *run.Submission.Language,
			TestCaseID:         run.TestCaseID,
			InputFile:          inputUrl,
			OutputFile:         outputUrl,
			CodeFile:           codeUrl,
			TestCaseUpdatedAt:  run.TestCase.UpdatedAt,
			MemoryLimit:        run.Problem.MemoryLimit,
			TimeLimit:          run.Problem.TimeLimit,
			CompileEnvironment: run.Problem.CompileEnvironment,
			CompareScriptName:  run.Problem.CompareScriptName,
		},
	}
}

func GetTask(c echo.Context) error {
	var poll bool
	var err error
	var run *models.Run
	if c.QueryParam("poll") == "1" {
		poll = true
	}

	taskLock.Lock()
	run = getRun()
	if run == nil {
		taskLock.Unlock()
		if poll {
			goto poll
		}
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	}
	run.Status = "JUDGING"
	run.JudgerName = c.Request().Header.Get("Judger-Name")
	err = base.DB.Save(&run).Error
	taskLock.Unlock()
	if err != nil {
		panic(errors.Wrap(err, "could not update run"))
	}
	return c.JSON(http.StatusOK, generateResponse(run))

poll:
	timeoutChan := time.After(viper.GetDuration("polling_timeout") * time.Second)
	timeout := false
	sub := base.Redis.Subscribe(c.Request().Context(), "runs")
	for {
		select {
		case <-sub.Channel():
			taskLock.Lock()
			run = getRun()
			if run == nil {
				taskLock.Unlock()
				break
			}
			run.Status = "JUDGING"
			run.JudgerName = c.Request().Header.Get("Judger-Name")
			err := base.DB.Save(&run).Error
			taskLock.Unlock()
			if err != nil {
				panic(errors.Wrap(err, "could not update run"))
			}
			break
		case <-c.Request().Context().Done():
			// context cancelled
			return nil
		case <-timeoutChan:
			timeout = true
		}
		if timeout {
			break
		}
		if run != nil {
			break
		}
	}
	if run == nil {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	}
	return c.JSON(http.StatusOK, generateResponse(run))
}
