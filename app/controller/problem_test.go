package controller_test

import (
	"bytes"
	"fmt"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/minio/minio-go"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

var inputTextBase64 = "aW5wdXQgdGV4dAo="
var outputTextBase64 = "b3V0cHV0IHRleHQK"
var newInputTextBase64 = "bmV3IGlucHV0IHRleHQK"
var newOutputTextBase64 = "bmV3IG91dHB1dCB0ZXh0"
var attachmentFileBase64 = "YXR0YWNobWVudCBmaWxlIGZvciB0ZXN0"
var newAttachmentFileBase64 = "bmV3IGF0dGFjaG1lbnQgZmlsZSBmb3IgdGVzdAo="

func getObjectContent(t *testing.T, bucketName, objectName string) (content []byte) {
	obj, err := base.Storage.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(obj)
	assert.Nil(t, err)
	return
}

func createProblemForTest(t *testing.T, name string, id int, attachmentFile *fileContent) (problem models.Problem, user models.User) {
	problem = models.Problem{
		Name:               fmt.Sprintf("problem_for_testing_%s_%d", name, id),
		Description:        fmt.Sprintf("a problem used to test API: %s(%d)", name, id),
		AttachmentFileName: "",
		Public:             true,
		Privacy:            false,
		MemoryLimit:        1024,
		TimeLimit:          1000,
		LanguageAllowed:    fmt.Sprintf("test_%s_language_allowed_%d", name, id),
		CompileEnvironment: fmt.Sprintf("test_%s_compile_environment_%d", name, id),
		CompareScriptID:    1,
	}
	if attachmentFile != nil {
		problem.AttachmentFileName = attachmentFile.fileName
	}
	assert.Nil(t, base.DB.Create(&problem).Error)
	user = createUserForTest(t, name, id)
	user.GrantRole("problem_creator", problem)
	if attachmentFile != nil {
		attachmentBytes, err := ioutil.ReadAll(attachmentFile.reader)
		assert.Nil(t, err)
		_, err = attachmentFile.reader.Seek(0, io.SeekStart)
		assert.Nil(t, err)
		_, err = base.Storage.PutObject("problems", fmt.Sprintf("%d/attachment", problem.ID), bytes.NewReader(attachmentBytes), int64(len(attachmentBytes)), minio.PutObjectOptions{})
		assert.Nil(t, err)
	}
	return
}

func createTestCaseForTest(t *testing.T, problem models.Problem, score uint, inputFile, outputFile *fileContent) (testCase models.TestCase) {
	var inputFileName, outputFileName string

	if inputFile != nil {
		inputFileName = inputFile.fileName
	}
	if outputFile != nil {
		outputFileName = outputFile.fileName
	}

	testCase = models.TestCase{
		Score:          score,
		InputFileName:  inputFileName,
		OutputFileName: outputFileName,
	}
	assert.Nil(t, base.DB.Model(&problem).Association("TestCases").Append(&testCase))

	if inputFile != nil {
		inputBytes, err := ioutil.ReadAll(inputFile.reader)
		assert.Nil(t, err)
		inputFile.reader = bytes.NewReader(inputBytes)
		_, err = base.Storage.PutObject("problems", fmt.Sprintf("%d/input/%d.in", problem.ID, testCase.ID), bytes.NewReader(inputBytes), int64(len(inputBytes)), minio.PutObjectOptions{})
		assert.Nil(t, err)
	}
	if outputFile != nil {
		outputBytes, err := ioutil.ReadAll(outputFile.reader)
		assert.Nil(t, err)
		outputFile.reader = bytes.NewReader(outputBytes)
		_, err = base.Storage.PutObject("problems", fmt.Sprintf("%d/output/%d.out", problem.ID, testCase.ID), bytes.NewReader(outputBytes), int64(len(outputBytes)), minio.PutObjectOptions{})
		assert.Nil(t, err)
	}

	return
}

func TestGetProblem(t *testing.T) {
	t.Parallel()

	// publicFalseProblem means a problem which "public" field is false
	publicFalseProblem := models.Problem{
		Name:               "test_get_problem_public_false",
		AttachmentFileName: "test_get_problem_public_false_attachment_file_name",
		LanguageAllowed:    "test_get_problem_public_false_language_allowed",
	}
	assert.Nil(t, base.DB.Create(&publicFalseProblem).Error)

	failTests := []failTest{
		{
			name:   "NonExistId",
			method: "GET",
			path:   base.Echo.Reverse("problem.getProblem", -1),
			req:    request.GetProblemRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			name:   "PublicFalseFail",
			method: "GET",
			path:   base.Echo.Reverse("problem.getProblem", publicFalseProblem.ID),
			req:    request.GetProblemRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
	}

	runFailTests(t, failTests, "GetProblem")

	successTests := []struct {
		name      string
		path      string
		req       request.GetProblemRequest
		problem   models.Problem
		isAdmin   bool
		testCases []models.TestCase
	}{
		{
			name: "AdminUserWithoutTestCases",
			path: "id",
			req:  request.GetProblemRequest{},
			problem: models.Problem{
				Name:               "test_get_problem_1",
				AttachmentFileName: "test_get_problem_1_attachment_file_name",
				LanguageAllowed:    "test_get_problem_1_language_allowed",
				Public:             true,
			},
			isAdmin:   true,
			testCases: nil,
		},
		{
			name: "NormalUserWithoutTestCases",
			path: "id",
			req:  request.GetProblemRequest{},
			problem: models.Problem{
				Name:               "test_get_problem_2",
				AttachmentFileName: "test_get_problem_2_attachment_file_name",
				LanguageAllowed:    "test_get_problem_2_language_allowed",
				Public:             true,
			},
			isAdmin:   false,
			testCases: nil,
		},
		{
			name: "PublicFalseSuccess",
			path: "id",
			req:  request.GetProblemRequest{},
			problem: models.Problem{
				Name:               "test_get_problem_3",
				AttachmentFileName: "test_get_problem_3_attachment_file_name",
				LanguageAllowed:    "test_get_problem_3_language_allowed",
				Public:             false,
			},
			isAdmin:   true,
			testCases: nil,
		},
		{
			name: "AdminUserWithTestCases",
			path: "id",
			req:  request.GetProblemRequest{},
			problem: models.Problem{
				Name:               "test_admin_get_problem_4",
				AttachmentFileName: "test_admin_get_problem_4_attachment_file_name",
				LanguageAllowed:    "test_admin_get_problem_4_language_allowed",
				Public:             true,
			},
			isAdmin: true,
			testCases: []models.TestCase{
				{
					Score:          100,
					InputFileName:  "test_admin_get_problem_4_test_case_1_input_file_name",
					OutputFileName: "test_admin_get_problem_4_test_case_1_output_file_name",
				},
				{
					Score:          100,
					InputFileName:  "test_admin_get_problem_4_test_case_2_input_file_name",
					OutputFileName: "test_admin_get_problem_4_test_case_2_output_file_name",
				},
			},
		},
		{
			name: "NormalUserWithTestCases",
			path: "id",
			req:  request.GetProblemRequest{},
			problem: models.Problem{
				Name:               "test_admin_get_problem_5",
				AttachmentFileName: "test_admin_get_problem_5_attachment_file_name",
				LanguageAllowed:    "test_admin_get_problem_5_language_allowed",
				Public:             true,
			},
			isAdmin: false,
			testCases: []models.TestCase{
				{
					Score:          100,
					InputFileName:  "test_admin_get_problem_5_test_case_1_input_file_name",
					OutputFileName: "test_admin_get_problem_5_test_case_1_output_file_name",
				},
				{
					Score:          100,
					InputFileName:  "test_admin_get_problem_5_test_case_2_input_file_name",
					OutputFileName: "test_admin_get_problem_5_test_case_2_output_file_name",
				},
			},
		},
	}

	t.Run("testGetProblemSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("testGetProblem"+test.name, func(t *testing.T) {
				t.Parallel()
				assert.Nil(t, base.DB.Create(&test.problem).Error)
				for j := range test.testCases {
					assert.Nil(t, base.DB.Model(&test.problem).Association("TestCases").Append(&test.testCases[j]))
				}
				user := createUserForTest(t, "get_problem", i)
				if test.isAdmin {
					user.GrantRole("problem_creator", test.problem)
				}
				httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("problem.getProblem", test.problem.ID), request.GetUserRequest{}, headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				}))
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				if test.isAdmin {
					resp := response.AdminGetProblemResponse{}
					expectResp := response.AdminGetProblemResponse{
						Message: "SUCCESS",
						Error:   nil,
						Data: struct {
							*resource.ProblemForAdmin `json:"problem"`
						}{
							resource.GetProblemForAdmin(&test.problem),
						},
					}
					mustJsonDecode(httpResp, &resp)
					assert.Equal(t, expectResp, resp)
				} else {
					resp := response.GetProblemResponse{}
					expectResp := response.GetProblemResponse{
						Message: "SUCCESS",
						Error:   nil,
						Data: struct {
							*resource.Problem `json:"problem"`
						}{
							resource.GetProblem(&test.problem),
						},
					}
					mustJsonDecode(httpResp, &resp)
					assert.Equal(t, expectResp, resp)
				}
			})
		}
	})
}

