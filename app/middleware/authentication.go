package middleware

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/base/utils"
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
