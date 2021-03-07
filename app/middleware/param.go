package middleware

import (
	"github.com/EduOJ/backend/app/response"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func ValidateParams(intParams ...string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, p := range intParams {
				param := c.Param(p)
				if param != "" {
					if _, err := strconv.Atoi(param); err != nil {
						return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
					}
				}
			}
			return next(c)
		}
	}
}
