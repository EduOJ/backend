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
	class := models.Class{}
	if err := base.DB.First(&class, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("CLASS_NOT_FOUND", nil))
		}
		panic(errors.Wrap(err, "could not get class while creating problem set"))
	}
	problemSet := models.ProblemSet{
		Name:        req.Name,
		Description: req.Description,
		Problems:    nil,
		Grades:      nil,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
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
			return c.JSON(http.StatusNotFound, response.ErrorResp("CLASS_NOT_FOUND", nil))
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
	if err := base.DB.Model(&sourceClass).Preload("Problems").Association("ProblemSets").
		Find(&sourceProblemSets, req.SourceProblemSetID); err != nil {
		panic(errors.Wrap(err, "could not find source problem set when cloning problem set"))
	}
	if len(sourceProblemSets) == 0 {
		return c.JSON(http.StatusNotFound, response.ErrorResp("SOURCE_PROBLEM_SET_NOT_FOUND", nil))
	}
	problemSet := models.ProblemSet{
		Name:        sourceProblemSets[0].Name,
		Description: sourceProblemSets[0].Description,
		Problems:    sourceProblemSets[0].Problems,
		Grades:      nil,
		StartTime:   sourceProblemSets[0].StartTime,
		EndTime:     sourceProblemSets[0].EndTime,
	}
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
	if err := base.DB.First(&class, c.Param("class_id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("CLASS_NOT_FOUND", nil))
		}
		panic(errors.Wrap(err, "could not get class for getting problem set"))
	}
	if err := base.DB.Model(&class).Preload("Problems").Preload("Grades").Association("ProblemSets").
		Find(&problemSets, "id = ?", c.Param("id")); err != nil {
		panic(errors.Wrap(err, "could not get problem set for getting problem set"))
	}
	if len(problemSets) == 0 {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if len(problemSets) > 1 {
		panic("there are two problem sets that have same id")
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
	var users []models.User
	if err := base.DB.Model(&class).Association("Students").Find(&users, user.ID); err != nil {
		panic(errors.Wrap(err, "could not check student in class for getting problem set"))
	}
	if len(users) == 0 {
		return c.JSON(http.StatusNotFound, response.ErrorResp("CLASS_NOT_FOUND", nil))
	}
	if time.Now().Before(problemSets[0].StartTime) || time.Now().After(problemSets[0].EndTime) {
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

func UpdateProblemSet(c echo.Context) error {
	req := request.UpdateProblemSetRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	class := models.Class{}
	var problemSets []models.ProblemSet
	if err := base.DB.First(&class, c.Param("class_id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("CLASS_NOT_FOUND", nil))
		}
		panic(errors.Wrap(err, "could not get class for updating problem set"))
	}
	if err := base.DB.Model(&class).Preload("Problems").Preload("Grades").Association("ProblemSets").
		Find(&problemSets, "id = ?", c.Param("id")); err != nil {
		panic(errors.Wrap(err, "could not get problem set for updating problem set"))
	}
	if len(problemSets) == 0 {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if len(problemSets) > 1 {
		panic("there are two problem sets that have same id")
	}

	problemSets[0].Name = req.Name
	problemSets[0].Description = req.Description
	problemSets[0].StartTime = req.StartTime
	problemSets[0].EndTime = req.EndTime
	utils.PanicIfDBError(base.DB.Save(&problemSets[0]), "could not update problem set for updating problem set")
	return c.JSON(http.StatusOK, response.UpdateProblemSetResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ProblemSetDetail `json:"problem_set"`
		}{
			resource.GetProblemSetDetail(&problemSets[0]),
		},
	})
}

func AddProblemsToSet(c echo.Context) error {
	req := request.AddProblemsToSetRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}

	class := models.Class{}
	var problemSets []models.ProblemSet
	if err := base.DB.First(&class, c.Param("class_id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("CLASS_NOT_FOUND", nil))
		}
		panic(errors.Wrap(err, "could not get class for adding problems in set"))
	}
	if err := base.DB.Model(&class).Preload("Problems").Preload("Grades").Association("ProblemSets").
		Find(&problemSets, "id = ?", c.Param("id")); err != nil {
		panic(errors.Wrap(err, "could not get problem set for adding problems in set"))
	}
	if len(problemSets) == 0 {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if len(problemSets) > 1 {
		panic("there are two problem sets that have same id")
	}

	if err := problemSets[0].AddProblems(req.ProblemIDs); err != nil {
		panic(errors.Wrap(err, "could not add problems for problem set"))
	}
	return c.JSON(http.StatusOK, response.AddProblemsToSetResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ProblemSetDetail `json:"problem_set"`
		}{
			resource.GetProblemSetDetail(&problemSets[0]),
		},
	})
}

func DeleteProblemsFromSet(c echo.Context) error {
	req := request.DeleteProblemsFromSetRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}

	class := models.Class{}
	var problemSets []models.ProblemSet
	if err := base.DB.First(&class, c.Param("class_id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("CLASS_NOT_FOUND", nil))
		}
		panic(errors.Wrap(err, "could not get class for deleting problems in set"))
	}
	if err := base.DB.Model(&class).Preload("Problems").Preload("Grades").Association("ProblemSets").
		Find(&problemSets, "id = ?", c.Param("id")); err != nil {
		panic(errors.Wrap(err, "could not get problem set for deleting problems in set"))
	}
	if len(problemSets) == 0 {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if len(problemSets) > 1 {
		panic("there are two problem sets that have same id")
	}

	if err := problemSets[0].DeleteProblems(req.ProblemIDs); err != nil {
		panic(errors.Wrap(err, "could not delete problems for problem set"))
	}
	return c.JSON(http.StatusOK, response.DeleteProblemsFromSetResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ProblemSetDetail `json:"problem_set"`
		}{
			resource.GetProblemSetDetail(&problemSets[0]),
		},
	})
}

func DeleteProblemSet(c echo.Context) error {
	class := models.Class{}
	var problemSets []models.ProblemSet
	if err := base.DB.First(&class, c.Param("class_id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("CLASS_NOT_FOUND", nil))
		}
		panic(errors.Wrap(err, "could not get class for deleting problems in set"))
	}
	if err := base.DB.Model(&class).Association("ProblemSets").
		Find(&problemSets, "id = ?", c.Param("id")); err != nil {
		panic(errors.Wrap(err, "could not get problem set for deleting problems in set"))
	}
	utils.PanicIfDBError(base.DB.Delete(&problemSets), "could not delete problem set for deleting problem set")

	return c.JSON(http.StatusOK, response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}
