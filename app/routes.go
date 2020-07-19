package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/controller"
	"github.com/leoleoasd/EduOJBackend/app/middleware"
)

func Register(e *echo.Echo) {
	e.Validator = &Validator{
		v: validator.New(),
	}
	e.Use(middleware.Recover)

	api := e.Group("/api")

	auth := api.Group("/auth", middleware.Auth)
	auth.POST("/login", controller.Login).Name = "auth.login"
	auth.POST("/register", controller.Register).Name = "auth.register"

	// TODO: routes.
}
