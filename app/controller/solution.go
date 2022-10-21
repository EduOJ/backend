package controller

import (
	"net/http"
	"strconv"

	"github.com/EduOJ/backend/app/request"
	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/app/response/resource"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/utils"
	"github.com/EduOJ/backend/database/models"
	"github.com/labstack/echo/v4"
)

func GetSolutions(c echo.Context) error {
	req := request.GetSolutionsRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}

	problemID, _ := strconv.ParseUint(req.ProblemID, 10, 64)
	query := base.DB.Model(&models.Solution{}).Order("problem_id ASC").Where("problem_id = ?", problemID)
	var solutions []*models.Solution

	err := query.Find(&solutions).Error
	if err != nil {
		if herr, ok := err.(utils.HttpError); ok {
			return herr.Response(c)
		}
		panic(err)
	}

	return c.JSON(http.StatusOK, response.GetSolutionsResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Solutions []resource.Solution `json:"solutions"`
		}{
			Solutions: resource.GetSolutions(solutions),
		},
	})
}

func CreateSolution(c echo.Context) error {
	req := request.CreateSolutionRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}

	solution := models.Solution{
		ProblemID:   req.ProblemID,
		Name:        req.Name,
		Author:      req.Author,
		Description: req.Description,
		Likes:       0,
	}
	utils.PanicIfDBError(base.DB.Create(&solution), "could not create solution")

	return c.JSON(http.StatusCreated, response.CreateSolutionResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.Solution `json:"solution"`
		}{
			resource.GetSolution(&solution),
		},
	})
}
