package controller

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"net/http"
	"sync"
)

var initCreator sync.Once
var problemCreator models.Role

func initCreatorRole() {
	utils.PanicIfDBError(base.DB.Where("name = ? and target = ?", "creator", "problem").First(&problemCreator), "could not get problem creator role")
}

func AdminCreateProblem(c echo.Context) error {
	initCreator.Do(initCreatorRole)
	req := request.AdminCreateProblemRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	var public, privacy bool
	if req.Public == nil {
		public = false
	}
	if req.Privacy == nil {
		privacy = true
	}
	problem := models.Problem{
		Name:               req.Name,
		Description:        req.Description,
		AttachmentFileName: req.AttachmentFileName,
		Public:             public,
		Privacy:            privacy,
		MemoryLimit:        req.MemoryLimit,
		TimeLimit:          req.TimeLimit,
		LanguageAllowed:    req.LanguageAllowed,
		CompileEnvironment: req.CompileEnvironment,
		CompareScriptID:    req.CompareScriptID,
	}
	utils.PanicIfDBError(base.DB.Create(&problem), "could not create problem")
	var user models.User
	if user, ok = c.Get("user").(models.User); !ok {
		panic("could not get user to grant role problem creator")
	}
	user.GrantRole(problemCreator, problem)
	return c.JSON(http.StatusCreated, response.AdminCreateProblemResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ProblemProfileForAdmin `json:"problem"`
		}{
			resource.GetProblemProfileForAdmin(&problem),
		},
	})
}

func AdminGetProblem(c echo.Context) error {
	problem, err := findProblem(c.Param("id"))
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
	}
	// TODO: load test cases
	return c.JSON(http.StatusOK, response.AdminGetProblemResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ProblemProfileForAdmin `json:"problem"`
		}{
			resource.GetProblemProfileForAdmin(problem),
		},
	})
}

func AdminGetProblems(c echo.Context) error {
	req := request.AdminGetProblemsRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}

	query, err := utils.Sorter(base.DB.Model(&models.Problem{}), req.OrderBy)
	if err != nil {
		if herr, ok := err.(utils.HttpError); ok {
			return herr.Response(c)
		}
		panic(err)
	}

	if req.Search != "" {
		query = query.Where("id like ? or name like ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	var problems []models.Problem
	total, prevUrl, nextUrl, err := utils.Paginator(query, req.Limit, req.Offset, c.Request().URL, &problems)
	return c.JSON(http.StatusOK, response.AdminGetProblemsResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Problems []resource.ProblemProfileForAdmin `json:"problems"`
			Total    int                               `json:"total"`
			Count    int                               `json:"count"`
			Offset   int                               `json:"offset"`
			Prev     *string                           `json:"prev"`
			Next     *string                           `json:"next"`
		}{
			Problems: resource.GetProblemProfileForAdminSlice(problems),
			Total:    total,
			Count:    len(problems),
			Offset:   req.Offset,
			Prev:     prevUrl,
			Next:     nextUrl,
		},
	})
}

func AdminUpdateProblem(c echo.Context) error {
	initCreator.Do(initCreatorRole)
	req := request.AdminUpdateProblemRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}
	problem, err := findProblem(c.Param("id"))
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
	}
	problem.Name = req.Name
	problem.Description = req.Description
	problem.AttachmentFileName = req.AttachmentFileName
	if req.Public != nil {
		problem.Public = *req.Public
	} else {
		problem.Public = false
	}
	if req.Privacy != nil {
		problem.Privacy = *req.Privacy
	} else {
		problem.Public = true
	}
	problem.MemoryLimit = req.MemoryLimit
	problem.TimeLimit = req.TimeLimit
	problem.LanguageAllowed = req.LanguageAllowed
	problem.CompileEnvironment = req.CompileEnvironment
	problem.CompareScriptID = req.CompareScriptID
	utils.PanicIfDBError(base.DB.Save(&problem), "could not update problem")
	return c.JSON(http.StatusOK, response.AdminUpdateProblemResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ProblemProfileForAdmin `json:"problem"`
		}{
			resource.GetProblemProfileForAdmin(problem),
		},
	})
}

func AdminDeleteProblem(c echo.Context) error {
	initCreator.Do(initCreatorRole)
	problem, err := findProblem(c.Param("id"))
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
	}
	utils.PanicIfDBError(base.DB.Delete(&models.UserHasRole{}, "role_id = ? and target_id = ?", problemCreator.ID, problem.ID), "could not delete problem")
	utils.PanicIfDBError(base.DB.Delete(&problem), "could not delete problem")
	return c.JSON(http.StatusNoContent, response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}

func findProblem(id string) (*models.Problem, error) {
	problem := models.Problem{}
	err := base.DB.Where("id = ?", id).First(&problem).Error
	if err != nil {
		err = base.DB.Where("name = ?", id).First(&problem).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, err
			} else {
				panic(errors.Wrap(err, "could not query problem"))
			}
		}
	}
	return &problem, nil
}