func TestGetProblems(t *testing.T) {
	t.Parallel()
	problem1 := models.Problem{
		Name:               "test_get_problems_1",
		AttachmentFileName: "test_get_problems_1_attachment_file_name",
		LanguageAllowed:    "test_get_problems_1_language_allowed",
		Public:             true,
	}
	problem2 := models.Problem{
		Name:               "test_get_problems_2",
		AttachmentFileName: "test_get_problems_2_attachment_file_name",
		LanguageAllowed:    "test_get_problems_2_language_allowed",
		Public:             true,
	}
	problem3 := models.Problem{
		Name:               "test_get_problems_3",
		AttachmentFileName: "test_get_problems_3_attachment_file_name",
		LanguageAllowed:    "test_get_problems_3_language_allowed",
		Public:             true,
	}
	problem4 := models.Problem{
		Name:               "test_get_problems_4",
		AttachmentFileName: "test_get_problems_4_attachment_file_name",
		LanguageAllowed:    "test_get_problems_4_language_allowed",
		Public:             false,
	}
	assert.Nil(t, base.DB.Create(&problem1).Error)
	assert.Nil(t, base.DB.Create(&problem2).Error)
	assert.Nil(t, base.DB.Create(&problem3).Error)
	assert.Nil(t, base.DB.Create(&problem4).Error)

	type respData struct {
		Problems []models.Problem `json:"problems"`
		Total    int              `json:"total"`
		Count    int              `json:"count"`
		Offset   int              `json:"offset"`
		Prev     *string          `json:"prev"`
		Next     *string          `json:"next"`
	}

	successTests := []struct {
		name     string
		req      request.GetProblemsRequest
		respData respData
		isAdmin  bool
	}{
		{
			name: "All",
			req: request.GetProblemsRequest{
				Search: "test_get_problems",
				Limit:  0,
				Offset: 0,
			},
			respData: respData{
				Problems: []models.Problem{
					problem1,
					problem2,
					problem3,
				},
				Total:  3,
				Count:  3,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
			isAdmin: false,
		},
		{
			name: "AllWithAdminPermission",
			req: request.GetProblemsRequest{
				Search: "test_get_problems",
				Limit:  0,
				Offset: 0,
			},
			respData: respData{
				Problems: []models.Problem{
					problem1,
					problem2,
					problem3,
					problem4,
				},
				Total:  4,
				Count:  4,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
			isAdmin: true,
		},
		{
			name: "NonExist",
			req: request.GetProblemsRequest{
				Search: "test_get_problems_non_exist",
				Limit:  0,
				Offset: 0,
			},
			respData: respData{
				Problems: []models.Problem{},
				Total:    0,
				Count:    0,
				Offset:   0,
				Prev:     nil,
				Next:     nil,
			},
			isAdmin: false,
		},
		{
			name: "Search",
			req: request.GetProblemsRequest{
				Search: "test_get_problems_2",
				Limit:  0,
				Offset: 0,
			},
			respData: respData{
				Problems: []models.Problem{
					problem2,
				},
				Total:  1,
				Count:  1,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
			isAdmin: false,
		},
		{
			name: "Paginator",
			req: request.GetProblemsRequest{
				Search: "test_get_problems",
				Limit:  2,
				Offset: 0,
			},
			respData: respData{
				Problems: []models.Problem{
					problem1,
					problem2,
				},
				Total:  3,
				Count:  2,
				Offset: 0,
				Prev:   nil,
				Next: getUrlStringPointer("problem.getProblems", map[string]string{
					"limit":  "2",
					"offset": "2",
				}),
			},
			isAdmin: false,
		},
	}

	t.Run("testGetProblemsSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("testGetProblems"+test.name, func(t *testing.T) {
				t.Parallel()
				user := createUserForTest(t, "get_problems", i)
				// assert.False(t,user.Can("manage_problem"))
				if test.isAdmin {
					user.GrantRole("admin")
				}
				httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("problem.getProblems"), test.req, headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				}))
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)

				if test.isAdmin {
					resp := response.AdminGetProblemsResponse{}
					mustJsonDecode(httpResp, &resp)
					expectResp := response.AdminGetProblemsResponse{
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
							Problems: resource.GetProblemForAdminSlice(test.respData.Problems),
							Total:    test.respData.Total,
							Count:    test.respData.Count,
							Offset:   test.respData.Offset,
							Prev:     test.respData.Prev,
							Next:     test.respData.Next,
						},
					}
					assert.Equal(t, expectResp, resp)
				} else {
					resp := response.GetProblemsResponse{}
					mustJsonDecode(httpResp, &resp)
					expectResp := response.GetProblemsResponse{
						Message: "SUCCESS",
						Error:   nil,
						Data: struct {
							Problems []resource.Problem `json:"problems"`
							Total    int                `json:"total"`
							Count    int                `json:"count"`
							Offset   int                `json:"offset"`
							Prev     *string            `json:"prev"`
							Next     *string            `json:"next"`
						}{
							Problems: resource.GetProblemSlice(test.respData.Problems),
							Total:    test.respData.Total,
							Count:    test.respData.Count,
							Offset:   test.respData.Offset,
							Prev:     test.respData.Prev,
							Next:     test.respData.Next,
						},
					}
					assert.Equal(t, expectResp, resp)
				}

			})
		}
	})
}

