package controller_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/EduOJ/backend/app/request"
	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/app/response/resource"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/utils"
	"github.com/EduOJ/backend/database/models"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	t.Parallel()

	failTests := []failTest{
		{
			name:   "NonExistId",
			method: "GET",
			path:   base.Echo.Reverse("user.getUser", -1),
			req:    request.GetUserRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
		{
			name:   "NonExistUsername",
			method: "GET",
			path:   base.Echo.Reverse("user.getUser", "test_get_non_existing_user"),
			req:    request.GetUserRequest{},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusNotFound,
			resp:       response.ErrorResp("NOT_FOUND", nil),
		},
	}

	runFailTests(t, failTests, "GetUser")

	classA := testClass{ID: 1}
	dummy := "test_class"
	testRole := models.Role{
		Name:   "testGetUserRole",
		Target: &dummy,
	}
	base.DB.Create(&testRole)
	testRole.AddPermission("testGetUserPerm")

	successTests := []struct {
		name       string
		path       string
		req        request.GetUserRequest
		user       models.User
		roleName   *string
		roleTarget models.HasRole
	}{
		{
			name: "WithId",
			path: "id",
			req:  request.GetUserRequest{},
			user: models.User{
				Username: "test_get_user_1",
				Nickname: "test_get_user_1_nick",
				Email:    "test_get_user_1@mail.com",
				Password: utils.HashPassword("test_get_user_1_pwd"),
			},
			roleName:   nil,
			roleTarget: nil,
		},
		{
			name: "WithUsername",
			path: base.Echo.Reverse("user.getUser", "test_get_user_2"),
			req:  request.GetUserRequest{},
			user: models.User{
				Username: "test_get_user_2",
				Nickname: "test_get_user_2_nick",
				Email:    "test_get_user_2@mail.com",
				Password: utils.HashPassword("test_get_user_2_pwd"),
			},
			roleName:   nil,
			roleTarget: nil,
		},
		{
			name: "WithRole",
			path: "id",
			req:  request.GetUserRequest{},
			user: models.User{
				Username: "test_get_user_3",
				Nickname: "test_get_user_3_nick",
				Email:    "test_get_user_3@mail.com",
				Password: utils.HashPassword("test_get_user_3_pwd"),
			},
			roleName:   &testRole.Name,
			roleTarget: classA,
		},
	}

	t.Run("testGetUserSuccess", func(t *testing.T) {
		t.Parallel()
		for _, test := range successTests {
			test := test
			t.Run("testGetUser"+test.name, func(t *testing.T) {
				t.Parallel()
				assert.NoError(t, base.DB.Create(&test.user).Error)
				if test.roleName != nil {
					test.user.GrantRole(*test.roleName, test.roleTarget)
				} else {
					test.user.LoadRoles()
				}
				if test.path == "id" {
					test.path = base.Echo.Reverse("user.getUser", test.user.ID)
				}
				httpResp := makeResp(makeReq(t, "GET", test.path, test.req, applyNormalUser))
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				databaseUser := models.User{}
				assert.NoError(t, base.DB.First(&databaseUser, test.user.ID).Error)
				assert.Equal(t, test.user.Username, databaseUser.Username)
				assert.Equal(t, test.user.Nickname, databaseUser.Nickname)
				assert.Equal(t, test.user.Email, databaseUser.Email)
				jsonEQ(t, response.GetUserResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data: struct {
						*resource.User `json:"user"`
					}{
						resource.GetUser(&test.user),
					},
				}, httpResp)
			})
		}
	})
}

