package controller_test

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/minio/minio-go"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAdminCreateProblem(t *testing.T) {
	t.Parallel()
	FailTests := []failTest{
		{
			name:   "WithoutParams",
			method: "POST",
			path:   "/api/admin/problem",
			req: request.AdminCreateProblemRequest{
				Name:               "",
				Description:        "",
				Public:             nil,
				Privacy:            nil,
				MemoryLimit:        0,
				TimeLimit:          0,
				LanguageAllowed:    "",
				CompileEnvironment: "",
				CompareScriptID:    0,
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusBadRequest,
			resp: response.ErrorResp("VALIDATION_ERROR", []interface{}{
				map[string]interface{}{
					"field":       "Name",
					"reason":      "required",
					"translation": "名称为必填字段",
				},
				map[string]interface{}{
					"field":       "Description",
					"reason":      "required",
					"translation": "介绍为必填字段",
				},
				map[string]interface{}{
					"field":       "Public",
					"reason":      "required",
					"translation": "是否公开为必填字段",
				},
				map[string]interface{}{
					"field":       "Privacy",
					"reason":      "required",
					"translation": "是否可见细节为必填字段",
				},
				map[string]interface{}{
					"field":       "MemoryLimit",
					"reason":      "required",
					"translation": "内存限制为必填字段",
				},
				map[string]interface{}{
					"field":       "TimeLimit",
					"reason":      "required",
					"translation": "运行时间限制为必填字段",
				},
				map[string]interface{}{
					"field":       "LanguageAllowed",
					"reason":      "required",
					"translation": "可用语言为必填字段",
				},
				map[string]interface{}{
					"field":       "CompareScriptID",
					"reason":      "required",
					"translation": "比较脚本编号为必填字段",
				},
			}),
		},
		{
			name:   "PermissionDenied",
			method: "POST",
			path:   "/api/admin/problem",
			req: request.AdminCreateProblemRequest{
				Name: "test_admin_create_problem_perm",
			},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, FailTests, "AdminCreateProblem")

	// TODO: specific public and privacy
	// TODO: with test case

	boolTrue := true
	boolFalse := false

	t.Run("testAdminCreateProblemSuccess", func(t *testing.T) {
		t.Parallel()
		user := createUserForTest(t, "admin_create_problem", 1)
		user.GrantRole("admin")
		req := request.AdminCreateProblemRequest{
			Name:            "test_admin_create_problem_1",
			Description:     "test_admin_create_problem_1_desc",
			MemoryLimit:     4294967296,
			TimeLimit:       1000,
			LanguageAllowed: "test_admin_create_problem_1_language_allowed",
			CompareScriptID: 1,
			Public:          &boolFalse,
			Privacy:         &boolTrue,
		}
		httpReq := makeReq(t, "POST", "/api/admin/problem", req, headerOption{
			"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
		})
		log.Debug(httpReq.Header.Get("Content-Type"))
		httpResp := makeResp(httpReq)
		assert.Equal(t, http.StatusCreated, httpResp.StatusCode)
		databaseProblem := models.Problem{}
		assert.Nil(t, base.DB.Where("name = ?", req.Name).First(&databaseProblem).Error)
		// request == database
		assert.Equal(t, req.Name, databaseProblem.Name)
		assert.Equal(t, req.Description, databaseProblem.Description)
		assert.Equal(t, req.MemoryLimit, databaseProblem.MemoryLimit)
		assert.Equal(t, req.TimeLimit, databaseProblem.TimeLimit)
		assert.Equal(t, req.LanguageAllowed, databaseProblem.LanguageAllowed)
		assert.Equal(t, req.CompareScriptID, databaseProblem.CompareScriptID)
		assert.False(t, databaseProblem.Public)
		assert.True(t, databaseProblem.Privacy)
		// response == database
		resp := response.AdminCreateProblemResponse{}
		mustJsonDecode(httpResp, &resp)
		jsonEQ(t, response.AdminUpdateProblemResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.ProblemForAdmin `json:"problem"`
			}{
				resource.GetProblemForAdmin(&databaseProblem),
			},
		}, resp)
		assert.True(t, user.HasRole("creator", databaseProblem))
	})
}

