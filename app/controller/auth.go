package controller

import (
	"github.com/go-playground/validator/v10"
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

func Login(c echo.Context) error {
	return nil
}

func Register(c echo.Context) error {
	req := new(request.UserRequest)
	if err := c.Bind(req); err != nil {
		log.Error(errors.Wrap(err, "could not bind request "), c)
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    -1,
			Message: "Internal error",
			Error:   nil,
			Data:    nil,
		})
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
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    1,
				Message: "validation error",
				Error:   validationErrors,
				Data:    nil,
			})
		}
		log.Error(errors.Wrap(err, "validate failed"), c)
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    -1,
			Message: "Internal error",
			Error:   nil,
			Data:    nil,
		})
	}
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Error(errors.Wrap(err, "could not hash user password"))
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    -1,
			Message: "Internal error",
			Error:   nil,
			Data:    nil,
		})
	}
	count := 0
	base.DB.Model(&models.User{}).Where("email = ? or username = ?", req.Email, req.Nickname).Count(&count)
	if base.DB.Error != nil {
		log.Error(errors.Wrap(err, "could not query user count"))
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    -1,
			Message: "Internal error",
			Error:   nil,
			Data:    nil,
		})
	}
	if count != 0 {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    2,
			Message: "duplicate username or email",
			Error:   nil,
			Data:    nil,
		})
	}
	user := models.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Password: hashed,
	}
	err = base.DB.Create(&user).Error
	if err != nil {
		log.Error(errors.Wrap(err, "could not create user"))
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    -1,
			Message: "Internal error",
			Error:   nil,
			Data:    nil,
		})
	}
	token := models.Token{
		Token: utils.RandStr(32),
		User:  user,
	}
	err = base.DB.Create(&token).Error
	if err != nil {
		log.Error(errors.Wrap(err, "could not create token for user"))
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    -1,
			Message: "Internal error",
			Error:   nil,
			Data:    nil,
		})
	}
	return c.JSON(http.StatusCreated, response.RegisterResponse{
		Code:    0,
		Message: "success",
		Error:   nil,
		Data: struct {
			models.User  `json:"user"`
			Token string `json:"token"`
		}{
			user,
			token.Token,
		},
	})
}
