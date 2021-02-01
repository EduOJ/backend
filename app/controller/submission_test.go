package controller_test

import (
	"fmt"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/minio/minio-go"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"
)

const (
	normalUser = iota
	problemCreator
	adminUser
	submitter
)

func createSubmissionForTest(t *testing.T, name string, id int, problem *models.Problem, user *models.User, code *fileContent, testCaseCount int) (submission models.Submission) {
	for i := 0; i < testCaseCount; i++ {
		createTestCaseForTest(t, *problem, testCaseData{
			Score:      uint(i),
			Sample:     i%3 == 0,
			InputFile:  newFileContent("input", "input_file", b64Encode(fmt.Sprintf("problem_%d_test_case_%d_input", problem.ID, i))),
			OutputFile: newFileContent("output", "output_file", b64Encode(fmt.Sprintf("problem_%d_test_case_%d_output", problem.ID, i))),
		})
	}
	submission = models.Submission{
		UserID:       user.ID,
		ProblemID:    problem.ID,
		ProblemSetId: 0,
		Language:     fmt.Sprintf("test_%s_language_allowed_%d", name, id),
		FileName:     fmt.Sprintf("test_%s_code_file_name_%d", name, id),
		Priority:     models.PriorityDefault,
		Judged:       false,
		Score:        0,
		Status:       "PENDING",
		Runs:         make([]models.Run, len(problem.TestCases)),
	}
	for i, testCase := range problem.TestCases {
		submission.Runs[i] = models.Run{
			UserID:             user.ID,
			ProblemID:          problem.ID,
			ProblemSetId:       0,
			TestCaseID:         testCase.ID,
			Sample:             testCase.Sample,
			Priority:           models.PriorityDefault,
			Judged:             false,
			Status:             "PENDING",
			MemoryUsed:         0,
			TimeUsed:           0,
			OutputStrippedHash: "",
		}
	}
	assert.Nil(t, base.DB.Create(&submission).Error)
	if code != nil {
		_, err := base.Storage.PutObject("submissions", fmt.Sprintf("%d/code", submission.ID), code.reader, code.size, minio.PutObjectOptions{})
		assert.Nil(t, err)
		_, err = code.reader.Seek(0, io.SeekStart)
		assert.Nil(t, err)
	}
	return
}