func TestAdminUpdateProblem(t *testing.T) {
	t.Parallel()
	problem1 := models.Problem{
		Name:               "test_admin_update_problem_1",
		AttachmentFileName: "test_admin_update_problem_1_attachment_file_name",
		LanguageAllowed:    "test_admin_update_problem_1_language_allowed",
	}
	assert.Nil(t, base.DB.Create(&problem1).Error)

	userWithProblem1Perm := models.User{
		Username: "test_admin_update_problem_user_p1",
		Nickname: "test_admin_update_problem_user_p1_nick",
		Email:    "test_admin_update_problem_user_p1@e.e",
		Password: utils.HashPassword("test_admin_update_problem_user_p1"),
	}
	assert.Nil(t, base.DB.Create(&userWithProblem1Perm).Error)
	userWithProblem1Perm.GrantRole("creator", problem1)

	FailTests := []failTest{
		{
			name:   "WithoutParams",
			method: "PUT",
			path:   fmt.Sprintf("/api/admin/problem/%d", problem1.ID),
			req: request.AdminUpdateProblemRequest{
				Name:               "",
				Description:        "",
				Public:             nil,
				Privacy:            nil,
				MemoryLimit:        0,
				TimeLimit:          0,
				LanguageAllowed:    "",
				CompileEnvironment: "",
				CompareScriptID:    0,
			},
			reqOptions: []reqOption{
				headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", userWithProblem1Perm.ID)},
				},
			},
			statusCode: http.StatusBadRequest,
			resp: response.ErrorResp("VALIDATION_ERROR", []interface{}{
				map[string]interface{}{
					"field":       "Name",
					"reason":      "required",
					"translation": "名称为必填字段",
				},
				map[string]interface{}{
					"field":       "Description",
					"reason":      "required",
					"translation": "介绍为必填字段",
				},
				map[string]interface{}{
					"field":       "Public",
					"reason":      "required",
					"translation": "是否公开为必填字段",
				},
				map[string]interface{}{
					"field":       "Privacy",
					"reason":      "required",
					"translation": "是否可见细节为必填字段",
				},
				map[string]interface{}{
					"field":       "MemoryLimit",
					"reason":      "required",
					"translation": "内存限制为必填字段",
				},
				map[string]interface{}{
					"field":       "TimeLimit",
					"reason":      "required",
					"translation": "运行时间限制为必填字段",
				},
				map[string]interface{}{
					"field":       "LanguageAllowed",
					"reason":      "required",
					"translation": "可用语言为必填字段",
				},
				map[string]interface{}{
					"field":       "CompareScriptID",
					"reason":      "required",
					"translation": "比较脚本编号为必填字段",
				},
			}),
		},
		{
			name:   "NonExistId",
			method: "PUT",
			path:   "/api/admin/problem/-1",
			req: request.AdminUpdateProblemRequest{
				Name:            "test_admin_update_problem_non_exist",
				LanguageAllowed: "test_admin_update_problem_non_exist_language_allowed",
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
		{
			name:   "PermissionDenied",
			method: "PUT",
			path:   fmt.Sprintf("/api/admin/problem/%d", problem1.ID),
			req: request.AdminUpdateProblemRequest{
				Name:            "test_admin_update_problem_prem",
				LanguageAllowed: "test_admin_update_problem_perm_language_allowed",
			},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, FailTests, "AdminUpdateProblem")

	boolTrue := true
	boolFalse := false

	successTests := []struct {
		name            string
		path            string
		originalProblem models.Problem
		expectedProblem models.Problem
		req             request.AdminUpdateProblemRequest
	}{
		{
			name: "WithSpecifiedPublicAndPrivacy",
			path: "id",
			originalProblem: models.Problem{
				Name:            "test_admin_update_problem_3",
				Description:     "test_admin_update_problem_3_desc",
				LanguageAllowed: "test_admin_update_problem_3_language_allowed",
				Public:          false,
				Privacy:         true,
				MemoryLimit:     1024,
				TimeLimit:       1000,
				CompareScriptID: 1,
			},
			expectedProblem: models.Problem{
				Name:            "test_admin_update_problem_30",
				Description:     "test_admin_update_problem_30_desc",
				LanguageAllowed: "test_admin_update_problem_30_language_allowed",
				Public:          true,
				Privacy:         false,
				MemoryLimit:     2048,
				TimeLimit:       2000,
				CompareScriptID: 2,
			},
			req: request.AdminUpdateProblemRequest{
				Name:            "test_admin_update_problem_30",
				Description:     "test_admin_update_problem_30_desc",
				LanguageAllowed: "test_admin_update_problem_30_language_allowed",
				Public:          &boolTrue,
				Privacy:         &boolFalse,
				MemoryLimit:     2048,
				TimeLimit:       2000,
				CompareScriptID: 2,
			},
		},
		// TODO: with test case
	}

	t.Run("testAdminUpdateProblemSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("testAdminUpdateProblem"+test.name, func(t *testing.T) {
				t.Parallel()
				assert.Nil(t, base.DB.Create(&test.originalProblem).Error)
				path := fmt.Sprintf("/api/admin/problem/%d", test.originalProblem.ID)
				user := createUserForTest(t, "admin_update_problem", i)
				user.GrantRole("creator", test.originalProblem)
				httpResp := makeResp(makeReq(t, "PUT", path, test.req, headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				}))
				databaseProblem := models.Problem{}
				assert.Nil(t, base.DB.First(&databaseProblem, test.originalProblem.ID).Error)
				// ignore other fields
				test.expectedProblem.ID = databaseProblem.ID
				test.expectedProblem.CreatedAt = databaseProblem.CreatedAt
				test.expectedProblem.UpdatedAt = databaseProblem.UpdatedAt
				test.expectedProblem.DeletedAt = databaseProblem.DeletedAt
				assert.Equal(t, test.expectedProblem, databaseProblem)
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				jsonEQ(t, response.AdminUpdateProblemResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data: struct {
						*resource.ProblemForAdmin `json:"problem"`
					}{
						resource.GetProblemForAdmin(&databaseProblem),
					},
				}, httpResp)
			})
		}
	})
}