func TestGetUsers(t *testing.T) {
	t.Parallel()
	user1 := models.User{
		Username: "test_get_users_01",
		Nickname: "c_test_get_users_1_nick",
		Email:    "1_test_get_users@e.com",
		Password: "test_get_users_1_passwd",
	}
	user2 := models.User{
		Username: "test_get_users_2",
		Nickname: "a0_test_get_users_2_nick",
		Email:    "2_test_get_users@e.com",
		Password: "test_get_users_2_passwd",
	}
	user3 := models.User{
		Username: "test_get_users_03",
		Nickname: "d0_test_get_users_3_nick",
		Email:    "3_test_get_users@f.com",
		Password: "test_get_users_3_passwd",
	}
	user4 := models.User{
		Username: "test_get_users_4",
		Nickname: "b_test_get_users_4_nick",
		Email:    "4_test_get_users@e.com",
		Password: "test_get_users_4_passwd",
	}
	DLUsers := make([]*models.User, 25) // DL: Default Limit
	assert.NoError(t, base.DB.Create(&user1).Error)
	assert.NoError(t, base.DB.Create(&user2).Error)
	assert.NoError(t, base.DB.Create(&user3).Error)
	assert.NoError(t, base.DB.Create(&user4).Error)

	for i := 0; i < 25; i++ {
		DLUsers[i] = &models.User{
			Username: fmt.Sprintf("test_DL_get_users_%d", i),
			Nickname: fmt.Sprintf("test_DL_get_users_n_%d", i),
			Email:    fmt.Sprintf("test_DL_get_users_%d@e.e", i),
			Password: fmt.Sprintf("test_DL_get_users_pwd_%d", i),
		}
		assert.NoError(t, base.DB.Create(&DLUsers[i]).Error)
	}

	type respData struct {
		Users  []resource.User `json:"users"`
		Total  int             `json:"total"`
		Count  int             `json:"count"`
		Offset int             `json:"offset"`
		Prev   *string         `json:"prev"`
		Next   *string         `json:"next"`
	}

	failTests := []failTest{
		{
			name:   "WithWrongOrderByPara",
			method: "GET",
			path:   base.Echo.Reverse("user.getUsers"),
			req: request.GetUsersRequest{
				Search:  "test_get_users",
				OrderBy: "wrongOrderByPara",
			},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("INVALID_ORDER", nil),
		},
		{
			name:   "OrderByNonExistingColumn",
			method: "GET",
			path:   base.Echo.Reverse("user.getUsers"),
			req: request.GetUsersRequest{
				Search:  "test_get_users",
				OrderBy: "nonExistingColumn.ASC",
			},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("INVALID_ORDER", nil),
		},
		{
			name:   "OrderByNonExistingOrder",
			method: "GET",
			path:   base.Echo.Reverse("user.getUsers"),
			req: request.GetUsersRequest{
				Search:  "test_get_users",
				OrderBy: "id.NonExistingOrder",
			},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("INVALID_ORDER", nil),
		},
	}

	runFailTests(t, failTests, "GetUsers")

	successTests := []struct {
		name     string
		req      request.GetUsersRequest
		respData respData
	}{
		{
			name: "All",
			req: request.GetUsersRequest{
				Search: "test_get_users",
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user1),
					*resource.GetUser(&user2),
					*resource.GetUser(&user3),
					*resource.GetUser(&user4),
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
			req: request.GetUsersRequest{
				Search: "test_get_users_non_exist",
			},
			respData: respData{
				Users: []resource.User{},
			},
		},
		{
			name: "SearchUsernameSingle",
			req: request.GetUsersRequest{
				Search: "test_get_users_2",
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user2),
				},
				Total:  1,
				Count:  1,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "SearchNicknameSingle",
			req: request.GetUsersRequest{
				Search: "test_get_users_3_nick",
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user3),
				},
				Total:  1,
				Count:  1,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "SearchEmailSingle",
			req: request.GetUsersRequest{
				Search: "4_test_get_users@e.com",
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user4),
				},
				Total:  1,
				Count:  1,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "SearchUsernameMultiple",
			req: request.GetUsersRequest{
				Search: "test_get_users_0",
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user1),
					*resource.GetUser(&user3),
				},
				Total:  2,
				Count:  2,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "SearchNicknameMultiple",
			req: request.GetUsersRequest{
				Search: "0_test_get_users_",
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user2),
					*resource.GetUser(&user3),
				},
				Total:  2,
				Count:  2,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "SearchEmailMultiple",
			req: request.GetUsersRequest{
				Search: "_test_get_users@e.com",
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user1),
					*resource.GetUser(&user2),
					*resource.GetUser(&user4),
				},
				Total:  3,
				Count:  3,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "Limit",
			req: request.GetUsersRequest{
				Search: "test_get_users",
				Limit:  2,
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user1),
					*resource.GetUser(&user2),
				},
				Total:  4,
				Count:  2,
				Offset: 0,
				Prev:   nil,
				Next: getUrlStringPointer("user.getUsers", map[string]string{
					"limit":  "2",
					"offset": "2",
				}),
			},
		},
		{
			name: "Offset",
			req: request.GetUsersRequest{
				Search: "test_get_users",
				Limit:  2,
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user1),
					*resource.GetUser(&user2),
				},
				Total:  4,
				Count:  2,
				Offset: 0,
				Prev:   nil,
				Next: getUrlStringPointer("user.getUsers", map[string]string{
					"limit":  "2",
					"offset": "2",
				}),
			},
		},
		{
			name: "LimitAndOffsetNext",
			req: request.GetUsersRequest{
				Search: "test_get_users",
				Limit:  2,
				Offset: 1,
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user2),
					*resource.GetUser(&user3),
				},
				Total:  4,
				Count:  2,
				Offset: 1,
				Prev:   nil,
				Next: getUrlStringPointer("user.getUsers", map[string]string{
					"limit":  "2",
					"offset": "3",
				}),
			},
		},
		{
			name: "LimitAndOffsetPrev",
			req: request.GetUsersRequest{
				Search: "test_get_users",
				Limit:  2,
				Offset: 2,
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user3),
					*resource.GetUser(&user4),
				},
				Total:  4,
				Count:  2,
				Offset: 2,
				Prev: getUrlStringPointer("user.getUsers", map[string]string{
					"limit":  "2",
					"offset": "0",
				}),
				Next: nil,
			},
		},
		{
			name: "LimitAndOffsetPrevNext",
			req: request.GetUsersRequest{
				Search: "test_get_users",
				Limit:  1,
				Offset: 2,
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user3),
				},
				Total:  4,
				Count:  1,
				Offset: 2,
				Prev: getUrlStringPointer("user.getUsers", map[string]string{
					"limit":  "1",
					"offset": "1",
				}),
				Next: getUrlStringPointer("user.getUsers", map[string]string{
					"limit":  "1",
					"offset": "3",
				}),
			},
		},
		{
			name: "OrderByIdDESC",
			req: request.GetUsersRequest{
				Search:  "test_get_users",
				OrderBy: "id.DESC",
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user4),
					*resource.GetUser(&user3),
					*resource.GetUser(&user2),
					*resource.GetUser(&user1),
				},
				Total:  4,
				Count:  4,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "OrderByUsernameASC",
			req: request.GetUsersRequest{
				Search:  "test_get_users",
				OrderBy: "username.ASC",
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user1),
					*resource.GetUser(&user3),
					*resource.GetUser(&user2),
					*resource.GetUser(&user4),
				},
				Total:  4,
				Count:  4,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "OrderByNicknameDESC",
			req: request.GetUsersRequest{
				Search:  "test_get_users",
				OrderBy: "nickname.DESC",
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user3),
					*resource.GetUser(&user1),
					*resource.GetUser(&user4),
					*resource.GetUser(&user2),
				},
				Total:  4,
				Count:  4,
				Offset: 0,
				Prev:   nil,
				Next:   nil,
			},
		},
		{
			name: "OrderByNicknameDESCWithLimitAndOffset",
			req: request.GetUsersRequest{
				Search:  "test_get_users",
				OrderBy: "nickname.DESC",
				Limit:   1,
				Offset:  2,
			},
			respData: respData{
				Users: []resource.User{
					*resource.GetUser(&user4),
				},
				Total:  4,
				Count:  1,
				Offset: 2,
				Prev: getUrlStringPointer("user.getUsers", map[string]string{
					"limit":  "1",
					"offset": "1",
				}),
				Next: getUrlStringPointer("user.getUsers", map[string]string{
					"limit":  "1",
					"offset": "3",
				}),
			},
		},
		{
			name: "DefaultLimit",
			req: request.GetUsersRequest{
				Search: "test_DL_get_users_",
			},
			respData: respData{
				Users:  resource.GetUserSlice(DLUsers[:20]),
				Total:  25,
				Count:  20,
				Offset: 0,
				Prev:   nil,
				Next: getUrlStringPointer("user.getUsers", map[string]string{
					"limit":  "20",
					"offset": "20",
				}),
			},
		},
	}

	t.Run("testGetUsersSuccess", func(t *testing.T) {
		t.Parallel()
		for _, test := range successTests {
			test := test
			t.Run("testGetUsers"+test.name, func(t *testing.T) {
				t.Parallel()
				httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("user.getUsers"), test.req, applyNormalUser))
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				resp := response.Response{}
				mustJsonDecode(httpResp, &resp)
				jsonEQ(t, response.GetUsersResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data:    test.respData,
				}, resp)
			})
		}
	})
}

