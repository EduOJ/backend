package middleware

import (
	"net/http"
	"time"

	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/log"
	"github.com/EduOJ/backend/base/utils"
	"github.com/EduOJ/backend/database/models"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return next(c)
		}
		token, err := utils.GetToken(tokenString)
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
		utils.PanicIfDBError(base.DB.Omit(clause.Associations).Save(&token), "could not update token")
		c.Set("user", token.User)
		return next(c)
	}
}

func Logged(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, ok := c.Get("user").(models.User)
		if !ok {
			return c.JSON(http.StatusUnauthorized, response.ErrorResp("AUTH_NEED_TOKEN", nil))
		}
		return next(c)
	}
}

func EmailVerified(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(models.User)
		// for some APIs, we allow guest access, but for logged user, we require email verification.
		if ok && viper.GetBool("email.need_verification") && !user.EmailVerified {
			return c.JSON(http.StatusUnauthorized, response.ErrorResp("AUTH_NEED_EMAIL_VERIFICATION", nil))
		}
		return next(c)
	}
}

// Using this middleware means the controller could accept access from guests.
// The only exception is role information, guest users don't have any permissions or roles.
// Any user information other than roles SHOULD NOT be read in controllers that use this middleware.
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
			DeletedAt: gorm.DeletedAt{},
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
