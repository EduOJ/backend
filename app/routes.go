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

	// TODO: replace hard-coded path with echo.Reverse with route name.

	auth := api.Group("/auth", middleware.Auth)
	auth.POST("/login", controller.Login).Name = "auth.login"
	auth.POST("/register", controller.Register).Name = "auth.register"
	auth.GET("/email_registered", controller.EmailRegistered).Name = "auth.emailRegistered"

	admin := api.Group("/admin", middleware.Logged)
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

	admin.POST("/problem",
		controller.AdminCreateProblem, middleware.HasPermission("create_problem")).Name = "admin.problem.createProblem"
	admin.GET("/problem/:id",
		controller.AdminGetProblem, middleware.HasPermission("get_problem", "problem")).Name = "admin.problem.getProblem"
	admin.GET("/problems",
		controller.AdminGetProblems, middleware.HasPermission("get_problems")).Name = "admin.problem.getProblems"
	admin.PUT("/problem/:id",
		controller.AdminUpdateProblem, middleware.HasPermission("update_problem", "problem")).Name = "admin.problem.updateProblem"
	admin.DELETE("/problem/:id",
		controller.AdminDeleteProblem, middleware.HasPermission("delete_problem", "problem")).Name = "admin.problem.deleteProblem"

	api.GET("/problem/:id",
		controller.GetProblem).Name = "problem.getProblem"
	api.GET("/problems",
		controller.GetProblems).Name = "problem.getProblems"

	api.GET("/problem/:id/attachment_file", controller.GetProblemAttachmentFile).Name = "problem.getProblemAttachmentFile"

	admin.POST("/problem/:id/test_case",
		controller.AdminCreateTestCase, middleware.HasPermission("create_test_case", "problem")).Name = "admin.problem.createTestCase"
	admin.PUT("/problem/:id/test_case/:test_case_id",
		controller.AdminUpdateTestCase, middleware.HasPermission("update_test_case", "problem")).Name = "admin.problem.updateTestCase"
	admin.DELETE("/problem/:id/test_case/:test_case_id",
		controller.AdminDeleteTestCase, middleware.HasPermission("delete_test_case", "problem")).Name = "admin.problem.deleteTestCase"
	admin.DELETE("/problem/:id/test_cases",
		controller.AdminDeleteTestCases, middleware.HasPermission("delete_test_case", "problem")).Name = "admin.problem.deleteTestCases"

	admin.GET("/problem/:id/test_case/:test_case_id/input_file",
		controller.AdminGetTestCaseInputFile, middleware.HasPermission("get_test_case_input_file", "problem")).Name = "admin.problem.getTestCaseInputFile"
	admin.GET("/problem/:id/test_case/:test_case_id/output_file",
		controller.AdminGetTestCaseOutputFile, middleware.HasPermission("get_test_case_output_file", "problem")).Name = "admin.problem.getTestCaseOutputFile"
}
