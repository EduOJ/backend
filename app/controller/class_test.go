package controller_test

import (
	"fmt"
	"github.com/EduOJ/backend/app/request"
	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/app/response/resource"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/utils"
	"github.com/EduOJ/backend/database/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"testing"
)

func checkInviteCode(t *testing.T, code string) {
	assert.Regexp(t, regexp.MustCompile("^[a-zA-Z2-9]{5}$"), code)
	var count int64
	assert.NoError(t, base.DB.Model(models.Class{}).Where("invite_code = ?", code).Count(&count).Error)
	assert.Equal(t, int64(1), count)
}

func createClassForTest(t *testing.T, name string, id int, managers, students []*models.User) models.Class {
	inviteCode := utils.GenerateInviteCode()
	class := models.Class{
		Name:        fmt.Sprintf("test_%s_%d_name", name, id),
		CourseName:  fmt.Sprintf("test_%s_%d_course_name", name, id),
		Description: fmt.Sprintf("test_%s_%d_description", name, id),
		InviteCode:  inviteCode,
		Managers:    managers,
		Students:    students,
	}
	assert.NoError(t, base.DB.Create(&class).Error)
	return class
}

func TestCreateClass(t *testing.T) {
	t.Parallel()

	class := createClassForTest(t, "test_create_class_permission_denied", 0, nil, nil)
	failTests := []failTest{
		{
			// testCreateClassWithoutParams
			name:       "WithoutParams",
			method:     "POST",
			path:       base.Echo.Reverse("class.createClass", -1),
			req:        request.CreateClassRequest{},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusBadRequest,
			resp: response.ErrorResp("VALIDATION_ERROR", []interface{}{
				map[string]interface{}{
					"field":       "Name",
					"reason":      "required",
					"translation": "名称为必填字段",
				},
				map[string]interface{}{
					"field":       "CourseName",
					"reason":      "required",
					"translation": "课程名称为必填字段",
				},
				map[string]interface{}{
					"field":       "Description",
					"reason":      "required",
					"translation": "描述为必填字段",
				},
			}),
		},
		{
			// testCreateClassPermissionDenied
			name:   "PermissionDenied",
			method: "POST",
			path:   base.Echo.Reverse("class.createClass", class.ID),
			req: request.CreateClassRequest{
				Name:        "test_create_class_permission_denied_name",
				CourseName:  "test_create_class_permission_denied_course_name",
				Description: "test_create_class_permission_denied_description",
			},
			reqOptions: []reqOption{applyNormalUser},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "CreateClass")

	user := createUserForTest(t, "test_create_class", 1)
	user.GrantRole("admin")
	httpResp := makeResp(makeReq(t, "POST", base.Echo.Reverse("class.createClass"), request.CreateClassRequest{
		Name:        "test_create_class_1_name",
		CourseName:  "test_create_class_1_course_name",
		Description: "test_create_class_1_description",
	}, applyUser(user)))
	assert.Equal(t, http.StatusCreated, httpResp.StatusCode)

	databaseClass := models.Class{}
	assert.NoError(t, base.DB.Preload("Managers").Preload("Students").
		First(&databaseClass, "name = ? ", "test_create_class_1_name").Error)
	checkInviteCode(t, databaseClass.InviteCode)
	assert.True(t, user.HasRole("class_creator", databaseClass))
	databaseClass.Managers[0].LoadRoles()
	databaseUser := models.User{}
	assert.NoError(t, base.DB.First(&databaseUser, user.ID).Error)
	databaseUser.LoadRoles()
	expectedClass := models.Class{
		ID:          databaseClass.ID,
		Name:        "test_create_class_1_name",
		CourseName:  "test_create_class_1_course_name",
		Description: "test_create_class_1_description",
		InviteCode:  databaseClass.InviteCode,
		Managers: []*models.User{
			&databaseUser,
		},
		Students:  []*models.User{},
		CreatedAt: databaseClass.CreatedAt,
		UpdatedAt: databaseClass.UpdatedAt,
		DeletedAt: gorm.DeletedAt{},
	}
	assert.Equal(t, expectedClass, databaseClass)
	resp := response.CreateClassResponse{}
	mustJsonDecode(httpResp, &resp)
	assert.Equal(t, response.CreateClassResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ClassDetail `json:"class"`
		}{
			resource.GetClassDetail(&expectedClass),
		},
	}, resp)
}