func TestGetUserMe(t *testing.T) {
	t.Parallel()
	classA := testClass{ID: 1}
	dummy := "test_class"
	testRole := models.Role{
		Name:   "testGetMeTestRole",
		Target: &dummy,
	}
	base.DB.Create(&testRole)
	testRole.AddPermission("all")

	successTests := []struct {
		name       string
		user       models.User
		roleName   *string
		roleTarget models.HasRole
	}{
		{
			name: "Success",
			user: models.User{
				Username: "test_get_me_4",
				Nickname: "test_get_me_4_nick",
				Email:    "test_get_me_4@mail.com",
				Password: utils.HashPassword("test_get_me_4_password"),
			},
			roleName:   nil,
			roleTarget: nil,
		},
		{
			name: "SuccessWithRole",
			user: models.User{
				Username: "test_get_me_5",
				Nickname: "test_get_me_5_nick",
				Email:    "test_get_me_5@mail.com",
				Password: utils.HashPassword("test_get_me_5_password"),
			},
			roleName:   &testRole.Name,
			roleTarget: classA,
		},
	}

	t.Run("testGetMeSuccess", func(t *testing.T) {
		t.Parallel()
		for _, test := range successTests {
			test := test
			t.Run("testGetMe"+test.name, func(t *testing.T) {
				t.Parallel()
				assert.NoError(t, base.DB.Create(&test.user).Error)
				if test.roleName != nil {
					test.user.GrantRole(*test.roleName, test.roleTarget)
				}
				test.user.LoadRoles()
				httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("user.getMe"), request.GetMeRequest{}, applyUser(test.user)))
				resp := response.GetMeResponse{}
				mustJsonDecode(httpResp, &resp)
				assert.Equal(t, test.user.Username, resp.Data.Username)
				assert.Equal(t, test.user.Nickname, resp.Data.Nickname)
				assert.Equal(t, test.user.Email, resp.Data.Email)
				assert.Equal(t, resource.GetRoleSlice(test.user.Roles), resp.Data.Roles)
				databaseUser := models.User{}
				assert.NoError(t, base.DB.First(&databaseUser, test.user.ID).Error)
				databaseUser.LoadRoles()
				assert.Equal(t, test.user.Username, databaseUser.Username)
				assert.Equal(t, test.user.Nickname, databaseUser.Nickname)
				assert.Equal(t, test.user.Email, databaseUser.Email)
				assert.Equal(t, test.user.Roles, databaseUser.Roles)
			})
		}
	})
}