func TestAdminDeleteProblem(t *testing.T) {
	t.Parallel()

	failTests := []failTest{
		{
			name:   "NonExistId",
			method: "DELETE",
			path:   "/api/admin/problem/-1",
			req:    request.AdminDeleteProblemRequest{},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
		{
			name:   "PermissionDenied",
			method: "DELETE",
			path:   "/api/admin/problem/-1",
			req:    request.AdminDeleteProblemRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "AdminDeleteProblem")

	successTests := []struct {
		name    string
		problem models.Problem
	}{
		{
			name: "success",
			problem: models.Problem{
				Name:               "test_admin_delete_problem_1",
				AttachmentFileName: "test_admin_delete_problem_1_attachment_file_name",
				LanguageAllowed:    "test_admin_delete_problem_1_language_allowed",
			},
		},
		// TODO: with test case
	}

	t.Run("testAdminDeleteProblemSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("testAdminDeleteProblem"+test.name, func(t *testing.T) {
				t.Parallel()
				assert.Nil(t, base.DB.Create(&test.problem).Error)
				user := createUserForTest(t, "admin_delete_problem", i)
				user.GrantRole("creator", test.problem)
				httpResp := makeResp(makeReq(t, "DELETE", fmt.Sprintf("/api/admin/problem/%d", test.problem.ID), request.AdminDeleteUserRequest{}, headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				}))
				resp := response.Response{}
				mustJsonDecode(httpResp, &resp)
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				assert.Equal(t, response.Response{
					Message: "SUCCESS",
					Error:   nil,
					Data:    nil,
				}, resp)
				assert.Equal(t, gorm.ErrRecordNotFound, base.DB.First(models.Problem{}, test.problem.ID).Error)
				assert.False(t, user.HasRole("creator", test.problem))
			})
		}
	})
}

