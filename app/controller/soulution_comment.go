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
	sc := []models.SolutionComment{}
	query := base.DB
	err := query.Where("solution_id = ?", c.Param("solutionId")).Find(&sc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not query solution comment"))
		}
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

func CreateSolutionComment(c echo.Context) error {
	req := request.CreateSolutionCommentRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	var FatherNode uint
	if req.IsRoot {
		FatherNode = 0
	} else {
		FatherNode = req.FatherNode
	}
	comment := models.SolutionComment{
		SolutionID:  req.SolutionID,
		FatherNode:  FatherNode,
		Description: req.Description,
		Speaker:     req.Speaker,
	}
	utils.PanicIfDBError(base.DB.Create(&comment), "could not create comment")
	return c.JSON(http.StatusOK, response.CreateSolutionCommentResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.SolutionComment `json:"solution_comment_create"`
		}{
			resource.GetSolutionComment(&comment),
		},
	})
}

func GetCommentTree(c echo.Context) error {
	solutionComments := []models.SolutionComment{}
	query := base.DB
	err := query.Where("SolutionID = ?", c.Param("id")).Find(&solutionComments).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not query solution comment"))
		}
	}

	commentNodes := make([]resource.SolutionCommentNode, len(solutionComments))
	for i, solutionComment := range solutionComments {
		commentNodes[i].ConvertCommentToNode(&solutionComment)
	}

	return c.JSON(http.StatusOK, response.GetCommentTreeResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.SolutionCommentTree `json:"solution_comment_tree"`
		}{
			resource.GetSolutionCommentTree(commentNodes),
		},
	})
}
