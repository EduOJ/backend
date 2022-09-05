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
			Solutions: resource.GetSolutionSummarySlice(solutions),
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