func TestAdminGetProblem(t *testing.T) {
	t.Parallel()

	failTests := []failTest{
		{
			name:   "NonExistId",
			method: "GET",
			path:   "/api/admin/problem/-1",
			req:    request.AdminGetProblemRequest{},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
		{
			name:   "PermissionDenied",
			method: "GET",
			path:   "/api/admin/problem/-1",
			req:    request.AdminGetProblemRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "AdminGetProblem")

	successTests := []struct {
		name       string
		path       string
		req        request.AdminGetProblemRequest
		problem    models.Problem
		role       models.Role
		roleTarget models.HasRole
	}{
		{
			name: "WithId",
			path: "id",
			req:  request.AdminGetProblemRequest{},
			problem: models.Problem{
				Name:               "test_admin_get_problem_1",
				AttachmentFileName: "test_admin_get_problem_1_attachment_file_name",
				LanguageAllowed:    "test_admin_get_problem_1_language_allowed",
			},
			role:       models.Role{},
			roleTarget: nil,
		},
		// TODO: with test case
	}

	t.Run("testAdminGetUserSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("testAdminGetProblem"+test.name, func(t *testing.T) {
				t.Parallel()
				assert.Nil(t, base.DB.Create(&test.problem).Error)
				user := createUserForTest(t, "admin_get_problem", i)
				user.GrantRole("creator", test.problem)
				httpResp := makeResp(makeReq(t, "GET", fmt.Sprintf("/api/admin/problem/%d", test.problem.ID), request.AdminGetUserRequest{}, headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				}))
				resp := response.AdminGetProblemResponse{}
				mustJsonDecode(httpResp, &resp)
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				assert.Equal(t, response.AdminGetProblemResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data: struct {
						*resource.ProblemForAdmin `json:"problem"`
					}{
						resource.GetProblemForAdmin(&test.problem),
					},
				}, resp)
			})
		}
	})
}

