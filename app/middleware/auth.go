package middleware

import 	"github.com/labstack/echo/v4"

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// do nothing.
		// we allow register/login when already logged in from api side.
		return next(c)
	}
}
