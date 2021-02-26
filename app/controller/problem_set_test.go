package controller_test

import (
	"fmt"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func TestCreateProblemSet(t *testing.T) {
	t.Parallel()

	class := createClassForTest(t, "create_problem_set_permission_denied", 0, nil, nil)

	failTests := []failTest{
		{
			name:       "WithoutParams",
			method:     "POST",
			path:       base.Echo.Reverse("problemSet.createProblemSet", class.ID),
			req:        request.CreateProblemSetRequest{},
			reqOptions: []reqOption{applyAdminUser},
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
					"translation": "描述为必填字段",
				},
				map[string]interface{}{
					"field":       "StartAt",
					"reason":      "required",
					"translation": "开始时间为必填字段",
				},
				map[string]interface{}{
					"field":       "EndAt",
					"reason":      "required",
					"translation": "结束时间为必填字段",
				},
			}),
		},
		{
			name:   "NonExistingClass",
			method: "POST",
			path:   base.Echo.Reverse("problemSet.createProblemSet", -1),
			req: request.CreateProblemSetRequest{
				Name:        "test_create_problem_set_non_existing_class_name",
				Description: "test_create_problem_set_non_existing_class_description",
				StartAt:     hashStringToTime("test_create_problem_set_non_existing_class_start_at"),
				EndAt:       hashStringToTime("test_create_problem_set_non_existing_class_end_at"),
			},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			name:   "PermissionDenied",
			method: "POST",
			path:   base.Echo.Reverse("problemSet.createProblemSet", class.ID),
			req: request.CreateProblemSetRequest{
				Name:        "test_create_problem_set_permission_denied_name",
				Description: "test_create_problem_set_permission_denied_description",
				StartAt:     hashStringToTime("test_create_problem_set_permission_denied_start_at"),
				EndAt:       hashStringToTime("test_create_problem_set_permission_denied_end_at"),
			},
			reqOptions: []reqOption{applyNormalUser},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "")

	t.Run("Success", func(t *testing.T) {
		user := createUserForTest(t, "create_problem_set_success", 0)
		class := createClassForTest(t, "create_problem_set_success", 0, nil, nil)
		user.GrantRole("class_creator", class)

		httpResp := makeResp(makeReq(t, "POST", base.Echo.Reverse("problemSet.createProblemSet", class.ID), request.CreateProblemSetRequest{
			Name:        "test_create_problem_set_success_name",
			Description: "test_create_problem_set_success_description",
			StartAt:     hashStringToTime("test_create_problem_set_success_start_at"),
			EndAt:       hashStringToTime("test_create_problem_set_success_end_at"),
		}, applyUser(user)))
		assert.Equal(t, http.StatusCreated, httpResp.StatusCode)

		databaseProblemSet := models.ProblemSet{}
		assert.NoError(t, base.DB.Preload("Problems").Preload("Scores").First(&databaseProblemSet, "name = ?", "test_create_problem_set_success_name").Error)
		expectedProblemSet := models.ProblemSet{
			ID:          databaseProblemSet.ID,
			ClassID:     class.ID,
			Name:        "test_create_problem_set_success_name",
			Description: "test_create_problem_set_success_description",
			Problems:    []models.Problem{},
			Scores:      []models.Grade{},
			StartAt:     hashStringToTime("test_create_problem_set_success_start_at"),
			EndAt:       hashStringToTime("test_create_problem_set_success_end_at"),
			CreatedAt:   databaseProblemSet.CreatedAt,
			UpdatedAt:   databaseProblemSet.UpdatedAt,
			DeletedAt:   gorm.DeletedAt{},
		}
		assert.Equal(t, expectedProblemSet, databaseProblemSet)

		resp := response.CreateProblemSetResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, response.CreateProblemSetResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.ProblemSetDetail `json:"problem_set"`
			}{
				resource.GetProblemSetDetail(&expectedProblemSet),
			},
		}, resp)
	})
}

func createProblemSetForTest(t *testing.T, name string, id int, class *models.Class, problems []models.Problem) *models.ProblemSet {
	problemSet := models.ProblemSet{
		Name:        fmt.Sprintf("test_%s_%d_name", name, id),
		Description: fmt.Sprintf("test_%s_%d_description", name, id),
		Problems:    []models.Problem{},
		Scores:      []models.Grade{},
		StartAt:     hashStringToTime(fmt.Sprintf("test_%s_%d_start_at", name, id)),
		EndAt:       hashStringToTime(fmt.Sprintf("test_%s_%d_end_at", name, id)),
	}
	assert.NoError(t, base.DB.Model(&class).Association("ProblemSets").Append(&problemSet))
	assert.NoError(t, base.DB.Model(&problemSet).Association("Problems").Append(problems))
	return &problemSet
}

