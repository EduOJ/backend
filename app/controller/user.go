package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"net/http"
)

func ChangePassword(c echo.Context) error {
	req := new(request.ChangePasswordRequest)
	//TODO: use bind and validate
	if err := c.Bind(req); err != nil {
		panic(err)
	}
	if err := c.Validate(req); err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			validationErrors := make([]response.ValidationError, len(e))
			for i, v := range e {
				validationErrors[i] = response.ValidationError{
					Field:  v.Field(),
					Reason: v.Tag(),
				}
			}
			return c.JSON(http.StatusBadRequest, response.ErrorResp("VALIDATION_ERROR", validationErrors))
		}
		log.Error(errors.Wrap(err, "validate failed"), c)
		return response.InternalErrorResp(c)
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

func GetUser(c echo.Context) error {

	user, err, ok := findUser(c)
	if !ok {
		return err
	}
	return c.JSON(http.StatusOK, response.GetUserResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			models.User `json:"user"`
		}{
			user,
		},
	})
}

func GetUsers(c echo.Context) error {
	req := new(request.GetUsersRequest)
	err, ok := utils.BindAndValidate(req, &c)
	if !ok {
		return err
	}
	var users []models.User
	if req.OrderBy != "" {
		err = base.DB.Where("username like ? and nickname like ?", "%"+req.Username+"%", "%"+req.Nickname+"%").Order(req.OrderBy).Find(&users).Error
	} else {
		err = base.DB.Where("username like ? and nickname like ?", "%"+req.Username+"%", "%"+req.Nickname+"%").Find(&users).Error
	}
	if err != nil {
		panic(errors.Wrap(err, "could not query users"))
	}
	if req.Offset > len(users) {
		return c.JSON(http.StatusNotFound, response.ErrorResp("GET_USERS_OFFSET_OUT_OF_BOUNDS", nil))
	}
	if req.Limit > 0 && req.Offset+req.Limit < len(users) {
		users = users[req.Offset : req.Offset+req.Limit]
	} else {
		users = users[req.Offset:]
	}
	return c.JSON(http.StatusOK, response.GetUsersResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Users  []models.User `json:"users"`
			Limit  int           `json:"limit"`
			Offset int           `json:"offset"`
			Prev   string        `json:"prev"`
			Next   string        `json:"next"`
		}{
			users,
			req.Limit,
			req.Offset,
			"",
			"", //TODO:fill this
		},
	})
}

func findUser(c echo.Context) (models.User, error, bool) {
	id := c.Param("id")
	user := models.User{}
	err := base.DB.Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		err = base.DB.Where("username = ?", id).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return user, c.JSON(http.StatusNotFound, response.ErrorResp("USER_NOT_FOUND", nil)), false
			} else {
				panic(errors.Wrap(err, "could not query username"))
			}
		}
	} else if err != nil {
		panic(errors.Wrap(err, "could not query id"))
	}
	return user, nil, true
}

