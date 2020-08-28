package controller_test

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func getUrlStringPointer(rawUrl string, paras map[string]string) *string {
	thisURL, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		panic(err)
	}
	q, err := url.ParseQuery(thisURL.RawQuery)
	if err != nil {
		panic(err)
	}
	for key := range paras {
		q.Add(key, paras[key])
	}
	thisURL.RawQuery = q.Encode()
	str := thisURL.String()
	return &str
}

func TestAdminCreateUser(t *testing.T) {
	t.Parallel()
	FailTests := []failTest{
		{
			name:   "WithoutParams",
			method: "POST",
			path:   "/api/admin/user",
			req: request.AdminCreateUserRequest{
				Username: "",
				Nickname: "",
				Email:    "",
				Password: "",
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusBadRequest,
			resp: response.ErrorResp("VALIDATION_ERROR", []interface{}{
				map[string]interface{}{
					"field":       "Username",
					"reason":      "required",
					"translation": "用户名为必填字段",
				},
				map[string]interface{}{
					"field":       "Nickname",
					"reason":      "required",
					"translation": "昵称为必填字段",
				},
				map[string]interface{}{
					"field":       "Email",
					"reason":      "required",
					"translation": "邮箱为必填字段",
				},
				map[string]interface{}{
					"field":       "Password",
					"reason":      "required",
					"translation": "密码为必填字段",
				},
			}),
		},
		{
			name:   "PermissionDenied",
			method: "POST",
			path:   "/api/admin/user",
			req: request.AdminCreateUserRequest{
				Username: "test_create_user_perm",
				Nickname: "test_create_user_perm",
				Email:    "test_create_user_perm@mail.com",
				Password: "test_create_user_perm",
			},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
		{
			name:   "ConflictEmail",
			method: "POST",
			path:   "/api/admin/user",
			req: request.AdminCreateUserRequest{
				Username: "test_create_user_1",
				Nickname: "test_create_user_1_nick",
				Email:    "test_create_user_conflict@mail.com",
				Password: "test_create_user_1_pwd",
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusConflict,
			resp:       response.ErrorResp("CONFLICT_EMAIL", nil),
		},
		{
			name:   "ConflictUsername",
			method: "POST",
			path:   "/api/admin/user",
			req: request.AdminCreateUserRequest{
				Username: "test_create_user_conflict",
				Nickname: "test_create_user_1_nick",
				Email:    "test_create_user_1@mail.com",
				Password: "test_create_user_1_pwd",
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusConflict,
			resp:       response.ErrorResp("CONFLICT_USERNAME", nil),
		},
	}

	dummyUserForConflict := models.User{
		Username: "test_create_user_conflict",
		Nickname: "test_create_user_conflict_nick",
		Email:    "test_create_user_conflict@mail.com",
		Password: utils.HashPassword("test_create_user_conflict_pwd"),
	}
	assert.Nil(t, base.DB.Create(&dummyUserForConflict).Error)
	runFailTests(t, FailTests, "AdminCreateUser")

	t.Run("testAdminCreateUserSuccess", func(t *testing.T) {
		t.Parallel()
		req := request.AdminCreateUserRequest{
			Username: "test_create_user_success_0",
			Nickname: "test_create_user_success_0",
			Email:    "test_create_user_success_0@mail.com",
			Password: "test_create_user_success_0",
		}
		httpResp := makeResp(makeReq(t, "POST", "/api/admin/user", req, applyAdminUser))
		assert.Equal(t, http.StatusCreated, httpResp.StatusCode)
		databaseUser := models.User{}
		assert.Nil(t, base.DB.Where("email = ?", req.Email).First(&databaseUser).Error)
		// request == database
		assert.Equal(t, req.Username, databaseUser.Username)
		assert.Equal(t, req.Nickname, databaseUser.Nickname)
		assert.Equal(t, req.Email, databaseUser.Email)
		assert.True(t, utils.VerifyPassword(req.Password, databaseUser.Password))
		// response == database
		jsonEQ(t, response.AdminUpdateUserResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*models.User `json:"user"`
			}{
				&databaseUser,
			},
		}, httpResp)
	})
}

func TestAdminUpdateUser(t *testing.T) {
	t.Parallel()
	user1 := models.User{
		Username: "test_update_user_1",
		Nickname: "test_update_user_1_rand_str",
		Email:    "test_update_user_1@mail.com",
		Password: utils.HashPassword("test_update_user_1_password"),
	}
	dummyUserForConflict := models.User{
		Username: "test_update_user_conflict",
		Nickname: "test_update_user_conflict_nick",
		Email:    "test_update_user_conflict@mail.com",
		Password: utils.HashPassword("test_update_user_conflict_pwd"),
	}
	assert.Nil(t, base.DB.Create(&user1).Error)
	assert.Nil(t, base.DB.Create(&dummyUserForConflict).Error)

	FailTests := []failTest{
		{
			name:   "WithoutParams",
			method: "PUT",
			path:   fmt.Sprintf("/api/admin/user/%d", user1.ID),
			req: request.AdminUpdateUserRequest{
				Username: "",
				Nickname: "",
				Email:    "",
				Password: "",
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusBadRequest,
			resp: response.ErrorResp("VALIDATION_ERROR", []interface{}{
				map[string]interface{}{
					"field":       "Username",
					"reason":      "required",
					"translation": "用户名为必填字段",
				},
				map[string]interface{}{
					"field":       "Nickname",
					"reason":      "required",
					"translation": "昵称为必填字段",
				},
				map[string]interface{}{
					"field":       "Email",
					"reason":      "required",
					"translation": "邮箱为必填字段",
				},
			}),
		},
		{
			name:   "NonExistId",
			method: "PUT",
			path:   "/api/admin/user/-1",
			req: request.AdminUpdateUserRequest{
				Username: "test_update_user_non_exist_1",
				Nickname: "test_update_user_non_exist_1_n",
				Email:    "test_update_user_non_exist_1@e.com",
				Password: "test_update_user_non_exist_1_p",
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			name:   "NonExistUsername",
			method: "PUT",
			path:   "/api/admin/user/test_put_non_existing_username",
			req: request.AdminUpdateUserRequest{
				Username: "test_update_user_non_exist_2",
				Nickname: "test_update_user_non_exist_2_n",
				Email:    "test_update_user_non_exist_2@e.com",
				Password: "test_update_user_non_exist_2_p",
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			name:   "ConflictEmail",
			method: "PUT",
			path:   fmt.Sprintf("/api/admin/user/%d", user1.ID),
			req: request.AdminUpdateUserRequest{
				Username: "test_update_user_1",
				Nickname: "test_update_user_1_nick",
				Email:    "test_update_user_conflict@mail.com",
				Password: "test_update_user_1_pwd",
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusConflict,
			resp:       response.ErrorResp("CONFLICT_EMAIL", nil),
		},
		{
			name:   "ConflictUsername",
			method: "PUT",
			path:   fmt.Sprintf("/api/admin/user/%d", user1.ID),
			req: request.AdminUpdateUserRequest{
				Username: "test_update_user_conflict",
				Nickname: "test_update_user_1_nick",
				Email:    "test_update_user_1@mail.com",
				Password: "test_update_user_1_pwd",
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusConflict,
			resp:       response.ErrorResp("CONFLICT_USERNAME", nil),
		},
		{
			name:   "PermissionDenied",
			method: "PUT",
			path:   fmt.Sprintf("/api/admin/user/%d", user1.ID),
			req: request.AdminUpdateUserRequest{
				Username: "test_update_user_perm",
				Nickname: "test_update_user_perm_nick",
				Email:    "test_update_user_perm@mail.com",
				Password: "test_update_user_perm_pwd",
			},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, FailTests, "AdminUpdateUser")

	successTests := []struct {
		name         string
		path         string
		originalUser models.User
		expectedUser models.User
		req          request.AdminUpdateUserRequest
	}{
		{
			name: "WithId",
			path: "id",
			originalUser: models.User{
				Username: "test_update_user_2",
				Nickname: "test_update_user_2_nick",
				Email:    "test_update_user_2@mail.com",
				Password: utils.HashPassword("test_update_user_2_pwd"),
			},
			expectedUser: models.User{
				Username: "test_update_user_20",
				Nickname: "test_update_user_20_nick",
				Email:    "test_update_user_20@mail.com",
				Password: "test_update_user_20_pwd",
			},
			req: request.AdminUpdateUserRequest{
				Username: "test_update_user_20",
				Nickname: "test_update_user_20_nick",
				Email:    "test_update_user_20@mail.com",
				Password: "test_update_user_20_pwd",
			},
		},
		{
			name: "WithUsername",
			path: "/api/admin/user/test_update_user_3",
			originalUser: models.User{
				Username: "test_update_user_3",
				Nickname: "test_update_user_3_nick",
				Email:    "test_update_user_3@mail.com",
				Password: utils.HashPassword("test_update_user_3_pwd"),
			},
			expectedUser: models.User{
				Username: "test_update_user_30",
				Nickname: "test_update_user_30_nick",
				Email:    "test_update_user_30@mail.com",
				Password: "test_update_user_30_pwd",
			},
			req: request.AdminUpdateUserRequest{
				Username: "test_update_user_30",
				Nickname: "test_update_user_30_nick",
				Email:    "test_update_user_30@mail.com",
				Password: "test_update_user_30_pwd",
			},
		},
		{
			name: "WithoutChangingPasswordEmpty",
			path: "id",
			originalUser: models.User{
				Username: "test_update_user_4",
				Nickname: "test_update_user_4_nick",
				Email:    "test_update_user_4@mail.com",
				Password: utils.HashPassword("test_update_user_4_pwd"),
			},
			expectedUser: models.User{
				Username: "test_update_user_40",
				Nickname: "test_update_user_40_nick",
				Email:    "test_update_user_40@mail.com",
				Password: "test_update_user_4_pwd",
			},
			req: request.AdminUpdateUserRequest{
				Username: "test_update_user_40",
				Nickname: "test_update_user_40_nick",
				Email:    "test_update_user_40@mail.com",
				Password: "",
			},
		},
		{
			name: "WithoutChangingPassword123456",
			path: "id",
			originalUser: models.User{
				Username: "test_update_user_5",
				Nickname: "test_update_user_5_nick",
				Email:    "test_update_user_5@mail.com",
				Password: utils.HashPassword("test_update_user_5_pwd"),
			},
			expectedUser: models.User{
				Username: "test_update_user_50",
				Nickname: "test_update_user_50_nick",
				Email:    "test_update_user_50@mail.com",
				Password: "test_update_user_5_pwd",
			},
			req: request.AdminUpdateUserRequest{
				Username: "test_update_user_50",
				Nickname: "test_update_user_50_nick",
				Email:    "test_update_user_50@mail.com",
				Password: "123456",
			},
		},
	}

	t.Run("testAdminUpdateUserSuccess", func(t *testing.T) {
		t.Parallel()
		for _, test := range successTests {
			t.Run("testAdminUpdateUser"+test.name, func(t *testing.T) {
				test := test
				t.Parallel()
				assert.Nil(t, base.DB.Create(&test.originalUser).Error)
				if test.path == "id" {
					test.path = fmt.Sprintf("/api/admin/user/%d", test.originalUser.ID)
				}
				httpResp := makeResp(makeReq(t, "PUT", test.path, test.req, applyAdminUser))
				databaseUser := models.User{}
				assert.Nil(t, base.DB.First(&databaseUser, test.originalUser.ID).Error)
				assert.Equal(t, test.expectedUser.Username, databaseUser.Username)
				assert.Equal(t, test.expectedUser.Nickname, databaseUser.Nickname)
				assert.Equal(t, test.expectedUser.Email, databaseUser.Email)
				utils.VerifyPassword(test.expectedUser.Password, databaseUser.Password)
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				jsonEQ(t, response.AdminUpdateUserResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data: struct {
						*models.User `json:"user"`
					}{
						&databaseUser,
					},
				}, httpResp)
			})
		}
	})
}

func TestAdminDeleteUser(t *testing.T) {
	t.Parallel()

	failTests := []failTest{
		{
			name:   "NonExistId",
			method: "DELETE",
			path:   "/api/admin/user/-1",
			req:    request.AdminDeleteUserRequest{},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			name:   "NonExistUsername",
			method: "DELETE",
			path:   "/api/admin/user/test_delete_non_existing_username",
			req:    request.AdminDeleteUserRequest{},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			name:   "NonExistPermissionDenied",
			method: "DELETE",
			path:   "/api/admin/user/-1",
			req:    request.AdminDeleteUserRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "testAdminDeleteUser")

	successTests := []struct {
		name string
		path string
		user models.User
	}{
		{
			name: "WithId",
			path: "id",
			user: models.User{
				Username: "test_delete_user_1",
				Nickname: "test_delete_user_1_nick",
				Email:    "test_delete_user_1@mail.com",
				Password: utils.HashPassword("test_delete_user_1_pwd"),
			},
		},
		{
			name: "WithUsername",
			path: "/api/admin/user/test_delete_user_2",
			user: models.User{
				Username: "test_delete_user_2",
				Nickname: "test_delete_user_2_nick",
				Email:    "test_delete_user_2@mail.com",
				Password: utils.HashPassword("test_delete_user_2_pwd"),
			},
		},
	}

	t.Run("testAdminDeleteUserSuccess", func(t *testing.T) {
		t.Parallel()
		for _, test := range successTests {
			t.Run("testAdminDeleteUser"+test.name, func(t *testing.T) {
				test := test
				t.Parallel()
				assert.Nil(t, base.DB.Create(&test.user).Error)
				if test.path == "id" {
					test.path = fmt.Sprintf("/api/admin/user/%d", test.user.ID)
				}
				httpResp := makeResp(makeReq(t, "DELETE", test.path, request.AdminDeleteUserRequest{}, applyAdminUser))
				resp := response.Response{}
				mustJsonDecode(httpResp, &resp)
				assert.Equal(t, http.StatusNoContent, httpResp.StatusCode)
				assert.Equal(t, response.Response{
					Message: "SUCCESS",
					Error:   nil,
					Data:    nil,
				}, resp)
				assert.Equal(t, gorm.ErrRecordNotFound, base.DB.First(models.User{}, test.user.ID).Error)
			})
		}
	})
}
func TestAdminGetUser(t *testing.T) {
	t.Parallel()

	failTests := []failTest{
		{
			name:   "NonExistId",
			method: "GET",
			path:   "/api/admin/user/-1",
			req:    request.AdminGetUserRequest{},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			name:   "NonExistUsername",
			method: "GET",
			path:   "/api/admin/user/test_get_non_existing_user",
			req:    request.AdminGetUserRequest{},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			name:   "PermissionDenied",
			method: "GET",
			path:   "/api/admin/user/-1",
			req:    request.AdminGetUserRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "testAdminGetUser")

	classA := testClass{ID: 1}
	dummy := "test_class"
	testRole := models.Role{
		Name:   "testAdminGetUserRole",
		Target: &dummy,
	}
	base.DB.Create(&testRole)
	testRole.AddPermission("testAdminGetUserPerm")

	successTests := []struct {
		name       string
		path       string
		req        request.AdminGetUserRequest
		user       models.User
		role       models.Role
		roleTarget models.HasRole
	}{
		{
			name: "WithId",
			path: "id",
			req:  request.AdminGetUserRequest{},
			user: models.User{
				Username: "test_get_user_1",
				Nickname: "test_get_user_1_nick",
				Email:    "test_get_user_1@mail.com",
				Password: utils.HashPassword("test_get_user_1_pwd"),
				Roles:    []models.UserHasRole{},
			},
			role:       models.Role{},
			roleTarget: nil,
		},
		{
			name: "WithUsername",
			path: "/api/admin/user/test_get_user_2",
			req:  request.AdminGetUserRequest{},
			user: models.User{
				Username: "test_get_user_2",
				Nickname: "test_get_user_2_nick",
				Email:    "test_get_user_2@mail.com",
				Password: utils.HashPassword("test_get_user_2_pwd"),
				Roles:    []models.UserHasRole{},
			},
			role:       models.Role{},
			roleTarget: nil,
		},
		{
			name: "WithRole",
			path: "id",
			req:  request.AdminGetUserRequest{},
			user: models.User{
				Username: "test_get_user_3",
				Nickname: "test_get_user_3_nick",
				Email:    "test_get_user_3@mail.com",
				Password: utils.HashPassword("test_get_user_3_pwd"),
			},
			role:       testRole,
			roleTarget: classA,
		},
	}

	t.Run("testAdminGetUserSuccess", func(t *testing.T) {
		t.Parallel()
		for _, test := range successTests {
			t.Run("testAdminGetUser"+test.name, func(t *testing.T) {
				test := test
				t.Parallel()
				assert.Nil(t, base.DB.Create(&test.user).Error)
				if test.role.Name != "" {
					test.user.GrantRole(test.role, test.roleTarget)
				}
				if test.path == "id" {
					test.path = fmt.Sprintf("/api/admin/user/%d", test.user.ID)
				}
				httpResp := makeResp(makeReq(t, "GET", test.path, test.req, applyAdminUser))
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				databaseUser := models.User{}
				assert.Nil(t, base.DB.First(&databaseUser, test.user.ID).Error)
				assert.Equal(t, test.user.Username, databaseUser.Username)
				assert.Equal(t, test.user.Nickname, databaseUser.Nickname)
				assert.Equal(t, test.user.Email, databaseUser.Email)
				jsonEQ(t, response.AdminGetUserResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data: struct {
						*models.User `json:"user"`
					}{
						User: &test.user,
					},
				}, httpResp)
			})
		}
	})
}

func TestAdminGetUsers(t *testing.T) {
	t.Parallel()
	user1 := models.User{
		Username: "test_admin_get_users_01",
		Nickname: "c_test_admin_get_users_1_nick",
		Email:    "1_test_admin_get_users@e.com",
		Password: "test_admin_get_users_1_passwd",
	}
	user2 := models.User{
		Username: "test_admin_get_users_2",
		Nickname: "a0_test_admin_get_users_2_nick",
		Email:    "2_test_admin_get_users@e.com",
		Password: "test_admin_get_users_2_passwd",
	}
	user3 := models.User{
		Username: "test_admin_get_users_03",
		Nickname: "d0_test_admin_get_users_3_nick",
		Email:    "3_test_admin_get_users@f.com",
		Password: "test_admin_get_users_3_passwd",
	}
	user4 := models.User{
		Username: "test_admin_get_users_4",
		Nickname: "b_test_admin_get_users_4_nick",
		Email:    "4_test_admin_get_users@e.com",
		Password: "test_admin_get_users_4_passwd",
	}
	DLUsers := make([]models.User, 25) // DL: Default Limit
	assert.Nil(t, base.DB.Create(&user1).Error)
	assert.Nil(t, base.DB.Create(&user2).Error)
	assert.Nil(t, base.DB.Create(&user3).Error)
	assert.Nil(t, base.DB.Create(&user4).Error)

	for i := 0; i < 25; i++ {
		DLUsers[i] = models.User{
			Username: fmt.Sprintf("test_DL_admin_get_users_%d", i),
			Nickname: fmt.Sprintf("test_DL_admin_get_users_n_%d", i),
			Email:    fmt.Sprintf("test_DL_admin_get_users_%d@e.e", i),
			Password: fmt.Sprintf("test_DL_admin_get_users_pwd_%d", i),
		}
		assert.Nil(t, base.DB.Create(&DLUsers[i]).Error)
	}

	type respData struct {
		Users  []models.User `json:"users"` // TODO:modify models.users
		Total  int           `json:"total"`
		Count  int           `json:"count"`
		Offset int           `json:"offset"`
		Prev   *string       `json:"prev"`
		Next   *string       `json:"next"`
	}

	baseUrl := "/api/admin/users"

	failTests := []failTest{
		{
			name:   "WithWrongOrderByPara",
			method: "GET",
			path:   "/api/admin/users",
			req: request.AdminGetUsersRequest{
				Search:  "test_admin_get_users",
				OrderBy: "wrongOrderByPara",
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("INVALID_ORDER", nil),
		},
		{
			name:   "OrderByNonExistingColumn",
			method: "GET",
			path:   "/api/admin/users",
			req: request.AdminGetUsersRequest{
				Search:  "test_admin_get_users",
				OrderBy: "nonExistingColumn.ASC",
			},
			reqOptions: []reqOption{
				applyAdminUser,
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("INVALID_ORDER", nil),
		},
		{
			name:   "OrderByNonExistingOrder",
			method: "GET",
			path:   "/api/admin/users",
			req: request.AdminGetUsersRequest{
				Search:  "test_admin_get_users",
				OrderBy: "id.NonExistingOrder",
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
			path:   "/api/admin/users",
			req:    request.AdminGetUsersRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusForbidden,
			resp:       response.ErrorResp("PERMISSION_DENIED", nil),
		},
	}

	runFailTests(t, failTests, "testAdminGetUsers")

	successTests := []struct {
		name     string
		req      request.AdminGetUsersRequest
		respData respData
	}{
		{
			name: "testAdminGetUsersAll",
			req: request.AdminGetUsersRequest{
				Search: "test_admin_get_users",
			},
			respData: respData{
				Users: []models.User{
					user1,
					user2,
					user3,
					user4,
				},
				Total:  4,
				Count:  4,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "testAdminGetUsersNonExist",
			req: request.AdminGetUsersRequest{
				Search: "test_admin_get_users_non_exist",
			},
			respData: respData{
				Users: []models.User{},
			},
		},
		{
			name: "testAdminGetUsersSearchUsernameSingle",
			req: request.AdminGetUsersRequest{
				Search: "test_admin_get_users_2",
			},
			respData: respData{
				Users: []models.User{
					user2,
				},
				Total:  1,
				Count:  1,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "testAdminGetUsersSearchNicknameSingle",
			req: request.AdminGetUsersRequest{
				Search: "test_admin_get_users_3_nick",
			},
			respData: respData{
				Users: []models.User{
					user3,
				},
				Total:  1,
				Count:  1,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "testAdminGetUsersSearchEmailSingle",
			req: request.AdminGetUsersRequest{
				Search: "4_test_admin_get_users@e.com",
			},
			respData: respData{
				Users: []models.User{
					user4,
				},
				Total:  1,
				Count:  1,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "testAdminGetUsersSearchUsernameMultiple",
			req: request.AdminGetUsersRequest{
				Search: "test_admin_get_users_0",
			},
			respData: respData{
				Users: []models.User{
					user1,
					user3,
				},
				Total:  2,
				Count:  2,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "testAdminGetUsersSearchNicknameMultiple",
			req: request.AdminGetUsersRequest{
				Search: "0_test_admin_get_users_",
			},
			respData: respData{
				Users: []models.User{
					user2,
					user3,
				},
				Total:  2,
				Count:  2,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "testAdminGetUsersSearchEmailMultiple",
			req: request.AdminGetUsersRequest{
				Search: "_test_admin_get_users@e.com",
			},
			respData: respData{
				Users: []models.User{
					user1,
					user2,
					user4,
				},
				Total:  3,
				Count:  3,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "testAdminGetUsersLimit",
			req: request.AdminGetUsersRequest{
				Search: "test_admin_get_users",
				Limit:  2,
			},
			respData: respData{
				Users: []models.User{
					user1,
					user2,
				},
				Total:  4,
				Count:  2,
				Offset: 0,
				Prev:   nil,
				Next: getUrlStringPointer(baseUrl, map[string]string{
					"limit":  "2",
					"offset": "2",
				}),
			},
		},
		{
			name: "testAdminGetUsersOffset",
			req: request.AdminGetUsersRequest{
				Search: "test_admin_get_users",
				Limit:  2,
			},
			respData: respData{
				Users: []models.User{
					user1,
					user2,
				},
				Total:  4,
				Count:  2,
				Offset: 0,
				Prev:   nil,
				Next: getUrlStringPointer(baseUrl, map[string]string{
					"limit":  "2",
					"offset": "2",
				}),
			},
		},
		{
			name: "testAdminGetUsersLimitAndOffsetNext",
			req: request.AdminGetUsersRequest{
				Search: "test_admin_get_users",
				Limit:  2,
				Offset: 1,
			},
			respData: respData{
				Users: []models.User{
					user2,
					user3,
				},
				Total:  4,
				Count:  2,
				Offset: 1,
				Prev:   nil,
				Next: getUrlStringPointer(baseUrl, map[string]string{
					"limit":  "2",
					"offset": "3",
				}),
			},
		},
		{
			name: "testAdminGetUsersLimitAndOffsetPrev",
			req: request.AdminGetUsersRequest{
				Search: "test_admin_get_users",
				Limit:  2,
				Offset: 2,
			},
			respData: respData{
				Users: []models.User{
					user3,
					user4,
				},
				Total:  4,
				Count:  2,
				Offset: 2,
				Prev: getUrlStringPointer(baseUrl, map[string]string{
					"limit":  "2",
					"offset": "0",
				}),
				Next: nil,
			},
		},
		{
			name: "testAdminGetUsersLimitAndOffsetPrevNext",
			req: request.AdminGetUsersRequest{
				Search: "test_admin_get_users",
				Limit:  1,
				Offset: 2,
			},
			respData: respData{
				Users: []models.User{
					user3,
				},
				Total:  4,
				Count:  1,
				Offset: 2,
				Prev: getUrlStringPointer(baseUrl, map[string]string{
					"limit":  "1",
					"offset": "1",
				}),
				Next: getUrlStringPointer(baseUrl, map[string]string{
					"limit":  "1",
					"offset": "3",
				}),
			},
		},
		{
			name: "testAdminGetUsersOrderByIdDESC",
			req: request.AdminGetUsersRequest{
				Search:  "test_admin_get_users",
				OrderBy: "id.DESC",
			},
			respData: respData{
				Users: []models.User{
					user4,
					user3,
					user2,
					user1,
				},
				Total:  4,
				Count:  4,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "testAdminGetUsersOrderByUsernameASC",
			req: request.AdminGetUsersRequest{
				Search:  "test_admin_get_users",
				OrderBy: "username.ASC",
			},
			respData: respData{
				Users: []models.User{
					user1,
					user3,
					user2,
					user4,
				},
				Total:  4,
				Count:  4,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "testAdminGetUsersOrderByNicknameDESC",
			req: request.AdminGetUsersRequest{
				Search:  "test_admin_get_users",
				OrderBy: "nickname.DESC",
			},
			respData: respData{
				Users: []models.User{
					user3,
					user1,
					user4,
					user2,
				},
				Total:  4,
				Count:  4,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "testAdminGetUsersOrderByNicknameDESCWithLimitAndOffset",
			req: request.AdminGetUsersRequest{
				Search:  "test_admin_get_users",
				OrderBy: "nickname.DESC",
				Limit:   1,
				Offset:  2,
			},
			respData: respData{
				Users: []models.User{
					user4,
				},
				Total:  4,
				Count:  1,
				Offset: 2,
				Prev: getUrlStringPointer(baseUrl, map[string]string{
					"limit":  "1",
					"offset": "1",
				}),
				Next: getUrlStringPointer(baseUrl, map[string]string{
					"limit":  "1",
					"offset": "3",
				}),
			},
		},
		{
			name: "testAdminGetUsersDefaultLimit",
			req: request.AdminGetUsersRequest{
				Search: "test_DL_admin_get_users_",
			},
			respData: respData{
				Users:  DLUsers[:20],
				Total:  25,
				Count:  20,
				Offset: 0,
				Prev:   nil,
				Next: getUrlStringPointer(baseUrl, map[string]string{
					"limit":  "20",
					"offset": "20",
				}),
			},
		},
	}

	t.Run("testAdminGetUsersSuccess", func(t *testing.T) {
		t.Parallel()
		for _, test := range successTests {
			t.Run(test.name, func(t *testing.T) {
				test := test
				t.Parallel()
				httpResp := makeResp(makeReq(t, "GET", "/api/admin/users", test.req, applyAdminUser))
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				resp := response.Response{}
				mustJsonDecode(httpResp, &resp)
				jsonEQ(t, response.AdminGetUsersResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data:    test.respData,
				}, resp)
			})
		}
	})
}