func TestGetProblemAttachmentFile(t *testing.T) {
	problemWithoutAttachmentFile := models.Problem{
		Name:               "test_get_problem_attachment_file_0",
		Description:        "test_get_problem_attachment_file_0_desc",
		AttachmentFileName: "",
		Public:             true,
		Privacy:            false,
		MemoryLimit:        1024,
		TimeLimit:          1000,
		LanguageAllowed:    "test_get_problem_attachment_file_0_language_allowed",
		CompileEnvironment: "test_get_problem_attachment_file_0_compile_environment",
		CompareScriptID:    1,
	}
	// publicFalseProblem means a problem which "public" field is false
	publicFalseProblem := models.Problem{
		Name:               "test_get_problem_attachment_file_1",
		Description:        "test_get_problem_attachment_file_1_desc",
		AttachmentFileName: "",
		Public:             false,
		Privacy:            false,
		MemoryLimit:        1024,
		TimeLimit:          1000,
		LanguageAllowed:    "test_get_problem_attachment_file_1_language_allowed",
		CompileEnvironment: "test_get_problem_attachment_file_1_compile_environment",
		CompareScriptID:    1,
	}
	assert.Nil(t, base.DB.Create(&problemWithoutAttachmentFile).Error)
	assert.Nil(t, base.DB.Create(&publicFalseProblem).Error)

	failTests := []failTest{
		{
			name:   "NonExistingProblem",
			method: "GET",
			path:   base.Echo.Reverse("problem.getProblemAttachmentFile", -1),
			req:    nil,
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("PROBLEM_NOT_FOUND", nil),
		},
		{
			name:   "ProblemWithoutAttachmentFile",
			method: "GET",
			path:   base.Echo.Reverse("problem.getProblemAttachmentFile", problemWithoutAttachmentFile.ID),
			req:    nil,
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			name:   "PublicFalse",
			method: "GET",
			path:   base.Echo.Reverse("problem.getProblemAttachmentFile", publicFalseProblem.ID),
			req:    nil,
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("PROBLEM_NOT_FOUND", nil),
		},
	}

	runFailTests(t, failTests, "GetProblemAttachmentFile")

	successTests := []struct {
		name                   string
		file                   *fileContent
		respContentDisposition string
	}{
		{
			name:                   "PDFFile",
			file:                   newFileContent("", "test_get_problem_attachment.pdf", "cGRmIGNvbnRlbnQK"),
			respContentDisposition: `inline; filename="test_get_problem_attachment.pdf"`,
		},
		{
			name:                   "NonPDFFile",
			file:                   newFileContent("", "test_get_problem_attachment.txt", "dHh0IGNvbnRlbnQK"),
			respContentDisposition: `attachment; filename="test_get_problem_attachment.txt"`,
		},
	}

	t.Run("testGetProblemAttachmentFileSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("testGetProblemAttachmentFile"+test.name, func(t *testing.T) {
				t.Parallel()
				problem, _ := createProblemForTest(t, "test_get_problem_attachment_file", i+2, test.file)
				httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("problem.getProblemAttachmentFile", problem.ID), nil, applyNormalUser))
				assert.Equal(t, test.respContentDisposition, httpResp.Header.Get("Content-Disposition"))
				assert.Equal(t, "public; max-age=31536000", httpResp.Header.Get("Cache-Control"))
				respBytes, err := ioutil.ReadAll(httpResp.Body)
				assert.Nil(t, err)
				fileBytes, err := ioutil.ReadAll(test.file.reader)
				assert.Equal(t, fileBytes, respBytes)
			})
		}
	})

}

