package controller_test

import (
	"encoding/json"
	"fmt"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
)

var initNormalUser sync.Once
var normalUser models.User

func initNormalUserFunc() {
	normalUser = models.User{
		Username: "test_user_normal_user",
		Nickname: "test_user_normal_nickname",
		Email:    "test_user_normal@mail.com",
		Password: "test_user_normal_password",
	}
	base.DB.Create(&normalUser)
}

func getNormalToken() (token models.Token) {
	initNormalUser.Do(initNormalUserFunc)
	token = models.Token{
		User:  normalUser,
		Token: utils.RandStr(32),
	}
	base.DB.Create(&token)
	return
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	token := getNormalToken()

	t.Run("getUserNonExistId", func(t *testing.T) {
		t.Parallel()
		resp := response.Response{}
		httpResp := makeResp(makeReq(t, "GET", "/api/user/-1", request.GetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
		assert.Equal(t, response.ErrorResp("NOT_FOUND", nil), resp)
	})
	t.Run("getUserNonExistUsername", func(t *testing.T) {
		t.Parallel()
		resp := response.Response{}
		httpResp := makeResp(makeReq(t, "GET", "/api/user/test_get_non_existing_user", request.GetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
		assert.Equal(t, response.ErrorResp("NOT_FOUND", nil), resp)
	})
	t.Run("getUserSuccessWithId", func(t *testing.T) {
		t.Parallel()
		user := models.User{
			Username: "test_get_user_4",
			Nickname: "test_get_user_4_rand_str",
			Email:    "test_get_user_4@mail.com",
			Password: utils.HashPassword("test_get_user_4_password"),
		}
		base.DB.Create(&user)
		resp := makeResp(makeReq(t, "GET", fmt.Sprintf("/api/user/%d", user.ID), request.GetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		jsonEQ(t, response.GetUserResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*models.User `json:"user"`
			}{
				&user,
			},
		}, resp)
	})
	t.Run("getUserSuccessWithUsername", func(t *testing.T) {
		user := models.User{
			Username: "test_get_user_5",
			Nickname: "test_get_user_5_rand_str",
			Email:    "test_get_user_5@mail.com",
			Password: utils.HashPassword("test_get_user_5_password"),
		}
		base.DB.Create(&user)
		t.Parallel()
		resp := makeResp(makeReq(t, "GET", "/api/user/test_get_user_5", request.GetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		jsonEQ(t, response.GetUserResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*models.User `json:"user"`
			}{
				&user,
			},
		}, resp)
	})
}

func TestGetUserMe(t *testing.T) {
	t.Parallel()

	t.Run("getUserSuccess", func(t *testing.T) {
		t.Parallel()
		user := models.User{
			Username: "test_get_user_me_1",
			Nickname: "test_get_user_me_1_rand_str",
			Email:    "test_get_user_me_1@mail.com",
			Password: utils.HashPassword("test_get_user_me_1_password"),
			Roles:    []models.UserHasRole{},
		}
		base.DB.Create(&user)
		token := models.Token{
			Token: utils.RandStr(32),
			User:  user,
		}
		base.DB.Create(&token)
		resp := makeResp(makeReq(t, "GET", "/api/user/me", request.GetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		jsonEQ(t, response.GetUserResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*models.User `json:"user"`
			}{
				&user,
			},
		}, resp)
	})
	t.Run("getUserSuccessWithRole", func(t *testing.T) {
		t.Parallel()
		classA := testClass{ID: 1}
		dummy := "test_class"
		adminRole := models.Role{
			Name:   "admin",
			Target: &dummy,
		}
		base.DB.Create(&adminRole)
		adminRole.AddPermission("all")
		user := models.User{
			Username: "test_get_user_me_2",
			Nickname: "test_get_user_me_2_rand_str",
			Email:    "test_get_user_me_2@mail.com",
			Password: utils.HashPassword("test_get_user_me_2_password"),
			Roles:    []models.UserHasRole{},
		}
		base.DB.Create(&user)
		user.GrantRole(adminRole, classA)
		token := models.Token{
			Token: utils.RandStr(32),
			User:  user,
		}
		base.DB.Create(&token)
		resp := makeResp(makeReq(t, "GET", "/api/user/me", request.GetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		jsonEQ(t, response.GetUserResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*models.User `json:"user"`
			}{
				&user,
			},
		}, resp)
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
	assert.Nil(t, base.DB.Create(&user1).Error)
	assert.Nil(t, base.DB.Create(&user2).Error)
	assert.Nil(t, base.DB.Create(&user3).Error)
	assert.Nil(t, base.DB.Create(&user4).Error)

	type respData struct {
		Users  []models.User `json:"users"` // TODO:modify models.users
		Total  int           `json:"total"`
		Count  int           `json:"count"`
		Offset int           `json:"offset"`
		Prev   *string       `json:"prev"`
		Next   *string       `json:"next"`
	}

	token := getNormalToken()
	baseUrl := "/api/users"

	t.Run("testGetUsersSuccess", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			req      request.GetUsersRequest
			respData respData
		}{
			{
				name: "testGetUsersAll",
				req: request.GetUsersRequest{
					Search: "test_get_users",
				},
				respData: respData{
					Users: []models.User{
						user1,
						user2,
						user3,
						user4,
					},
					Total: 4,
					Count: 4,
				},
			},
			{
				name: "testGetUsersNonExist",
				req: request.GetUsersRequest{
					Search: "test_get_users_non_exist",
				},
				respData: respData{
					Users: []models.User{},
				},
			},
			{
				name: "testGetUsersSearchUsernameSingle",
				req: request.GetUsersRequest{
					Search: "test_get_users_2",
				},
				respData: respData{
					Users: []models.User{
						user2,
					},
					Total: 1,
					Count: 1,
				},
			},
			{
				name: "testGetUsersSearchNicknameSingle",
				req: request.GetUsersRequest{
					Search: "test_get_users_3_nick",
				},
				respData: respData{
					Users: []models.User{
						user3,
					},
					Total: 1,
					Count: 1,
				},
			},
			{
				name: "testGetUsersSearchEmailSingle",
				req: request.GetUsersRequest{
					Search: "4_test_get_users@e.com",
				},
				respData: respData{
					Users: []models.User{
						user4,
					},
					Total: 1,
					Count: 1,
				},
			},
			{
				name: "testGetUsersSearchUsernameMultiple",
				req: request.GetUsersRequest{
					Search: "test_get_users_0",
				},
				respData: respData{
					Users: []models.User{
						user1,
						user3,
					},
					Total: 2,
					Count: 2,
				},
			},
			{
				name: "testGetUsersSearchNicknameMultiple",
				req: request.GetUsersRequest{
					Search: "0_test_get_users_",
				},
				respData: respData{
					Users: []models.User{
						user2,
						user3,
					},
					Total: 2,
					Count: 2,
				},
			},
			{
				name: "testGetUsersSearchEmailMultiple",
				req: request.GetUsersRequest{
					Search: "_test_get_users@e.com",
				},
				respData: respData{
					Users: []models.User{
						user1,
						user2,
						user4,
					},
					Total: 3,
					Count: 3,
				},
			},
			{
				name: "testGetUsersLimit",
				req: request.GetUsersRequest{
					Search: "test_get_users",
					Limit:  2,
				},
				respData: respData{
					Users: []models.User{
						user1,
						user2,
					},
					Total: 4,
					Count: 2,
					Next: getUrlStringPointer(baseUrl, map[string]string{
						"limit":  "2",
						"offset": "2",
					}),
				},
			},
			{
				name: "testGetUsersOffset",
				req: request.GetUsersRequest{
					Search: "test_get_users",
					Limit:  2,
				},
				respData: respData{
					Users: []models.User{
						user1,
						user2,
					},
					Total: 4,
					Count: 2,
					Next: getUrlStringPointer(baseUrl, map[string]string{
						"limit":  "2",
						"offset": "2",
					}),
				},
			},
			{
				name: "testGetUsersLimitAndOffsetNext",
				req: request.GetUsersRequest{
					Search: "test_get_users",
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
					Next: getUrlStringPointer(baseUrl, map[string]string{
						"limit":  "2",
						"offset": "3",
					}),
				},
			},
			{
				name: "testGetUsersLimitAndOffsetPrev",
				req: request.GetUsersRequest{
					Search: "test_get_users",
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
				},
			},
			{
				name: "testGetUsersLimitAndOffsetPrevNext",
				req: request.GetUsersRequest{
					Search: "test_get_users",
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
				name: "testGetUsersOrderByIdDESC",
				req: request.GetUsersRequest{
					Search:  "test_get_users",
					OrderBy: "id.DESC",
				},
				respData: respData{
					Users: []models.User{
						user4,
						user3,
						user2,
						user1,
					},
					Total: 4,
					Count: 4,
				},
			},
			{
				name: "testGetUsersOrderByUsernameASC",
				req: request.GetUsersRequest{
					Search:  "test_get_users",
					OrderBy: "username.ASC",
				},
				respData: respData{
					Users: []models.User{
						user1,
						user3,
						user2,
						user4,
					},
					Total: 4,
					Count: 4,
				},
			},
			{
				name: "testGetUsersOrderByNicknameDESC",
				req: request.GetUsersRequest{
					Search:  "test_get_users",
					OrderBy: "nickname.DESC",
				},
				respData: respData{
					Users: []models.User{
						user3,
						user1,
						user4,
						user2,
					},
					Total: 4,
					Count: 4,
				},
			},
			{
				name: "testGetUsersOrderByNicknameDESCWithLimitAndOffset",
				req: request.GetUsersRequest{
					Search:  "test_get_users",
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
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				test := test
				t.Parallel()
				httpResp := makeResp(makeReq(t, "GET", "/api/users", test.req, headerOption{
					"Authorization": {token.Token},
				}))
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
	t.Run("testGetUsersWithWrongOrderByPara", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "GET", "/api/users", request.GetUsersRequest{
			Search:  "test_get_users",
			OrderBy: "wrongOrderByPara",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		jsonEQ(t, response.Response{
			Message: "INVALID_ORDER",
			Error:   nil,
			Data:    nil,
		}, resp)
	})
	t.Run("testGetUsersOrderByNonExistingColumn", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "GET", "/api/users", request.GetUsersRequest{
			Search:  "test_get_users",
			OrderBy: "nonExistingColumn.ASC",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		jsonEQ(t, response.Response{
			Message: "INVALID_ORDER",
			Error:   nil,
			Data:    nil,
		}, resp)
	})
	t.Run("testGetUsersOrderByNonExistingOrder", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "GET", "/api/users", request.GetUsersRequest{
			Search:  "test_get_users",
			OrderBy: "id.NonExistingOrder",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		jsonEQ(t, response.Response{
			Message: "INVALID_ORDER",
			Error:   nil,
			Data:    nil,
		}, resp)
	})
	t.Run("testGetUsersDefaultLimit", func(t *testing.T) {
		t.Parallel()
		// DL: default limit
		users := make([]models.User, 25)
		for i := 0; i < 25; i++ {
			users[i] = models.User{
				Username: fmt.Sprintf("test_DL_get_users_%d", i),
				Nickname: fmt.Sprintf("test_DL_get_users_n_%d", i),
				Email:    fmt.Sprintf("test_DL_get_users_%d@e.e", i),
				Password: fmt.Sprintf("test_DL_get_users_pwd_%d", i),
			}
			base.DB.Create(&users[i])
		}
		httpResp := makeResp(makeReq(t, "GET", "/api/users", request.GetUsersRequest{
			Search: "test_DL_get_users_",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		jsonEQ(t, response.GetUsersResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: respData{
				Users: users[:20],
				Count: 20,
				Total: 25,
				Next: getUrlStringPointer(baseUrl, map[string]string{
					"limit":  "20",
					"offset": "20",
				}),
			},
		}, resp)
	})
}

func TestUpdateUserMe(t *testing.T) {
	t.Parallel()
	t.Run("testUpdateUserMeWithoutParams", func(t *testing.T) {
		t.Parallel()
		user := models.User{
			Username: "test_put_user_me_1",
			Nickname: "test_put_user_me_1_rand_str",
			Email:    "test_put_user_me_1@mail.com",
			Password: utils.HashPassword("test_put_user_me_1_password"),
		}
		base.DB.Create(&user)
		token := models.Token{
			Token: utils.RandStr(32),
			User:  user,
		}
		base.DB.Create(&token)
		httpResp := makeResp(makeReq(t, "PUT", "/api/user/me", request.UpdateUserRequest{
			Username: "",
			Nickname: "",
			Email:    "",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)
		assert.Equal(t, response.Response{
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
		}, resp)
	})
	t.Run("testUpdateUserMeWithParams", func(t *testing.T) {
		t.Parallel()
		user2 := models.User{
			Username: "test_put_user_me_2",
			Nickname: "test_put_user_me_2_rand_str",
			Email:    "test_put_user_me_2@mail.com",
			Password: utils.HashPassword("test_put_user_me_2_password"),
		}
		user3 := models.User{
			Username: "test_put_user_me_3",
			Nickname: "test_put_user_me_3_rand_str",
			Email:    "test_put_user_me_3@mail.com",
			Password: utils.HashPassword("test_put_user_me_3_password"),
		}
		base.DB.Create(&user2)
		base.DB.Create(&user3)
		user2.LoadRoles()
		user3.LoadRoles()
		user2Token := models.Token{
			Token: utils.RandStr(32),
			User:  user2,
		}
		base.DB.Create(&user2Token)
		t.Run("testUpdateUserMeSuccess", func(t *testing.T) {
			t.Parallel()
			user := models.User{
				Username: "test_put_user_me_4",
				Nickname: "test_put_user_me_4_rand_str",
				Email:    "test_put_user_me_4@mail.com",
				Password: utils.HashPassword("test_put_user_me_4_password"),
			}
			base.DB.Create(&user)
			user.LoadRoles()
			token := models.Token{
				Token: utils.RandStr(32),
				User:  user,
			}
			base.DB.Create(&token)
			respResponse := makeResp(makeReq(t, "PUT", "/api/user/me", request.UpdateUserRequest{
				Username: "test_put_user_me_success_0",
				Nickname: "test_put_user_me_success_0",
				Email:    "test_put_user_me_success_0@e.com",
			}, headerOption{
				"Authorization": {token.Token},
			}))
			assert.Equal(t, http.StatusOK, respResponse.StatusCode)
			resp := response.UpdateUserResponse{}
			respBytes, err := ioutil.ReadAll(respResponse.Body)
			assert.Equal(t, nil, err)
			err = json.Unmarshal(respBytes, &resp)
			assert.Equal(t, nil, err)
			databaseUser := models.User{}
			err = base.DB.Where("id = ?", user.ID).First(&databaseUser).Error
			assert.Equal(t, nil, err)
			databaseUser.LoadRoles()
			jsonEQ(t, resp.Data.User, databaseUser)
		})
		t.Run("testUpdateUserMeDuplicateEmail", func(t *testing.T) {
			t.Parallel()
			resp := response.Response{}
			httpResp := makeResp(makeReq(t, "PUT", "/api/user/me", request.UpdateUserRequest{
				Username: "test_put_user_me_2",
				Nickname: "test_put_user_me_2_rand_str",
				Email:    "test_put_user_me_3@mail.com",
			}, headerOption{
				"Authorization": {user2Token.Token},
			}))
			mustJsonDecode(httpResp, &resp)
			assert.Equal(t, http.StatusConflict, httpResp.StatusCode)
			assert.Equal(t, response.ErrorResp("DUPLICATE_EMAIL", nil), resp)
		})
		t.Run("testUpdateUserMeDuplicateUsername", func(t *testing.T) {
			t.Parallel()
			resp := response.Response{}
			httpResp := makeResp(makeReq(t, "PUT", "/api/user/me", request.UpdateUserRequest{
				Username: "test_put_user_me_3",
				Nickname: "test_put_user_me_2_rand_str",
				Email:    "test_put_user_me_2@mail.com",
			}, headerOption{
				"Authorization": {user2Token.Token},
			}))
			mustJsonDecode(httpResp, &resp)
			assert.Equal(t, http.StatusConflict, httpResp.StatusCode)
			assert.Equal(t, response.ErrorResp("DUPLICATE_USERNAME", nil), resp)
		})
	})
}

func TestChangePassword(t *testing.T) {
	t.Parallel()

	token := getNormalToken()

	t.Run("testChangePasswordWithoutParams", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/user/change_password", request.ChangePasswordRequest{
			OldPassword: "",
			NewPassword: "",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)
		jsonEQ(t, response.Response{
			Message: "VALIDATION_ERROR",
			Error: []map[string]string{
				{
					"field":  "OldPassword",
					"reason": "required",
				},
				{
					"field":  "NewPassword",
					"reason": "required",
				},
			},
			Data: nil,
		}, httpResp)
	})

	t.Run("testChangePasswordSuccess", func(t *testing.T) {
		t.Parallel()
		user := models.User{
			Username: "test_change_passwd_1",
			Nickname: "test_change_passwd_1_rand_str",
			Email:    "test_change_passwd_1@mail.com",
			Password: utils.HashPassword("test_change_passwd_old_passwd"),
		}
		assert.Nil(t, base.DB.Create(&user).Error)
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
		assert.Nil(t, base.DB.Create(&token1).Error)
		assert.Nil(t, base.DB.Create(&token2).Error)
		assert.Nil(t, base.DB.Create(&token3).Error)
		httpResp := makeResp(makeReq(t, "POST", "/api/user/change_password", request.ChangePasswordRequest{
			OldPassword: "test_change_passwd_old_passwd",
			NewPassword: "test_change_passwd_new_passwd",
		}, headerOption{
			"Authorization": {token1.Token},
		}))
		var tokens []models.Token
		var updatedUser models.User
		assert.Nil(t, base.DB.Preload("User").Where("user_id = ?", user.ID).Find(&tokens).Error)
		token1, _ = utils.GetToken(token1.Token)
		assert.Equal(t, []models.Token{
			token1,
		}, tokens)

		assert.Nil(t, base.DB.First(&updatedUser, user.ID).Error)
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
			Nickname: "test_change_passwd_2_rand_str",
			Email:    "test_change_passwd_2@mail.com",
			Password: utils.HashPassword("test_change_passwd_old_passwd"),
		}
		assert.Nil(t, base.DB.Create(&user).Error)
		mainToken := models.Token{
			Token: utils.RandStr(32),
			User:  user,
		}
		assert.Nil(t, base.DB.Create(&mainToken).Error)
		httpResp := makeResp(makeReq(t, "POST", "/api/user/change_password", request.ChangePasswordRequest{
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