func TestUpdateUserMe(t *testing.T) {
	t.Parallel()

	user1 := models.User{
		Username: "test_update_me_1",
		Nickname: "test_update_me_1_nick",
		Email:    "test_update_me_1@mail.com",
		Password: utils.HashPassword("test_update_me_1_password"),
	}
	user2 := models.User{
		Username: "test_update_me_2",
		Nickname: "test_update_me_2_nick",
		Email:    "test_update_me_2@mail.com",
		Password: utils.HashPassword("test_update_me_2_password"),
	}
	user3 := models.User{
		Username: "test_update_me_3",
		Nickname: "test_update_me_3_nick",
		Email:    "test_update_me_3@mail.com",
		Password: utils.HashPassword("test_update_me_3_password"),
	}
	dummyUserForConflict := models.User{
		Username: "test_update_me_conflict",
		Nickname: "test_update_me_conflict_nick",
		Email:    "test_update_me_conflict@mail.com",
		Password: utils.HashPassword("test_update_me_conflict_pwd"),
	}
	assert.NoError(t, base.DB.Create(&user1).Error)
	assert.NoError(t, base.DB.Create(&user2).Error)
	assert.NoError(t, base.DB.Create(&user3).Error)
	assert.NoError(t, base.DB.Create(&dummyUserForConflict).Error)

	failTests := []failTest{
		{
			name:   "WithoutParams",
			method: "PUT",
			path:   base.Echo.Reverse("user.getMe"),
			req: request.UpdateMeRequest{
				Username: "",
				Nickname: "",
				Email:    "",
			},
			reqOptions: []reqOption{
				applyUser(user1),
			},
			statusCode: http.StatusBadRequest,
			resp: response.Response{
				Message: "VALIDATION_ERROR",
				Error: []interface{}{
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
				},
				Data: nil,
			},
		},
		{
			name:   "ConflictUsername",
			method: "PUT",
			path:   base.Echo.Reverse("user.getMe"),
			req: request.UpdateMeRequest{
				Username: "test_update_me_conflict",
				Nickname: "test_update_me_2_nick",
				Email:    "test_update_me_2@mail.com",
			},
			reqOptions: []reqOption{
				applyUser(user2),
			},
			statusCode: http.StatusConflict,
			resp:       response.ErrorResp("CONFLICT_USERNAME", nil),
		},
		{
			name:   "ConflictEmail",
			method: "PUT",
			path:   base.Echo.Reverse("user.getMe"),
			req: request.UpdateMeRequest{
				Username: "test_update_me_3",
				Nickname: "test_update_me_3_nick",
				Email:    "test_update_me_conflict@mail.com",
			},
			reqOptions: []reqOption{
				applyUser(user3),
			},
			statusCode: http.StatusConflict,
			resp:       response.ErrorResp("CONFLICT_EMAIL", nil),
		},
	}

	runFailTests(t, failTests, "UpdateMe")

	classA := testClass{ID: 1}
	dummy := "test_class"
	testRole := models.Role{
		Name:   "testUpdateMeTestRole",
		Target: &dummy,
	}
	base.DB.Create(&testRole)
	testRole.AddPermission("all")

	successTests := []struct {
		name       string
		user       models.User
		req        request.UpdateMeRequest
		roleName   *string
		roleTarget models.HasRole
	}{
		{
			name: "Success",
			user: models.User{
				Username: "test_update_me_4",
				Nickname: "test_update_me_4_nick",
				Email:    "test_update_me_4@mail.com",
				Password: utils.HashPassword("test_update_me_4_password"),
			},
			req: request.UpdateMeRequest{
				Username: "test_update_me_success_4",
				Nickname: "test_update_me_success_4",
				Email:    "test_update_me_success_4@e.com",
			},
			roleName:   nil,
			roleTarget: nil,
		},
		{
			name: "SuccessWithRole",
			user: models.User{
				Username: "test_update_me_5",
				Nickname: "test_update_me_5_nick",
				Email:    "test_update_me_5@mail.com",
				Password: utils.HashPassword("test_update_me_5_password"),
			},
			req: request.UpdateMeRequest{
				Username: "test_update_me_success_5",
				Nickname: "test_update_me_success_5",
				Email:    "test_update_me_success_5@e.com",
			},
			roleName:   &testRole.Name,
			roleTarget: classA,
		},
	}

	t.Run("testUpdateMeSuccess", func(t *testing.T) {
		t.Parallel()
		for _, test := range successTests {
			test := test
			t.Run("testUpdateMe"+test.name, func(t *testing.T) {
				t.Parallel()
				assert.NoError(t, base.DB.Create(&test.user).Error)
				assert.NoError(t, base.DB.Create(&models.EmailVerificationToken{
					User:  &test.user,
					Email: test.user.Email,
					Token: "test_token" + test.name,
					Used:  false,
				}).Error)
				if test.roleName != nil {
					test.user.GrantRole(*test.roleName, test.roleTarget)
				}
				test.user.LoadRoles()
				httpResp := makeResp(makeReq(t, "PUT", base.Echo.Reverse("user.getMe"), test.req, applyUser(test.user)))
				resp := response.UpdateMeResponse{}
				mustJsonDecode(httpResp, &resp)
				assert.Equal(t, test.req.Username, resp.Data.Username)
				assert.Equal(t, test.req.Nickname, resp.Data.Nickname)
				assert.Equal(t, test.req.Email, resp.Data.Email)
				assert.Equal(t, resource.GetRoleSlice(test.user.Roles), resp.Data.Roles)
				databaseUser := models.User{}
				assert.NoError(t, base.DB.First(&databaseUser, test.user.ID).Error)
				databaseUser.LoadRoles()
				assert.Equal(t, test.req.Username, databaseUser.Username)
				assert.Equal(t, test.req.Nickname, databaseUser.Nickname)
				assert.Equal(t, test.req.Email, databaseUser.Email)
				assert.Equal(t, test.user.Roles, databaseUser.Roles)

				if test.req.Email != test.user.Email {
					assert.False(t, databaseUser.EmailVerified)
					count := int64(0)
					assert.NoError(t, base.DB.Model(&models.EmailVerificationToken{}).Where("user_id = ? and email != ?", databaseUser.ID, databaseUser.Email).Count(&count).Error)
					assert.Zero(t, count)
					assert.NoError(t, base.DB.Model(&models.EmailVerificationToken{}).Where("user_id = ? and email = ?", databaseUser.ID, databaseUser.Email).Count(&count).Error)
					assert.Equal(t, int64(1), count)
				}
			})
		}
	})
}

