package controller_test

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/minio/minio-go"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func checkObject(t *testing.T, bucketName, objectName string) (content []byte, found bool) {
	obj, err := base.Storage.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(obj)
	if err != nil && err.Error() == "The specified key does not exist." {
		assert.Equal(t, []byte{}, content)
		found = false
		return
	}
	assert.Nil(t, err)
	found = true
	return
}

func TestAdminCreateProblem(t *testing.T) {
	t.Parallel()
	FailTests := []failTest{
		{
			name:   "WithoutParams",
			method: "POST",
			path:   base.Echo.Reverse("admin.problem.createProblem"),
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
					"translation": "题目描述为必填字段",
				},
				map[string]interface{}{
					"field":       "Public",
					"reason":      "required",
					"translation": "是否公开为必填字段",
				},
				map[string]interface{}{
					"field":       "Privacy",
					"reason":      "required",
					"translation": "是否可以看到详细评测结果为必填字段",
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
			path:   base.Echo.Reverse("admin.problem.createProblem"),
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

	boolTrue := true
	boolFalse := false

	successTests := []struct {
		name       string
		req        request.AdminCreateProblemRequest
		attachment *fileContent
	}{
		{
			name: "SuccessWithoutAttachment",
			req: request.AdminCreateProblemRequest{
				Name:            "test_admin_create_problem_1",
				Description:     "test_admin_create_problem_1_desc",
				MemoryLimit:     4294967296,
				TimeLimit:       1000,
				LanguageAllowed: "test_admin_create_problem_1_language_allowed",
				CompareScriptID: 1,
				Public:          &boolFalse,
				Privacy:         &boolTrue,
			},
			attachment: nil,
		},
		{
			name: "SuccessWithAttachment",
			req: request.AdminCreateProblemRequest{
				Name:            "test_admin_create_problem_2",
				Description:     "test_admin_create_problem_2_desc",
				MemoryLimit:     4294967296,
				TimeLimit:       1000,
				LanguageAllowed: "test_admin_create_problem_2_language_allowed",
				CompareScriptID: 2,
				Public:          &boolTrue,
				Privacy:         &boolFalse,
			},
			attachment: newFileContent("attachment_file", "test_admin_create_problem_attachment_file", "YXR0YWNobWVudCBmaWxlIGZvciB0ZXN0"),
		},
	}
	t.Run("testAdminCreateProblemSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("testAdminCreateProblem"+test.name, func(t *testing.T) {
				t.Parallel()
				user := createUserForTest(t, "admin_create_problem", i)
				user.GrantRole("admin")
				var data interface{}
				if test.attachment != nil {
					data = addFieldContentSlice([]reqContent{
						test.attachment,
					}, map[string]string{
						"name":              test.req.Name,
						"description":       test.req.Description,
						"memory_limit":      fmt.Sprint(test.req.MemoryLimit),
						"time_limit":        fmt.Sprint(test.req.TimeLimit),
						"language_allowed":  test.req.LanguageAllowed,
						"compare_script_id": fmt.Sprint(test.req.CompareScriptID),
						"public":            fmt.Sprint(*test.req.Public),
						"privacy":           fmt.Sprint(*test.req.Privacy),
					})
				} else {
					data = test.req
				}
				httpReq := makeReq(t, "POST", base.Echo.Reverse("admin.problem.createProblem"), data, headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				})
				httpResp := makeResp(httpReq)
				assert.Equal(t, http.StatusCreated, httpResp.StatusCode)
				databaseProblem := models.Problem{}
				assert.Nil(t, base.DB.Where("name = ?", test.req.Name).First(&databaseProblem).Error)
				// request == database
				assert.Equal(t, test.req.Name, databaseProblem.Name)
				assert.Equal(t, test.req.Description, databaseProblem.Description)
				assert.Equal(t, test.req.MemoryLimit, databaseProblem.MemoryLimit)
				assert.Equal(t, test.req.TimeLimit, databaseProblem.TimeLimit)
				assert.Equal(t, test.req.LanguageAllowed, databaseProblem.LanguageAllowed)
				assert.Equal(t, test.req.CompareScriptID, databaseProblem.CompareScriptID)
				assert.Equal(t, *test.req.Public, databaseProblem.Public)
				assert.Equal(t, *test.req.Privacy, databaseProblem.Privacy)
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
				if test.attachment != nil {
					storageContent, found := checkObject(t, "problems", fmt.Sprintf("%d/attachment", databaseProblem.ID))
					assert.True(t, found)
					expectedContent, err := ioutil.ReadAll(test.attachment.reader)
					assert.Nil(t, err)
					assert.Equal(t, expectedContent, storageContent)
					assert.Equal(t, test.attachment.fileName, databaseProblem.AttachmentFileName)
				} else {
					assert.Equal(t, "", databaseProblem.AttachmentFileName)
				}
			})
		}
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
			path:   base.Echo.Reverse("admin.problem.updateProblem", problem1.ID),
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
					"translation": "题目描述为必填字段",
				},
				map[string]interface{}{
					"field":       "Public",
					"reason":      "required",
					"translation": "是否公开为必填字段",
				},
				map[string]interface{}{
					"field":       "Privacy",
					"reason":      "required",
					"translation": "是否可以看到详细评测结果为必填字段",
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
			path:   base.Echo.Reverse("admin.problem.updateProblem", -1),
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
			path:   base.Echo.Reverse("admin.problem.updateProblem", problem1.ID),
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
		name               string
		path               string
		originalProblem    models.Problem
		expectedProblem    models.Problem
		req                request.AdminUpdateProblemRequest
		updatedAttachment  *fileContent
		originalAttachment *fileContent
		testCases          []models.TestCase
	}{
		{
			name: "WithoutAttachmentAndTestCase",
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
			updatedAttachment: nil,
			testCases:         nil,
		},
		{
			name: "WithAddingAttachment",
			path: "id",
			originalProblem: models.Problem{
				Name:               "test_admin_update_problem_4",
				Description:        "test_admin_update_problem_4_desc",
				LanguageAllowed:    "test_admin_update_problem_4_language_allowed",
				Public:             true,
				Privacy:            true,
				MemoryLimit:        1024,
				TimeLimit:          1000,
				CompareScriptID:    1,
				AttachmentFileName: "",
			},
			expectedProblem: models.Problem{
				Name:               "test_admin_update_problem_40",
				Description:        "test_admin_update_problem_40_desc",
				LanguageAllowed:    "test_admin_update_problem_40_language_allowed",
				Public:             false,
				Privacy:            false,
				MemoryLimit:        2048,
				TimeLimit:          2000,
				CompareScriptID:    2,
				AttachmentFileName: "test_admin_update_problem_attachment_40",
			},
			req: request.AdminUpdateProblemRequest{
				Name:            "test_admin_update_problem_40",
				Description:     "test_admin_update_problem_40_desc",
				LanguageAllowed: "test_admin_update_problem_40_language_allowed",
				Public:          &boolFalse,
				Privacy:         &boolFalse,
				MemoryLimit:     2048,
				TimeLimit:       2000,
				CompareScriptID: 2,
			},
			updatedAttachment: newFileContent("attachment_file", "test_admin_update_problem_attachment_40", "bmV3IGF0dGFjaG1lbnQgZmlsZSBmb3IgdGVzdA"),
			testCases:         nil,
		},
		{
			name: "WithChangingAttachment",
			path: "id",
			originalProblem: models.Problem{
				Name:               "test_admin_update_problem_5",
				Description:        "test_admin_update_problem_5_desc",
				LanguageAllowed:    "test_admin_update_problem_5_language_allowed",
				Public:             true,
				Privacy:            true,
				MemoryLimit:        1024,
				TimeLimit:          1000,
				CompareScriptID:    1,
				AttachmentFileName: "test_admin_update_problem_attachment_5",
			},
			expectedProblem: models.Problem{
				Name:               "test_admin_update_problem_50",
				Description:        "test_admin_update_problem_50_desc",
				LanguageAllowed:    "test_admin_update_problem_50_language_allowed",
				Public:             false,
				Privacy:            false,
				MemoryLimit:        2048,
				TimeLimit:          2000,
				CompareScriptID:    2,
				AttachmentFileName: "test_admin_update_problem_attachment_50",
			},
			req: request.AdminUpdateProblemRequest{
				Name:            "test_admin_update_problem_50",
				Description:     "test_admin_update_problem_50_desc",
				LanguageAllowed: "test_admin_update_problem_50_language_allowed",
				Public:          &boolFalse,
				Privacy:         &boolFalse,
				MemoryLimit:     2048,
				TimeLimit:       2000,
				CompareScriptID: 2,
			},
			originalAttachment: newFileContent("attachment_file", "test_admin_update_problem_attachment_5", "YXR0YWNobWVudCBmaWxlIGZvciB0ZXN0"),
			updatedAttachment:  newFileContent("attachment_file", "test_admin_update_problem_attachment_50", "bmV3IGF0dGFjaG1lbnQgZmlsZSBmb3IgdGVzdA"),
			testCases:          nil,
		},
		{
			name: "WithoutChangingAttachment",
			path: "id",
			originalProblem: models.Problem{
				Name:               "test_admin_update_problem_6",
				Description:        "test_admin_update_problem_6_desc",
				LanguageAllowed:    "test_admin_update_problem_6_language_allowed",
				Public:             true,
				Privacy:            true,
				MemoryLimit:        1024,
				TimeLimit:          1000,
				CompareScriptID:    1,
				AttachmentFileName: "test_admin_update_problem_attachment_6",
			},
			expectedProblem: models.Problem{
				Name:               "test_admin_update_problem_60",
				Description:        "test_admin_update_problem_60_desc",
				LanguageAllowed:    "test_admin_update_problem_60_language_allowed",
				Public:             false,
				Privacy:            false,
				MemoryLimit:        2048,
				TimeLimit:          2000,
				CompareScriptID:    2,
				AttachmentFileName: "test_admin_update_problem_attachment_6",
			},
			req: request.AdminUpdateProblemRequest{
				Name:            "test_admin_update_problem_60",
				Description:     "test_admin_update_problem_60_desc",
				LanguageAllowed: "test_admin_update_problem_60_language_allowed",
				Public:          &boolFalse,
				Privacy:         &boolFalse,
				MemoryLimit:     2048,
				TimeLimit:       2000,
				CompareScriptID: 2,
			},
			originalAttachment: newFileContent("attachment_file", "test_admin_update_problem_attachment_6", "YXR0YWNobWVudCBmaWxlIGZvciB0ZXN0"),
			updatedAttachment:  nil,
			testCases:          nil,
		},
		{
			name: "WithTestCase",
			path: "id",
			originalProblem: models.Problem{
				Name:            "test_admin_update_problem_7",
				Description:     "test_admin_update_problem_7_desc",
				LanguageAllowed: "test_admin_update_problem_7_language_allowed",
				Public:          false,
				Privacy:         true,
				MemoryLimit:     1024,
				TimeLimit:       1000,
				CompareScriptID: 1,
			},
			expectedProblem: models.Problem{
				Name:            "test_admin_update_problem_70",
				Description:     "test_admin_update_problem_70_desc",
				LanguageAllowed: "test_admin_update_problem_70_language_allowed",
				Public:          true,
				Privacy:         false,
				MemoryLimit:     2048,
				TimeLimit:       2000,
				CompareScriptID: 2,
			},
			req: request.AdminUpdateProblemRequest{
				Name:            "test_admin_update_problem_70",
				Description:     "test_admin_update_problem_70_desc",
				LanguageAllowed: "test_admin_update_problem_70_language_allowed",
				Public:          &boolTrue,
				Privacy:         &boolFalse,
				MemoryLimit:     2048,
				TimeLimit:       2000,
				CompareScriptID: 2,
			},
			updatedAttachment: nil,
			testCases: []models.TestCase{
				{
					Score:          100,
					InputFileName:  "test_admin_update_problem_7_test_case_1_input_file_name",
					OutputFileName: "test_admin_update_problem_7_test_case_1_output_file_name",
				},
				{
					Score:          100,
					InputFileName:  "test_admin_update_problem_7_test_case_2_input_file_name",
					OutputFileName: "test_admin_update_problem_7_test_case_2_output_file_name",
				},
			},
		},
	}

	t.Run("testAdminUpdateProblemSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("testAdminUpdateProblem"+test.name, func(t *testing.T) {
				t.Parallel()
				assert.Nil(t, base.DB.Create(&test.originalProblem).Error)
				for j := range test.testCases {
					assert.Nil(t, base.DB.Model(&test.originalProblem).Association("TestCases").Append(&test.testCases[j]).Error)
				}
				if test.originalAttachment != nil {
					b, err := ioutil.ReadAll(test.originalAttachment.reader)
					assert.Nil(t, err)
					_, err = base.Storage.PutObject("problems", fmt.Sprintf("%d/attachment", test.originalProblem.ID), bytes.NewReader(b), int64(len(b)), minio.PutObjectOptions{})
					assert.Nil(t, err)
					test.originalAttachment.reader = bytes.NewReader(b)
				}
				path := base.Echo.Reverse("admin.problem.updateProblem", test.originalProblem.ID)
				user := createUserForTest(t, "admin_update_problem", i)
				user.GrantRole("creator", test.originalProblem)
				var data interface{}
				if test.updatedAttachment != nil {
					data = addFieldContentSlice([]reqContent{
						test.updatedAttachment,
					}, map[string]string{
						"name":              test.req.Name,
						"description":       test.req.Description,
						"memory_limit":      fmt.Sprint(test.req.MemoryLimit),
						"time_limit":        fmt.Sprint(test.req.TimeLimit),
						"language_allowed":  test.req.LanguageAllowed,
						"compare_script_id": fmt.Sprint(test.req.CompareScriptID),
						"public":            fmt.Sprint(*test.req.Public),
						"privacy":           fmt.Sprint(*test.req.Privacy),
					})
				} else {
					data = test.req
				}
				httpResp := makeResp(makeReq(t, "PUT", path, data, headerOption{
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
				assert.Nil(t, base.DB.Set("gorm:auto_preload", true).Model(databaseProblem).Related(&databaseProblem.TestCases).Error)
				if test.testCases != nil {
					jsonEQ(t, test.testCases, databaseProblem.TestCases)
				} else {
					assert.Equal(t, []models.TestCase{}, databaseProblem.TestCases)
				}
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
				if test.updatedAttachment != nil || test.originalAttachment != nil {
					storageContent, found := checkObject(t, "problems", fmt.Sprintf("%d/attachment", databaseProblem.ID))
					assert.True(t, found)
					var err error
					var expectedContent []byte
					if test.updatedAttachment != nil {
						expectedContent, err = ioutil.ReadAll(test.updatedAttachment.reader)
					} else {
						expectedContent, err = ioutil.ReadAll(test.originalAttachment.reader)
					}
					assert.Nil(t, err)
					assert.Equal(t, expectedContent, storageContent)
				}
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
			path:   base.Echo.Reverse("admin.problem.deleteProblem", -1),
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
			path:   base.Echo.Reverse("admin.problem.deleteProblem", -1),
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
		name               string
		problem            models.Problem
		originalAttachment *fileContent
		testCases          []struct {
			testcase   models.TestCase
			inputFile  *fileContent
			outputFile *fileContent
		}
	}{
		{
			name: "SuccessWithoutAttachmentAndTestCases",
			problem: models.Problem{
				Name:               "test_admin_delete_problem_1",
				AttachmentFileName: "",
				LanguageAllowed:    "test_admin_delete_problem_1_language_allowed",
			},
			originalAttachment: nil,
			testCases:          nil,
		},
		{
			name: "SuccessWithAttachment",
			problem: models.Problem{
				Name:               "test_admin_delete_problem_2",
				AttachmentFileName: "test_admin_delete_problem_attachment_2",
				LanguageAllowed:    "test_admin_delete_problem_2_language_allowed",
			},
			originalAttachment: newFileContent("attachment_file", "test_admin_delete_problem_attachment_2", "YXR0YWNobWVudCBmaWxlIGZvciB0ZXN0"),
			testCases:          nil,
		},
		{
			name: "SuccessWithTestCases",
			problem: models.Problem{
				Name:               "test_admin_delete_problem_3",
				AttachmentFileName: "",
				LanguageAllowed:    "test_admin_delete_problem_3_language_allowed",
			},
			originalAttachment: nil,
			testCases: []struct {
				testcase   models.TestCase
				inputFile  *fileContent
				outputFile *fileContent
			}{
				{
					testcase: models.TestCase{
						Score:          100,
						InputFileName:  "test_admin_delete_problem_3_test_case_1_input_file_name",
						OutputFileName: "test_admin_delete_problem_3_test_case_1_output_file_name",
					},
					inputFile:  newFileContent("input_file", "test_admin_delete_problem_3_test_case_1.in", "aW5wdXQgdGV4dAo="),
					outputFile: newFileContent("output_file", "test_admin_delete_problem_3_test_case_1.out", "b3V0cHV0IHRleHQK"),
				},
				{
					testcase: models.TestCase{
						Score:          100,
						InputFileName:  "test_admin_delete_problem_3_test_case_2_input_file_name",
						OutputFileName: "test_admin_delete_problem_3_test_case_2_output_file_name",
					},
					inputFile:  newFileContent("input_file", "test_admin_delete_problem_3_test_case_2.in", "aW5wdXQgdGV4dAo="),
					outputFile: newFileContent("output_file", "test_admin_delete_problem_3_test_case_2.out", "b3V0cHV0IHRleHQK"),
				},
			},
		},
		{
			name: "SuccessWithAttachmentAndTestCases",
			problem: models.Problem{
				Name:               "test_admin_delete_problem_4",
				AttachmentFileName: "test_admin_delete_problem_attachment_4",
				LanguageAllowed:    "test_admin_delete_problem_4_language_allowed",
			},
			originalAttachment: newFileContent("attachment_file", "test_admin_delete_problem_attachment_4", "YXR0YWNobWVudCBmaWxlIGZvciB0ZXN0"),
			testCases: []struct {
				testcase   models.TestCase
				inputFile  *fileContent
				outputFile *fileContent
			}{
				{
					testcase: models.TestCase{
						Score:          100,
						InputFileName:  "test_admin_delete_problem_4_test_case_1_input_file_name",
						OutputFileName: "test_admin_delete_problem_4_test_case_1_output_file_name",
					},
					inputFile:  newFileContent("input_file", "test_admin_delete_problem_4_test_case_1.in", "aW5wdXQgdGV4dAo="),
					outputFile: newFileContent("output_file", "test_admin_delete_problem_4_test_case_1.out", "b3V0cHV0IHRleHQK"),
				},
				{
					testcase: models.TestCase{
						Score:          100,
						InputFileName:  "test_admin_delete_problem_4_test_case_2_input_file_name",
						OutputFileName: "test_admin_delete_problem_4_test_case_2_output_file_name",
					},
					inputFile:  newFileContent("input_file", "test_admin_delete_problem_4_test_case_2.in", "aW5wdXQgdGV4dAo="),
					outputFile: newFileContent("output_file", "test_admin_delete_problem_4_test_case_2.out", "b3V0cHV0IHRleHQK"),
				},
			},
		},
	}

	t.Run("testAdminDeleteProblemSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("testAdminDeleteProblem"+test.name, func(t *testing.T) {
				t.Parallel()
				assert.Nil(t, base.DB.Create(&test.problem).Error)
				for j := range test.testCases {
					createTestCaseForTest(t, test.problem, test.testCases[j].testcase.Score, test.testCases[j].inputFile, test.testCases[j].outputFile)
				}
				if test.originalAttachment != nil {
					b, err := ioutil.ReadAll(test.originalAttachment.reader)
					assert.Nil(t, err)
					_, err = base.Storage.PutObject("problems", fmt.Sprintf("%d/attachment", test.problem.ID), bytes.NewReader(b), int64(len(b)), minio.PutObjectOptions{})
					assert.Nil(t, err)
					test.originalAttachment.reader = bytes.NewReader(b)
				}
				user := createUserForTest(t, "admin_delete_problem", i)
				user.GrantRole("creator", test.problem)
				httpResp := makeResp(makeReq(t, "DELETE", base.Echo.Reverse("admin.problem.deleteProblem", test.problem.ID), request.AdminDeleteProblemRequest{}, headerOption{
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
				if test.originalAttachment != nil {
					_, found := checkObject(t, "problems", fmt.Sprintf("%d/attachment", test.problem.ID))
					assert.False(t, found)
				}
				for j := range test.testCases {
					assert.Equal(t, gorm.ErrRecordNotFound, base.DB.First(models.TestCase{}, test.testCases[j].testcase.ID).Error)
					_, found := checkObject(t, "problems", fmt.Sprintf("%d/input/%s", test.problem.ID, test.testCases[j].inputFile.fileName))
					assert.False(t, found)
					_, found = checkObject(t, "problems", fmt.Sprintf("%d/output/%s", test.problem.ID, test.testCases[j].outputFile.fileName))
					assert.False(t, found)
				}
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
			path:   base.Echo.Reverse("admin.problem.getProblem", -1),
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
			path:   base.Echo.Reverse("admin.problem.getProblem", -1),
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
		name      string
		path      string
		req       request.AdminGetProblemRequest
		problem   models.Problem
		testCases []models.TestCase
	}{
		{
			name: "WithoutTestCases",
			path: "id",
			req:  request.AdminGetProblemRequest{},
			problem: models.Problem{
				Name:               "test_admin_get_problem_1",
				AttachmentFileName: "test_admin_get_problem_1_attachment_file_name",
				LanguageAllowed:    "test_admin_get_problem_1_language_allowed",
			},
			testCases: nil,
		},
		{
			name: "WithTestCases",
			path: "id",
			req:  request.AdminGetProblemRequest{},
			problem: models.Problem{
				Name:               "test_admin_get_problem_2",
				AttachmentFileName: "test_admin_get_problem_2_attachment_file_name",
				LanguageAllowed:    "test_admin_get_problem_2_language_allowed",
			},
			testCases: []models.TestCase{
				{
					Score:          100,
					InputFileName:  "test_admin_get_problem_2_test_case_1_input_file_name",
					OutputFileName: "test_admin_get_problem_2_test_case_1_output_file_name",
				},
				{
					Score:          100,
					InputFileName:  "test_admin_get_problem_2_test_case_2_input_file_name",
					OutputFileName: "test_admin_get_problem_2_test_case_2_output_file_name",
				},
			},
		},
	}

	t.Run("testAdminGetProblemSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("testAdminGetProblem"+test.name, func(t *testing.T) {
				t.Parallel()
				assert.Nil(t, base.DB.Create(&test.problem).Error)
				for j := range test.testCases {
					assert.Nil(t, base.DB.Model(&test.problem).Association("TestCases").Append(&test.testCases[j]).Error)
				}
				user := createUserForTest(t, "admin_get_problem", i)
				user.GrantRole("creator", test.problem)
				httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("admin.problem.getProblem", test.problem.ID), request.AdminGetUserRequest{}, headerOption{
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

	failTests := []failTest{
		{
			name:   "NotAllowedColumn",
			method: "GET",
			path:   base.Echo.Reverse("admin.problem.getProblems"),
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
			path:   base.Echo.Reverse("admin.problem.getProblems"),
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
				Next: getUrlStringPointer("admin.problem.getProblems", map[string]string{
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
				httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("admin.problem.getProblems"), test.req, applyAdminUser))
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
		inputFile.reader = bytes.NewReader(inputBytes)
		_, err = base.Storage.PutObject("problems", fmt.Sprintf("%d/input/%s", problem.ID, inputFile.fileName), bytes.NewReader(inputBytes), int64(len(inputBytes)), minio.PutObjectOptions{})
		assert.Nil(t, err)
		inputFileName = inputFile.fileName
	}
	if outputFile != nil {
		utils.MustCreateBucket("problems")
		outputBytes, err := ioutil.ReadAll(outputFile.reader)
		assert.Nil(t, err)
		outputFile.reader = bytes.NewReader(outputBytes)
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
			name:   "NonExistingProblem",
			method: "POST",
			path:   base.Echo.Reverse("admin.problem.createTestCase", -1),
			req: addFieldContentSlice([]reqContent{
				newFileContent("input_file", "test_admin_create_test_case_non_existing_problem.in", "aW5wdXQgdGV4dAo="),
				newFileContent("output_file", "test_admin_create_test_case_non_existing_problem.out", "b3V0cHV0IHRleHQK"),
			}, map[string]string{
				"score": "100",
			}),
			reqOptions: []reqOption{
				headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				},
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
		{
			name:   "LackInputFile",
			method: "POST",
			path:   base.Echo.Reverse("admin.problem.createTestCase", problem.ID),
			req: addFieldContentSlice([]reqContent{
				newFileContent("output_file", "test_admin_create_test_case_lack_input_file.out", "b3V0cHV0IHRleHQK"),
			}, map[string]string{
				"score": "100",
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
			path:   base.Echo.Reverse("admin.problem.createTestCase", problem.ID),
			req: addFieldContentSlice([]reqContent{
				newFileContent("input_file", "test_admin_create_test_case_lack_output_file.in", "aW5wdXQgdGV4dAo="),
			}, map[string]string{
				"score": "100",
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
			path:   base.Echo.Reverse("admin.problem.createTestCase", problem.ID),
			req: addFieldContentSlice([]reqContent{}, map[string]string{
				"score": "100",
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
			path:   base.Echo.Reverse("admin.problem.createTestCase", problem.ID),
			req: addFieldContentSlice([]reqContent{
				newFileContent("input_file", "test_admin_create_test_case_permission_denied.in", "aW5wdXQgdGV4dAo="),
				newFileContent("output_file", "test_admin_create_test_case_permission_denied.out", "b3V0cHV0IHRleHQK"),
			}, map[string]string{
				"score": "100",
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
		t.Parallel()
		req := makeReq(t, "POST", base.Echo.Reverse("admin.problem.createTestCase", problem.ID), addFieldContentSlice([]reqContent{
			newFileContent("input_file", "test_admin_create_test_case_success.in", "aW5wdXQgdGV4dAo="),
			newFileContent("output_file", "test_admin_create_test_case_success.out", "b3V0cHV0IHRleHQK"),
		}, map[string]string{
			"score": "100",
		}), headerOption{
			"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
		})
		//req.Header.Set("Set-User-For-Test", fmt.Sprintf("%d", user.ID))
		httpResp := makeResp(req)
		expectedTestCase := models.TestCase{
			ProblemID:      problem.ID,
			Score:          100,
			InputFileName:  "test_admin_create_test_case_success.in",
			OutputFileName: "test_admin_create_test_case_success.out",
		}
		databaseTestCase := models.TestCase{}
		assert.Nil(t, base.DB.Where("problem_id = ? and input_file_name = ?", problem.ID, "test_admin_create_test_case_success.in").First(&databaseTestCase).Error)
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
			path:   base.Echo.Reverse("admin.problem.getTestCaseInputFile", -1, 1),
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
			path:   base.Echo.Reverse("admin.problem.getTestCaseInputFile", problem.ID, -1),
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
			path:   base.Echo.Reverse("admin.problem.getTestCaseInputFile", problem.ID, 1),
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
		newFileContent("", "test_admin_get_test_case_input_file_success.in", "aW5wdXQgdGV4dAo="),
		//newFileContent("","test_admin_get_test_case_input_file_success.out","b3V0cHV0IHRleHQK"),
		nil,
	)

	req := makeReq(t, "GET", base.Echo.Reverse("admin.problem.getTestCaseInputFile", problem.ID, testCase.ID), request.AdminGetTestCaseInputFileRequest{}, headerOption{
		"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
	})
	httpResp := makeResp(req)

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "input text\n", string(respBytes))
}

func TestAdminGetTestCaseOutputFile(t *testing.T) {
	problem, user := createProblemForTest(t, "admin_get_test_case_output_file", 0)
	failTests := []failTest{
		{
			name:   "NonExistingProblem",
			method: "GET",
			path:   base.Echo.Reverse("admin.problem.getTestCaseOutputFile", -1, 1),
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
			path:   base.Echo.Reverse("admin.problem.getTestCaseOutputFile", problem.ID, -1),
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
			path:   base.Echo.Reverse("admin.problem.getTestCaseOutputFile", problem.ID, 1),
			req:    request.AdminGetTestCaseOutputFileRequest{},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "AdminGetTestCaseOutputFile")

	testCase := createTestCaseForTest(t, problem, 52,
		nil,
		newFileContent("", "test_admin_get_test_case_output_file_success.out", "b3V0cHV0IHRleHQK"),
	)

	req := makeReq(t, "GET", base.Echo.Reverse("admin.problem.getTestCaseOutputFile", problem.ID, testCase.ID), request.AdminGetTestCaseOutputFileRequest{}, headerOption{
		"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
	})
	httpResp := makeResp(req)

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "output text\n", string(respBytes))
}

func TestAdminUpdateTestCase(t *testing.T) {
	problem, user := createProblemForTest(t, "admin_update_test_case", 0)

	failTests := []failTest{
		{
			name:   "NonExistingProblem",
			method: "PUT",
			path:   base.Echo.Reverse("admin.problem.updateTestCase", -1, 1),
			req: request.AdminUpdateTestCaseRequest{
				Score: 100,
			},
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
			method: "PUT",
			path:   base.Echo.Reverse("admin.problem.updateTestCase", problem.ID, -1),
			req: request.AdminUpdateTestCaseRequest{
				Score: 100,
			},
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
			method: "PUT",
			path:   base.Echo.Reverse("admin.problem.updateTestCase", problem.ID, 1),
			req: request.AdminUpdateTestCaseRequest{
				Score: 100,
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "AdminUpdateTestCase")

	successTests := []struct {
		name               string
		originalScore      uint
		updatedScore       uint
		originalInputFile  *fileContent
		originalOutputFile *fileContent
		updatedInputFile   *fileContent
		updatedOutputFile  *fileContent
		expectedTestCase   models.TestCase
	}{
		{
			name:               "SuccessWithoutUpdatingFile",
			originalScore:      0,
			updatedScore:       100,
			originalInputFile:  newFileContent("input_file", "test_update_test_case_1.in", "aW5wdXQgdGV4dAo="),
			originalOutputFile: newFileContent("output_file", "test_update_test_case_1.out", "b3V0cHV0IHRleHQK"),
			updatedInputFile:   nil,
			updatedOutputFile:  nil,
			expectedTestCase: models.TestCase{
				ProblemID:      problem.ID,
				Score:          100,
				InputFileName:  "test_update_test_case_1.in",
				OutputFileName: "test_update_test_case_1.out",
			},
		},
		{
			name:               "SuccessWithUpdatingInputFile",
			originalScore:      0,
			updatedScore:       100,
			originalInputFile:  newFileContent("input_file", "test_update_test_case_2.in", "aW5wdXQgdGV4dAo="),
			originalOutputFile: newFileContent("output_file", "test_update_test_case_2.out", "b3V0cHV0IHRleHQK"),
			updatedInputFile:   newFileContent("input_file", "test_update_test_case_20.in", "bmV3IGlucHV0IHRleHQ"),
			updatedOutputFile:  nil,
			expectedTestCase: models.TestCase{
				ProblemID:      problem.ID,
				Score:          100,
				InputFileName:  "test_update_test_case_20.in",
				OutputFileName: "test_update_test_case_2.out",
			},
		},
		{
			name:               "SuccessWithUpdatingOutputFile",
			originalScore:      0,
			updatedScore:       100,
			originalInputFile:  newFileContent("input_file", "test_update_test_case_3.in", "aW5wdXQgdGV4dAo="),
			originalOutputFile: newFileContent("output_file", "test_update_test_case_3.out", "b3V0cHV0IHRleHQK"),
			updatedInputFile:   nil,
			updatedOutputFile:  newFileContent("output_file", "test_update_test_case_30.out", "bmV3IG91dHB1dCB0ZXh0"),
			expectedTestCase: models.TestCase{
				ProblemID:      problem.ID,
				Score:          100,
				InputFileName:  "test_update_test_case_3.in",
				OutputFileName: "test_update_test_case_30.out",
			},
		},
		{
			name:               "SuccessWithUpdatingBothFile",
			originalScore:      0,
			updatedScore:       100,
			originalInputFile:  newFileContent("input_file", "test_update_test_case_4.in", "aW5wdXQgdGV4dAo="),
			originalOutputFile: newFileContent("output_file", "test_update_test_case_4.out", "b3V0cHV0IHRleHQK"),
			updatedInputFile:   newFileContent("input_file", "test_update_test_case_40.in", "bmV3IGlucHV0IHRleHQ"),
			updatedOutputFile:  newFileContent("output_file", "test_update_test_case_40.out", "bmV3IG91dHB1dCB0ZXh0"),
			expectedTestCase: models.TestCase{
				ProblemID:      problem.ID,
				Score:          100,
				InputFileName:  "test_update_test_case_40.in",
				OutputFileName: "test_update_test_case_40.out",
			},
		},
	}

	t.Run("testAdminUpdateTestCaseSuccess", func(t *testing.T) {
		t.Parallel()
		for _, test := range successTests {
			test := test
			t.Run("testAdminUpdateTestCase"+test.name, func(t *testing.T) {
				t.Parallel()
				testCase := createTestCaseForTest(t, problem, test.originalScore, test.originalInputFile, test.originalOutputFile)
				var reqContentSlice []reqContent
				if test.updatedInputFile != nil {
					reqContentSlice = append(reqContentSlice, test.updatedInputFile)
				}
				if test.updatedOutputFile != nil {
					reqContentSlice = append(reqContentSlice, test.updatedOutputFile)
				}
				req := makeReq(t, "PUT", base.Echo.Reverse("admin.problem.updateTestCase", problem.ID, testCase.ID), addFieldContentSlice(
					reqContentSlice, map[string]string{
						"score": fmt.Sprintf("%d", test.updatedScore),
					}), headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				})
				httpResp := makeResp(req)
				databaseTestCase := models.TestCase{}
				base.DB.First(&databaseTestCase, testCase.ID)
				assert.Equal(t, test.expectedTestCase.ProblemID, databaseTestCase.ProblemID)
				assert.Equal(t, test.expectedTestCase.Score, databaseTestCase.Score)
				assert.Equal(t, test.expectedTestCase.InputFileName, databaseTestCase.InputFileName)
				assert.Equal(t, test.expectedTestCase.OutputFileName, databaseTestCase.OutputFileName)

				var expectedInputContent []byte
				var err error
				if test.updatedInputFile == nil || test.originalInputFile.fileName == test.updatedInputFile.fileName {
					if test.updatedInputFile == nil {
						expectedInputContent, err = ioutil.ReadAll(test.originalInputFile.reader)
					} else {
						expectedInputContent, err = ioutil.ReadAll(test.updatedInputFile.reader)
					}
					assert.Nil(t, err)
				} else {
					_, found := checkObject(t, "problems", fmt.Sprintf("%d/input/%s", problem.ID, test.originalInputFile.fileName))
					assert.False(t, found)
					expectedInputContent, err = ioutil.ReadAll(test.updatedInputFile.reader)
					assert.Nil(t, err)
				}
				storageInputContent, found := checkObject(t, "problems", fmt.Sprintf("%d/input/%s", problem.ID, databaseTestCase.InputFileName))
				assert.True(t, found)
				assert.Equal(t, expectedInputContent, storageInputContent)

				var expectedOutputContent []byte
				if test.updatedOutputFile == nil || test.originalOutputFile.fileName == test.updatedOutputFile.fileName {
					if test.updatedOutputFile == nil {
						expectedOutputContent, err = ioutil.ReadAll(test.originalOutputFile.reader)
					} else {
						expectedOutputContent, err = ioutil.ReadAll(test.updatedOutputFile.reader)
					}
					assert.Nil(t, err)
				} else {
					_, found := checkObject(t, "problems", fmt.Sprintf("%d/output/%s", problem.ID, test.originalOutputFile.fileName))
					assert.False(t, found)
					expectedOutputContent, err = ioutil.ReadAll(test.updatedOutputFile.reader)
					assert.Nil(t, err)
				}
				storageOutputContent, found := checkObject(t, "problems", fmt.Sprintf("%d/output/%s", problem.ID, databaseTestCase.OutputFileName))
				assert.True(t, found)
				assert.Equal(t, expectedOutputContent, storageOutputContent)

				resp := response.AdminUpdateTestCaseResponse{}
				mustJsonDecode(httpResp, &resp)
				assert.Equal(t, response.AdminUpdateTestCaseResponse{
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
	})
}

func TestAdminDeleteTestCase(t *testing.T) {
	problem, user := createProblemForTest(t, "admin_delete_test_case", 0)

	failTests := []failTest{
		{
			name:   "NonExistingProblem",
			method: "DELETE",
			path:   base.Echo.Reverse("admin.problem.deleteTestCases", -1, 1),
			req: request.AdminUpdateTestCaseRequest{
				Score: 100,
			},
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
			method: "DELETE",
			path:   base.Echo.Reverse("admin.problem.deleteTestCase", problem.ID, -1),
			req: request.AdminUpdateTestCaseRequest{
				Score: 100,
			},
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
			method: "DELETE",
			path:   base.Echo.Reverse("admin.problem.deleteTestCases", problem.ID, 1),
			req: request.AdminUpdateTestCaseRequest{
				Score: 100,
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "AdminDeleteTestCase")

	t.Run("testAdminDeleteTestCaseSuccess", func(t *testing.T) {
		t.Parallel()
		testCase := createTestCaseForTest(t, problem, 72,
			newFileContent("input_file", "test_delete_test_case_0.in", "aW5wdXQgdGV4dAo="),
			newFileContent("output_file", "test_delete_test_case_0.out", "b3V0cHV0IHRleHQK"),
		)

		req := makeReq(t, "DELETE", base.Echo.Reverse("admin.problem.deleteTestCases", problem.ID, testCase.ID), request.AdminDeleteTestCaseRequest{}, headerOption{
			"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
		})
		httpResp := makeResp(req)
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.Equal(t, response.Response{
			Message: "SUCCESS",
			Error:   nil,
			Data:    nil,
		}, resp)
		databaseTestcase := models.TestCase{}
		err := base.DB.First(&databaseTestcase, testCase.ID).Error
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestAdminDeleteTestCases(t *testing.T) {
	t.Parallel()
	problem, user := createProblemForTest(t, "admin_delete_test_cases", 0)

	failTests := []failTest{
		{
			name:   "NonExistingProblem",
			method: "DELETE",
			path:   base.Echo.Reverse("admin.problem.deleteTestCases", -1),
			req: request.AdminUpdateTestCaseRequest{
				Score: 100,
			},
			reqOptions: []reqOption{
				headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				},
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
		{
			name:   "PermissionDenied",
			method: "DELETE",
			path:   base.Echo.Reverse("admin.problem.deleteTestCases", problem.ID),
			req: request.AdminUpdateTestCaseRequest{
				Score: 100,
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "AdminDeleteTestCases")

	t.Run("testAdminDeleteTestCasesSuccess", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 5; i++ {
			createTestCaseForTest(t, problem, 0,
				newFileContent("input_file", fmt.Sprintf("test_delete_test_cases_%d.in", i), "aW5wdXQgdGV4dAo="),
				newFileContent("output_file", fmt.Sprintf("test_delete_test_cases_%d.out", i), "b3V0cHV0IHRleHQK"),
			)
		}
		req := makeReq(t, "DELETE", base.Echo.Reverse("admin.problem.deleteTestCases", problem.ID), request.AdminDeleteTestCasesRequest{}, headerOption{
			"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
		})
		httpResp := makeResp(req)
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		var databaseTestCases []models.TestCase
		assert.Nil(t, base.DB.Find(&databaseTestCases, "problem_id = ?", problem.ID).Error)
		assert.Equal(t, []models.TestCase{}, databaseTestCases)
	})
}
