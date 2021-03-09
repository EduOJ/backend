package app

import (
	"github.com/EduOJ/backend/app/controller"
	"github.com/EduOJ/backend/app/middleware"
	"github.com/EduOJ/backend/base/log"
	"github.com/EduOJ/backend/base/utils"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"net/http"
	"net/http/pprof"
)

func Register(e *echo.Echo) {
	utils.InitOrigin()
	e.Use(middleware.Recover)
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     utils.Origins,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
	}))
	api := e.Group("/api", middleware.Authentication)

	auth := api.Group("/auth", middleware.Auth)
	auth.POST("/login", controller.Login).Name = "auth.login"
	auth.GET("/login/webauthn", controller.BeginLogin).Name = "auth.webauthn.beginLogin"
	auth.POST("/login/webauthn", controller.FinishLogin).Name = "auth.webauthn.finishLogin"
	auth.POST("/register", controller.Register).Name = "auth.register"
	auth.GET("/email_registered", controller.EmailRegistered).Name = "auth.emailRegistered"

	admin := api.Group("/admin", middleware.Logged)

	admin.POST("/user",
		controller.AdminCreateUser, middleware.Logged, middleware.HasPermission(middleware.UnscopedPermission{P: "manage_user"})).Name = "admin.user.createUser"
	admin.PUT("/user/:id",
		controller.AdminUpdateUser, middleware.Logged, middleware.HasPermission(middleware.UnscopedPermission{P: "manage_user"})).Name = "admin.user.updateUser"
	admin.DELETE("/user/:id",
		controller.AdminDeleteUser, middleware.Logged, middleware.HasPermission(middleware.UnscopedPermission{P: "manage_user"})).Name = "admin.user.deleteUser"
	admin.GET("/user/:id",
		controller.AdminGetUser, middleware.Logged, middleware.HasPermission(middleware.UnscopedPermission{P: "read_user"})).Name = "admin.user.getUser"
	admin.GET("/users",
		controller.AdminGetUsers, middleware.Logged, middleware.HasPermission(middleware.UnscopedPermission{P: "read_user"})).Name = "admin.user.getUsers"

	api.GET("/user/me", controller.GetMe, middleware.Logged).Name = "user.getMe"
	api.PUT("/user/me", controller.UpdateMe, middleware.Logged).Name = "user.updateMe"
	api.GET("/user/:id", controller.GetUser).Name = "user.getUser"
	api.GET("/user/me/managing_classes",
		controller.GetClassesIManage, middleware.Logged).Name = "user.getClassesIManage"
	api.GET("/user/me/taking_classes",
		controller.GetClassesITake, middleware.Logged).Name = "user.getClassesITake"
	api.GET("/user/:id/problem_info", controller.GetUserProblemInfo).Name = "user.getUserProblemInfo"
	api.GET("/users", controller.GetUsers).Name = "user.getUsers"

	api.GET("/webauthn/register", controller.BeginRegistration, middleware.Logged).Name = "webauthn.BeginRegister"
	api.POST("/webauthn/register", controller.FinishRegistration, middleware.Logged).Name = "webauthn.FinishRegister"

	api.POST("/user/change_password", controller.ChangePassword, middleware.Logged).Name = "user.changePassword"

	api.GET("/image/:id", controller.GetImage).Name = "image.getImage"
	api.POST("/image", controller.CreateImage, middleware.Logged).Name = "image.createImage"

	adminProblem := admin.Group("/problem", middleware.ValidateParams(map[string]string{
		"id":           "NOT_FOUND",
		"test_case_id": "TEST_CASE_NOT_FOUND",
	}))
	problem := api.Group("/problem", middleware.ValidateParams(map[string]string{
		"id":           "NOT_FOUND",
		"test_case_id": "TEST_CASE_NOT_FOUND",
	}))
	adminProblem.POST("",
		controller.CreateProblem, middleware.Logged, middleware.HasPermission(middleware.UnscopedPermission{P: "create_problem"})).Name = "problem.createProblem"
	adminProblem.PUT("/:id",
		controller.UpdateProblem, middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "update_problem", T: "problem"},
			B: middleware.UnscopedPermission{P: "update_problem"},
		})).Name = "problem.updateProblem"
	adminProblem.DELETE("/:id",
		controller.DeleteProblem, middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "delete_problem", T: "problem"},
			B: middleware.UnscopedPermission{P: "delete_problem"},
		})).Name = "problem.deleteProblem"

	problem.GET("/random", controller.GetRandomProblem, middleware.AllowGuest).Name = "problem.getRandomProblem"
	problem.GET("/:id", controller.GetProblem, middleware.AllowGuest).Name = "problem.getProblem"
	problem.GET("s", controller.GetProblems, middleware.AllowGuest).Name = "problem.getProblems"

	problem.GET("/:id/attachment_file", controller.GetProblemAttachmentFile, middleware.AllowGuest).Name = "problem.getProblemAttachmentFile"

	adminProblem.POST("/:id/test_case",
		controller.CreateTestCase,
		middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "update_problem", T: "problem"},
			B: middleware.UnscopedPermission{P: "update_problem"},
		})).Name = "problem.createTestCase"
	adminProblem.PUT("/:id/test_case/:test_case_id",
		controller.UpdateTestCase,
		middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "update_problem", T: "problem"},
			B: middleware.UnscopedPermission{P: "update_problem"},
		})).Name = "problem.updateTestCase"
	adminProblem.DELETE("/:id/test_case/all",
		controller.DeleteTestCases,
		middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "update_problem", T: "problem"},
			B: middleware.UnscopedPermission{P: "update_problem"},
		})).Name = "problem.deleteTestCases"
	adminProblem.DELETE("/:id/test_case/:test_case_id",
		controller.DeleteTestCase,
		middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "update_problem", T: "problem"},
			B: middleware.UnscopedPermission{P: "update_problem"},
		})).Name = "problem.deleteTestCase"

	problem.GET("/:id/test_case/:test_case_id/input_file",
		controller.GetTestCaseInputFile,
		middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.CustomPermission{
				F: middleware.IsTestCaseSample,
			},
			B: middleware.OrPermission{
				A: middleware.ScopedPermission{P: "read_problem_secret", T: "problem"},
				B: middleware.UnscopedPermission{P: "read_problem_secret"},
			},
		})).Name = "problem.getTestCaseInputFile"
	problem.GET("/:id/test_case/:test_case_id/output_file",
		controller.GetTestCaseOutputFile,
		middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.CustomPermission{
				F: middleware.IsTestCaseSample,
			},
			B: middleware.OrPermission{
				A: middleware.ScopedPermission{P: "read_problem_secret", T: "problem"},
				B: middleware.UnscopedPermission{P: "read_problem_secret"},
			},
		})).Name = "problem.getTestCaseOutputFile"

	api.POST("/problem/:pid/submission", controller.CreateSubmission,
		middleware.ValidateParams(map[string]string{
			"pid": "PROBLEM_NOT_FOUND",
		}), middleware.Logged).Name = "submission.createSubmission"

	submission := api.Group("/submission", middleware.ValidateParams(map[string]string{
		"id":     "NOT_FOUND",
		"run_id": "RUN_NOT_FOUND",
	}))
	submission.GET("/:id", controller.GetSubmission, middleware.Logged).Name = "submission.getSubmission"
	submission.GET("s", controller.GetSubmissions, middleware.Logged).Name = "submission.getSubmissions"

	submission.GET("/:id/code", controller.GetSubmissionCode, middleware.Logged).Name = "submission.getSubmissionCode"
	submission.GET("/:submission_id/run/:id/output", controller.GetRunOutput, middleware.Logged).Name = "submission.getRunOutput"
	submission.GET("/:submission_id/run/:id/input", controller.GetRunInput, middleware.Logged).Name = "submission.getRunInput"
	submission.GET("/:submission_id/run/:id/compiler_output", controller.GetRunCompilerOutput, middleware.Logged).Name = "submission.getRunCompilerOutput"
	submission.GET("/:submission_id/run/:id/comparer_output", controller.GetRunComparerOutput, middleware.Logged).Name = "submission.getRunComparerOutput"

	api.POST("/problem_set/:problem_set_id/problem/:pid/submission",
		controller.ProblemSetCreateSubmission, middleware.ValidateParams(map[string]string{
			"problem_set_id": "PROBLEM_SET_NOT_FOUND",
			"pid":            "PROBLEM_NOT_FOUND",
		}), middleware.Logged,
		middleware.HasPermission(middleware.CustomPermission{F: middleware.ProblemSetStarted})).Name = "problemSet.createSubmission"
	problemSetSubmission := api.Group("/problem_set",
		middleware.ValidateParams(map[string]string{
			"id":             "NOT_FOUND",
			"problem_set_id": "PROBLEM_SET_NOT_FOUND",
			"submission_id":  "SUBMISSION_NOT_FOUND",
		}),
		middleware.Logged,
		middleware.HasPermission(middleware.OrPermission{
			A: middleware.OrPermission{
				A: middleware.ScopedPermission{P: "read_answers", T: "problem_set"},
				B: middleware.UnscopedPermission{P: "read_answers"},
			},
			B: middleware.CustomPermission{F: middleware.ProblemSetStarted},
		}))
	problemSetSubmission.GET("/:problem_set_id/submission/:id",
		controller.ProblemSetGetSubmission).Name = "problemSet.getSubmission"
	problemSetSubmission.GET("/:problem_set_id/submissions",
		controller.ProblemSetGetSubmissions).Name = "problemSet.getSubmissions"
	problemSetSubmission.GET("/:problem_set_id/submission/:id/code",
		controller.ProblemSetGetSubmissionCode).Name = "problemSet.getSubmissionCode"
	problemSetSubmission.GET("/:problem_set_id/submission/:submission_id/run/:id/output",
		controller.ProblemSetGetRunOutput).Name = "problemSet.getRunOutput"
	problemSetSubmission.GET("/:problem_set_id/submission/:submission_id/run/:id/input",
		controller.ProblemSetGetRunInput).Name = "problemSet.getRunInput"
	problemSetSubmission.GET("/:problem_set_id/submission/:submission_id/run/:id/compiler_output",
		controller.ProblemSetGetRunCompilerOutput).Name = "problemSet.getRunCompilerOutput"
	problemSetSubmission.GET("/:problem_set_id/submission/:submission_id/run/:id/comparer_output",
		controller.ProblemSetGetRunComparerOutput).Name = "problemSet.getRunComparerOutput"

	admin.GET("/logs",
		controller.AdminGetLogs, middleware.Logged, middleware.HasPermission(middleware.UnscopedPermission{P: "read_logs"})).Name = "admin.getLogs"

	judger := e.Group("/judger", middleware.Judger)

	judger.GET("/script/:name", controller.GetScript).Name = "judger.getScript"
	judger.GET("/task", controller.GetTask).Name = "judger.getTask"
	judger.PUT("/run/:id", controller.UpdateRun, middleware.ValidateParams(map[string]string{
		"id": "NOT_FOUND",
	})).Name = "judger.updateRun"

	class := api.Group("/class", middleware.ValidateParams(map[string]string{
		"id":       "NOT_FOUND",
		"class_id": "CLASS_NOT_FOUND",
	}))
	class.POST("",
		controller.CreateClass, middleware.Logged, middleware.HasPermission(middleware.UnscopedPermission{P: "manage_class"})).Name = "class.createClass"
	class.GET("/:id",
		controller.GetClass, middleware.Logged).Name = "class.getClass"
	class.PUT("/:id",
		controller.UpdateClass, middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "manage_class", T: "class"},
			B: middleware.UnscopedPermission{P: "manage_class"},
		})).Name = "class.updateClass"
	class.PUT("/:id/invite_code",
		controller.RefreshInviteCode, middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "manage_class", T: "class"},
			B: middleware.UnscopedPermission{P: "manage_class"},
		})).Name = "class.refreshInviteCode"
	class.POST("/:id/students",
		controller.AddStudents, middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "manage_class", T: "class"},
			B: middleware.UnscopedPermission{P: "manage_class"},
		})).Name = "class.addStudents"
	class.DELETE("/:id/students",
		controller.DeleteStudents, middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "manage_class", T: "class"},
			B: middleware.UnscopedPermission{P: "manage_class"},
		})).Name = "class.deleteStudents"
	class.POST("/:id/join", controller.JoinClass, middleware.Logged).Name = "class.joinClass"
	class.DELETE("/:id",
		controller.DeleteClass, middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "manage_class", T: "class"},
			B: middleware.UnscopedPermission{P: "manage_class"},
		})).Name = "class.deleteClass"

	class.POST("/:id/problem_set",
		controller.CreateProblemSet, middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "manage_problem_sets", T: "class"},
			B: middleware.UnscopedPermission{P: "manage_problem_sets"},
		})).Name = "problemSet.createProblemSet"
	class.POST("/:id/problem_set/clone",
		controller.CloneProblemSet, middleware.Logged, middleware.HasPermission(middleware.AndPermission{
			A: middleware.OrPermission{
				A: middleware.ScopedPermission{P: "manage_problem_sets", T: "class"},
				B: middleware.UnscopedPermission{P: "manage_problem_sets"},
			},
			B: middleware.OrPermission{
				A: middleware.ScopedPermission{P: "clone_problem_sets", T: "class"},
				B: middleware.UnscopedPermission{P: "clone_problem_sets"},
			},
		})).Name = "problemSet.cloneProblemSet"
	class.GET("/:class_id/problem_set/:id",
		controller.GetProblemSet, middleware.Logged).Name = "problemSet.getProblemSet"
	class.PUT("/:class_id/problem_set/:id",
		controller.UpdateProblemSet, middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "manage_problem_sets", T: "class", IdFieldName: "class_id"},
			B: middleware.UnscopedPermission{P: "manage_problem_sets"},
		})).Name = "problemSet.updateProblemSet"
	class.POST("/:class_id/problem_set/:id/problems",
		controller.AddProblemsToSet, middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "manage_problem_sets", T: "class", IdFieldName: "class_id"},
			B: middleware.UnscopedPermission{P: "manage_problem_sets"},
		})).Name = "problemSet.addProblemsToSet"
	class.DELETE("/:class_id/problem_set/:id/problems",
		controller.DeleteProblemsFromSet, middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "manage_problem_sets", T: "class", IdFieldName: "class_id"},
			B: middleware.UnscopedPermission{P: "manage_problem_sets"},
		})).Name = "problemSet.deleteProblemsFromSet"
	class.DELETE("/:class_id/problem_set/:id",
		controller.DeleteProblemSet, middleware.Logged, middleware.HasPermission(middleware.OrPermission{
			A: middleware.ScopedPermission{P: "manage_problem_sets", T: "class", IdFieldName: "class_id"},
			B: middleware.UnscopedPermission{P: "manage_problem_sets"},
		})).Name = "problemSet.deleteProblemSet"

	if viper.GetBool("debug") {
		log.Debugf("Adding pprof handlers. SHOULD NOT BE USED IN PRODUCTION")
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
