package controller

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

func AdminCreateProblem(c echo.Context) error {
	file, err := c.FormFile("attachment_file")
	if err != nil && err != http.ErrMissingFile && err.Error() != "request Content-Type isn't multipart/form-data" {
		panic(errors.Wrap(err, "could not read file"))
	}

	req := request.AdminCreateProblemRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	var public, privacy bool
	if req.Public == nil {
		public = false
	} else {
		public = *req.Public
	}
	if req.Privacy == nil {
		privacy = true
	} else {
		privacy = *req.Privacy
	}

	problem := models.Problem{
		Name:               req.Name,
		Description:        req.Description,
		Public:             public,
		Privacy:            privacy,
		MemoryLimit:        req.MemoryLimit,
		TimeLimit:          req.TimeLimit,
		LanguageAllowed:    req.LanguageAllowed,
		CompileEnvironment: req.CompileEnvironment,
		CompareScriptID:    req.CompareScriptID,
	}
	if file != nil {
		problem.AttachmentFileName = file.Filename
	}
	utils.PanicIfDBError(base.DB.Create(&problem), "could not create problem")

	if file != nil { // TODO: use MustPutObject
		src, err := file.Open()
		if err != nil {
			panic(err)
		}
		defer func() {
			err := src.Close()
			if err != nil {
				panic(errors.Wrap(err, "could not close file reader"))
			}
		}()
		_, err = base.Storage.PutObjectWithContext(c.Request().Context(), "problems", fmt.Sprintf("%d/attachment", problem.ID), src, file.Size, minio.PutObjectOptions{})
		if err != nil {
			panic(errors.Wrap(err, "could write attachment file to s3 storage."))
		}
	}

	var user models.User
	if user, ok = c.Get("user").(models.User); !ok {
		panic("could not get user to grant role problem creator")
	}
	user.GrantRole("creator", problem)
	return c.JSON(http.StatusCreated, response.AdminCreateProblemResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ProblemForAdmin `json:"problem"`
		}{
			resource.GetProblemForAdmin(&problem),
		},
	})
}

// TODO: add test for file operation

func AdminGetProblem(c echo.Context) error {
	problem, err := utils.FindProblem(c.Param("id"), false)
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
	}
	// TODO: load test cases
	return c.JSON(http.StatusOK, response.AdminGetProblemResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ProblemForAdmin `json:"problem"`
		}{
			resource.GetProblemForAdmin(problem),
		},
	})
}

func AdminGetProblems(c echo.Context) error {
	req := request.AdminGetProblemsRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}

	query, err := utils.Sorter(base.DB.Model(&models.Problem{}), req.OrderBy, "id")
	if err != nil {
		if herr, ok := err.(utils.HttpError); ok {
			return herr.Response(c)
		}
		panic(err)
	}

	if req.Search != "" {
		id, _ := strconv.ParseUint(req.Search, 10, 64)
		query = query.Where("id = ? or name like ?", id, "%"+req.Search+"%")
	}

	var problems []models.Problem
	total, prevUrl, nextUrl, err := utils.Paginator(query, req.Limit, req.Offset, c.Request().URL, &problems)
	return c.JSON(http.StatusOK, response.AdminGetProblemsResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Problems []resource.ProblemForAdmin `json:"problems"`
			Total    int                        `json:"total"`
			Count    int                        `json:"count"`
			Offset   int                        `json:"offset"`
			Prev     *string                    `json:"prev"`
			Next     *string                    `json:"next"`
		}{
			Problems: resource.GetProblemForAdminSlice(problems),
			Total:    total,
			Count:    len(problems),
			Offset:   req.Offset,
			Prev:     prevUrl,
			Next:     nextUrl,
		},
	})
}

