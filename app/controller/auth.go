package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"net/http"
)

var tokenCountMaximum int

func Login(c echo.Context) error {
	req := new(request.LoginRequest)
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
	t := base.DB.HasTable("users")
	_ = t
	user := models.User{}
	t = base.DB.HasTable("users")
	err := base.DB.Where("email = ? or username = ?", req.UsernameOrEmail, req.UsernameOrEmail).First(&user).Error
	t = base.DB.HasTable("users")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, response.ErrorResp(2, "wrong username or email", nil))
		} else {
			panic(errors.Wrap(err, "could not query username or email"))
		}
	}
	if !utils.VerifyPassword(req.Password, user.Password) {
		return c.JSON(http.StatusForbidden, response.ErrorResp(3, "wrong password", nil))
	}
	token := models.Token{
		Token:      utils.RandStr(32),
		User:       user,
		RememberMe: req.RememberMe,
	}
	utils.PanicIfDBError(base.DB.Create(&token), "could not create token for user")
	if tokenCountMaximum == 0 {
		tokenCountMaximum = config.MustGet("auth.token_count_maximum", 10).Value().(int)
	}
	var lastPassToken models.Token
	utils.PanicIfDBError(base.DB.Preload("User").Model(models.Token{}).Where("user_id = ?", &user.ID).Order("updated_at desc").Limit(tokenCountMaximum).Find(&lastPassToken), "could not find and order tokens for user")
	utils.PanicIfDBError(base.DB.Preload("User").Model(models.Token{}).Delete(models.Token{}, "user_id = ? AND updated_at < ?", user.ID, lastPassToken.UpdatedAt), "could not delete tokens for user")
	return c.JSON(http.StatusOK, response.RegisterResponse{
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
	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count), "could not query user count")
	if count != 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResp(2, "duplicate email", nil))
	}
	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count), "could not query user count")
	if count != 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResp(3, "duplicate username", nil))
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