func TestCreateProblem(t *testing.T) {
	t.Parallel()
	FailTests := []failTest{
		{
			name:   "WithoutParams",
			method: "POST",
			path:   base.Echo.Reverse("problem.createProblem"),
			req: request.CreateProblemRequest{
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
					"translation": "评测脚本为必填字段",
				},
			}),
		},
		{
			name:   "PermissionDenied",
			method: "POST",
			path:   base.Echo.Reverse("problem.createProblem"),
			req: request.CreateProblemRequest{
				Name: "test_create_problem_perm",
			},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, FailTests, "CreateProblem")

	boolTrue := true
	boolFalse := false

	successTests := []struct {
		name       string
		req        request.CreateProblemRequest
		attachment *fileContent
	}{
		{
			name: "SuccessWithoutAttachment",
			req: request.CreateProblemRequest{
				Name:            "test_create_problem_1",
				Description:     "test_create_problem_1_desc",
				MemoryLimit:     4294967296,
				TimeLimit:       1000,
				LanguageAllowed: "test_create_problem_1_language_allowed",
				CompareScriptID: 1,
				Public:          &boolFalse,
				Privacy:         &boolTrue,
			},
			attachment: nil,
		},
		{
			name: "SuccessWithAttachment",
			req: request.CreateProblemRequest{
				Name:            "test_create_problem_2",
				Description:     "test_create_problem_2_desc",
				MemoryLimit:     4294967296,
				TimeLimit:       1000,
				LanguageAllowed: "test_create_problem_2_language_allowed",
				CompareScriptID: 2,
				Public:          &boolTrue,
				Privacy:         &boolFalse,
			},
			attachment: newFileContent("attachment_file", "test_create_problem_attachment_file", attachmentFileBase64),
		},
	}
	t.Run("TestCreateProblemSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("TestCreateProblem"+test.name, func(t *testing.T) {
				t.Parallel()
				user := createUserForTest(t, "create_problem", i)
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
				httpReq := makeReq(t, "POST", base.Echo.Reverse("problem.createProblem"), data, headerOption{
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
				resp := response.CreateProblemResponse{}
				mustJsonDecode(httpResp, &resp)
				jsonEQ(t, response.UpdateProblemResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data: struct {
						*resource.ProblemForAdmin `json:"problem"`
					}{
						resource.GetProblemForAdmin(&databaseProblem),
					},
				}, resp)
				assert.True(t, user.HasRole("problem_creator", databaseProblem))
				if test.attachment != nil {
					storageContent := getObjectContent(t, "problems", fmt.Sprintf("%d/attachment", databaseProblem.ID))
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

func TestUpdateProblem(t *testing.T) {
	t.Parallel()
	boolTrue := true
	boolFalse := false
	problem1 := models.Problem{
		Name:               "test_update_problem_1",
		AttachmentFileName: "test_update_problem_1_attachment_file_name",
		LanguageAllowed:    "test_update_problem_1_language_allowed",
	}
	assert.Nil(t, base.DB.Create(&problem1).Error)

	userWithProblem1Perm := models.User{
		Username: "test_update_problem_user_p1",
		Nickname: "test_update_problem_user_p1_nick",
		Email:    "test_update_problem_user_p1@e.e",
		Password: utils.HashPassword("test_update_problem_user_p1"),
	}
	assert.Nil(t, base.DB.Create(&userWithProblem1Perm).Error)
	userWithProblem1Perm.GrantRole("problem_creator", problem1)

	FailTests := []failTest{
		{
			name:   "WithoutParams",
			method: "PUT",
			path:   base.Echo.Reverse("problem.updateProblem", problem1.ID),
			req: request.UpdateProblemRequest{
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
					"translation": "评测脚本为必填字段",
				},
			}),
		},
		{
			name:   "NonExistId",
			method: "PUT",
			path:   base.Echo.Reverse("problem.updateProblem", -1),
			req: request.UpdateProblemRequest{
				Name:               "test_update_problem_non_exist",
				Description:        "test_update_problem_non_exist_desc",
				Public:             &boolFalse,
				Privacy:            &boolFalse,
				MemoryLimit:        1024,
				TimeLimit:          1000,
				LanguageAllowed:    "test_update_problem_non_exist_language_allowed",
				CompileEnvironment: "test_update_problem_non_exist_compile_environment",
				CompareScriptID:    1,
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			name:   "PermissionDenied",
			method: "PUT",
			path:   base.Echo.Reverse("problem.updateProblem", problem1.ID),
			req: request.UpdateProblemRequest{
				Name:            "test_update_problem_prem",
				LanguageAllowed: "test_update_problem_perm_language_allowed",
			},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, FailTests, "UpdateProblem")

	successTests := []struct {
		name               string
		path               string
		originalProblem    models.Problem
		expectedProblem    models.Problem
		req                request.UpdateProblemRequest
		updatedAttachment  *fileContent
		originalAttachment *fileContent
		testCases          []models.TestCase
	}{
		{
			name: "WithoutAttachmentAndTestCase",
			path: "id",
			originalProblem: models.Problem{
				Name:            "test_update_problem_3",
				Description:     "test_update_problem_3_desc",
				LanguageAllowed: "test_update_problem_3_language_allowed",
				Public:          false,
				Privacy:         true,
				MemoryLimit:     1024,
				TimeLimit:       1000,
				CompareScriptID: 1,
			},
			expectedProblem: models.Problem{
				Name:            "test_update_problem_30",
				Description:     "test_update_problem_30_desc",
				LanguageAllowed: "test_update_problem_30_language_allowed",
				Public:          true,
				Privacy:         false,
				MemoryLimit:     2048,
				TimeLimit:       2000,
				CompareScriptID: 2,
			},
			req: request.UpdateProblemRequest{
				Name:            "test_update_problem_30",
				Description:     "test_update_problem_30_desc",
				LanguageAllowed: "test_update_problem_30_language_allowed",
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
				Name:               "test_update_problem_4",
				Description:        "test_update_problem_4_desc",
				LanguageAllowed:    "test_update_problem_4_language_allowed",
				Public:             true,
				Privacy:            true,
				MemoryLimit:        1024,
				TimeLimit:          1000,
				CompareScriptID:    1,
				AttachmentFileName: "",
			},
			expectedProblem: models.Problem{
				Name:               "test_update_problem_40",
				Description:        "test_update_problem_40_desc",
				LanguageAllowed:    "test_update_problem_40_language_allowed",
				Public:             false,
				Privacy:            false,
				MemoryLimit:        2048,
				TimeLimit:          2000,
				CompareScriptID:    2,
				AttachmentFileName: "test_update_problem_attachment_40",
			},
			req: request.UpdateProblemRequest{
				Name:            "test_update_problem_40",
				Description:     "test_update_problem_40_desc",
				LanguageAllowed: "test_update_problem_40_language_allowed",
				Public:          &boolFalse,
				Privacy:         &boolFalse,
				MemoryLimit:     2048,
				TimeLimit:       2000,
				CompareScriptID: 2,
			},
			updatedAttachment: newFileContent("attachment_file", "test_update_problem_attachment_40", newAttachmentFileBase64),
			testCases:         nil,
		},
		{
			name: "WithChangingAttachment",
			path: "id",
			originalProblem: models.Problem{
				Name:               "test_update_problem_5",
				Description:        "test_update_problem_5_desc",
				LanguageAllowed:    "test_update_problem_5_language_allowed",
				Public:             true,
				Privacy:            true,
				MemoryLimit:        1024,
				TimeLimit:          1000,
				CompareScriptID:    1,
				AttachmentFileName: "test_update_problem_attachment_5",
			},
			expectedProblem: models.Problem{
				Name:               "test_update_problem_50",
				Description:        "test_update_problem_50_desc",
				LanguageAllowed:    "test_update_problem_50_language_allowed",
				Public:             false,
				Privacy:            false,
				MemoryLimit:        2048,
				TimeLimit:          2000,
				CompareScriptID:    2,
				AttachmentFileName: "test_update_problem_attachment_50",
			},
			req: request.UpdateProblemRequest{
				Name:            "test_update_problem_50",
				Description:     "test_update_problem_50_desc",
				LanguageAllowed: "test_update_problem_50_language_allowed",
				Public:          &boolFalse,
				Privacy:         &boolFalse,
				MemoryLimit:     2048,
				TimeLimit:       2000,
				CompareScriptID: 2,
			},
			originalAttachment: newFileContent("attachment_file", "test_update_problem_attachment_5", attachmentFileBase64),
			updatedAttachment:  newFileContent("attachment_file", "test_update_problem_attachment_50", newAttachmentFileBase64),
			testCases:          nil,
		},
		{
			name: "WithoutChangingAttachment",
			path: "id",
			originalProblem: models.Problem{
				Name:               "test_update_problem_6",
				Description:        "test_update_problem_6_desc",
				LanguageAllowed:    "test_update_problem_6_language_allowed",
				Public:             true,
				Privacy:            true,
				MemoryLimit:        1024,
				TimeLimit:          1000,
				CompareScriptID:    1,
				AttachmentFileName: "test_update_problem_attachment_6",
			},
			expectedProblem: models.Problem{
				Name:               "test_update_problem_60",
				Description:        "test_update_problem_60_desc",
				LanguageAllowed:    "test_update_problem_60_language_allowed",
				Public:             false,
				Privacy:            false,
				MemoryLimit:        2048,
				TimeLimit:          2000,
				CompareScriptID:    2,
				AttachmentFileName: "test_update_problem_attachment_6",
			},
			req: request.UpdateProblemRequest{
				Name:            "test_update_problem_60",
				Description:     "test_update_problem_60_desc",
				LanguageAllowed: "test_update_problem_60_language_allowed",
				Public:          &boolFalse,
				Privacy:         &boolFalse,
				MemoryLimit:     2048,
				TimeLimit:       2000,
				CompareScriptID: 2,
			},
			originalAttachment: newFileContent("attachment_file", "test_update_problem_attachment_6", attachmentFileBase64),
			updatedAttachment:  nil,
			testCases:          nil,
		},
		{
			name: "WithTestCase",
			path: "id",
			originalProblem: models.Problem{
				Name:            "test_update_problem_7",
				Description:     "test_update_problem_7_desc",
				LanguageAllowed: "test_update_problem_7_language_allowed",
				Public:          false,
				Privacy:         true,
				MemoryLimit:     1024,
				TimeLimit:       1000,
				CompareScriptID: 1,
			},
			expectedProblem: models.Problem{
				Name:            "test_update_problem_70",
				Description:     "test_update_problem_70_desc",
				LanguageAllowed: "test_update_problem_70_language_allowed",
				Public:          true,
				Privacy:         false,
				MemoryLimit:     2048,
				TimeLimit:       2000,
				CompareScriptID: 2,
			},
			req: request.UpdateProblemRequest{
				Name:            "test_update_problem_70",
				Description:     "test_update_problem_70_desc",
				LanguageAllowed: "test_update_problem_70_language_allowed",
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
					InputFileName:  "test_update_problem_7_test_case_1_input_file_name",
					OutputFileName: "test_update_problem_7_test_case_1_output_file_name",
				},
				{
					Score:          100,
					InputFileName:  "test_update_problem_7_test_case_2_input_file_name",
					OutputFileName: "test_update_problem_7_test_case_2_output_file_name",
				},
			},
		},
	}

	t.Run("TestUpdateProblemSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("TestUpdateProblem"+test.name, func(t *testing.T) {
				t.Parallel()
				assert.Nil(t, base.DB.Create(&test.originalProblem).Error)
				for j := range test.testCases {
					assert.Nil(t, base.DB.Model(&test.originalProblem).Association("TestCases").Append(&test.testCases[j]))
				}
				if test.originalAttachment != nil {
					b, err := ioutil.ReadAll(test.originalAttachment.reader)
					assert.Nil(t, err)
					_, err = base.Storage.PutObject("problems", fmt.Sprintf("%d/attachment", test.originalProblem.ID), bytes.NewReader(b), int64(len(b)), minio.PutObjectOptions{})
					assert.Nil(t, err)
					test.originalAttachment.reader = bytes.NewReader(b)
				}
				path := base.Echo.Reverse("problem.updateProblem", test.originalProblem.ID)
				user := createUserForTest(t, "update_problem", i)
				user.GrantRole("problem_creator", test.originalProblem)
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
				assert.Nil(t, base.DB.Set("gorm:auto_preload", true).Model(databaseProblem).Association("TestCases").Find(&databaseProblem.TestCases))
				if test.testCases != nil {
					jsonEQ(t, test.testCases, databaseProblem.TestCases)
				} else {
					assert.Equal(t, []models.TestCase{}, databaseProblem.TestCases)
				}
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				jsonEQ(t, response.UpdateProblemResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data: struct {
						*resource.ProblemForAdmin `json:"problem"`
					}{
						resource.GetProblemForAdmin(&databaseProblem),
					},
				}, httpResp)
				if test.updatedAttachment != nil || test.originalAttachment != nil {
					storageContent := getObjectContent(t, "problems", fmt.Sprintf("%d/attachment", databaseProblem.ID))
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

func TestDeleteProblem(t *testing.T) {
	t.Parallel()

	failTests := []failTest{
		{
			name:   "NonExistId",
			method: "DELETE",
			path:   base.Echo.Reverse("problem.deleteProblem", -1),
			req:    request.DeleteProblemRequest{},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			name:   "PermissionDenied",
			method: "DELETE",
			path:   base.Echo.Reverse("problem.deleteProblem", -1),
			req:    request.DeleteProblemRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "DeleteProblem")

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
				Name:               "test_delete_problem_1",
				AttachmentFileName: "",
				LanguageAllowed:    "test_delete_problem_1_language_allowed",
			},
			originalAttachment: nil,
			testCases:          nil,
		},
		{
			name: "SuccessWithAttachment",
			problem: models.Problem{
				Name:               "test_delete_problem_2",
				AttachmentFileName: "test_delete_problem_attachment_2",
				LanguageAllowed:    "test_delete_problem_2_language_allowed",
			},
			originalAttachment: newFileContent("attachment_file", "test_delete_problem_attachment_2", attachmentFileBase64),
			testCases:          nil,
		},
		{
			name: "SuccessWithTestCases",
			problem: models.Problem{
				Name:               "test_delete_problem_3",
				AttachmentFileName: "",
				LanguageAllowed:    "test_delete_problem_3_language_allowed",
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
						InputFileName:  "test_delete_problem_3_test_case_1_input_file_name",
						OutputFileName: "test_delete_problem_3_test_case_1_output_file_name",
					},
					inputFile:  newFileContent("input_file", "test_delete_problem_3_test_case_1.in", inputTextBase64),
					outputFile: newFileContent("output_file", "test_delete_problem_3_test_case_1.out", outputTextBase64),
				},
				{
					testcase: models.TestCase{
						Score:          100,
						InputFileName:  "test_delete_problem_3_test_case_2_input_file_name",
						OutputFileName: "test_delete_problem_3_test_case_2_output_file_name",
					},
					inputFile:  newFileContent("input_file", "test_delete_problem_3_test_case_2.in", inputTextBase64),
					outputFile: newFileContent("output_file", "test_delete_problem_3_test_case_2.out", outputTextBase64),
				},
			},
		},
		{
			name: "SuccessWithAttachmentAndTestCases",
			problem: models.Problem{
				Name:               "test_delete_problem_4",
				AttachmentFileName: "test_delete_problem_attachment_4",
				LanguageAllowed:    "test_delete_problem_4_language_allowed",
			},
			originalAttachment: newFileContent("attachment_file", "test_delete_problem_attachment_4", attachmentFileBase64),
			testCases: []struct {
				testcase   models.TestCase
				inputFile  *fileContent
				outputFile *fileContent
			}{
				{
					testcase: models.TestCase{
						Score:          100,
						InputFileName:  "test_delete_problem_4_test_case_1_input_file_name",
						OutputFileName: "test_delete_problem_4_test_case_1_output_file_name",
					},
					inputFile:  newFileContent("input_file", "test_delete_problem_4_test_case_1.in", inputTextBase64),
					outputFile: newFileContent("output_file", "test_delete_problem_4_test_case_1.out", outputTextBase64),
				},
				{
					testcase: models.TestCase{
						Score:          100,
						InputFileName:  "test_delete_problem_4_test_case_2_input_file_name",
						OutputFileName: "test_delete_problem_4_test_case_2_output_file_name",
					},
					inputFile:  newFileContent("input_file", "test_delete_problem_4_test_case_2.in", inputTextBase64),
					outputFile: newFileContent("output_file", "test_delete_problem_4_test_case_2.out", outputTextBase64),
				},
			},
		},
	}

	t.Run("TestDeleteProblemSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("TestDeleteProblem"+test.name, func(t *testing.T) {
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
				user := createUserForTest(t, "delete_problem", i)
				user.GrantRole("problem_creator", test.problem)
				httpResp := makeResp(makeReq(t, "DELETE", base.Echo.Reverse("problem.deleteProblem", test.problem.ID), request.DeleteProblemRequest{}, headerOption{
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
				assert.False(t, user.HasRole("problem_creator", test.problem))
				for j := range test.testCases {
					assert.Equal(t, gorm.ErrRecordNotFound, base.DB.First(models.TestCase{}, test.testCases[j].testcase.ID).Error)
				}
			})
		}
	})
}