func AdminUpdateProblem(c echo.Context) error {
	req := request.AdminUpdateProblemRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}
	problem, err := utils.FindProblem(c.Param("id"), false)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		}
		panic(err)
	}
	problem.Name = req.Name
	problem.Description = req.Description

	file, err := c.FormFile("attachment_file")
	if err != nil && err != http.ErrMissingFile && err.Error() != "request Content-Type isn't multipart/form-data" {
		panic(errors.Wrap(err, "could not read file"))
	}
	if file != nil { // TODO: use MustPutObject
		problem.AttachmentFileName = file.Filename
		src, err := file.Open()
		if err != nil {
			panic(err)
		}
		defer func() {
			err := src.Close()
			if err != nil {
				panic(errors.Wrap(err, "could not close file reader"))
			}
		}()
		_, err = base.Storage.PutObjectWithContext(c.Request().Context(), "problems", fmt.Sprintf("%d/attachment", problem.ID), src, file.Size, minio.PutObjectOptions{})
		if err != nil {
			panic(errors.Wrap(err, "could write attachment file to s3 storage."))
		}
	}

	if req.Public != nil {
		problem.Public = *req.Public
	} else {
		problem.Public = false
	}
	if req.Privacy != nil {
		problem.Privacy = *req.Privacy
	} else {
		problem.Privacy = true
	}
	problem.MemoryLimit = req.MemoryLimit
	problem.TimeLimit = req.TimeLimit
	problem.LanguageAllowed = req.LanguageAllowed
	problem.CompileEnvironment = req.CompileEnvironment
	problem.CompareScriptID = req.CompareScriptID
	utils.PanicIfDBError(base.DB.Save(&problem), "could not update problem")
	return c.JSON(http.StatusOK, response.AdminUpdateProblemResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ProblemForAdmin `json:"problem"`
		}{
			resource.GetProblemForAdmin(problem),
		},
	})
}

// TODO: add test for file operation

func AdminDeleteProblem(c echo.Context) error {
	problem, err := utils.FindProblem(c.Param("id"), false)
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
	}

	var roles []models.Role
	if err := base.DB.Where("target = ?", "problem").Find(&roles).Error; err != gorm.ErrRecordNotFound && err != nil {
		panic(errors.Wrap(err, "could not find roles"))
	}
	roleIds := make([]uint, len(roles))
	for i, role := range roles {
		roleIds[i] = role.ID
	}
	utils.PanicIfDBError(base.DB.Delete(&models.UserHasRole{}, "role_id IN (?) and target_id = ?", roleIds, problem.ID), "could not delete user has roles")
	utils.PanicIfDBError(base.DB.Delete(&problem), "could not delete problem")
	return c.JSON(http.StatusOK, response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}

func AdminCreateTestCase(c echo.Context) error {

	problem, err := utils.FindProblem(c.Param("id"), false)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, response.ErrorResp("WRONG_PROBLEM", nil))
		}
		panic(err)
	}

	inputFile, err := c.FormFile("input_file")
	if err != nil && err != http.ErrMissingFile && err.Error() != "request Content-Type isn't multipart/form-data" {
		panic(errors.Wrap(err, "could not read input file"))
	}
	outputFile, err := c.FormFile("output_file")
	if err != nil && err != http.ErrMissingFile && err.Error() != "request Content-Type isn't multipart/form-data" {
		panic(errors.Wrap(err, "could not read output file"))
	}

	if inputFile == nil || outputFile == nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("LACK_FILE", nil)) // TODO: code name ?
	}

	req := request.AdminCreateTestCaseRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}

	testCase := models.TestCase{
		Score:          req.Score,
		InputFileName:  inputFile.Filename,
		OutputFileName: outputFile.Filename,
	}

	if err := base.DB.Model(&problem).Association("TestCases").Append(&testCase).Error; err != nil {
		panic(errors.Wrap(err, "could not create test case"))
	}

	utils.MustPutObject(inputFile, c.Request().Context(), "problems", fmt.Sprintf("%d/input/%s", problem.ID, inputFile.Filename))
	utils.MustPutObject(outputFile, c.Request().Context(), "problems", fmt.Sprintf("%d/output/%s", problem.ID, outputFile.Filename))
	// TODO: add perm

	return c.JSON(http.StatusCreated, response.AdminCreateTestCaseResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.TestCaseForAdmin `json:"test_case"`
		}{
			resource.GetTestCaseForAdmin(&testCase),
		},
	})
}

func AdminGetTestCaseInputFile(c echo.Context) error {
	testCase, problem, err := utils.FindTestCase(c.Param("id"), c.Param("test_case_id"), false)
	if problem == nil {
		return c.JSON(http.StatusNotFound, response.ErrorResp("PROBLEM_NOT_FOUND", nil))
	} else if testCase == nil {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", err))
	}

	c.Response().Header().Set("Access-Control-Allow-Origin", strings.Join(utils.Origins, ", "))
	c.Response().Header().Set("Cache-Control", "public; max-age=31536000")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, testCase.InputFileName))

	return c.Stream(http.StatusOK, "", utils.MustGetObject("problems", fmt.Sprintf("%d/input/%s", problem.ID, testCase.InputFileName)))
}

