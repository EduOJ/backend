package controller

import (
	"fmt"
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
	"net/url"
	"strings"
)

func AdminCreateProblem(c echo.Context) error {
	req := request.AdminCreateProblemRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	count := 0
	utils.PanicIfDBError(base.DB.Model(&models.Problem{}).Where("name = ?", req.Name).Count(&count), "could not query problem count")
	if count != 0 {
		return c.JSON(http.StatusConflict, response.ErrorResp("CONFLICT_NAME", nil))
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
	var problems []models.Problem
	var total int

	query := base.DB.Model(&models.Problem{})
	if req.OrderBy != "" {
		order := strings.SplitN(req.OrderBy, ".", 2)
		if len(order) != 2 {
			return c.JSON(http.StatusBadRequest, response.ErrorResp("INVALID_ORDER", nil))
		}
		if !utils.Contain(order[0], []string{"id", "name"}) {
			return c.JSON(http.StatusBadRequest, response.ErrorResp("INVALID_ORDER", nil))
		}
		if !utils.Contain(order[1], []string{"ASC", "DESC"}) {
			return c.JSON(http.StatusBadRequest, response.ErrorResp("INVALID_ORDER", nil))
		}
		query = query.Order(strings.Join(order, " "))
	}
	if req.Search != "" {
		query = query.Where("id like ? or name like ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if req.Limit == 0 {
		req.Limit = 20 // Default limit
	}
	err := query.Limit(req.Limit).Offset(req.Offset).Find(&problems).Error
	if err != nil {
		panic(errors.Wrap(err, "could not query problems"))
	}
	err = query.Count(&total).Error
	if err != nil {
		panic(errors.Wrap(err, "could not query count of problems"))
	}

	var nextUrlStr *string
	var prevUrlStr *string

	if req.Offset-req.Limit >= 0 {
		prevURL := c.Request().URL
		q, err := url.ParseQuery(prevURL.RawQuery)
		if err != nil {
			panic(errors.Wrap(err, "could not parse query for url"))
		}
		q.Set("offset", fmt.Sprint(req.Offset-req.Limit))
		q.Set("limit", fmt.Sprint(req.Limit))
		prevURL.RawQuery = q.Encode()
		temp := prevURL.String()
		prevUrlStr = &temp
	} else {
		prevUrlStr = nil
	}
	if req.Offset+len(problems) < total {
		nextURL := c.Request().URL
		q, err := url.ParseQuery(nextURL.RawQuery)
		if err != nil {
			panic(errors.Wrap(err, "could not parse query for url"))
		}
		q.Set("offset", fmt.Sprint(req.Offset+req.Limit))
		q.Set("limit", fmt.Sprint(req.Limit))
		nextURL.RawQuery = q.Encode()
		temp := nextURL.String()
		nextUrlStr = &temp
	} else {
		nextUrlStr = nil
	}
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
			resource.GetProblemProfileForAdminSlice(problems),
			total,
			len(problems),
			req.Offset,
			prevUrlStr,
			nextUrlStr,
		},
	})
}

func AdminUpdateProblem(c echo.Context) error {
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
	count := 0
	utils.PanicIfDBError(base.DB.Model(&models.Problem{}).Where("name = ?", req.Name).Count(&count), "could not query problem count")
	if count > 1 || (count == 1 && problem.Name != req.Name) {
		return c.JSON(http.StatusConflict, response.ErrorResp("CONFLICT_NAME", nil))
	}
	problem.Name = req.Name
	problem.Description = req.Description
	problem.AttachmentFileName = req.AttachmentFileName
	problem.Public = *req.Public
	problem.Privacy = *req.Privacy
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
	problem, err := findProblem(c.Param("id"))
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
	}
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
