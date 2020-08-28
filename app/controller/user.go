package controller

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

func GetUser(c echo.Context) error {

	user, err := utils.FindUser(c.Param("id"))
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
	}
	return c.JSON(http.StatusOK, response.GetUserResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*response.UserProfile `json:"user"`
		}{
			response.GetUserProfile(user),
		},
	})
}

func GetMe(c echo.Context) error {
	var user models.User
	var ok bool
	if user, ok = c.Get("user").(models.User); !ok {
		panic("could not convert my user into type models.User")
	}
	if !user.RoleLoaded {
		user.LoadRoles()
	}
	return c.JSON(http.StatusOK, response.GetMeResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*response.UserProfileForMe `json:"user"`
		}{
			response.GetUserProfileForMe(&user),
		},
	})
}

func GetUsers(c echo.Context) error {
	req := new(request.GetUsersRequest)
	if err, ok := utils.BindAndValidate(req, c); !ok {
		return err
	}
	var users []models.User
	var total int

	query := base.DB.Model(&models.User{})
	if req.OrderBy != "" {
		order := strings.SplitN(req.OrderBy, ".", 2)
		if len(order) != 2 {
			return c.JSON(http.StatusBadRequest, response.ErrorResp("INVALID_ORDER", nil))
		}
		if !utils.Contain(order[0], []string{"username", "id", "nickname"}) {
			return c.JSON(http.StatusBadRequest, response.ErrorResp("INVALID_ORDER", nil))
		}
		if !utils.Contain(order[1], []string{"ASC", "DESC"}) {
			return c.JSON(http.StatusBadRequest, response.ErrorResp("INVALID_ORDER", nil))
		}
		query = query.Order(strings.Join(order, " "))
	}
	if req.Search != "" {
		query = query.Where("id like ? or username like ? or email like ? or nickname like ?", "%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if req.Limit == 0 {
		req.Limit = 20 // Default limit
	}
	err := query.Limit(req.Limit).Offset(req.Offset).Find(&users).Error
	if err != nil {
		panic(errors.Wrap(err, "could not query users"))
	}
	err = query.Count(&total).Error
	if err != nil {
		panic(errors.Wrap(err, "could not query count of users"))
	}

	var nextUrlStr *string
	var prevUrlStr *string

	if req.Offset-req.Limit >= 0 {
		prevURL := c.Request().URL
		q, err := url.ParseQuery(prevURL.RawQuery)
		if err != nil {
			panic(errors.Wrap(err, "could not parse query for url"))
		}
		q.Set("offset", fmt.Sprint(req.Offset-req.Limit))
		q.Set("limit", fmt.Sprint(req.Limit))
		prevURL.RawQuery = q.Encode()
		temp := prevURL.String()
		prevUrlStr = &temp
	} else {
		prevUrlStr = nil
	}
	if req.Offset+len(users) < total {
		nextURL := c.Request().URL
		q, err := url.ParseQuery(nextURL.RawQuery)
		if err != nil {
			panic(errors.Wrap(err, "could not parse query for url"))
		}
		q.Set("offset", fmt.Sprint(req.Offset+req.Limit))
		q.Set("limit", fmt.Sprint(req.Limit))
		nextURL.RawQuery = q.Encode()
		temp := nextURL.String()
		nextUrlStr = &temp
	} else {
		nextUrlStr = nil
	}

	return c.JSON(http.StatusOK, response.GetUsersResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Users  []response.UserProfile `json:"users"`
			Total  int                    `json:"total"`
			Count  int                    `json:"count"`
			Offset int                    `json:"offset"`
			Prev   *string                `json:"prev"`
			Next   *string                `json:"next"`
		}{
			response.GetUserProfileSlice(users),
			total,
			len(users),
			req.Offset,
			prevUrlStr,
			nextUrlStr,
		},
	})
}

func UpdateMe(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		panic("could not convert my user into type models.User")
	}
	if !user.RoleLoaded {
		user.LoadRoles()
	}
	req := request.UpdateMeRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	count := 0
	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count), "could not query user count")
	if count > 1 || (count == 1 && user.Email != req.Email) {
		return c.JSON(http.StatusConflict, response.ErrorResp("CONFLICT_EMAIL", nil))
	}
	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count), "could not query user count")
	if count > 1 || (count == 1 && user.Username != req.Username) {
		return c.JSON(http.StatusConflict, response.ErrorResp("CONFLICT_USERNAME", nil))
	}
	user.Username = req.Username
	user.Nickname = req.Nickname
	user.Email = req.Email
	utils.PanicIfDBError(base.DB.Save(&user), "could not update user")
	return c.JSON(http.StatusOK, response.UpdateMeResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*response.UserProfileForMe `json:"user"`
		}{
			response.GetUserProfileForMe(&user),
		},
	})
}

func ChangePassword(c echo.Context) error {
	req := request.ChangePasswordRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	user, ok := c.Get("user").(models.User)
	if !ok {
		panic("could not get user from context")
	}
	if !utils.VerifyPassword(req.OldPassword, user.Password) {
		return c.JSON(http.StatusForbidden, response.ErrorResp("WRONG_PASSWORD", nil))
	}
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		panic("could not get tokenString from request header")
	}
	utils.PanicIfDBError(base.DB.Where("user_id = ? and token != ?", user.ID, tokenString).Delete(models.Token{}), "could not remove token")
	user.Password = utils.HashPassword(req.NewPassword)
	base.DB.Save(&user)
	return c.JSON(http.StatusOK, response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}
