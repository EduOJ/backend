package controller

import (
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
	"time"
)

func CreateProblemSet(c echo.Context) error {
	req := request.CreateProblemSetRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	problemSet := models.ProblemSet{
		Name:        req.Name,
		Description: req.Description,
		Problems:    nil,
		Scores:      nil,
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
	}
	class := models.Class{}
	if err := base.DB.First(&class, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		}
		panic(errors.Wrap(err, "could not get class while creating problem set"))
	}
	if err := base.DB.Model(&class).Association("ProblemSets").Append(&problemSet); err != nil {
		panic(errors.Wrap(err, "could not add problem set for class when creating problem set"))
	}
	return c.JSON(http.StatusCreated, response.CreateProblemSetResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ProblemSetDetail `json:"problem_set"`
		}{
			resource.GetProblemSetDetail(&problemSet),
		},
	})
}

func CloneProblemSet(c echo.Context) error {
	req := request.CloneProblemSetRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	class := models.Class{}
	if err := base.DB.First(&class, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		}
		panic(errors.Wrap(err, "could not get class while cloning problem set"))
	}
	sourceClass := models.Class{}
	var sourceProblemSets []models.ProblemSet
	if err := base.DB.First(&sourceClass, req.SourceClassID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("SOURCE_CLASS_NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find source class when cloning problem set"))
		}
	}
	// TODO: maybe Preload("ProblemSets") + for instead of Association?
	if err := base.DB.Model(&sourceClass).Association("ProblemSets").
		Find(&sourceProblemSets, req.SourceProblemSetID); err != nil {
		panic(errors.Wrap(err, "could not find source problem set when cloning problem set"))
	}
	if len(sourceProblemSets) == 0 {
		return c.JSON(http.StatusNotFound, response.ErrorResp("SOURCE_PROBLEM_SET_NOT_FOUND", nil))
	}
	problemSet := sourceProblemSets[0]
	if err := base.DB.Model(&class).Association("ProblemSets").Append(&problemSet); err != nil {
		panic(errors.Wrap(err, "could not add problem set for class when cloning problem set"))
	}
	return c.JSON(http.StatusCreated, response.CloneProblemSetResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ProblemSetDetail `json:"problem_set"`
		}{
			resource.GetProblemSetDetail(&problemSet),
		},
	})
}

func GetProblemSet(c echo.Context) error {
	class := models.Class{}
	var problemSets []models.ProblemSet
	if err := base.DB.Preload("Problems").First(&class, c.Param("class_id")).
		Association("ProblemSets").Find(&problemSets, c.Param("id")); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("CLASS_NOT_FOUND", nil))
		}
		panic(errors.Wrap(err, "could not get class or problem set while getting problem set"))
	}
	if len(problemSets) == 0 {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	}
	user := c.Get("user").(models.User)
	if user.Can("manage_problem_sets", class) || user.Can("manage_problem_sets") {
		return c.JSON(http.StatusOK, response.GetProblemSetResponseForAdmin{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.ProblemSetDetail `json:"problem_set"`
			}{
				resource.GetProblemSetDetail(&problemSets[0]),
			},
		})
	}
	// TODO: maybe Preload("Students") + for instead of Association?
	if err := base.DB.Model(&class).Association("Students").Find(&models.User{}, user.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("CLASS_NOT_FOUND", nil))
		}
		panic(errors.Wrap(err, "could not check student in class for getting problem set"))
	}
	if time.Now().Before(problemSets[0].StartAt) || time.Now().After(problemSets[0].EndAt) {
		problemSets[0].Problems = nil
	}
	return c.JSON(http.StatusOK, response.GetProblemSetResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ProblemSet `json:"problem_set"`
		}{
			resource.GetProblemSet(&problemSets[0]),
		},
	})
}