func TestGetClass(t *testing.T) {
	t.Parallel()

	class := createClassForTest(t, "test_get_class_permission_denied", 0, nil, nil)
	failTests := []failTest{
		{
			// testGetClassNonExist
			name:       "NonExist",
			method:     "GET",
			path:       base.Echo.Reverse("class.getClass", -1),
			req:        request.GetClassRequest{},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			// testGetClassPermissionDenied
			name:       "PermissionDenied",
			method:     "GET",
			path:       base.Echo.Reverse("class.getClass", class.ID),
			req:        request.GetClassRequest{},
			reqOptions: []reqOption{applyNormalUser},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "GetClass")

	t.Run("Admin", func(t *testing.T) {
		t.Parallel()

		user := createUserForTest(t, "test_get_class_admin", 0)
		class := createClassForTest(t, "test_get_class_admin", 0, nil, nil)
		problem1 := createProblemForTest(t, "test_get_class_admin", 1, nil, user)
		problem2 := createProblemForTest(t, "test_get_class_admin", 2, nil, user)
		problemSet1 := createProblemSetForTest(t, "test_get_class_admin", 1, &class, []models.Problem{problem1})
		problemSet2 := createProblemSetForTest(t, "test_get_class_admin", 1, &class, []models.Problem{problem1, problem2})
		class.ProblemSets = []*models.ProblemSet{problemSet1, problemSet2}
		assert.NoError(t, base.DB.Save(&class).Error)
		user.GrantRole("class_creator", class)
		httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("class.getClass", class.ID),
			request.GetClassRequest{}, applyUser(user)))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		resp := response.GetClassResponseForAdmin{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, response.GetClassResponseForAdmin{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.ClassDetail `json:"class"`
			}{
				resource.GetClassDetail(&class),
			},
		}, resp)
	})

	t.Run("Student", func(t *testing.T) {
		t.Parallel()

		user := createUserForTest(t, "test_get_class_student", 0)
		class := createClassForTest(t, "test_get_class_student", 0, nil, []*models.User{&user})
		problem1 := createProblemForTest(t, "test_get_class_student", 1, nil, user)
		problem2 := createProblemForTest(t, "test_get_class_student", 2, nil, user)
		problemSet1 := createProblemSetForTest(t, "test_get_class_student", 1, &class, []models.Problem{problem1})
		problemSet2 := createProblemSetForTest(t, "test_get_class_student", 1, &class, []models.Problem{problem1, problem2})
		class.ProblemSets = []*models.ProblemSet{problemSet1, problemSet2}
		assert.NoError(t, base.DB.Save(&class).Error)
		httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("class.getClass", class.ID), request.GetClassRequest{}, applyUser(user)))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		resp := response.GetClassResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, response.GetClassResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.Class `json:"class"`
			}{
				resource.GetClass(&class),
			},
		}, resp)
	})
}