func TestUpdateEmail(t *testing.T) {
	t.Parallel()

	user := models.User{
		Username:      "test_update_email_1",
		Nickname:      "test_update_email_1_nick",
		Email:         "test_update_email_1@mail.com",
		EmailVerified: true,
	}
	assert.NoError(t, base.DB.Create(&user).Error)

	failTests := []failTest{
		{
			name:   "ModifyVerifiedEmail",
			method: "PUT",
			path:   base.Echo.Reverse("user.updateEmail"),
			req: request.UpdateEmailRequest{
				Email: "test_update_email_1_new@mail.com",
			},
			reqOptions: []reqOption{
				applyUser(user),
			},
			statusCode: http.StatusNotAcceptable,
			resp: response.Response{
				Message: "EMAIL_VERIFIED",
				Error:   nil,
				Data:    nil,
			},
		},
	}

	runFailTests(t, failTests, "UpdateEmail")

	successTests := []struct {
		name string
		user models.User
		req  request.UpdateEmailRequest
	}{
		{
			name: "Success",
			user: models.User{
				Username:      "test_update_email_2",
				Nickname:      "test_update_email_2_nick",
				Email:         "test_update_email_2@mail.com",
				EmailVerified: false,
				Password:      utils.HashPassword("test_update_me_2_password"),
			},
			req: request.UpdateEmailRequest{
				Email: "test_update_email_2_new@mail.com",
			},
		},
	}

	t.Run("testUpdateEmailSuccess", func(t *testing.T) {
		t.Parallel()
		for _, test := range successTests {
			test := test
			t.Run("testUpdateEmail"+test.name, func(t *testing.T) {
				t.Parallel()
				assert.NoError(t, base.DB.Create(&test.user).Error)
				assert.NoError(t, base.DB.Create(&models.EmailVerificationToken{
					User:  &test.user,
					Email: test.user.Email,
					Token: "test_token" + test.name,
					Used:  false,
				}).Error)
				httpResp := makeResp(makeReq(t, "PUT", base.Echo.Reverse("user.updateEmail"), test.req, applyUser(test.user)))
				resp := response.UpdateEmailResponse{}
				mustJsonDecode(httpResp, &resp)
				assert.Equal(t, test.req.Email, resp.Data.Email)
				databaseUser := models.User{}
				assert.NoError(t, base.DB.First(&databaseUser, test.user.ID).Error)
				assert.Equal(t, test.req.Email, databaseUser.Email)
				assert.False(t, databaseUser.EmailVerified)
				count := int64(0)
				assert.NoError(t, base.DB.Model(&models.EmailVerificationToken{}).Where("user_id = ?", databaseUser.ID).Count(&count).Error)
				assert.Equal(t, count, int64(1))
			})
		}
	})
}

