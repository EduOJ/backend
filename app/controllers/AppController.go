package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/base"
	"net/http"
)

func Recv(c echo.Context) error {
	base.Redis.Set(c.Request().Context(), "123", 123, 0)
	v := base.Redis.Get(c.Request().Context(), "123")
	return c.String(http.StatusOK, v.Val())
}

func Send(c echo.Context) error {
	base.Redis.Set(c.Request().Context(), "123", 123, 0)
	v := base.Redis.Get(c.Request().Context(), "123")
	return c.String(http.StatusOK, v.Val())
}