func TestUpdateClass(t *testing.T) {
	t.Parallel()

	class := createClassForTest(t, "test_update_class_permission_denied", 0, nil, nil)
	failTests := []failTest{
		{
			// testUpdateClassWithoutParams
			name:       "WithoutParams",
			method:     "PUT",
			path:       base.Echo.Reverse("class.updateClass", -1),
			req:        request.UpdateClassRequest{},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusBadRequest,
			resp: response.ErrorResp("VALIDATION_ERROR", []interface{}{
				map[string]interface{}{
					"field":       "Name",
					"reason":      "required",
					"translation": "名称为必填字段",
				},
				map[string]interface{}{
					"field":       "CourseName",
					"reason":      "required",
					"translation": "课程名称为必填字段",
				},
				map[string]interface{}{
					"field":       "Description",
					"reason":      "required",
					"translation": "描述为必填字段",
				},
			}),
		},
		{
			// testUpdateClassNonExist
			name:   "NonExist",
			method: "PUT",
			path:   base.Echo.Reverse("class.updateClass", -1),
			req: request.UpdateClassRequest{
				Name:        "test_update_class_non_exist_name",
				CourseName:  "test_update_class_non_exist_course_name",
				Description: "test_update_class_non_exist_description",
			},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			// testUpdateClassPermissionDenied
			name:   "PermissionDenied",
			method: "PUT",
			path:   base.Echo.Reverse("class.updateClass", class.ID),
			req: request.UpdateClassRequest{
				Name:        "test_update_class_permission_denied_name",
				CourseName:  "test_update_class_permission_denied_course_name",
				Description: "test_update_class_permission_denied_description",
			},
			reqOptions: []reqOption{applyNormalUser},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "UpdateClass")

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		class := createClassForTest(t, "update_class", 0, nil, nil)
		user := createUserForTest(t, "update_class", 0)
		user.GrantRole("class_creator", class)

		httpResp := makeResp(makeReq(t, "PUT", base.Echo.Reverse("class.updateClass", class.ID), request.UpdateClassRequest{
			Name:        "test_update_class_00_name",
			CourseName:  "test_update_class_00_course_name",
			Description: "test_update_class_00_description",
		}, applyUser(user)))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)

		databaseClass := models.Class{}
		assert.NoError(t, base.DB.Preload("Managers").Preload("Students").
			First(&databaseClass, class.ID).Error)
		checkInviteCode(t, databaseClass.InviteCode)
		expectedClass := models.Class{
			ID:          databaseClass.ID,
			Name:        "test_update_class_00_name",
			CourseName:  "test_update_class_00_course_name",
			Description: "test_update_class_00_description",
			InviteCode:  databaseClass.InviteCode,
			Managers:    []*models.User{},
			Students:    []*models.User{},
			CreatedAt:   databaseClass.CreatedAt,
			UpdatedAt:   databaseClass.UpdatedAt,
			DeletedAt:   gorm.DeletedAt{},
		}
		assert.Equal(t, expectedClass, databaseClass)
		resp := response.UpdateClassResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, response.UpdateClassResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.ClassDetail `json:"class"`
			}{
				resource.GetClassDetail(&expectedClass),
			},
		}, resp)
	})
}

func TestRefreshInviteCode(t *testing.T) {
	t.Parallel()

	class := createClassForTest(t, "test_refresh_invite_code_permission_denied", 0, nil, nil)
	failTests := []failTest{
		{
			// testUpdateClassNonExist
			name:       "NonExist",
			method:     "PUT",
			path:       base.Echo.Reverse("class.refreshInviteCode", -1),
			req:        request.RefreshInviteCodeRequest{},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			// testUpdateClassPermissionDenied
			name:       "PermissionDenied",
			method:     "PUT",
			path:       base.Echo.Reverse("class.refreshInviteCode", class.ID),
			req:        request.RefreshInviteCodeRequest{},
			reqOptions: []reqOption{applyNormalUser},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}
	runFailTests(t, failTests, "RefreshInviteCode")

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		class := createClassForTest(t, "refresh_invite_code", 0, nil, nil)
		originalInviteCode := class.InviteCode

		user := createUserForTest(t, "refresh_invite_code", 0)
		user.GrantRole("class_creator", class)

		httpResp := makeResp(makeReq(t, "PUT", base.Echo.Reverse("class.refreshInviteCode", class.ID),
			request.RefreshInviteCodeRequest{}, applyUser(user)))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)

		databaseClass := models.Class{}
		assert.NoError(t, base.DB.Preload("Managers").Preload("Students").
			First(&databaseClass, class.ID).Error)
		checkInviteCode(t, databaseClass.InviteCode)
		assert.NotEqual(t, originalInviteCode, databaseClass.InviteCode)
		expectedClass := models.Class{
			ID:          databaseClass.ID,
			Name:        "test_refresh_invite_code_0_name",
			CourseName:  "test_refresh_invite_code_0_course_name",
			Description: "test_refresh_invite_code_0_description",
			InviteCode:  databaseClass.InviteCode,
			Managers:    []*models.User{},
			Students:    []*models.User{},
			CreatedAt:   databaseClass.CreatedAt,
			UpdatedAt:   databaseClass.UpdatedAt,
			DeletedAt:   gorm.DeletedAt{},
		}
		assert.Equal(t, expectedClass, databaseClass)
		resp := response.RefreshInviteCodeResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, response.RefreshInviteCodeResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.ClassDetail `json:"class"`
			}{
				resource.GetClassDetail(&expectedClass),
			},
		}, resp)
	})
}

