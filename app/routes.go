package app

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/leoleoasd/EduOJBackend/app/controller"
	"github.com/leoleoasd/EduOJBackend/app/middleware"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"net/http"
	"net/http/pprof"
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
	admin.POST("/user",
		controller.AdminCreateUser, middleware.HasPermission("manage_user")).Name = "admin.user.createUser"
	admin.PUT("/user/:id",
		controller.AdminUpdateUser, middleware.HasPermission("manage_user")).Name = "admin.user.updateUser"
	admin.DELETE("/user/:id",
		controller.AdminDeleteUser, middleware.HasPermission("manage_user")).Name = "admin.user.deleteUser"
	admin.GET("/user/:id",
		controller.AdminGetUser, middleware.HasPermission("read_user")).Name = "admin.user.getUser"
	admin.GET("/users",
		controller.AdminGetUsers, middleware.HasPermission("read_user")).Name = "admin.user.getUsers"

	api.GET("/user/me", controller.GetMe, middleware.Logged).Name = "user.getMe"
	api.PUT("/user/me", controller.UpdateMe, middleware.Logged).Name = "user.updateMe"
	api.GET("/user/:id", controller.GetUser).Name = "user.getUser"
	api.GET("/users", controller.GetUsers).Name = "user.getUsers"

	api.POST("/user/change_password", controller.ChangePassword, middleware.Logged).Name = "user.changePassword"

	api.GET("/image/:id", controller.GetImage).Name = "image.getImage"
	api.POST("/image", controller.CreateImage, middleware.Logged).Name = "image.createImage"

	admin.POST("/problem",
		controller.AdminCreateProblem, middleware.HasPermission("manage_problem")).Name = "admin.problem.createProblem"
	admin.GET("/problem/:id",
		controller.AdminGetProblem, middleware.HasPermission("read_problem", "problem")).Name = "admin.problem.getProblem" // TODO: merge this API into GetProblem
	admin.GET("/problems",
		controller.AdminGetProblems, middleware.HasPermission("read_problem")).Name = "admin.problem.getProblems" // TODO: merge this API into GetProblems
	admin.PUT("/problem/:id",
		controller.AdminUpdateProblem, middleware.HasPermission("manage_problem", "problem")).Name = "admin.problem.updateProblem"
	admin.DELETE("/problem/:id",
		controller.AdminDeleteProblem, middleware.HasPermission("manage_problem", "problem")).Name = "admin.problem.deleteProblem"

	api.GET("/problem/:id", controller.GetProblem).Name = "problem.getProblem"
	api.GET("/problems", controller.GetProblems).Name = "problem.getProblems"

	api.GET("/problem/:id/attachment_file", controller.GetProblemAttachmentFile).Name = "problem.getProblemAttachmentFile"

	admin.POST("/problem/:id/test_case",
		controller.AdminCreateTestCase,
		middleware.HasPermission("manage_problem", "problem")).Name = "admin.problem.createTestCase"
	admin.PUT("/problem/:id/test_case/:test_case_id",
		controller.AdminUpdateTestCase,
		middleware.HasPermission("manage_problem", "problem")).Name = "admin.problem.updateTestCase"
	admin.DELETE("/problem/:id/test_case/all",
		controller.AdminDeleteTestCases,
		middleware.HasPermission("manage_problem", "problem")).Name = "admin.problem.deleteTestCases"
	admin.DELETE("/problem/:id/test_case/:test_case_id",
		controller.AdminDeleteTestCase,
		middleware.HasPermission("manage_problem", "problem")).Name = "admin.problem.deleteTestCase"

	admin.GET("/problem/:id/test_case/:test_case_id/input_file",
		controller.AdminGetTestCaseInputFile,
		middleware.HasPermission("read_problem", "problem")).Name = "admin.problem.getTestCaseInputFile"
	admin.GET("/problem/:id/test_case/:test_case_id/output_file",
		controller.AdminGetTestCaseOutputFile,
		middleware.HasPermission("read_problem", "problem")).Name = "admin.problem.getTestCaseOutputFile"

	admin.GET("/logs",
		controller.AdminGetLogs, middleware.HasPermission("read_logs")).Name = "admin.getLogs"

	if config.MustGet("debug", false).Value().(bool) {
		log.Debugf("Adding pprof handlers. SHOULD NOT BE USED UNDER PRODUCTION")
		e.Any("/debug/pprof/", func(c echo.Context) error {
			pprof.Index(c.Response().Writer, c.Request())
			return nil
		})
		e.Any("/debug/pprof/*", func(c echo.Context) error {
			pprof.Index(c.Response().Writer, c.Request())
			return nil
		})
		e.Any("/debug/pprof/cmdline", func(c echo.Context) error {
			pprof.Cmdline(c.Response().Writer, c.Request())
			return nil
		})
		e.Any("/debug/pprof/profile", func(c echo.Context) error {
			pprof.Profile(c.Response().Writer, c.Request())
			return nil
		})
		e.Any("/debug/pprof/symbol", func(c echo.Context) error {
			pprof.Symbol(c.Response().Writer, c.Request())
			return nil
		})
		e.Any("/debug/pprof/trace", func(c echo.Context) error {
			pprof.Trace(c.Response().Writer, c.Request())
			return nil
		})
	}
}
