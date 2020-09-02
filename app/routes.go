package app

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/leoleoasd/EduOJBackend/app/controller"
	"github.com/leoleoasd/EduOJBackend/app/middleware"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"net/http"
)

func Register(e *echo.Echo) {
	utils.InitOrigin()
	e.Use(middleware.Recover)
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: utils.Origins,
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	api := e.Group("/api", middleware.Authentication)

	auth := api.Group("/auth", middleware.Auth)
	auth.POST("/login", controller.Login).Name = "auth.login"
	auth.POST("/register", controller.Register).Name = "auth.register"
	auth.GET("/email_registered", controller.EmailRegistered).Name = "auth.emailRegistered"

	admin := api.Group("/admin", middleware.Logged)
	// TODO: add HasPermission
	admin.POST("/user",
		controller.AdminCreateUser, middleware.HasPermission("create_user")).Name = "admin.user.createUser"
	admin.PUT("/user/:id",
		controller.AdminUpdateUser, middleware.HasPermission("update_user")).Name = "admin.user.updateUser"
	admin.DELETE("/user/:id",
		controller.AdminDeleteUser, middleware.HasPermission("delete_user")).Name = "admin.user.deleteUser"
	admin.GET("/user/:id",
		controller.AdminGetUser, middleware.HasPermission("get_user")).Name = "admin.user.getUser"
	admin.GET("/users",
		controller.AdminGetUsers, middleware.HasPermission("get_users")).Name = "admin.getUsers"

	api.GET("/user/me", controller.GetMe, middleware.Logged).Name = "user.getMe"
	api.PUT("/user/me", controller.UpdateMe, middleware.Logged).Name = "user.updateMe"
	api.GET("/user/:id", controller.GetUser).Name = "user.getUser"
	api.GET("/users", controller.GetUsers).Name = "user.getUsers"

	api.POST("/user/change_password", controller.ChangePassword, middleware.Logged).Name = "user.changePassword"

	api.GET("/image/:id", controller.GetImage).Name = "image.getImage"
	api.POST("/image", controller.CreateImage, middleware.Logged).Name = "image.create"
	// TODO: routes.
}
