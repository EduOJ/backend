package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/pkg/errors"
)

func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if xx := recover(); xx != nil {
				if err, ok := xx.(error); ok {
					log.Error(errors.Wrap(err, "controller panics"))
				} else {
					log.Error("controller panics: ", xx)
				}
				response.InternalErrorResp(c)
			}
		}()
		return next(c)
	}
}