func TestAddAndDeleteStudents(t *testing.T) {
	t.Parallel()

	class := createClassForTest(t, "test_add_and_delete_student_permission_denied", 0, nil, nil)
	addStudentsFailTests := []failTest{
		{
			// testAddStudentsWithoutParams
			name:       "WithoutParams",
			method:     "POST",
			path:       base.Echo.Reverse("class.addStudents", -1),
			req:        request.AddStudentsRequest{},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusBadRequest,
			resp: response.ErrorResp("VALIDATION_ERROR", []interface{}{
				map[string]interface{}{
					"field":       "UserIds",
					"reason":      "required",
					"translation": "用户ID数组为必填字段",
				},
			}),
		},
		{
			// testAddStudentsNonExist
			name:   "NonExist",
			method: "POST",
			path:   base.Echo.Reverse("class.addStudents", -1),
			req: request.AddStudentsRequest{
				UserIds: []uint{0},
			},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			// testAddStudentsPermissionDenied
			name:   "PermissionDenied",
			method: "POST",
			path:   base.Echo.Reverse("class.addStudents", class.ID),
			req: request.AddStudentsRequest{
				UserIds: []uint{0},
			},
			reqOptions: []reqOption{applyNormalUser},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}
	runFailTests(t, addStudentsFailTests, "AddStudents")

	deleteStudentsFailTests := []failTest{
		{
			// testDeleteStudentsWithoutParams
			name:       "WithoutParams",
			method:     "DELETE",
			path:       base.Echo.Reverse("class.deleteStudents", -1),
			req:        request.DeleteStudentsRequest{},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusBadRequest,
			resp: response.ErrorResp("VALIDATION_ERROR", []interface{}{
				map[string]interface{}{
					"field":       "UserIds",
					"reason":      "required",
					"translation": "用户ID数组为必填字段",
				},
			}),
		},
		{
			// testDeleteStudentsNonExist
			name:   "NonExist",
			method: "DELETE",
			path:   base.Echo.Reverse("class.deleteStudents", -1),
			req: request.DeleteStudentsRequest{
				UserIds: []uint{0},
			},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			// testDeleteStudentsPermissionDenied
			name:   "PermissionDenied",
			method: "DELETE",
			path:   base.Echo.Reverse("class.deleteStudents", class.ID),
			req: request.DeleteStudentsRequest{
				UserIds: []uint{0},
			},
			reqOptions: []reqOption{applyNormalUser},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}
	runFailTests(t, deleteStudentsFailTests, "DeleteStudents")

	user1 := createUserForTest(t, "test_add_and_delete_students_user", 1)
	user2 := createUserForTest(t, "test_add_and_delete_students_user", 2)
	user3 := createUserForTest(t, "test_add_and_delete_students_user", 3)
	user4 := createUserForTest(t, "test_add_and_delete_students_user", 4)

	t.Run("AddStudentSuccess", func(t *testing.T) {
		t.Parallel()

		class := createClassForTest(t, "add_students_success", 0, nil, []*models.User{
			&user1,
			&user2,
		})
		user := createUserForTest(t, "add_students_success", 0)
		user.GrantRole("class_creator", class)
		httpResp := makeResp(makeReq(t, "POST", base.Echo.Reverse("class.addStudents", class.ID),
			request.AddStudentsRequest{
				UserIds: []uint{
					user2.ID,
					user3.ID,
					user4.ID,
				},
			}, applyUser(user)))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		databaseClass := models.Class{}
		assert.NoError(t, base.DB.Preload("Managers").Preload("Students").
			First(&databaseClass, class.ID).Error)
		expectedClass := models.Class{
			ID:          databaseClass.ID,
			Name:        "test_add_students_success_0_name",
			CourseName:  "test_add_students_success_0_course_name",
			Description: "test_add_students_success_0_description",
			InviteCode:  databaseClass.InviteCode,
			Managers:    []*models.User{},
			Students: []*models.User{
				&user1,
				&user2,
				&user3,
				&user4,
			},
			CreatedAt: databaseClass.CreatedAt,
			UpdatedAt: databaseClass.UpdatedAt,
			DeletedAt: gorm.DeletedAt{},
		}
		assert.Equal(t, expectedClass, databaseClass)
		resp := response.AddStudentsResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, response.AddStudentsResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.ClassDetail `json:"class"`
			}{
				resource.GetClassDetail(&expectedClass),
			},
		}, resp)
	})
	t.Run("DeleteStudentSuccess", func(t *testing.T) {
		t.Parallel()

		class := createClassForTest(t, "delete_students_success", 0, nil, []*models.User{
			&user1,
			&user2,
			&user3,
		})
		user := createUserForTest(t, "delete_students_success", 0)
		user.GrantRole("class_creator", class)
		httpResp := makeResp(makeReq(t, "DELETE", base.Echo.Reverse("class.deleteStudents", class.ID),
			request.DeleteStudentsRequest{
				UserIds: []uint{
					user2.ID,
					user3.ID,
					user4.ID,
				},
			}, applyUser(user)))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		databaseClass := models.Class{}
		assert.NoError(t, base.DB.Preload("Managers").Preload("Students").
			First(&databaseClass, class.ID).Error)
		expectedClass := models.Class{
			ID:          databaseClass.ID,
			Name:        "test_delete_students_success_0_name",
			CourseName:  "test_delete_students_success_0_course_name",
			Description: "test_delete_students_success_0_description",
			InviteCode:  databaseClass.InviteCode,
			Managers:    []*models.User{},
			Students: []*models.User{
				&user1,
			},
			CreatedAt: databaseClass.CreatedAt,
			UpdatedAt: databaseClass.UpdatedAt,
			DeletedAt: gorm.DeletedAt{},
		}
		assert.Equal(t, expectedClass, databaseClass)
		resp := response.DeleteStudentsResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, response.DeleteStudentsResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.ClassDetail `json:"class"`
			}{
				resource.GetClassDetail(&expectedClass),
			},
		}, resp)
	})
}

