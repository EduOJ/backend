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
			authConf, err := config.Get("database")
			sessionTimeoutInt := 168
			if err != nil || authConf == nil {
				log.Warning("Cannot read auth config")
			} else {
				sessionTimeoutInt = authConf.MustGet("session_timeout", 168).Value().(int)
			}
			sessionTimeout = time.Second * time.Duration(sessionTimeoutInt*3600)
			if err != nil {
				log.Error(errors.Wrap(err, "ParseDuration failed"), c)
				return response.InternalErrorResp(c)
			}
		}
		if time.Now().Add(sessionTimeout).Before(token.UpdatedAt) {
			base.DB.Delete(&token)
			return c.JSON(http.StatusRequestTimeout, response.ErrorResp(1, "session expired", nil))
		}
		token.UpdatedAt = time.Now()
		utils.PanicIfDBError(base.DB.Save(&token), "could not update token")
		c.Set("token", token)
		log.Debug(token.UpdatedAt.String())
		return next(c)
	}
}

func LoginCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("token")
		if token == nil {
			return c.JSON(http.StatusUnauthorized, response.ErrorResp(1, "Unauthorized", nil))
		}
		return next(c)
	}
}
