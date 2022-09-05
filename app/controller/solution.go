package controller

import (
	"errors"
	"fmt"
	"github.com/EduOJ/backend/app/request"
	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/app/response/resource"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/utils"
	"github.com/EduOJ/backend/database/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

// TODO
func GetSolution(c echo.Context) error {
	user := c.Get("user").(models.User)
	solution, err := utils.FindSolution(c.Param("id"), &user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
	}
	if user.Can("read_solution_secrets", solution) || user.Can("read_solution_secrets") {
		return c.JSON(http.StatusOK, response.GetSolutionResponseForAdmin{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.SolutionForAdmin `json:"problem"`
			}{
				resource.GetSolutonForAdmin(solution),
			},
		})
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

// TODO
func GetSolutions(c echo.Context) error {
	req := request.GetSolutionsRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}

	query := base.DB.Model(&models.Solution{}).Order("id ASC").Omit("Description") // Force order by id asc.

	user := c.Get("user").(models.User)
	isAdmin := user.Can("manage_solution")
	if !isAdmin {
		query = query.Where("public = ?", true)
	}

	if req.Search != "" {
		id, _ := strconv.ParseUint(req.Search, 10, 64)
		query = query.Where("id = ? or name like ?", id, "%"+req.Search+"%")
	}

	where := (*gorm.DB)(nil)

	if req.Tags != "" {
		tags := strings.Split(req.Tags, ",")
		query = query.Where("id in (?)",
			base.DB.Table("tags").
				Where("name in (?)", tags).
				Group("solution_id").
				Having("count(*) = ?", len(tags)).
				Select("solution_id"),
		)
	}

	if where != nil {
		query = query.Where(where)
	}

	var solutions []*models.Solution
	total, prevUrl, nextUrl, err := utils.Paginator(query.WithContext(c.Request().Context()).Preload("Tags"), req.Limit, req.Offset, c.Request().URL, &solutions)
	if err != nil {
		if herr, ok := err.(utils.HttpError); ok {
			return herr.Response(c)
		}
		panic(err)
	}

	if isAdmin {
		return c.JSON(http.StatusOK, response.GetSolutionsResponseForAdmin{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				Solutions []resource.SolutionSummaryForAdmin `json:"solutions"`
				Total     int                                `json:"total"`
				Count     int                                `json:"count"`
				Offset    int                                `json:"offset"`
				Prev      *string                            `json:"prev"`
				Next      *string                            `json:"next"`
			}{
				Solutions: resource.GetSolutionSummaryForAdminSlice(solutions, passed),
				Total:     total,
				Count:     len(solutions),
				Offset:    req.Offset,
				Prev:      prevUrl,
				Next:      nextUrl,
			},
		})
	}
	return c.JSON(http.StatusOK, response.GetSolutionsResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Solutions []resource.SolutionSummary `json:"solutions"`
			Total     int                        `json:"total"`
			Count     int                        `json:"count"`
			Offset    int                        `json:"offset"`
			Prev      *string                    `json:"prev"`
			Next      *string                    `json:"next"`
		}{
			Solutions: resource.GetSolutionSummarySlice(solutions, passed),
			Total:     total,
			Count:     len(solutions),
			Offset:    req.Offset,
			Prev:      prevUrl,
			Next:      nextUrl,
		},
	})
}

// TODO
func CreateSolution(c echo.Context) error {
	file, err := c.FormFile("attachment_file")
	if err != nil && err != http.ErrMissingFile && err.Error() != "request Content-Type isn't multipart/form-data" {
		panic(errors.Wrap(err, "could not read file"))
	}

	req := request.CreateSolutionRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	var public, privacy bool
	if req.Public == nil {
		public = false
	} else {
		public = *req.Public
	}
	if req.Privacy == nil {
		privacy = false
	} else {
		privacy = *req.Privacy
	}

	solution := models.Solution{
		Name:        req.Name,
		Description: req.Description,
	}
	utils.PanicIfDBError(base.DB.Create(&solution), "could not create solution")
	user := c.Get("user").(models.User)
	user.GrantRole("solution_creator", solution)

	if file != nil {
		utils.MustPutObject(file, c.Request().Context(), "solutions", fmt.Sprintf("%d/attachment", solution.ID))
	}

	return c.JSON(http.StatusCreated, response.CreateSolutionResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.SolutionForAdmin(&solution)
		},
	})

}
