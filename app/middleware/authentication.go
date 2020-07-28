package middleware

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

var sessionTimeout time.Duration

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("token")
		if tokenString == "" {
			return next(c)
		}
		token, err := utils.GetToken(tokenString)
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusUnauthorized, response.ErrorResp(1, "Unauthorized", nil))
		}
		if err != nil {
			log.Error(errors.Wrap(err, "fail to get user from token"), c)
			return response.InternalErrorResp(c)
		}
		if sessionTimeout == 0 {
			sessionTimeoutInt := config.MustGet("auth.session_timeout", 168).Value().(int)
			sessionTimeout = -1 * time.Second * time.Duration(sessionTimeoutInt*3600)
		}
		if time.Now().Add(sessionTimeout).After(token.UpdatedAt) {
			base.DB.Delete(&token)
			return c.JSON(http.StatusRequestTimeout, response.ErrorResp(1, "session expired", nil))
		}
		token.UpdatedAt = time.Now()
		utils.PanicIfDBError(base.DB.Save(&token), "could not update token")
		c.Set("user", token.User)
		return next(c)
		//TODO:remember me
		//TODO:delete earliest token if one user have too much token
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
