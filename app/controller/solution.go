package controller

import (
	"net/http"

	"github.com/EduOJ/backend/app/request"
	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/app/response/resource"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/utils"
	"github.com/EduOJ/backend/database/models"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func GetSolution(c echo.Context) error {
	solution, err := utils.FindSolution(c.Param("id"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
	}
	return c.JSON(http.StatusOK, response.GetSolutionResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.Solution `json:"solution"`
		}{
			resource.GetSolution(solution),
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
		Likes:       req.Likes,
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