func TestCreateSubmission(t *testing.T) {
	// publicFalseProblem means a problem which "public" field is false
	publicFalseProblem, _ := createProblemForTest(t, "test_create_submission_public_false", 0, nil)
	assert.Nil(t, base.DB.Model(&publicFalseProblem).Update("public", false).Error)
	assert.Nil(t, base.DB.Model(&publicFalseProblem).Update("language_allowed", "test_language,golang").Error)
	failTests := []failTest{
		{
			// testCreateSubmissionNonExistingProblem
			name:   "NonExistingProblem",
			method: "POST",
			path:   base.Echo.Reverse("submission.createSubmission", -1),
			req: addFieldContentSlice([]reqContent{
				newFileContent("code", "code_file_name", b64Encode("test code content")),
			}, map[string]string{"language": "test_language"}),
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			// testCreateSubmissionPublicFalseProblem
			name:   "PublicFalseProblem",
			method: "POST",
			path:   base.Echo.Reverse("submission.createSubmission", publicFalseProblem.ID),
			req: addFieldContentSlice([]reqContent{
				newFileContent("code", "code_file_name", b64Encode("test code content")),
			}, map[string]string{"language": "test_language"}),
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
		{
			// testCreateSubmissionWithoutCode
			name:   "WithoutCode",
			method: "POST",
			path:   base.Echo.Reverse("submission.createSubmission", publicFalseProblem.ID),
			req: request.CreateSubmissionRequest{
				Language: "test_language",
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("INVALID_FILE", nil),
		},
		{
			// testCreateSubmissionInvalidLanguage
			name:   "InvalidLanguage",
			method: "POST",
			path:   base.Echo.Reverse("submission.createSubmission", publicFalseProblem.ID),
			req: addFieldContentSlice([]reqContent{
				newFileContent("code", "code_file_name", b64Encode("test code content")),
			}, map[string]string{"language": "invalid_language"}),
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("INVALID_LANGUAGE", nil),
		},
	}

	// testCreateSubmissionFail
	runFailTests(t, failTests, "CreateSubmission")

	successfulTests := []struct {
		name          string
		testCaseCount int
		problemPublic bool
		requestUser   int // 0->normalUser / 1->problemCreator / 2->adminUser
		response      resource.SubmissionDetail
	}{
		// testCreateSubmissionWithoutTestCases
		{
			name:          "WithoutTestCases",
			testCaseCount: 0,
			problemPublic: true,
			requestUser:   normalUser,
		},
		// testCreateSubmissionPublicProblem
		{
			name:          "PublicProblem",
			testCaseCount: 1,
			problemPublic: true,
			requestUser:   normalUser,
		},
		// testCreateSubmissionCreator
		{
			name:          "Creator",
			testCaseCount: 2,
			problemPublic: true,
			requestUser:   problemCreator,
		},
		// testCreateSubmissionAdmin
		{
			name:          "Admin",
			testCaseCount: 5,
			problemPublic: true,
			requestUser:   adminUser,
		},
	}
	t.Run("testCreateSubmissionSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successfulTests {
			i := i
			test := test
			t.Run("testCreateSubmission"+test.name, func(t *testing.T) {
				t.Parallel()
				problem, creator := createProblemForTest(t, "test_create_submission", i, nil)
				assert.Nil(t, base.DB.Model(&problem).Update("language_allowed", "test_language,golang").Error)
				for j := 0; j < test.testCaseCount; j++ {
					createTestCaseForTest(t, problem, testCaseData{
						Score:  0,
						Sample: true,
						InputFile: newFileContent("input", "input_file",
							b64Encode(fmt.Sprintf("test_create_submission_%d_test_case_%d_input_content", i, j))),
						OutputFile: newFileContent("output", "output_file",
							b64Encode(fmt.Sprintf("test_create_submission_%d_test_case_%d_output_content", i, j))),
					})
				}
				problem.LoadTestCases()
				var applyUser reqOption
				switch test.requestUser {
				case normalUser:
					applyUser = applyNormalUser
				case problemCreator:
					applyUser = headerOption{
						"Set-User-For-Test": {fmt.Sprintf("%d", creator.ID)},
					}
				case adminUser:
					applyUser = applyAdminUser
				default:
					t.Fail()
				}
				req := makeReq(t, "POST", base.Echo.Reverse("submission.createSubmission", problem.ID),
					addFieldContentSlice([]reqContent{
						newFileContent("code", "code_file_name", b64Encode(fmt.Sprintf("test_create_submission_%d_code", i))),
					}, map[string]string{
						"language": "test_language",
					}), applyUser)
				httpResp := makeResp(req)
				resp := response.CreateSubmissionResponse{}
				mustJsonDecode(httpResp, &resp)
				responseSubmission := *resp.Data.SubmissionDetail
				databaseSubmission := models.Submission{}
				reqUserID, err := strconv.ParseUint(req.Header.Get("Set-User-For-Test"), 10, 64)
				assert.Nil(t, err)
				assert.Nil(t, base.DB.Preload("Runs").First(&databaseSubmission, "problem_id = ? and user_id = ?", problem.ID, reqUserID).Error)
				databaseSubmissionDetail := resource.GetSubmissionDetail(&databaseSubmission)
				databaseRunData := map[uint]struct {
					ID        uint
					CreatedAt time.Time
				}{}
				for _, run := range databaseSubmission.Runs {
					databaseRunData[run.TestCaseID] = struct {
						ID        uint
						CreatedAt time.Time
					}{
						ID:        run.ID,
						CreatedAt: run.CreatedAt,
					}
				}
				expectedRunSlice := make([]resource.Run, test.testCaseCount)
				for i, testCase := range problem.TestCases {
					expectedRunSlice[i] = resource.Run{
						ID:           databaseRunData[testCase.ID].ID,
						UserID:       uint(reqUserID),
						ProblemID:    problem.ID,
						ProblemSetId: 0,
						TestCaseID:   testCase.ID,
						Sample:       testCase.Sample,
						SubmissionID: databaseSubmission.ID,
						Priority:     127,
						Judged:       false,
						Status:       "PENDING",
						MemoryUsed:   0,
						TimeUsed:     0,
						CreatedAt:    databaseRunData[testCase.ID].CreatedAt,
					}
				}
				expectedSubmission := resource.SubmissionDetail{
					ID:           databaseSubmissionDetail.ID,
					UserID:       uint(reqUserID),
					ProblemID:    problem.ID,
					ProblemSetId: 0,
					Language:     "test_language",
					FileName:     "code_file_name",
					Priority:     127,
					Judged:       false,
					Score:        0,
					Status:       "PENDING",
					Runs:         expectedRunSlice,
					CreatedAt:    databaseSubmission.CreatedAt,
				}
				assert.Equal(t, &expectedSubmission, databaseSubmissionDetail)
				assert.Equal(t, expectedSubmission, responseSubmission)

				storageContent := string(getObjectContent(t, "submissions", fmt.Sprintf("%d/code", databaseSubmissionDetail.ID)))
				expectedContent := fmt.Sprintf("test_create_submission_%d_code", i)
				assert.Equal(t, expectedContent, storageContent)
			})
		}

	})
}

