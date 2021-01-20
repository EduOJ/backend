package middleware

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return next(c)
		}
		token, err := utils.GetToken(tokenString)
		if err == gorm.ErrRecordNotFound {
			return next(c)
		}
		if err != nil {
			log.Error(errors.Wrap(err, "fail to get user from token"), c)
			return response.InternalErrorResp(c)
		}
		if utils.IsTokenExpired(token) {
			base.DB.Delete(&token)
			return c.JSON(http.StatusRequestTimeout, response.ErrorResp("AUTH_SESSION_EXPIRED", nil))
		}
		token.UpdatedAt = time.Now()
		utils.PanicIfDBError(base.DB.Save(&token), "could not update token")
		c.Set("user", token.User)
		return next(c)
	}
}

func Logged(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		if user == nil {
			return c.JSON(http.StatusUnauthorized, response.ErrorResp("AUTH_NEED_TOKEN", nil))
		}
		return next(c)
	}
}

// Using this middleware means the controller could accept guest users, whose most information doesn't exist.
// The only exception is roles information, guest users don't have any permissions or roles.
// Any User information other than roles SHOULD NOT be read in controllers that use this middleware.
func AllowGuest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		guestUser := models.User{
			ID:        0,
			Username:  "guest_user",
			Nickname:  "guest_user_nick",
			Email:     "guest_user@email.com",
			Password:  "guest_user_pwd",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: nil,
			// The above content is to filled for debug, and should not be used in formal applications.

			Roles:      []models.UserHasRole{},
			RoleLoaded: true,
		}
		if c.Get("user") == nil {
			c.Set("user", guestUser)
		}
		return next(c)
	}
}
