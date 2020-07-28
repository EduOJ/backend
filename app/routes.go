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

	api := e.Group("/api", middleware.Authentication)

	auth := api.Group("/auth", middleware.Auth)
	auth.POST("/login", controller.Login).Name = "auth.login"
	auth.POST("/register", controller.Register).Name = "auth.register"

	loginCheck := api.Group("/", middleware.LoginCheck)

	_ = loginCheck
	// TODO: routes.
}