func TestCreateTestCase(t *testing.T) {
	problem, user := createProblemForTest(t, "create_test_case", 0, nil)

	failTests := []failTest{
		{
			name:   "NonExistingProblem",
			method: "POST",
			path:   base.Echo.Reverse("problem.createTestCase", -1),
			req: addFieldContentSlice([]reqContent{
				newFileContent("input_file", "test_create_test_case_non_existing_problem.in", inputTextBase64),
				newFileContent("output_file", "test_create_test_case_non_existing_problem.out", outputTextBase64),
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
			path:   base.Echo.Reverse("problem.createTestCase", problem.ID),
			req: addFieldContentSlice([]reqContent{
				newFileContent("output_file", "test_create_test_case_lack_input_file.out", outputTextBase64),
			}, map[string]string{
				"score": "100",
			}),
			reqOptions: []reqOption{
				headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				},
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("INVALID_FILE", nil),
		},
		{
			name:   "LackOutputFile",
			method: "POST",
			path:   base.Echo.Reverse("problem.createTestCase", problem.ID),
			req: addFieldContentSlice([]reqContent{
				newFileContent("input_file", "test_create_test_case_lack_output_file.in", inputTextBase64),
			}, map[string]string{
				"score": "100",
			}),
			reqOptions: []reqOption{
				headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				},
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("INVALID_FILE", nil),
		},
		{
			name:   "LackBothFile",
			method: "POST",
			path:   base.Echo.Reverse("problem.createTestCase", problem.ID),
			req: addFieldContentSlice([]reqContent{}, map[string]string{
				"score": "100",
			}),
			reqOptions: []reqOption{
				headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				},
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("INVALID_FILE", nil),
		},
		{
			name:   "PermissionDenied",
			method: "POST",
			path:   base.Echo.Reverse("problem.createTestCase", problem.ID),
			req: addFieldContentSlice([]reqContent{
				newFileContent("input_file", "test_create_test_case_permission_denied.in", inputTextBase64),
				newFileContent("output_file", "test_create_test_case_permission_denied.out", outputTextBase64),
			}, map[string]string{
				"score": "100",
			}),
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "CreateTestCase")

	t.Run("TestCreateTestCaseSuccess", func(t *testing.T) {
		t.Parallel()
		req := makeReq(t, "POST", base.Echo.Reverse("problem.createTestCase", problem.ID), addFieldContentSlice([]reqContent{
			newFileContent("input_file", "test_create_test_case_success.in", inputTextBase64),
			newFileContent("output_file", "test_create_test_case_success.out", outputTextBase64),
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
			InputFileName:  "test_create_test_case_success.in",
			OutputFileName: "test_create_test_case_success.out",
		}
		databaseTestCase := models.TestCase{}
		assert.Nil(t, base.DB.Where("problem_id = ? and input_file_name = ?", problem.ID, "test_create_test_case_success.in").First(&databaseTestCase).Error)
		assert.Equal(t, expectedTestCase.ProblemID, databaseTestCase.ProblemID)
		assert.Equal(t, expectedTestCase.InputFileName, databaseTestCase.InputFileName)
		assert.Equal(t, expectedTestCase.OutputFileName, databaseTestCase.OutputFileName)
		assert.Equal(t, expectedTestCase.Score, databaseTestCase.Score)
		resp := response.CreateTestCaseResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusCreated, httpResp.StatusCode)
		assert.Equal(t, response.CreateTestCaseResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.TestCaseForAdmin `json:"test_case"`
			}{
				resource.GetTestCaseForAdmin(&databaseTestCase),
			},
		}, resp)
		assert.Equal(t, "input text\n", string(getObjectContent(t, "problems", fmt.Sprintf("%d/input/%d.in", problem.ID, databaseTestCase.ID))))
		assert.Equal(t, "output text\n", string(getObjectContent(t, "problems", fmt.Sprintf("%d/output/%d.out", problem.ID, databaseTestCase.ID))))
	})
}

func TestGetTestCaseInputFile(t *testing.T) {
	problem, user := createProblemForTest(t, "get_test_case_input_file", 0, nil)
	failTests := []failTest{
		{
			name:   "NonExistingProblem",
			method: "GET",
			path:   base.Echo.Reverse("problem.getTestCaseInputFile", -1, 1),
			req:    request.GetTestCaseInputFileRequest{},
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
			path:   base.Echo.Reverse("problem.getTestCaseInputFile", problem.ID, -1),
			req:    request.GetTestCaseInputFileRequest{},
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
			path:   base.Echo.Reverse("problem.getTestCaseInputFile", problem.ID, 1),
			req:    request.GetTestCaseInputFileRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "GetTestCaseInputFile")

	testCase := createTestCaseForTest(t, problem, 51,
		newFileContent("", "test_get_test_case_input_file_success.in", inputTextBase64),
		nil,
	)

	req := makeReq(t, "GET", base.Echo.Reverse("problem.getTestCaseInputFile", problem.ID, testCase.ID), request.GetTestCaseInputFileRequest{}, headerOption{
		"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
	})
	httpResp := makeResp(req)

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "input text\n", string(respBytes))
}

func TestGetTestCaseOutputFile(t *testing.T) {
	problem, user := createProblemForTest(t, "get_test_case_output_file", 0, nil)
	failTests := []failTest{
		{
			name:   "NonExistingProblem",
			method: "GET",
			path:   base.Echo.Reverse("problem.getTestCaseOutputFile", -1, 1),
			req:    request.GetTestCaseOutputFileRequest{},
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
			path:   base.Echo.Reverse("problem.getTestCaseOutputFile", problem.ID, -1),
			req:    request.GetTestCaseOutputFileRequest{},
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
			path:   base.Echo.Reverse("problem.getTestCaseOutputFile", problem.ID, 1),
			req:    request.GetTestCaseOutputFileRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "GetTestCaseOutputFile")

	testCase := createTestCaseForTest(t, problem, 52,
		nil,
		newFileContent("", "test_get_test_case_output_file_success.out", outputTextBase64),
	)

	req := makeReq(t, "GET", base.Echo.Reverse("problem.getTestCaseOutputFile", problem.ID, testCase.ID), request.GetTestCaseOutputFileRequest{}, headerOption{
		"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
	})
	httpResp := makeResp(req)

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "output text\n", string(respBytes))
}

func TestUpdateTestCase(t *testing.T) {
	problem, user := createProblemForTest(t, "update_test_case", 0, nil)

	failTests := []failTest{
		{
			name:   "NonExistingProblem",
			method: "PUT",
			path:   base.Echo.Reverse("problem.updateTestCase", -1, 1),
			req: request.UpdateTestCaseRequest{
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
			path:   base.Echo.Reverse("problem.updateTestCase", problem.ID, -1),
			req: request.UpdateTestCaseRequest{
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
			path:   base.Echo.Reverse("problem.updateTestCase", problem.ID, 1),
			req: request.UpdateTestCaseRequest{
				Score: 100,
			},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "UpdateTestCase")

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
			originalInputFile:  newFileContent("input_file", "test_update_test_case_1.in", inputTextBase64),
			originalOutputFile: newFileContent("output_file", "test_update_test_case_1.out", outputTextBase64),
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
			originalInputFile:  newFileContent("input_file", "test_update_test_case_2.in", inputTextBase64),
			originalOutputFile: newFileContent("output_file", "test_update_test_case_2.out", outputTextBase64),
			updatedInputFile:   newFileContent("input_file", "test_update_test_case_20.in", newInputTextBase64),
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
			originalInputFile:  newFileContent("input_file", "test_update_test_case_3.in", inputTextBase64),
			originalOutputFile: newFileContent("output_file", "test_update_test_case_3.out", outputTextBase64),
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
			originalInputFile:  newFileContent("input_file", "test_update_test_case_4.in", inputTextBase64),
			originalOutputFile: newFileContent("output_file", "test_update_test_case_4.out", outputTextBase64),
			updatedInputFile:   newFileContent("input_file", "test_update_test_case_40.in", newInputTextBase64),
			updatedOutputFile:  newFileContent("output_file", "test_update_test_case_40.out", newOutputTextBase64),
			expectedTestCase: models.TestCase{
				ProblemID:      problem.ID,
				Score:          100,
				InputFileName:  "test_update_test_case_40.in",
				OutputFileName: "test_update_test_case_40.out",
			},
		},
	}

	t.Run("TestUpdateTestCaseSuccess", func(t *testing.T) {
		t.Parallel()
		for _, test := range successTests {
			test := test
			t.Run("TestUpdateTestCase"+test.name, func(t *testing.T) {
				t.Parallel()
				testCase := createTestCaseForTest(t, problem, test.originalScore, test.originalInputFile, test.originalOutputFile)
				var reqContentSlice []reqContent
				if test.updatedInputFile != nil {
					reqContentSlice = append(reqContentSlice, test.updatedInputFile)
				}
				if test.updatedOutputFile != nil {
					reqContentSlice = append(reqContentSlice, test.updatedOutputFile)
				}
				req := makeReq(t, "PUT", base.Echo.Reverse("problem.updateTestCase", problem.ID, testCase.ID), addFieldContentSlice(
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

				var expectedInputFileReader io.Reader
				if test.updatedInputFile != nil {
					expectedInputFileReader = test.updatedInputFile.reader
				} else {
					expectedInputFileReader = test.originalInputFile.reader
				}
				expectedInputContent, err := ioutil.ReadAll(expectedInputFileReader)
				assert.Nil(t, err)
				assert.Equal(t, expectedInputContent, getObjectContent(t, "problems", fmt.Sprintf("%d/input/%d.in", problem.ID, databaseTestCase.ID)))

				var expectedOutputFileReader io.Reader
				if test.updatedOutputFile != nil {
					expectedOutputFileReader = test.updatedOutputFile.reader
				} else {
					expectedOutputFileReader = test.originalOutputFile.reader
				}
				expectedOutputContent, err := ioutil.ReadAll(expectedOutputFileReader)
				assert.Nil(t, err)
				assert.Equal(t, expectedOutputContent, getObjectContent(t, "problems", fmt.Sprintf("%d/output/%d.out", problem.ID, databaseTestCase.ID)))

				resp := response.UpdateTestCaseResponse{}
				mustJsonDecode(httpResp, &resp)
				assert.Equal(t, response.UpdateTestCaseResponse{
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

func TestDeleteTestCase(t *testing.T) {
	problem, user := createProblemForTest(t, "delete_test_case", 0, nil)

	failTests := []failTest{
		{
			name:   "NonExistingProblem",
			method: "DELETE",
			path:   base.Echo.Reverse("problem.deleteTestCases", -1, 1),
			req:    request.DeleteProblemRequest{},
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
			path:   base.Echo.Reverse("problem.deleteTestCase", problem.ID, -1),
			req:    request.DeleteProblemRequest{},
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
			path:   base.Echo.Reverse("problem.deleteTestCases", problem.ID, 1),
			req:    request.DeleteProblemRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "DeleteTestCase")

	t.Run("TestDeleteTestCaseSuccess", func(t *testing.T) {
		t.Parallel()
		testCase := createTestCaseForTest(t, problem, 72,
			newFileContent("input_file", "test_delete_test_case_0.in", inputTextBase64),
			newFileContent("output_file", "test_delete_test_case_0.out", outputTextBase64),
		)

		req := makeReq(t, "DELETE", base.Echo.Reverse("problem.deleteTestCase", problem.ID, testCase.ID), request.DeleteTestCaseRequest{}, headerOption{
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

func TestDeleteTestCases(t *testing.T) {
	t.Parallel()
	problem, user := createProblemForTest(t, "delete_test_cases", 0, nil)

	failTests := []failTest{
		{
			name:   "NonExistingProblem",
			method: "DELETE",
			path:   base.Echo.Reverse("problem.deleteTestCases", -1),
			req:    request.DeleteTestCasesRequest{},
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
			path:   base.Echo.Reverse("problem.deleteTestCases", problem.ID),
			req:    request.DeleteTestCasesRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "DeleteTestCases")

	t.Run("TestDeleteTestCasesSuccess", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 5; i++ {
			createTestCaseForTest(t, problem, 0,
				newFileContent("input_file", fmt.Sprintf("test_delete_test_cases_%d.in", i), inputTextBase64),
				newFileContent("output_file", fmt.Sprintf("test_delete_test_cases_%d.out", i), outputTextBase64),
			)
		}
		req := makeReq(t, "DELETE", base.Echo.Reverse("problem.deleteTestCases", problem.ID), request.DeleteTestCasesRequest{}, headerOption{
			"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
		})
		httpResp := makeResp(req)
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		var databaseTestCases []models.TestCase
		assert.Nil(t, base.DB.Find(&databaseTestCases, "problem_id = ?", problem.ID).Error)
		assert.Equal(t, 0, len(databaseTestCases))
	})
}