func TestAdminGetProblems(t *testing.T) {
	t.Parallel()
	problem1 := models.Problem{
		Name:               "test_admin_get_problems_1",
		AttachmentFileName: "test_admin_get_problems_1_attachment_file_name",
		LanguageAllowed:    "test_admin_get_problems_1_language_allowed",
	}
	problem2 := models.Problem{
		Name:               "test_admin_get_problems_2",
		AttachmentFileName: "test_admin_get_problems_2_attachment_file_name",
		LanguageAllowed:    "test_admin_get_problems_2_language_allowed",
	}
	problem3 := models.Problem{
		Name:               "test_admin_get_problems_3",
		AttachmentFileName: "test_admin_get_problems_3_attachment_file_name",
		LanguageAllowed:    "test_admin_get_problems_3_language_allowed",
	}
	problem4 := models.Problem{
		Name:               "test_admin_get_problems_4",
		AttachmentFileName: "test_admin_get_problems_4_attachment_file_name",
		LanguageAllowed:    "test_admin_get_problems_4_language_allowed",
	}
	assert.Nil(t, base.DB.Create(&problem1).Error)
	assert.Nil(t, base.DB.Create(&problem2).Error)
	assert.Nil(t, base.DB.Create(&problem3).Error)
	assert.Nil(t, base.DB.Create(&problem4).Error)

	type respData struct {
		Problems []resource.ProblemForAdmin `json:"problems"`
		Total    int                        `json:"total"`
		Count    int                        `json:"count"`
		Offset   int                        `json:"offset"`
		Prev     *string                    `json:"prev"`
		Next     *string                    `json:"next"`
	}

	requestUrl := "/api/admin/problems"

	failTests := []failTest{
		{
			name:   "NotAllowedColumn",
			method: "GET",
			path:   requestUrl,
			req: request.AdminGetProblemsRequest{
				Search:  "test_admin_get_problems",
				OrderBy: "name.ASC",
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("INVALID_ORDER", nil),
		},
		{
			name:   "PermissionDenied",
			method: "GET",
			path:   requestUrl,
			req:    request.AdminGetProblemsRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "AdminGetProblems")

	successTests := []struct {
		name     string
		req      request.AdminGetProblemsRequest
		respData respData
	}{
		{
			name: "All",
			req: request.AdminGetProblemsRequest{
				Search:  "test_admin_get_problems",
				Limit:   0,
				Offset:  0,
				OrderBy: "",
			},
			respData: respData{
				Problems: []resource.ProblemForAdmin{
					*resource.GetProblemForAdmin(&problem1),
					*resource.GetProblemForAdmin(&problem2),
					*resource.GetProblemForAdmin(&problem3),
					*resource.GetProblemForAdmin(&problem4),
				},
				Total:  4,
				Count:  4,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "NonExist",
			req: request.AdminGetProblemsRequest{
				Search:  "test_admin_get_problems_non_exist",
				Limit:   0,
				Offset:  0,
				OrderBy: "",
			},
			respData: respData{
				Problems: []resource.ProblemForAdmin{},
				Total:    0,
				Count:    0,
				Offset:   0,
				Prev:     nil,
				Next:     nil,
			},
		},
		{
			name: "Search",
			req: request.AdminGetProblemsRequest{
				Search:  "test_admin_get_problems_2",
				Limit:   0,
				Offset:  0,
				OrderBy: "",
			},
			respData: respData{
				Problems: []resource.ProblemForAdmin{
					*resource.GetProblemForAdmin(&problem2),
				},
				Total:  1,
				Count:  1,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "Sort",
			req: request.AdminGetProblemsRequest{
				Search:  "test_admin_get_problems",
				Limit:   0,
				Offset:  0,
				OrderBy: "id.DESC",
			},
			respData: respData{
				Problems: []resource.ProblemForAdmin{
					*resource.GetProblemForAdmin(&problem4),
					*resource.GetProblemForAdmin(&problem3),
					*resource.GetProblemForAdmin(&problem2),
					*resource.GetProblemForAdmin(&problem1),
				},
				Total:  4,
				Count:  4,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "Paginator",
			req: request.AdminGetProblemsRequest{
				Search:  "test_admin_get_problems",
				Limit:   2,
				Offset:  1,
				OrderBy: "",
			},
			respData: respData{
				Problems: []resource.ProblemForAdmin{
					*resource.GetProblemForAdmin(&problem2),
					*resource.GetProblemForAdmin(&problem3),
				},
				Total:  4,
				Count:  2,
				Offset: 1,
				Prev:   nil,
				Next: getUrlStringPointer(requestUrl, map[string]string{
					"limit":  "2",
					"offset": "3",
				}),
			},
		},
	}

	t.Run("testAdminGetProblemsSuccess", func(t *testing.T) {
		t.Parallel()
		for _, test := range successTests {
			test := test
			t.Run("testAdminGetProblems"+test.name, func(t *testing.T) {
				t.Parallel()
				httpResp := makeResp(makeReq(t, "GET", requestUrl, test.req, applyAdminUser))
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				resp := response.Response{}
				mustJsonDecode(httpResp, &resp)
				jsonEQ(t, response.AdminGetProblemsResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data:    test.respData,
				}, resp)
			})
		}
	})
}

func createProblemForTest(t *testing.T, name string, index int) (problem models.Problem, user models.User) {
	problem = models.Problem{
		Name:               fmt.Sprintf("problem_for_testing_%s_%d", name, index),
		Description:        fmt.Sprintf("a problem used to test API: %s(%d)", name, index),
		AttachmentFileName: "",
		Public:             true,
		Privacy:            false,
		MemoryLimit:        1024,
		TimeLimit:          1000,
		LanguageAllowed:    fmt.Sprintf("test_%s_language_allowed_%d", name, index),
		CompileEnvironment: fmt.Sprintf("test_%s_compile_environment_%d", name, index),
		CompareScriptID:    1,
	}
	assert.Nil(t, base.DB.Create(&problem).Error)
	user = createUserForTest(t, name, index)
	user.GrantRole("creator", problem)
	return
}

func createTestCaseForTest(t *testing.T, problem models.Problem, score uint, inputFile, outputFile *fileContent) (testCase models.TestCase) {
	var inputFileName, outputFileName string

	if inputFile != nil {
		utils.MustCreateBucket("problems")
		inputBytes, err := ioutil.ReadAll(inputFile.reader)

		assert.Nil(t, err)
		_, err = base.Storage.PutObject("problems", fmt.Sprintf("%d/input/%s", problem.ID, inputFile.fileName), bytes.NewReader(inputBytes), int64(len(inputBytes)), minio.PutObjectOptions{})
		assert.Nil(t, err)
		inputFileName = inputFile.fileName
	}
	if outputFile != nil {
		utils.MustCreateBucket("problems")
		outputBytes, err := ioutil.ReadAll(outputFile.reader)
		assert.Nil(t, err)
		_, err = base.Storage.PutObject("problems", fmt.Sprintf("%d/output/%s", problem.ID, outputFile.fileName), bytes.NewReader(outputBytes), int64(len(outputBytes)), minio.PutObjectOptions{})
		assert.Nil(t, err)
		outputFileName = outputFile.fileName
	}

	testCase = models.TestCase{
		Score:          score,
		InputFileName:  inputFileName,
		OutputFileName: outputFileName,
	}
	assert.Nil(t, base.DB.Model(&problem).Association("TestCases").Append(&testCase).Error)
	return
}

func TestAdminCreateTestCase(t *testing.T) {
	problem, user := createProblemForTest(t, "admin_create_test_case", 0)

	failTests := []failTest{
		{
			name:   "LackInputFile",
			method: "POST",
			path:   fmt.Sprintf("/api/admin/problem/%d/test_case", problem.ID),
			req: addFieldContentSlice([]reqContent{
				newFileContent("output_file", "test_admin_create_test_case_lack_input_file.out", "b3V0cHV0IHRleHQ"),
			}, map[string]string{
				"score": "41",
			}),
			reqOptions: []reqOption{
				headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				},
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("LACK_FILE", nil),
		},
		{
			name:   "LackOutputFile",
			method: "POST",
			path:   fmt.Sprintf("/api/admin/problem/%d/test_case", problem.ID),
			req: addFieldContentSlice([]reqContent{
				newFileContent("input_file", "test_admin_create_test_case_lack_output_file.in", "b3V0cHV0IHRleHQ"),
			}, map[string]string{
				"score": "42",
			}),
			reqOptions: []reqOption{
				headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				},
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("LACK_FILE", nil),
		},
		{
			name:   "LackBothFile",
			method: "POST",
			path:   fmt.Sprintf("/api/admin/problem/%d/test_case", problem.ID),
			req: addFieldContentSlice([]reqContent{}, map[string]string{
				"score": "43",
			}),
			reqOptions: []reqOption{
				headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				},
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("LACK_FILE", nil),
		},
		{
			name:   "PermissionDenied",
			method: "POST",
			path:   fmt.Sprintf("/api/admin/problem/%d/test_case", problem.ID),
			req: addFieldContentSlice([]reqContent{
				newFileContent("input_file", "test_admin_create_test_case_permission_denied.in", "b3V0cHV0IHRleHQ"),
				newFileContent("output_file", "test_admin_create_test_case_permission_denied.out", "b3V0cHV0IHRleHQ"),
			}, map[string]string{
				"score": "44",
			}),
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "AdminCreateTestCase")

	t.Run("testAdminCreateTestCaseSuccess", func(t *testing.T) {
		req := makeReq(t, "POST", fmt.Sprintf("/api/admin/problem/%d/test_case", problem.ID), addFieldContentSlice([]reqContent{
			newFileContent("input_file", "test_admin_create_test_case_success.in", "aW5wdXQgdGV4dA"),
			newFileContent("output_file", "test_admin_create_test_case_success.out", "b3V0cHV0IHRleHQ"),
		}, map[string]string{
			"score": "20",
		}), headerOption{
			"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
		})
		//req.Header.Set("Set-User-For-Test", fmt.Sprintf("%d", user.ID))
		httpResp := makeResp(req)
		expectedTestCase := models.TestCase{
			ProblemID:      problem.ID,
			Score:          20,
			InputFileName:  "test_admin_create_test_case_success.in",
			OutputFileName: "test_admin_create_test_case_success.out",
		}
		databaseTestCase := models.TestCase{}
		assert.Nil(t, base.DB.Where("problem_id = ? and score = ?", problem.ID, 20).First(&databaseTestCase).Error)
		assert.Equal(t, expectedTestCase.ProblemID, databaseTestCase.ProblemID)
		assert.Equal(t, expectedTestCase.InputFileName, databaseTestCase.InputFileName)
		assert.Equal(t, expectedTestCase.OutputFileName, databaseTestCase.OutputFileName)
		assert.Equal(t, expectedTestCase.Score, databaseTestCase.Score)
		resp := response.AdminCreateTestCaseResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusCreated, httpResp.StatusCode)
		assert.Equal(t, response.AdminCreateTestCaseResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.TestCaseForAdmin `json:"test_case"`
			}{
				resource.GetTestCaseForAdmin(&databaseTestCase),
			},
		}, resp)
	})
}

func TestAdminGetTestCaseInputFile(t *testing.T) {
	problem, user := createProblemForTest(t, "admin_get_test_case_input_file", 0)
	failTests := []failTest{
		{
			name:   "NonExistingProblem",
			method: "GET",
			path:   "/api/admin/problem/-1/test_case/1/input_file",
			req:    request.AdminGetTestCaseInputFileRequest{},
			reqOptions: []reqOption{
				headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				},
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
		{
			name:   "NonExistingTestCase",
			method: "GET",
			path:   fmt.Sprintf("/api/admin/problem/%d/test_case/-1/input_file", problem.ID),
			req:    request.AdminGetTestCaseInputFileRequest{},
			reqOptions: []reqOption{
				headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				},
			},
			statusCode: http.StatusNotFound,
			resp: response.ErrorResp("NOT_FOUND", map[string]interface{}{
				"Err":  map[string]interface{}{},
				"Func": "ParseUint",
				"Num":  "-1",
			}),
		},
		{
			name:   "PermissionDenied",
			method: "GET",
			path:   fmt.Sprintf("/api/admin/problem/%d/test_case/1/input_file", problem.ID),
			req:    request.AdminGetTestCaseInputFileRequest{},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "AdminGetTestCaseInputFile")

	testCase := createTestCaseForTest(t, problem, 51,
		newFileContent("", "test_admin_get_test_case_input_file_success.in", "aW5wdXQgdGV4dA"),
		//newFileContent("","test_admin_get_test_case_input_file_success.out","b3V0cHV0IHRleHQ"),
		nil,
	)

	req := makeReq(t, "GET", fmt.Sprintf("/api/admin/problem/%d/test_case/%d/input_file", problem.ID, testCase.ID), request.AdminGetTestCaseInputFileRequest{}, headerOption{
		"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
	})
	httpResp := makeResp(req)

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "input text", string(respBytes))
}

func TestAdminGetTestCaseOutputFile(t *testing.T) {
	problem, user := createProblemForTest(t, "admin_get_test_case_output_file", 0)
	failTests := []failTest{
		{
			name:   "NonExistingProblem",
			method: "GET",
			path:   "/api/admin/problem/-1/test_case/1/output_file",
			req:    request.AdminGetTestCaseOutputFileRequest{},
			reqOptions: []reqOption{
				headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				},
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
		{
			name:   "NonExistingTestCase",
			method: "GET",
			path:   fmt.Sprintf("/api/admin/problem/%d/test_case/-1/output_file", problem.ID),
			req:    request.AdminGetTestCaseOutputFileRequest{},
			reqOptions: []reqOption{
				headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				},
			},
			statusCode: http.StatusNotFound,
			resp: response.ErrorResp("NOT_FOUND", map[string]interface{}{
				"Err":  map[string]interface{}{},
				"Func": "ParseUint",
				"Num":  "-1",
			}),
		},
		{
			name:   "PermissionDenied",
			method: "GET",
			path:   fmt.Sprintf("/api/admin/problem/%d/test_case/1/output_file", problem.ID),
			req:    request.AdminGetTestCaseOutputFileRequest{},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "AdminGetTestCaseOutputFile")

	testCase := createTestCaseForTest(t, problem, 51,
		nil,
		newFileContent("", "test_admin_get_test_case_output_file_success.out", "b3V0cHV0IHRleHQ"),
	)

	req := makeReq(t, "GET", fmt.Sprintf("/api/admin/problem/%d/test_case/%d/output_file", problem.ID, testCase.ID), request.AdminGetTestCaseOutputFileRequest{}, headerOption{
		"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
	})
	httpResp := makeResp(req)

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "output text", string(respBytes))
}