func TestCloneProblemSet(t *testing.T) {
	t.Parallel()

	class1 := createClassForTest(t, "clone_problem_set_fail", 1, nil, nil)
	class2 := createClassForTest(t, "clone_problem_set_fail", 2, nil, nil)
	problemSet1 := createProblemSetForTest(t, "clone_problem_set_fail", 1, &class1, nil)
	problemSet2 := createProblemSetForTest(t, "clone_problem_set_fail", 2, &class2, nil)

	failTests := []failTest{
		{
			name:       "WithoutParams",
			method:     "POST",
			path:       base.Echo.Reverse("problemSet.cloneProblemSet", class1.ID),
			req:        request.CloneProblemSetRequest{},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusBadRequest,
			resp: response.ErrorResp("VALIDATION_ERROR", []interface{}{
				map[string]interface{}{
					"field":       "SourceClassID",
					"reason":      "required",
					"translation": "复制源班级ID为必填字段",
				},
				map[string]interface{}{
					"field":       "SourceProblemSetID",
					"reason":      "required",
					"translation": "复制源题目组ID为必填字段",
				},
			}),
		},
		{
			name:   "NonExistingClass",
			method: "POST",
			path:   base.Echo.Reverse("problemSet.cloneProblemSet", -1),
			req: request.CloneProblemSetRequest{
				SourceClassID:      class1.ID,
				SourceProblemSetID: problemSet1.ID,
			},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			name:   "NonExistingSourceClass",
			method: "POST",
			path:   base.Echo.Reverse("problemSet.cloneProblemSet", class1.ID),
			req: request.CloneProblemSetRequest{
				SourceClassID:      9999999, // non-existing class id
				SourceProblemSetID: problemSet1.ID,
			},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("SOURCE_CLASS_NOT_FOUND", nil),
		},
		{
			name:   "NonExistingSourceProblemSet",
			method: "POST",
			path:   base.Echo.Reverse("problemSet.cloneProblemSet", class1.ID),
			req: request.CloneProblemSetRequest{
				SourceClassID:      class1.ID,
				SourceProblemSetID: problemSet2.ID,
			},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("SOURCE_PROBLEM_SET_NOT_FOUND", nil),
		},
		{
			name:   "PermissionDenied",
			method: "POST",
			path:   base.Echo.Reverse("problemSet.cloneProblemSet", class1.ID),
			req: request.CloneProblemSetRequest{
				SourceClassID:      class1.ID,
				SourceProblemSetID: problemSet1.ID,
			},
			reqOptions: []reqOption{applyNormalUser},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "")

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		user := createUserForTest(t, "clone_problem_set_success", 0)
		class := createClassForTest(t, "clone_problem_set_success", 0, nil, nil)
		user.GrantRole("class_creator", class)

		sourceClass := createClassForTest(t, "clone_problem_set_success_source", 0, nil, nil)
		problem1 := createProblemForTest(t, "clone_problem_set_success_source", 1, nil, user)
		problem2 := createProblemForTest(t, "clone_problem_set_success_source", 2, nil, user)
		sourceProblemSet := createProblemSetForTest(t, "clone_problem_set_success_source", 0, &sourceClass, []models.Problem{problem1, problem2})
		assert.NoError(t, sourceProblemSet.UpdateGrade(models.Submission{
			UserID:    user.ID,
			ProblemID: problem1.ID,
			Score:     10,
		}))
		assert.NoError(t, sourceProblemSet.UpdateGrade(models.Submission{
			UserID:    user.ID,
			ProblemID: problem2.ID,
			Score:     20,
		}))

		httpResp := makeResp(makeReq(t, "POST", base.Echo.Reverse("problemSet.cloneProblemSet", class.ID), request.CloneProblemSetRequest{
			SourceClassID:      sourceClass.ID,
			SourceProblemSetID: sourceProblemSet.ID,
		}, applyUser(user)))
		assert.Equal(t, http.StatusCreated, httpResp.StatusCode)

		databaseProblemSet := models.ProblemSet{}
		assert.NoError(t, base.DB.Preload("Problems").Preload("Scores").
			First(&databaseProblemSet, "name = ? and class_id = ?", "test_clone_problem_set_success_source_0_name", class.ID).Error)
		expectedProblemSet := models.ProblemSet{
			ID:          databaseProblemSet.ID,
			ClassID:     class.ID,
			Name:        "test_clone_problem_set_success_source_0_name",
			Description: "test_clone_problem_set_success_source_0_description",
			Problems: []models.Problem{
				problem1,
				problem2,
			},
			Scores:    []models.Grade{},
			StartAt:   hashStringToTime("test_clone_problem_set_success_source_0_start_at"),
			EndAt:     hashStringToTime("test_clone_problem_set_success_source_0_end_at"),
			CreatedAt: databaseProblemSet.CreatedAt,
			UpdatedAt: databaseProblemSet.UpdatedAt,
			DeletedAt: gorm.DeletedAt{},
		}
		assert.Equal(t, expectedProblemSet, databaseProblemSet)

		resp := response.CloneProblemSetResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, response.CloneProblemSetResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.ProblemSetDetail `json:"problem_set"`
			}{
				resource.GetProblemSetDetail(&expectedProblemSet),
			},
		}, resp)
	})
}
