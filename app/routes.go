package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/controller"
	adminController "github.com/leoleoasd/EduOJBackend/app/controller/admin"
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

	loginCheck := api.Group("", middleware.LoginCheck)

	admin := api.Group("/admin")
	admin.POST("/user", adminController.PostUser)
	admin.PUT("/user/:id", adminController.PutUser)
	admin.DELETE("/user/:id", adminController.DeleteUser)
	admin.GET("/user/:id", adminController.GetUser)
	admin.GET("/users", adminController.GetUsers)

	api.GET("/user/:id", controller.Todo)
	api.GET("/users", controller.Todo)

	loginCheck.POST("/user/change_password", controller.ChangePassword)

	// TODO: routes.
}
