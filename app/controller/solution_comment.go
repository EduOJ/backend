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

func GetSolutionComments(c echo.Context) error {
	sc, err := utils.FindSolutionComments(c.Param("id"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
	}
	return c.JSON(http.StatusOK, response.GetSolutionCommentsResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			SolutionComments []resource.SolutionComment `json:"solution_comments"`
		}{
			resource.GetSolutionComments(sc),
		},
	})
}

func CreateSolutions(c echo.Context) error {
	req := request.CreateSolutionCommentRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}

	comment := models.SolutionComment{
		SolutionID:  req.SolutionID,
		FatherNode:  req.FatherNode,
		Description: req.Description,
		Speaker:     req.Speaker,
	}
	utils.PanicIfDBError(base.DB.Create(&comment), "could not create comment")

	return c.JSON(http.StatusOK, response.CreateSolutionCommentResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.SolutionComment `json:"solution_comment"`
		}{
			resource.GetSolutionComment(&comment),
		},
	})

}
