package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func CreateSubmission(c echo.Context) error {
	req := request.CreateSubmissionRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}

	problem := models.Problem{}
	if err := base.DB.First(&problem, c.Param("pid")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find problem"))
		}
	}
	problem.LoadTestCases()

	user := c.Get("user").(models.User)

	if !problem.Public && !user.Can("read_problem", problem) && !user.Can("read_problem") {
		return c.JSON(http.StatusForbidden, response.ErrorResp("PERMISSION_DENIED", nil))
	}

	languageAllow := false
	for _, language := range strings.Split(problem.LanguageAllowed, ",") {
		if language == req.Language {
			languageAllow = true
			break
		}
	}
	if !languageAllow {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("INVALID_LANGUAGE", nil))
	}

	file, err := c.FormFile("code")
	if err != nil && err != http.ErrMissingFile && err.Error() != "request Content-Type isn't multipart/form-data" {
		panic(errors.Wrap(err, "could not read file"))
	}
	if file == nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("INVALID_FILE", nil))
	}

	// TODO: save file

	priority := models.PriorityDefault

	submission := models.Submission{
		UserID:       user.ID,
		ProblemID:    problem.ID,
		ProblemSetId: 0,
		Language:     req.Language,
		FileName:     file.Filename,
		Priority:     priority,
		Judged:       false,
		Score:        0,
		Status:       "PENDING",
		Runs:         make([]models.Run, len(problem.TestCases)),
	}
	for i, testCase := range problem.TestCases {
		submission.Runs[i] = models.Run{
			UserID:             user.ID,
			ProblemID:          problem.ID,
			ProblemSetId:       0,
			TestCaseID:         testCase.ID,
			Sample:             testCase.Sample,
			Priority:           priority,
			Judged:             false,
			Status:             "PENDING",
			MemoryUsed:         0,
			TimeUsed:           0,
			OutputStrippedHash: "",
		}
	}
	utils.PanicIfDBError(base.DB.Create(&submission), "could not create submission and runs")

	utils.MustPutObject(file, c.Request().Context(), "submissions", fmt.Sprintf("%d/code", submission.ID))

	return c.JSON(http.StatusCreated, response.CreateSubmissionResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.SubmissionDetail `json:"submission"`
		}{
			resource.GetSubmissionDetail(&submission),
		},
	})
}

func GetSubmission(c echo.Context) error {
	user := c.Get("user").(models.User)
	submission := models.Submission{}
	if err := base.DB.Preload("Problem").First(&submission, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if user.Can("read_problem_secret") {
				return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
			} else {
				return c.JSON(http.StatusForbidden, response.ErrorResp("PERMISSION_DENIED", nil))
			}
		} else {
			panic(errors.Wrap(err, "could not find problem"))
		}
	}
	if user.ID != submission.UserID && !user.Can("read_problem_secret", submission.Problem) && !user.Can("read_problem_secret") {
		return c.JSON(http.StatusForbidden, response.ErrorResp("PERMISSION_DENIED", nil))
	}
	return c.JSON(http.StatusOK, response.GetSubmissionResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.SubmissionDetail `json:"submission"`
		}{
			resource.GetSubmissionDetail(&submission),
		},
	})
}

func GetSubmissions(c echo.Context) error {
	req := request.GetSubmissionsRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}

	query := base.DB.Model(&models.Submission{}).Where("problem_set_id = ?", 0).Order("id DESC") // Force order by id desc.

	if req.ProblemId != 0 {
		query = query.Where("problem_id = ?", req.ProblemId)
	}
	if req.UserId != 0 {
		query = query.Where("user_id = ?", req.UserId)
	}

	var submissions []models.Submission
	total, prevUrl, nextUrl, err := utils.Paginator(query, req.Limit, req.Offset, c.Request().URL, &submissions)
	if err != nil {
		if herr, ok := err.(utils.HttpError); ok {
			return herr.Response(c)
		}
		panic(err)
	}
	return c.JSON(http.StatusOK, response.GetSubmissionsResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Submissions []resource.Submission `json:"submissions"`
			Total       int                   `json:"total"`
			Count       int                   `json:"count"`
			Offset      int                   `json:"offset"`
			Prev        *string               `json:"prev"`
			Next        *string               `json:"next"`
		}{
			Submissions: resource.GetSubmissionSlice(submissions),
			Total:       total,
			Count:       len(submissions),
			Offset:      req.Offset,
			Prev:        prevUrl,
			Next:        nextUrl,
		},
	})
}