func TestChangePassword(t *testing.T) {
	t.Parallel()

	t.Run("testChangePasswordWithoutParams", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", base.Echo.Reverse("user.changePassword"), request.ChangePasswordRequest{
			OldPassword: "",
			NewPassword: "",
		}, applyNormalUser))
		assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)
		jsonEQ(t, response.Response{
			Message: "VALIDATION_ERROR",
			Error: []map[string]string{
				{
					"field":       "OldPassword",
					"reason":      "required",
					"translation": "旧密码为必填字段",
				},
				{
					"field":       "NewPassword",
					"reason":      "required",
					"translation": "新密码为必填字段",
				},
			},
			Data: nil,
		}, httpResp)
	})

	t.Run("testChangePasswordSuccess", func(t *testing.T) {
		t.Parallel()
		user := models.User{
			Username: "test_change_passwd_1",
			Nickname: "test_change_passwd_1_nick",
			Email:    "test_change_passwd_1@mail.com",
			Password: utils.HashPassword("test_change_passwd_old_passwd"),
		}
		assert.NoError(t, base.DB.Create(&user).Error)
		token1 := models.Token{
			Token: utils.RandStr(32),
			User:  user,
		}
		token2 := models.Token{
			Token: utils.RandStr(32),
			User:  user,
		}
		token3 := models.Token{
			Token: utils.RandStr(32),
			User:  user,
		}
		assert.NoError(t, base.DB.Create(&token1).Error)
		assert.NoError(t, base.DB.Create(&token2).Error)
		assert.NoError(t, base.DB.Create(&token3).Error)
		httpResp := makeResp(makeReq(t, "POST", base.Echo.Reverse("user.changePassword"), request.ChangePasswordRequest{
			OldPassword: "test_change_passwd_old_passwd",
			NewPassword: "test_change_passwd_new_passwd",
		}, headerOption{
			"Authorization": {token1.Token},
		}))
		var tokens []models.Token
		var updatedUser models.User
		assert.NoError(t, base.DB.Preload("User").Where("user_id = ?", user.ID).Find(&tokens).Error)
		token1, _ = utils.GetToken(token1.Token)
		assert.Equal(t, []models.Token{
			token1,
		}, tokens)

		assert.NoError(t, base.DB.First(&updatedUser, user.ID).Error)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		jsonEQ(t, response.Response{
			Message: "SUCCESS",
			Error:   nil,
			Data:    nil,
		}, httpResp)
		assert.True(t, utils.VerifyPassword("test_change_passwd_new_passwd", updatedUser.Password))
	})

	t.Run("testChangePasswordWithWrongPassword", func(t *testing.T) {
		t.Parallel()
		user := models.User{
			Username: "test_change_passwd_2",
			Nickname: "test_change_passwd_2_nick",
			Email:    "test_change_passwd_2@mail.com",
			Password: utils.HashPassword("test_change_passwd_old_passwd"),
		}
		assert.NoError(t, base.DB.Create(&user).Error)
		mainToken := models.Token{
			Token: utils.RandStr(32),
			User:  user,
		}
		assert.NoError(t, base.DB.Create(&mainToken).Error)
		httpResp := makeResp(makeReq(t, "POST", base.Echo.Reverse("user.changePassword"), request.ChangePasswordRequest{
			OldPassword: "test_change_passwd_wrong",
			NewPassword: "test_change_passwd_new_passwd",
		}, headerOption{
			"Authorization": {mainToken.Token},
		}))
		assert.Equal(t, http.StatusForbidden, httpResp.StatusCode)
		jsonEQ(t, response.Response{
			Message: "WRONG_PASSWORD",
			Error:   nil,
			Data:    nil,
		}, httpResp)
	})
}