func TestGetSubmission(t *testing.T) {

	failTests := []failTest{
		{
			// testGetSubmissionNormalUserNonExisting
			name:   "NormalUserNonExisting",
			method: "GET",
			path:   base.Echo.Reverse("submission.getSubmission", -1),
			req:    request.GetSubmissionRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
		{
			// testGetSubmissionAdminUserNonExisting
			name:   "AdminUserNonExisting",
			method: "GET",
			path:   base.Echo.Reverse("submission.getSubmission", -1),
			req:    request.GetSubmissionRequest{},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			// testGetSubmissionPermissionDenied
			name:   "PermissionDenied",
			method: "GET",
			path:   base.Echo.Reverse("submission.getSubmission", -1),
			req:    request.GetSubmissionRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	// testGetSubmissionFail
	runFailTests(t, failTests, "GetSubmission")

	successTests := []struct {
		name          string
		code          *fileContent
		testCaseCount int
		requestUser   uint
	}{
		{
			// testGetSubmissionWithoutTestCases
			name:          "WithoutTestCases",
			code:          newFileContent("code", "code_file_name", b64Encode("test_get_submission_code_0")),
			testCaseCount: 0,
			requestUser:   adminUser,
		},
		{
			// testGetSubmissionAdminUser
			name:          "AdminUser",
			code:          newFileContent("code", "code_file_name", b64Encode("test_get_submission_code_1")),
			testCaseCount: 2,
			requestUser:   adminUser,
		},
		{
			// testGetSubmissionSubmitter
			name:          "Submitter",
			code:          newFileContent("code", "code_file_name", b64Encode("test_get_submission_code_2")),
			testCaseCount: 2,
			requestUser:   submitter,
		},
	}

	t.Run("testGetSubmissionSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			t.Run("testGetSubmission"+test.name, func(t *testing.T) {
				t.Parallel()
				problem, user := createProblemForTest(t, "get_submission", i, nil)
				submission := createSubmissionForTest(t, "get_submission", i, &problem, &user, test.code, test.testCaseCount)
				var applyUser reqOption
				switch test.requestUser {
				case adminUser:
					applyUser = applyAdminUser
				case submitter:
					applyUser = headerOption{
						"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
					}
				default:
					t.Fail()
				}
				httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("submission.getSubmission", submission.ID),
					request.GetSubmissionRequest{}, applyUser))
				resp := response.GetSubmissionResponse{}
				mustJsonDecode(httpResp, &resp)
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				expectedSubmissionDetail := resource.GetSubmissionDetail(&submission)
				expectedSubmissionDetail.CreatedAt.UTC()
				assert.Equal(t, response.GetSubmissionResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data: struct {
						*resource.SubmissionDetail `json:"submission"`
					}{
						expectedSubmissionDetail,
					},
				}, resp)
			})
		}
	})
}
