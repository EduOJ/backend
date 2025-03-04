package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
		Likes:       "",
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

func GetLikes(c echo.Context) error {
	req := request.LikesRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	solution := models.Solution{}
	query := base.DB.Model(&models.Solution{})
	err = query.Where("id = ?", req.SolutionId).Find(&solution).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not query solution"))
		}
	}
	// err = query.Where("id = ?", req.SolutionId).Update("likes", "").Error
	// fmt.Println(solution.Likes)
	// "1,2,3,4"
	likeList := strings.Split(solution.Likes, ",")
	// likeList := strings.Split("2,3,4,", ",")
	// fmt.Println(likeList)
	count := len(likeList) - 1
	isLike := false
	// fmt.Println(count)
	switch req.IsLike {
	case 1:
		newLikes := solution.Likes + fmt.Sprint(req.UserId) + ","
		query.Where("id = ?", req.SolutionId).Update("likes", newLikes)
		count = count + 1
		isLike = true
		// fmt.Println(newLikes)
	case -1:
		newLikes := strings.Replace(solution.Likes, fmt.Sprint(req.UserId)+",", "", -1)
		query.Where("id = ?", req.SolutionId).Update("likes", newLikes)
		count = count - 1
		isLike = false
		// fmt.Println(newLikes)

	default:
		for i := 0; i < count; i++ {
			// fmt.Print(i)
			// fmt.Print(": ")
			// fmt.Println(likeList[i])
			if fmt.Sprint(req.UserId) == likeList[i] {
				isLike = true
			}
		}
	}

	likes := &models.Likes{
		Count:  count,
		IsLike: isLike,
	}

	return c.JSON(http.StatusOK, response.GetLikesResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Likes resource.Likes `json:"likes"`
		}{
			Likes: *resource.GetLikes(likes),
		},
	})
}
