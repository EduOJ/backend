package controller_test

import (
	"fmt"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetProblem(t *testing.T) {
	t.Parallel()

	privateProblem := models.Problem{
		Name:               "test_get_problem_private",
		AttachmentFileName: "test_get_problem_private_attachment_file_name",
		LanguageAllowed:    "test_get_problem_private_language_allowed",
	}
	assert.Nil(t, base.DB.Create(&privateProblem).Error)

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
			name:   "Private",
			method: "GET",
			path:   base.Echo.Reverse("problem.getProblem", privateProblem.ID),
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
		name       string
		path       string
		req        request.GetProblemRequest
		problem    models.Problem
		role       models.Role
		roleTarget models.HasRole
	}{
		{
			name: "Success",
			path: "id",
			req:  request.GetProblemRequest{},
			problem: models.Problem{
				Name:               "test_get_problem_1",
				AttachmentFileName: "test_get_problem_1_attachment_file_name",
				LanguageAllowed:    "test_get_problem_1_language_allowed",
				Public:             true,
			},
			role:       models.Role{},
			roleTarget: nil,
		},
		// TODO: with test case
	}

	t.Run("testGetUserSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("testGetProblem"+test.name, func(t *testing.T) {
				t.Parallel()
				assert.Nil(t, base.DB.Create(&test.problem).Error)
				user := createUserForTest(t, "get_problem", i)
				user.GrantRole("creator", test.problem)
				httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("problem.getProblem", test.problem.ID), request.GetUserRequest{}, headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				}))
				resp := response.GetProblemResponse{}
				mustJsonDecode(httpResp, &resp)
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				assert.Equal(t, response.GetProblemResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data: struct {
						*resource.Problem `json:"problem"`
					}{
						resource.GetProblem(&test.problem),
					},
				}, resp)
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

	failTests := []failTest{
		{
			name:   "NotAllowedColumn",
			method: "GET",
			path:   base.Echo.Reverse("problem.getProblems"),
			req: request.GetProblemsRequest{
				Search:  "test_get_problems",
				OrderBy: "name.ASC",
			},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("INVALID_ORDER", nil),
		},
	}

	runFailTests(t, failTests, "GetProblems")

	successTests := []struct {
		name     string
		req      request.GetProblemsRequest
		respData respData
	}{
		{
			name: "All",
			req: request.GetProblemsRequest{
				Search:  "test_get_problems",
				Limit:   0,
				Offset:  0,
				OrderBy: "",
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
		},
		{
			name: "NonExist",
			req: request.GetProblemsRequest{
				Search:  "test_get_problems_non_exist",
				Limit:   0,
				Offset:  0,
				OrderBy: "",
			},
			respData: respData{
				Problems: []resource.Problem{},
				Total:    0,
				Count:    0,
				Offset:   0,
				Prev:     nil,
				Next:     nil,
			},
		},
		{
			name: "Search",
			req: request.GetProblemsRequest{
				Search:  "test_get_problems_2",
				Limit:   0,
				Offset:  0,
				OrderBy: "",
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
		},
		{
			name: "Sort",
			req: request.GetProblemsRequest{
				Search:  "test_get_problems",
				Limit:   0,
				Offset:  0,
				OrderBy: "id.DESC",
			},
			respData: respData{
				Problems: []resource.Problem{
					*resource.GetProblem(&problem3),
					*resource.GetProblem(&problem2),
					*resource.GetProblem(&problem1),
				},
				Total:  3,
				Count:  3,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "Paginator",
			req: request.GetProblemsRequest{
				Search:  "test_get_problems",
				Limit:   2,
				Offset:  0,
				OrderBy: "",
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
		},
	}

	t.Run("testGetProblemsSuccess", func(t *testing.T) {
		t.Parallel()
		for _, test := range successTests {
			test := test
			t.Run("testGetProblems"+test.name, func(t *testing.T) {
				t.Parallel()
				httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("problem.getProblems"), test.req, applyNormalUser))
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
