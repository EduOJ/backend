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
	"net/http"
)

func GetProblem(c echo.Context) error {
	problem, err := findProblem(c.Param("id"), true)
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
	}
	// TODO: load test cases
	return c.JSON(http.StatusOK, response.GetProblemResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ProblemProfile `json:"problem"`
		}{
			resource.GetProblemProfile(problem),
		},
	})
}

func GetProblems(c echo.Context) error {
	req := request.GetProblemsRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}

	query, err := utils.Sorter(base.DB.Model(&models.Problem{}), req.OrderBy, "id")
	if err != nil {
		if herr, ok := err.(utils.HttpError); ok {
			return herr.Response(c)
		}
		panic(err)
	}

	query = query.Where("public = ?", true)

	if req.Search != "" {
		query = query.Where("id like ? or name like ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	var problems []models.Problem
	total, prevUrl, nextUrl, err := utils.Paginator(query, req.Limit, req.Offset, c.Request().URL, &problems)
	return c.JSON(http.StatusOK, response.GetProblemsResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Problems []resource.ProblemProfile `json:"problems"`
			Total    int                       `json:"total"`
			Count    int                       `json:"count"`
			Offset   int                       `json:"offset"`
			Prev     *string                   `json:"prev"`
			Next     *string                   `json:"next"`
		}{
			Problems: resource.GetProblemProfileSlice(problems),
			Total:    total,
			Count:    len(problems),
			Offset:   req.Offset,
			Prev:     prevUrl,
			Next:     nextUrl,
		},
	})
}