func TestGetClassesIManageAndTake(t *testing.T) {
	t.Parallel()

	users := []models.User{
		createUserForTest(t, "test_get_classes_i_manage_or_take", 0),
		createUserForTest(t, "test_get_classes_i_manage_or_take", 1),
		createUserForTest(t, "test_get_classes_i_manage_or_take", 2),
		createUserForTest(t, "test_get_classes_i_manage_or_take", 3),
	}

	class1 := createClassForTest(t, "test_get_classes_i_manage_or_take", 1, []*models.User{
		&users[1],
	}, []*models.User{
		&users[2],
		&users[3],
	})
	class2 := createClassForTest(t, "test_get_classes_i_manage_or_take", 2, []*models.User{
		&users[3],
	}, []*models.User{
		&users[2],
		&users[3],
	})
	class3 := createClassForTest(t, "test_get_classes_i_manage_or_take", 3, []*models.User{
		&users[1],
		&users[3],
	}, []*models.User{})
	createClassForTest(t, "test_get_classes_i_manage_or_take", 4, []*models.User{}, []*models.User{})

	for i := range users {
		assert.NoError(t, base.DB.First(&users[i], users[i].ID).Error)
	}

	createProblemSetForTest(t, "test_get_classes_i_manage_or_take_1", 1, &class1, nil)
	createProblemSetForTest(t, "test_get_classes_i_manage_or_take_1", 2, &class1, nil)
	createProblemSetForTest(t, "test_get_classes_i_manage_or_take_2", 1, &class2, nil)
	createProblemSetForTest(t, "test_get_classes_i_manage_or_take_2", 2, &class2, nil)
	createProblemSetForTest(t, "test_get_classes_i_manage_or_take_3", 1, &class3, nil)

	class1.Students = nil
	class1.Managers = nil
	class2.Students = nil
	class2.Managers = nil
	class3.Students = nil
	class3.Managers = nil

	manageClasses := map[int][]models.Class{
		0: {},
		1: {
			class1,
			class3,
		},
		2: {},
		3: {
			class2,
			class3,
		},
	}
	takeClasses := map[int][]models.Class{
		0: {},
		1: {},
		2: {
			class1,
			class2,
		},
		3: {
			class1,
			class2,
		},
	}

	t.Run("GetClassesIManage", func(t *testing.T) {
		t.Parallel()
		for i, classes := range manageClasses {
			i := i
			classes := classes
			for i := range classes {
				classes[i].ProblemSets = []*models.ProblemSet{}
			}
			t.Run(fmt.Sprintf("User%d", i), func(t *testing.T) {
				t.Parallel()

				httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("user.getClassesIManage"),
					request.GetClassesIManageRequest{}, applyUser(users[i])))
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				resp := response.GetClassesIManageResponse{}
				mustJsonDecode(httpResp, &resp)
				assert.Equal(t, response.GetClassesIManageResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data: struct {
						Classes []resource.Class `json:"classes"`
					}{
						resource.GetClassSlice(classes),
					},
				}, resp)
			})
		}
	})
	t.Run("GetClassesITake", func(t *testing.T) {
		t.Parallel()
		for i, classes := range takeClasses {
			i := i
			classes := classes
			t.Run(fmt.Sprintf("User%d", i), func(t *testing.T) {
				t.Parallel()

				httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("user.getClassesITake"),
					request.GetClassesIManageRequest{}, applyUser(users[i])))
				assert.Equal(t, http.StatusOK, httpResp.StatusCode)
				resp := response.GetClassesIManageResponse{}
				mustJsonDecode(httpResp, &resp)
				assert.Equal(t, response.GetClassesIManageResponse{
					Message: "SUCCESS",
					Error:   nil,
					Data: struct {
						Classes []resource.Class `json:"classes"`
					}{
						resource.GetClassSlice(classes),
					},
				}, resp)
			})
		}
	})
}

func TestGetUserProblemInfo(t *testing.T) {
	t.Parallel()

	user := createUserForTest(t, "get_user_problem_info", 1)
	otherUser := createUserForTest(t, "get_user_problem_info", 2)
	problem1 := createProblemForTest(t, "get_user_problem_info", 1, nil, user)
	problem2 := createProblemForTest(t, "get_user_problem_info", 2, nil, user)
	problem3 := createProblemForTest(t, "get_user_problem_info", 3, nil, user)

	submissionPassed1 := createSubmissionForTest(t, "get_user_problem_info_1_passed", 1, &problem1, &user, nil, 0)
	submissionPassed2 := createSubmissionForTest(t, "get_user_problem_info_1_passed", 2, &problem1, &otherUser, nil, 0)
	submissionPassed3 := createSubmissionForTest(t, "get_user_problem_info_2_passed", 3, &problem2, &user, nil, 0)
	submissionPassed4 := createSubmissionForTest(t, "get_user_problem_info_2_passed", 4, &problem2, &user, nil, 0)
	submissionPassed1.Status = "ACCEPTED"
	submissionPassed2.Status = "ACCEPTED"
	submissionPassed3.Status = "ACCEPTED"
	submissionPassed4.Status = "ACCEPTED"
	assert.NoError(t, base.DB.Save(&submissionPassed1).Error)
	assert.NoError(t, base.DB.Save(&submissionPassed2).Error)
	assert.NoError(t, base.DB.Save(&submissionPassed3).Error)
	assert.NoError(t, base.DB.Save(&submissionPassed4).Error)

	submissionFailed1 := createSubmissionForTest(t, "get_user_problem_info_1_failed", 1, &problem1, &otherUser, nil, 0)
	submissionFailed2 := createSubmissionForTest(t, "get_user_problem_info_2_failed", 2, &problem2, &user, nil, 0)
	submissionFailed3 := createSubmissionForTest(t, "get_user_problem_info_3_failed", 3, &problem3, &user, nil, 0)
	submissionFailed4 := createSubmissionForTest(t, "get_user_problem_info_3_failed", 4, &problem3, &user, nil, 0)
	submissionFailed1.Status = "WRONG_ANSWER"
	submissionFailed2.Status = "TIME_LIMIT_EXCEEDED"
	submissionFailed3.Status = "RUNTIME_ERROR"
	submissionFailed4.Status = "PENDING"
	assert.NoError(t, base.DB.Save(&submissionFailed1).Error)
	assert.NoError(t, base.DB.Save(&submissionFailed2).Error)
	assert.NoError(t, base.DB.Save(&submissionFailed3).Error)
	assert.NoError(t, base.DB.Save(&submissionFailed4).Error)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("user.getUserProblemInfo", user.ID),
			request.GetUserProblemInfoRequest{}, applyNormalUser))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		resp := response.GetUserProblemInfoResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, response.GetUserProblemInfoResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				TriedCount  int `json:"tried_count"`
				PassedCount int `json:"passed_count"`
				Rank        int `json:"rank"`
			}{
				TriedCount:  1,
				PassedCount: 2,
				Rank:        0,
			},
		}, resp)
	})

	t.Run("NonExistingUser", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("user.getUserProblemInfo", -1),
			request.GetUserProblemInfoRequest{}, applyNormalUser))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		resp := response.GetUserProblemInfoResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, response.GetUserProblemInfoResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				TriedCount  int `json:"tried_count"`
				PassedCount int `json:"passed_count"`
				Rank        int `json:"rank"`
			}{
				TriedCount:  0,
				PassedCount: 0,
				Rank:        0,
			},
		}, resp)
	})
}

