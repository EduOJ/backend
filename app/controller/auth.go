package controller

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"net/http"
)

func Login(c echo.Context) error {
	req := new(request.LoginRequest)
	if err, ok := utils.BindAndValidate(req, &c); !ok {
		return err
	}
	user := models.User{}
	err := base.DB.Where("email = ? or username = ?", req.UsernameOrEmail, req.UsernameOrEmail).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, response.ErrorResp("WRONG_USERNAME", nil))
		} else {
			panic(errors.Wrap(err, "could not query username or email"))
		}
	}
	if !utils.VerifyPassword(req.Password, user.Password) {
		return c.JSON(http.StatusForbidden, response.ErrorResp("WRONG_PASSWORD", nil))
	}
	token := models.Token{
		Token:      utils.RandStr(32),
		User:       user,
		RememberMe: req.RememberMe,
	}
	utils.PanicIfDBError(base.DB.Create(&token), "could not create token for users")
	return c.JSON(http.StatusOK, response.RegisterResponse{
		Message: "SUCCESS",
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

func Register(c echo.Context) error {
	req := new(request.RegisterRequest)
	err, ok := utils.BindAndValidate(req, &c)
	if !ok {
		return err
	}
	hashed := utils.HashPassword(req.Password)
	count := 0
	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count), "could not query user count")
	if count != 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("DUPLICATE_EMAIL", nil))
	}
	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count), "could not query user count")
	if count != 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("DUPLICATE_USERNAME", nil))
	}
	user := models.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Password: hashed,
	}
	utils.PanicIfDBError(base.DB.Create(&user), "could not create user")
	token := models.Token{
		Token: utils.RandStr(32),
		User:  user,
	}
	utils.PanicIfDBError(base.DB.Create(&token), "could not create token for user")
	return c.JSON(http.StatusCreated, response.RegisterResponse{
		Message: "SUCCESS",
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
