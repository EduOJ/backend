package controller

import (
	"net/http"

	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/app/response/resource"
	"github.com/EduOJ/backend/base/utils"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func GetCommentTree(c echo.Context) error {
	solutionComments, err := utils.FindSolutionComments(c.Param("solution_id"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
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
