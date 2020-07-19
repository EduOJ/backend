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
	req := new(request.RegisterRequest)
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
			return c.JSON(http.StatusBadRequest, response.ErrorResp(1, "validation error", validationErrors))
		}
		log.Error(errors.Wrap(err, "validate failed"), c)
		return response.InternalErrorResp(c)
	}
	hashed := utils.HashPassword(req.Password)
	count := 0
	panicIfDBError(base.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count), "could not query user count")
	if count != 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResp(2, "duplicate email", nil))
	}
	panicIfDBError(base.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count), "could not query user count")
	if count != 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResp(3, "duplicate username", nil))
	}
	user := models.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Password: hashed,
	}
	panicIfDBError(base.DB.Create(&user), "could not create user")
	token := models.Token{
		Token: utils.RandStr(32),
		User:  user,
	}
	panicIfDBError(base.DB.Create(&token), "could not create token for user")
	return c.JSON(http.StatusCreated, response.RegisterResponse{
		Code:    0,
		Message: "success",
		Error:   nil,
		Data: struct {
			models.User `json:"user"`
			Token       string `json:"token"`
		}{
			user,
			token.Token,
		},
	})
}
