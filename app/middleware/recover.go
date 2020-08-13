package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/pkg/errors"
	"net/http"
	"runtime/debug"
)

func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		defer func() {
			if xx := recover(); xx != nil {
				if err, ok := xx.(error); ok {
					log.Error(errors.Wrap(err, "controller panics"))
				} else {
					log.Error("controller panics: ", xx)
				}
				// TODO: show call stacks to admins.
				if config.MustGet("debug", false).Value().(bool) {
					stack := debug.Stack()
					err = c.JSON(http.StatusInternalServerError, response.ErrorResp("INTERNAL_ERROR", string(stack)))
				} else {
					err = response.InternalErrorResp(c)
				}
			}
		}()
		err = next(c)
		return
	}
}