func AdminGetTestCaseOutputFile(c echo.Context) error {
	testCase, problem, err := utils.FindTestCase(c.Param("id"), c.Param("test_case_id"), false)
	if problem == nil {
		return c.JSON(http.StatusNotFound, response.ErrorResp("PROBLEM_NOT_FOUND", nil))
	} else if testCase == nil {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", err))
	}

	c.Response().Header().Set("Access-Control-Allow-Origin", strings.Join(utils.Origins, ", "))
	c.Response().Header().Set("Cache-Control", "public; max-age=31536000")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, testCase.OutputFileName))

	return c.Stream(http.StatusOK, "", utils.MustGetObject("problems", fmt.Sprintf("%d/output/%s", problem.ID, testCase.OutputFileName)))
}

func AdminUpdateTestCase(c echo.Context) error {
	testCase, problem, err := utils.FindTestCase(c.Param("id"), c.Param("test_case_id"), false)
	if problem == nil {
		return c.JSON(http.StatusNotFound, response.ErrorResp("PROBLEM_NOT_FOUND", nil))
	} else if testCase == nil {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", err))
	}

	req := request.AdminUpdateTestCaseRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}

	inputFile, err := c.FormFile("input_file")
	if err != nil && err != http.ErrMissingFile && err.Error() != "request Content-Type isn't multipart/form-data" {
		panic(errors.Wrap(err, "could not read input file"))
	}
	outputFile, err := c.FormFile("output_file")
	if err != nil && err != http.ErrMissingFile && err.Error() != "request Content-Type isn't multipart/form-data" {
		panic(errors.Wrap(err, "could not read output file"))
	}

	if inputFile != nil {
		if err := base.Storage.RemoveObject("problems", fmt.Sprintf("%d/input/%s", problem.ID, testCase.InputFileName)); err != nil {
			panic(errors.Wrap(err, "could not remove object"))
		}
		utils.MustPutObject(inputFile, c.Request().Context(), "problems", fmt.Sprintf("%d/input/%s", problem.ID, inputFile.Filename))
		testCase.InputFileName = inputFile.Filename
	}
	if outputFile != nil {
		if err := base.Storage.RemoveObject("problems", fmt.Sprintf("%d/output/%s", problem.ID, testCase.OutputFileName)); err != nil {
			panic(errors.Wrap(err, "could not remove object"))
		}
		utils.MustPutObject(outputFile, c.Request().Context(), "problems", fmt.Sprintf("%d/output/%s", problem.ID, outputFile.Filename))
		testCase.OutputFileName = outputFile.Filename
	}

	testCase.Score = req.Score
	utils.PanicIfDBError(base.DB.Save(&testCase), "could not update testCase")

	return c.JSON(http.StatusOK, response.AdminUpdateTestCaseResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.TestCaseForAdmin `json:"test_case"`
		}{
			resource.GetTestCaseForAdmin(testCase),
		},
	})
}

func AdminDeleteTestCase(c echo.Context) error {
	testCase, problem, err := utils.FindTestCase(c.Param("id"), c.Param("test_case_id"), false)
	if problem == nil {
		return c.JSON(http.StatusNotFound, response.ErrorResp("PROBLEM_NOT_FOUND", nil))
	} else if testCase == nil {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", err))
	}

	if err := base.Storage.RemoveObject("problems", fmt.Sprintf("%d/input/%s", problem.ID, testCase.InputFileName)); err != nil {
		panic(errors.Wrap(err, "could not remove object"))
	}
	if err := base.Storage.RemoveObject("problems", fmt.Sprintf("%d/output/%s", problem.ID, testCase.OutputFileName)); err != nil {
		panic(errors.Wrap(err, "could not remove object"))
	}

	utils.PanicIfDBError(base.DB.Delete(&testCase), "could not remove test case")

	return c.JSON(http.StatusOK, response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}

func AdminDeleteTestCases(c echo.Context) error {
	problem, err := utils.FindProblem(c.Param("id"), false)
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusNotFound, response.ErrorResp("PROBLEM_NOT_FOUND", nil))
	} else if err != nil {
		panic(err)
	}

	testCaseIds := make([]uint, len(problem.TestCases))

	for i, testCase := range problem.TestCases {
		if err := base.Storage.RemoveObject("problems", fmt.Sprintf("%d/input/%s", problem.ID, testCase.InputFileName)); err != nil {
			panic(errors.Wrap(err, "could not remove object"))
		}
		if err := base.Storage.RemoveObject("problems", fmt.Sprintf("%d/output/%s", problem.ID, testCase.OutputFileName)); err != nil {
			panic(errors.Wrap(err, "could not remove object"))
		}
		testCaseIds[i] = testCase.ID
	}
	utils.PanicIfDBError(base.DB.Delete(&models.TestCase{}, "id IN (?)", testCaseIds), "could not delete test cases")

	return c.JSON(http.StatusOK, response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}
