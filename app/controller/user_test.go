package controller_test

import (
	"fmt"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
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
		t.Parallel()
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
		t.Parallel()
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