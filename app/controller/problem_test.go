package controller_test

import (
	"fmt"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

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
					assert.Nil(t, base.DB.Model(&test.problem).Association("TestCases").Append(&test.testCases[j]).Error)
				}
				user := createUserForTest(t, "get_problem", i)
				if test.isAdmin {
					user.GrantRole("creator", test.problem)
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
		Problems []resource.Problem `json:"problems"`
		Total    int                `json:"total"`
		Count    int                `json:"count"`
		Offset   int                `json:"offset"`
		Prev     *string            `json:"prev"`
		Next     *string            `json:"next"`
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
				Problems: []resource.Problem{
					*resource.GetProblem(&problem1),
					*resource.GetProblem(&problem2),
					*resource.GetProblem(&problem3),
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
				Problems: []resource.Problem{
					*resource.GetProblem(&problem1),
					*resource.GetProblem(&problem2),
					*resource.GetProblem(&problem3),
				},
				Total:  3,
				Count:  3,
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
				Problems: []resource.Problem{},
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
				Problems: []resource.Problem{
					*resource.GetProblem(&problem2),
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
				Problems: []resource.Problem{
					*resource.GetProblem(&problem1),
					*resource.GetProblem(&problem2),
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
				if test.isAdmin {
					user.GrantRole("creator")
				}
				httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("problem.getProblems"), test.req, headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				}))
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				resp := response.Response{}
				mustJsonDecode(httpResp, &resp)
				jsonEQ(t, response.GetProblemsResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data:    test.respData,
				}, resp)
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