func TestVerifyEmail(t *testing.T) {
	t.Parallel()
	user := models.User{
		Username: "test_verify_email_username",
		Nickname: "test_verify_email_nickname",
		Email:    "test_verify_email@e.com",
		Password: "test_verify_email_passwd",
	}
	base.DB.Create(&user)
	sucUser := models.User{
		Username: "test_verify_email2_username",
		Nickname: "test_verify_email2_nickname",
		Email:    "test_verify_email2@e.com",
		Password: "test_verify_email2_passwd",
	}
	base.DB.Create(&sucUser)
	code := models.EmailVerificationToken{
		User:  &user,
		Email: user.Email,
		Token: "QwE12",
		Used:  false,
	}
	base.DB.Create(&code)
	sucCode := models.EmailVerificationToken{
		User:  &sucUser,
		Email: sucUser.Email,
		Token: "QwE1X",
		Used:  false,
	}
	base.DB.Create(&sucCode)
	oldCode := models.EmailVerificationToken{
		User:  &user,
		Email: user.Email,
		Token: "QwE21",
		Used:  false,
	}
	base.DB.Create(&oldCode)
	oldCode.CreatedAt = time.Now().Add(-100 * time.Minute)
	base.DB.Save(&oldCode)
	usedCode := models.EmailVerificationToken{
		User:  &user,
		Email: user.Email,
		Token: "QwA21",
		Used:  true,
	}
	base.DB.Create(&usedCode)
	failTests := []failTest{
		{
			name:   "WithoutParams",
			method: "POST",
			path:   base.Echo.Reverse("auth.email.verify"),
			req:    request.VerifyEmailRequest{},
			reqOptions: []reqOption{
				applyUser(user),
			},
			statusCode: http.StatusBadRequest,
			resp: response.Response{
				Message: "VALIDATION_ERROR",
				Error: []interface{}{
					map[string]interface{}{
						"field":       "Token",
						"reason":      "required",
						"translation": "验证码为必填字段",
					},
				},
				Data: nil,
			},
		},
		{
			name:   "OldCode",
			method: "POST",
			path:   base.Echo.Reverse("auth.email.verify"),
			req:    request.VerifyEmailRequest{Token: oldCode.Token},
			reqOptions: []reqOption{
				applyUser(user),
			},
			statusCode: http.StatusRequestTimeout,
			resp: response.Response{
				Message: "CODE_EXPIRED",
				Error:   nil,
				Data:    nil,
			},
		},
		{
			name:   "UsedCode",
			method: "POST",
			path:   base.Echo.Reverse("auth.email.verify"),
			req:    request.VerifyEmailRequest{Token: usedCode.Token},
			reqOptions: []reqOption{
				applyUser(user),
			},
			statusCode: http.StatusRequestTimeout,
			resp: response.Response{
				Message: "CODE_USED",
				Error:   nil,
				Data:    nil,
			},
		},

		{
			name:   "WRONG_CODE",
			method: "POST",
			path:   base.Echo.Reverse("auth.email.verify"),
			req:    request.VerifyEmailRequest{Token: "QWERT"},
			reqOptions: []reqOption{
				applyUser(user),
			},
			statusCode: http.StatusUnauthorized,
			resp: response.Response{
				Message: "WRONG_CODE",
				Error:   nil,
				Data:    nil,
			},
		},
	}

	runFailTests(t, failTests, "VerifyEmail")

	t.Run("VerifyEmailSuccess", func(t *testing.T) {
		t.Parallel()

		httpResp := makeResp(makeReq(t, "POST", base.Echo.Reverse("auth.email.verify"),
			request.VerifyEmailRequest{Token: sucCode.Token}, applyUser(sucUser)))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		resp := response.EmailVerificationResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, response.EmailVerificationResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data:    nil,
		}, resp)
		base.DB.Find(&sucUser, sucUser.ID)
		assert.True(t, sucUser.EmailVerified)
		base.DB.Find(&sucCode, sucCode.ID)
		assert.True(t, sucCode.Used)
	})
}
