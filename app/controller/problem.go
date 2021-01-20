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
	"net/http"
	"strconv"
	"strings"
)

func GetProblem(c echo.Context) error {
	var user models.User
	var ok bool
	if user, ok = c.Get("user").(models.User); !ok {
		panic("could not convert my user into type models.User")
	}
	problem, err := utils.FindProblem(c.Param("id"), &user)
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
	}
	// TODO: load test cases
	if user.Can("read_problem", problem) {
		return c.JSON(http.StatusOK, response.AdminGetProblemResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.ProblemForAdmin `json:"problem"`
			}{
				resource.GetProblemForAdmin(problem),
			},
		})
	}
	return c.JSON(http.StatusOK, response.GetProblemResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.Problem `json:"problem"`
		}{
			resource.GetProblem(problem),
		},
	})
}

func GetProblems(c echo.Context) error { // TODO: add test for admin check
	req := request.GetProblemsRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}

	query := base.DB.Model(&models.Problem{}).Order("id ASC")

	var user models.User
	var ok bool
	if user, ok = c.Get("user").(models.User); !ok {
		panic("could not convert my user into type models.User")
	}
	isAdmin := user.Can("read_problem")
	query = query.Where("public = ?", !isAdmin)

	if req.Search != "" {
		id, _ := strconv.ParseUint(req.Search, 10, 64)
		query = query.Where("id = ? or name like ?", id, "%"+req.Search+"%")
	}

	var problems []models.Problem
	total, prevUrl, nextUrl, err := utils.Paginator(query, req.Limit, req.Offset, c.Request().URL, &problems)
	if err != nil {
		if herr, ok := err.(utils.HttpError); ok {
			return herr.Response(c)
		}
		panic(err)
	}
	if isAdmin {
		return c.JSON(http.StatusOK, response.AdminGetProblemsResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				Problems []resource.ProblemForAdmin `json:"problems"`
				Total    int                        `json:"total"`
				Count    int                        `json:"count"`
				Offset   int                        `json:"offset"`
				Prev     *string                    `json:"prev"`
				Next     *string                    `json:"next"`
			}{
				Problems: resource.GetProblemForAdminSlice(problems),
				Total:    total,
				Count:    len(problems),
				Offset:   req.Offset,
				Prev:     prevUrl,
				Next:     nextUrl,
			},
		})
	}
	return c.JSON(http.StatusOK, response.GetProblemsResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Problems []resource.Problem `json:"problems"`
			Total    int                `json:"total"`
			Count    int                `json:"count"`
			Offset   int                `json:"offset"`
			Prev     *string            `json:"prev"`
			Next     *string            `json:"next"`
		}{
			Problems: resource.GetProblemSlice(problems),
			Total:    total,
			Count:    len(problems),
			Offset:   req.Offset,
			Prev:     prevUrl,
			Next:     nextUrl,
		},
	})
}

func GetProblemAttachmentFile(c echo.Context) error { // TODO: use MustGetObject
	var user models.User
	var ok bool
	if user, ok = c.Get("user").(models.User); !ok {
		panic("could not convert my user into type models.User")
	}
	problem, err := utils.FindProblem(c.Param("id"), &user)
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusNotFound, response.ErrorResp("PROBLEM_NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
	}
	if problem.AttachmentFileName == "" {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	}
	object := utils.MustGetObject("problems", fmt.Sprintf("%d/attachment", problem.ID))
	contentType := "application/octet-stream"
	if strings.HasSuffix(problem.AttachmentFileName, ".pdf") {
		// If file is a pdf, render it in browser.
		contentType = "application/pdf"
		c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, problem.AttachmentFileName))
	} else {
		// If not, download it as a file.
		c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, problem.AttachmentFileName))
	}
	c.Response().Header().Set("Cache-Control", "public; max-age=31536000")

	return c.Stream(http.StatusOK, contentType, object)
}
