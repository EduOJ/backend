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
			return c.JSON(http.StatusUnauthorized, response.ErrorResp(1, "Token not found", nil))
		}
		if err != nil {
			log.Error(errors.Wrap(err, "fail to get user from token"), c)
			return response.InternalErrorResp(c)
		}
		if utils.IsTokenExpired(token) {
			base.DB.Delete(&token)
			return c.JSON(http.StatusRequestTimeout, response.ErrorResp(1, "session expired", nil))
		}
		token.UpdatedAt = time.Now()
		utils.PanicIfDBError(base.DB.Save(&token), "could not update token")
		c.Set("user", token.User)
		return next(c)
	}
}

func LoginCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		if user == nil {
			return c.JSON(http.StatusUnauthorized, response.ErrorResp(1, "Unauthorized", nil))
		}
		return next(c)
	}
}
