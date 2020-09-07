package controller_test

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func createUserForTest(t *testing.T, name string, index int) (user models.User) {
	user = models.User{
		Username: fmt.Sprintf("test_%s_user_%d", name, index),
		Nickname: fmt.Sprintf("test_%s_user_%d_nick", name, index),
		Email:    fmt.Sprintf("test_%s_user_%d@e.e", name, index),
		Password: utils.HashPassword(fmt.Sprintf("test_%s_user_%d_pwd", name, index)),
	}
	assert.Nil(t, base.DB.Create(&user).Error)
	return
}

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

	t.Run("testAdminCreateProblemSuccess", func(t *testing.T) {
		t.Parallel()
		user := models.User{
			Username: "test_admin_create_problem_user_1",
			Nickname: "test_admin_create_problem_user_1_nick",
			Email:    "test_admin_create_problem_user_1@e.e",
			Password: utils.HashPassword("test_admin_create_problem_user_1"),
		}
		assert.Nil(t, base.DB.Create(&user).Error)
		user.GrantRole("admin")
		req := request.AdminCreateProblemRequest{
			Name:        "test_admin_create_problem_1",
			Description: "test_admin_create_problem_1_desc",
			//AttachmentFileName: "test_admin_create_problem_1_attachment_file_name",
			MemoryLimit:     4294967296,
			TimeLimit:       1000,
			LanguageAllowed: "test_admin_create_problem_1_language_allowed",
			CompareScriptID: 1,
		}
		httpResp := makeResp(makeReq(t, "POST", "/api/admin/problem", req, headerOption{
			"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
		}))
		assert.Equal(t, http.StatusCreated, httpResp.StatusCode)
		databaseProblem := models.Problem{}
		assert.Nil(t, base.DB.Where("name = ?", req.Name).First(&databaseProblem).Error)
		// request == database
		assert.Equal(t, req.Name, databaseProblem.Name)
		assert.False(t, databaseProblem.Public)
		assert.True(t, databaseProblem.Privacy)
		// response == database
		resp := response.AdminCreateProblemResponse{}
		mustJsonDecode(httpResp, &resp)
		jsonEQ(t, response.AdminUpdateProblemResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.ProblemProfileForAdmin `json:"problem"`
			}{
				resource.GetProblemProfileForAdmin(&databaseProblem),
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
				AttachmentFileName: "",
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
					"field":       "AttachmentFileName",
					"reason":      "required",
					"translation": "附件名称为必填字段",
				},
				map[string]interface{}{
					"field":       "LanguageAllowed",
					"reason":      "required",
					"translation": "可用语言为必填字段",
				},
			}),
		},
		{
			name:   "NonExistId",
			method: "PUT",
			path:   "/api/admin/problem/-1",
			req: request.AdminUpdateProblemRequest{
				Name:               "test_admin_update_problem_non_exist",
				AttachmentFileName: "test_admin_update_problem_non_exist_attachment_file_name",
				LanguageAllowed:    "test_admin_update_problem_non_exist_language_allowed",
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
				Name:               "test_admin_update_problem_prem",
				AttachmentFileName: "test_admin_update_problem_perm_attachment_file_name",
				LanguageAllowed:    "test_admin_update_problem_perm_language_allowed",
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
			name: "WithDefaultPublicAndPrivacy",
			path: "id",
			originalProblem: models.Problem{
				Name:               "test_admin_update_problem_2",
				AttachmentFileName: "test_admin_update_problem_2_attachment_file_name",
				LanguageAllowed:    "test_admin_update_problem_2_language_allowed",
				Public:             true,
				Privacy:            false,
			},
			expectedProblem: models.Problem{
				Name:               "test_admin_update_problem_20",
				AttachmentFileName: "test_admin_update_problem_20_attachment_file_name",
				LanguageAllowed:    "test_admin_update_problem_20_language_allowed",
				Public:             false,
				Privacy:            true,
			},
			req: request.AdminUpdateProblemRequest{
				Name:               "test_admin_update_problem_20",
				AttachmentFileName: "test_admin_update_problem_20_attachment_file_name",
				LanguageAllowed:    "test_admin_update_problem_20_language_allowed",
			},
		},
		{
			name: "WithSpecifiedPublicAndPrivacy",
			path: "id",
			originalProblem: models.Problem{
				Name:               "test_admin_update_problem_3",
				AttachmentFileName: "test_admin_update_problem_3_attachment_file_name",
				LanguageAllowed:    "test_admin_update_problem_3_language_allowed",
				Public:             false,
				Privacy:            true,
			},
			expectedProblem: models.Problem{
				Name:               "test_admin_update_problem_30",
				AttachmentFileName: "test_admin_update_problem_30_attachment_file_name",
				LanguageAllowed:    "test_admin_update_problem_30_language_allowed",
				Public:             true,
				Privacy:            false,
			},
			req: request.AdminUpdateProblemRequest{
				Name:               "test_admin_update_problem_30",
				AttachmentFileName: "test_admin_update_problem_30_attachment_file_name",
				LanguageAllowed:    "test_admin_update_problem_30_language_allowed",
				Public:             &boolTrue,
				Privacy:            &boolFalse,
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
						*resource.ProblemProfileForAdmin `json:"problem"`
					}{
						resource.GetProblemProfileForAdmin(&databaseProblem),
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
						*resource.ProblemProfileForAdmin `json:"problem"`
					}{
						resource.GetProblemProfileForAdmin(&test.problem),
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
		Problems []resource.ProblemProfileForAdmin `json:"problems"`
		Total    int                               `json:"total"`
		Count    int                               `json:"count"`
		Offset   int                               `json:"offset"`
		Prev     *string                           `json:"prev"`
		Next     *string                           `json:"next"`
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
				Problems: []resource.ProblemProfileForAdmin{
					*resource.GetProblemProfileForAdmin(&problem1),
					*resource.GetProblemProfileForAdmin(&problem2),
					*resource.GetProblemProfileForAdmin(&problem3),
					*resource.GetProblemProfileForAdmin(&problem4),
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
				Problems: []resource.ProblemProfileForAdmin{},
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
				Problems: []resource.ProblemProfileForAdmin{
					*resource.GetProblemProfileForAdmin(&problem2),
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
				Problems: []resource.ProblemProfileForAdmin{
					*resource.GetProblemProfileForAdmin(&problem4),
					*resource.GetProblemProfileForAdmin(&problem3),
					*resource.GetProblemProfileForAdmin(&problem2),
					*resource.GetProblemProfileForAdmin(&problem1),
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
				Problems: []resource.ProblemProfileForAdmin{
					*resource.GetProblemProfileForAdmin(&problem2),
					*resource.GetProblemProfileForAdmin(&problem3),
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