func TestJoinClass(t *testing.T) {
	t.Parallel()
	user := createUserForTest(t, "test_join_class_already_in_class", 0)
	class1 := createClassForTest(t, "test_join_class_wrong_invite_code", 1, nil, nil)
	class2 := createClassForTest(t, "test_join_class_already_in_class", 2, nil, []*models.User{&user})
	failTests := []failTest{
		{
			// testJoinClassWithoutParams
			name:       "WithoutParams",
			method:     "POST",
			path:       base.Echo.Reverse("class.joinClass", -1),
			req:        request.JoinClassRequest{},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusBadRequest,
			resp: response.ErrorResp("VALIDATION_ERROR", []interface{}{
				map[string]interface{}{
					"field":       "InviteCode",
					"reason":      "required",
					"translation": "邀请码为必填字段",
				},
			}),
		},
		{
			// testJoinClassNonExist
			name:   "NonExist",
			method: "POST",
			path:   base.Echo.Reverse("class.joinClass", -1),
			req: request.JoinClassRequest{
				InviteCode: utils.GenerateInviteCode(),
			},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			// testJoinClassWrongInviteCode
			name:   "WrongInviteCode",
			method: "POST",
			path:   base.Echo.Reverse("class.joinClass", class1.ID),
			req: request.JoinClassRequest{
				InviteCode: utils.GenerateInviteCode(),
			},
			reqOptions: []reqOption{applyNormalUser},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("WRONG_INVITE_CODE", nil),
		},
		{
			// testJoinClassAlreadyInClass
			name:   "AlreadyInClass",
			method: "POST",
			path:   base.Echo.Reverse("class.joinClass", class2.ID),
			req: request.JoinClassRequest{
				InviteCode: class2.InviteCode,
			},
			reqOptions: []reqOption{applyUser(user)},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("ALREADY_IN_CLASS", nil),
		},
	}

	runFailTests(t, failTests, "JoinClass")

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		class := createClassForTest(t, "join_class", 0, nil, nil)
		user := createUserForTest(t, "join_class", 0)
		user.GrantRole("class_creator", class)

		httpResp := makeResp(makeReq(t, "POST", base.Echo.Reverse("class.joinClass", class.ID), request.JoinClassRequest{
			InviteCode: class.InviteCode,
		}, applyUser(user)))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)

		databaseClass := models.Class{}
		assert.NoError(t, base.DB.Preload("Managers").Preload("Students").
			First(&databaseClass, class.ID).Error)
		checkInviteCode(t, databaseClass.InviteCode)
		user.LoadRoles()
		databaseClass.Students[0].LoadRoles()
		expectedClass := models.Class{
			ID:          databaseClass.ID,
			Name:        "test_join_class_0_name",
			CourseName:  "test_join_class_0_course_name",
			Description: "test_join_class_0_description",
			InviteCode:  databaseClass.InviteCode,
			Managers:    []*models.User{},
			Students: []*models.User{
				&user,
			},
			CreatedAt: databaseClass.CreatedAt,
			UpdatedAt: databaseClass.UpdatedAt,
			DeletedAt: gorm.DeletedAt{},
		}
		assert.Equal(t, expectedClass, databaseClass)
		resp := response.JoinClassResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, response.JoinClassResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.Class `json:"class"`
			}{
				resource.GetClass(&expectedClass),
			},
		}, resp)
	})
}

func TestDeleteClass(t *testing.T) {
	t.Parallel()

	class := createClassForTest(t, "test_delete_class_permission_denied", 0, nil, nil)
	failTests := []failTest{
		{
			// testDeleteClassNonExistingClass
			name:       "NonExistingClass",
			method:     "DELETE",
			path:       base.Echo.Reverse("class.deleteClass", -1),
			req:        request.DeleteClassRequest{},
			reqOptions: []reqOption{applyAdminUser},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			// testDeleteClassPermissionDenied
			name:       "PermissionDenied",
			method:     "DELETE",
			path:       base.Echo.Reverse("class.deleteClass", class.ID),
			req:        request.DeleteClassRequest{},
			reqOptions: []reqOption{applyNormalUser},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "DeleteClass")

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		user := createUserForTest(t, "delete_class", 0)
		class := createClassForTest(t, "delete_class", 0, []*models.User{&user}, []*models.User{&user})
		user.GrantRole("class_creator", class)
		createProblemSetForTest(t, "delete_class", 0, &class, nil)

		httpResp := makeResp(makeReq(t, "DELETE", base.Echo.Reverse("class.deleteClass", class.ID),
			request.DeleteClassRequest{}, applyUser(user)))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)

		databaseClass := models.Class{}
		err := base.DB.First(&databaseClass, class.ID).Error
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		databaseUser := models.User{}
		err = base.DB.First(&databaseUser, user.ID).Error
		assert.NoError(t, err)
		assert.Empty(t, databaseUser.ClassesTaking)
		assert.Empty(t, databaseUser.ClassesManaging)
		err = base.DB.First(&models.ProblemSet{}, "class_id = ?", class.ID).Error
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
