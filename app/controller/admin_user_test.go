package controller_test

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"testing"
)

var initAdminUser sync.Once
var adminUser models.User

func initAdminUserFunc() {
	adminRole := models.Role{
		Name:   "globalAdmin",
		Target: nil,
	}
	base.DB.Create(&adminRole)
	adminRole.AddPermission("create_user")
	adminRole.AddPermission("update_user")
	adminRole.AddPermission("delete_user")
	adminUser = models.User{
		Username: "test_user_admin_user",
		Nickname: "test_user_admin_nickname",
		Email:    "test_user_admin@mail.com",
		Password: "test_user_admin_password",
	}
	base.DB.Create(&adminUser)
	adminUser.GrantRole(adminRole)
}

func getAdminToken() (token models.Token) {
	initAdminUser.Do(initAdminUserFunc)
	token = models.Token{
		User:  adminUser,
		Token: utils.RandStr(32),
	}
	base.DB.Create(&token)
	return
}

func getUrlStringPointer(url string, paras map[string]string) *string {
	s := fmt.Sprintf("%s?", url)
	keys := make([]string, 0, len(paras))
	for key := range paras {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for index, key := range keys {
		if index != 0 {
			s += "&"
		}
		s += fmt.Sprintf("%s=%s", key, paras[key])
	}
	return &s
}

func TestAdminCreateUser(t *testing.T) {
	t.Parallel()
	token := getAdminToken()

	t.Run("testAdminCreateUserWithoutParams", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/admin/user", request.AdminCreateUserRequest{
			Username: "",
			Nickname: "",
			Email:    "",
			Password: "",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
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
				map[string]interface{}{
					"field":       "Password",
					"reason":      "required",
					"translation": "密码为必填字段",
				},
			},
			Data: nil,
		}, resp)
	})
	t.Run("testAdminCreateUserSuccess", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/admin/user", request.AdminCreateUserRequest{
			Username: "test_post_user_success_0",
			Nickname: "test_post_user_success_0",
			Email:    "test_post_user_success_0@mail.com",
			Password: "test_post_user_success_0",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusCreated, httpResp.StatusCode)
		resp := response.AdminCreateUserResponse{}
		respBytes, err := ioutil.ReadAll(httpResp.Body)
		assert.Equal(t, nil, err)
		mustJsonDecode(respBytes, &resp)
		user := models.User{}
		err = base.DB.Where("email = ?", "test_post_user_success_0@mail.com").First(&user).Error
		assert.Equal(t, nil, err)
		jsonEQ(t, resp, response.AdminCreateUserResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*models.User `json:"user"`
			}{
				&user,
			},
		})

		resp2 := response.Response{}
		httpResp = makeResp(makeReq(t, "POST", "/api/admin/user", request.AdminCreateUserRequest{
			Username: "test_post_user_success_0",
			Nickname: "test_post_user_success_0",
			Email:    "test_post_user_success_0@mail.com",
			Password: "test_post_user_success_0",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		mustJsonDecode(httpResp, &resp2)
		assert.Equal(t, http.StatusConflict, httpResp.StatusCode)
		assert.Equal(t, response.ErrorResp("CONFLICT_EMAIL", nil), resp2)
		httpResp = makeResp(makeReq(t, "POST", "/api/admin/user", request.AdminCreateUserRequest{
			Username: "test_post_user_success_0",
			Nickname: "test_post_user_success_0",
			Email:    "test_post_user_success_1@mail.com",
			Password: "test_post_user_success_0",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		mustJsonDecode(httpResp, &resp2)
		assert.Equal(t, http.StatusConflict, httpResp.StatusCode)
		assert.Equal(t, response.ErrorResp("CONFLICT_USERNAME", nil), resp2)
	})
}

func TestAdminUpdateUser(t *testing.T) {
	t.Parallel()
	token := getAdminToken()

	t.Run("testAdminUpdateUserWithoutParams", func(t *testing.T) {
		t.Parallel()
		user := models.User{
			Username: "test_put_user_1",
			Nickname: "test_put_user_1_rand_str",
			Email:    "test_put_user_1@mail.com",
			Password: utils.HashPassword("test_put_user_1_password"),
		}
		base.DB.Create(&user)
		httpResp := makeResp(makeReq(t, "PUT", fmt.Sprintf("/api/admin/user/%d", user.ID), request.AdminUpdateUserRequest{
			Username: "",
			Nickname: "",
			Email:    "",
			Password: "",
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
				map[string]interface{}{
					"field":       "Password",
					"reason":      "required",
					"translation": "密码为必填字段",
				},
			},
			Data: nil,
		}, resp)
	})
	t.Run("testAdminUpdateUserNonExist", func(t *testing.T) {
		t.Parallel()
		t.Run("testAdminUpdateUserNonExistId", func(t *testing.T) {
			t.Parallel()
			resp := response.Response{}
			httpResp := makeResp(makeReq(t, "PUT", "/api/admin/user/-1", request.AdminUpdateUserRequest{
				Username: "test_put_user_non_exist",
				Nickname: "test_put_user_non_exist_nick",
				Email:    "test_put_user_non_exist@e.com",
				Password: "test_put_user_non_exist_passwd",
			}, headerOption{
				"Authorization": {token.Token},
			}))
			mustJsonDecode(httpResp, &resp)
			assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
			assert.Equal(t, response.ErrorResp("NOT_FOUND", nil), resp)
		})
		t.Run("testAdminUpdateUserNonExistUsername", func(t *testing.T) {
			t.Parallel()
			resp := response.Response{}
			httpResp := makeResp(makeReq(t, "PUT", "/api/admin/user/test_put_non_existing_user", request.AdminUpdateUserRequest{
				Username: "test_put_user_non_exist",
				Nickname: "test_put_user_non_exist_nick",
				Email:    "test_put_user_non_exist@e.com",
				Password: "test_put_user_non_exist_passwd",
			}, headerOption{
				"Authorization": {token.Token},
			}))
			mustJsonDecode(httpResp, &resp)
			assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
			assert.Equal(t, response.ErrorResp("NOT_FOUND", nil), resp)
		})
	})
	t.Run("testAdminUpdateUserWithParams", func(t *testing.T) {
		t.Parallel()
		user2 := models.User{
			Username: "test_put_user_2",
			Nickname: "test_put_user_2_rand_str",
			Email:    "test_put_user_2@mail.com",
			Password: utils.HashPassword("test_put_user_2_password"),
		}
		base.DB.Create(&user2)
		user3 := models.User{
			Username: "test_put_user_3",
			Nickname: "test_put_user_3_rand_str",
			Email:    "test_put_user_3@mail.com",
			Password: utils.HashPassword("test_put_user_3_password"),
		}
		base.DB.Create(&user3)
		t.Run("testAdminUpdateUserSuccessWithId", func(t *testing.T) {
			t.Parallel()
			user := models.User{
				Username: "test_put_user_4",
				Nickname: "test_put_user_4_rand_str",
				Email:    "test_put_user_4@mail.com",
				Password: utils.HashPassword("test_put_user_4_password"),
			}
			base.DB.Create(&user)
			respResponse := makeResp(makeReq(t, "PUT", fmt.Sprintf("/api/admin/user/%d", user.ID), request.AdminUpdateUserRequest{
				Username: "test_putUserSuccess_0",
				Nickname: "test_putUserSuccess_0",
				Email:    "test_putUserSuccess_0@mail.com",
				Password: "test_putUserSuccess_0",
			}, headerOption{
				"Authorization": {token.Token},
			}))
			assert.Equal(t, http.StatusOK, respResponse.StatusCode)
			resp := response.AdminUpdateUserResponse{}
			respBytes, err := ioutil.ReadAll(respResponse.Body)
			assert.Equal(t, nil, err)
			err = json.Unmarshal(respBytes, &resp)
			assert.Equal(t, nil, err)
			databaseUser := models.User{}
			err = base.DB.Where("id = ?", user.ID).First(&databaseUser).Error
			assert.Equal(t, nil, err)
			jsonEQ(t, resp.Data.User, databaseUser)
		})
		t.Run("testAdminUpdateUserSuccessWithUsername", func(t *testing.T) {
			t.Parallel()
			user := models.User{
				Username: "test_put_user_5",
				Nickname: "test_put_user_5_rand_str",
				Email:    "test_put_user_5@mail.com",
				Password: utils.HashPassword("test_put_user_5_password"),
			}
			base.DB.Create(&user)
			respResponse := makeResp(makeReq(t, "PUT", "/api/admin/user/test_put_user_5", request.AdminUpdateUserRequest{
				Username: "test_putUserSuccess_1",
				Nickname: "test_putUserSuccess_1",
				Email:    "test_putUserSuccess_1@mail.com",
				Password: "test_putUserSuccess_1",
			}, headerOption{
				"Authorization": {token.Token},
			}))
			assert.Equal(t, http.StatusOK, respResponse.StatusCode)
			resp := response.AdminUpdateUserResponse{}
			respBytes, err := ioutil.ReadAll(respResponse.Body)
			assert.Equal(t, nil, err)
			err = json.Unmarshal(respBytes, &resp)
			assert.Equal(t, nil, err)
			databaseUser := models.User{}
			err = base.DB.Where("id = ?", user.ID).First(&databaseUser).Error
			assert.Equal(t, nil, err)
			jsonEQ(t, resp.Data.User, databaseUser)
		})
		t.Run("testAdminUpdateUserDuplicateEmail", func(t *testing.T) {
			t.Parallel()
			resp := response.Response{}
			httpResp := makeResp(makeReq(t, "PUT", fmt.Sprintf("/api/admin/user/%d", user2.ID), request.AdminUpdateUserRequest{
				Username: "test_put_user_2",
				Nickname: "test_put_user_2_rand_str",
				Email:    "test_put_user_3@mail.com",
				Password: "test_put_user_2_password",
			}, headerOption{
				"Authorization": {token.Token},
			}))
			mustJsonDecode(httpResp, &resp)
			assert.Equal(t, http.StatusConflict, httpResp.StatusCode)
			assert.Equal(t, response.ErrorResp("CONFLICT_EMAIL", nil), resp)
		})
		t.Run("testAdminUpdateUserDuplicateUsername", func(t *testing.T) {
			t.Parallel()
			resp := response.Response{}
			httpResp := makeResp(makeReq(t, "PUT", fmt.Sprintf("/api/admin/user/%d", user2.ID), request.AdminUpdateUserRequest{
				Username: "test_put_user_3",
				Nickname: "test_put_user_2_rand_str",
				Email:    "test_put_user_2@mail.com",
				Password: "test_put_user_2_password",
			}, headerOption{
				"Authorization": {token.Token},
			}))
			mustJsonDecode(httpResp, &resp)
			assert.Equal(t, http.StatusConflict, httpResp.StatusCode)
			assert.Equal(t, response.ErrorResp("CONFLICT_USERNAME", nil), resp)
		})
	})
}

func TestAdminDeleteUser(t *testing.T) {
	t.Parallel()
	token := getAdminToken()
	t.Run("deleteUserNonExistId", func(t *testing.T) {
		t.Parallel()
		resp := response.Response{}
		httpResp := makeResp(makeReq(t, "DELETE", "/api/admin/user/-1", request.AdminDeleteUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
		assert.Equal(t, response.ErrorResp("NOT_FOUND", nil), resp)
	})
	t.Run("deleteUserNonExistUsername", func(t *testing.T) {
		t.Parallel()
		resp := response.Response{}
		httpResp := makeResp(makeReq(t, "DELETE", "/api/admin/user/test_delete_non_existing_user", request.AdminDeleteUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
		assert.Equal(t, response.ErrorResp("NOT_FOUND", nil), resp)
	})
	t.Run("deleteUserSuccessWithId", func(t *testing.T) {
		t.Parallel()
		user := models.User{
			Username: "test_delete_user_1",
			Nickname: "test_delete_user_1_rand_str",
			Email:    "test_delete_user_1@mail.com",
			Password: utils.HashPassword("test_delete_user_1_password"),
		}
		base.DB.Create(&user)
		resp := response.Response{}
		httpResp := makeResp(makeReq(t, "DELETE", fmt.Sprintf("/api/admin/user/%d", user.ID), request.AdminDeleteUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusNoContent, httpResp.StatusCode)
		assert.Equal(t, response.Response{
			Message: "SUCCESS",
			Error:   nil,
			Data:    nil,
		}, resp)
		databaseUser := models.User{}
		err := base.DB.Where("id = ?", user.ID).First(&databaseUser).Error
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
	t.Run("deleteUserSuccessWithUsername", func(t *testing.T) {
		t.Parallel()
		user := models.User{
			Username: "test_delete_user_2",
			Nickname: "test_delete_user_2_rand_str",
			Email:    "test_delete_user_2@mail.com",
			Password: utils.HashPassword("test_delete_user_2_password"),
		}
		base.DB.Create(&user)
		resp := response.Response{}
		httpResp := makeResp(makeReq(t, "DELETE", "/api/admin/user/test_delete_user_2", request.AdminDeleteUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusNoContent, httpResp.StatusCode)
		assert.Equal(t, response.Response{
			Message: "SUCCESS",
			Error:   nil,
			Data:    nil,
		}, resp)
		databaseUser := models.User{}
		err := base.DB.Where("id = ?", user.ID).First(&databaseUser).Error
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
func TestAdminGetUser(t *testing.T) {
	t.Parallel()
	token := getAdminToken()
	t.Run("getUserNonExistId", func(t *testing.T) {
		t.Parallel()
		resp := response.Response{}
		httpResp := makeResp(makeReq(t, "GET", "/api/admin/user/-1", request.AdminGetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
		assert.Equal(t, response.ErrorResp("NOT_FOUND", nil), resp)
	})
	t.Run("getUserNonExistUsername", func(t *testing.T) {
		t.Parallel()
		resp := response.Response{}
		httpResp := makeResp(makeReq(t, "GET", "/api/admin/user/test_get_non_existing_user", request.AdminGetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
		assert.Equal(t, response.ErrorResp("NOT_FOUND", nil), resp)
	})
	t.Run("getUserSuccessWithId", func(t *testing.T) {
		t.Parallel()
		user := models.User{
			Username: "test_get_user_1",
			Nickname: "test_get_user_1_rand_str",
			Email:    "test_get_user_1@mail.com",
			Password: utils.HashPassword("test_get_user_1_password"),
			Roles:    []models.UserHasRole{},
		}
		base.DB.Create(&user)
		resp := makeResp(makeReq(t, "GET", fmt.Sprintf("/api/admin/user/%d", user.ID), request.AdminGetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		jsonEQ(t, response.AdminGetUserResponse{
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
		t.Parallel()
		user := models.User{
			Username: "test_get_user_2",
			Nickname: "test_get_user_2_rand_str",
			Email:    "test_get_user_2@mail.com",
			Password: utils.HashPassword("test_get_user_2_password"),
			Roles:    []models.UserHasRole{},
		}
		base.DB.Create(&user)
		resp := makeResp(makeReq(t, "GET", "/api/admin/user/test_get_user_2", request.AdminGetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		jsonEQ(t, response.AdminGetUserResponse{
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
			Username: "test_get_user_3",
			Nickname: "test_get_user_3_rand_str",
			Email:    "test_get_user_3@mail.com",
			Password: utils.HashPassword("test_get_user_3_password"),
			Roles:    []models.UserHasRole{},
		}
		base.DB.Create(&user)
		user.GrantRole(adminRole, classA)
		resp := makeResp(makeReq(t, "GET", fmt.Sprintf("/api/admin/user/%d", user.ID), request.AdminGetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		jsonEQ(t, response.AdminGetUserResponse{
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

	token := getAdminToken()
	baseUrl := "/api/admin/users"

	t.Run("testAdminGetUsersSuccess", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
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
					Total: 4,
					Count: 4,
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
					Total: 1,
					Count: 1,
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
					Total: 1,
					Count: 1,
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
					Total: 1,
					Count: 1,
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
					Total: 2,
					Count: 2,
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
					Total: 2,
					Count: 2,
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
					Total: 3,
					Count: 3,
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
					Total: 4,
					Count: 2,
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
					Total: 4,
					Count: 2,
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
					Total: 4,
					Count: 4,
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
					Total: 4,
					Count: 4,
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
					Total: 4,
					Count: 4,
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
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				test := test
				t.Parallel()
				httpResp := makeResp(makeReq(t, "GET", "/api/admin/users", test.req, headerOption{
					"Authorization": {token.Token},
				}))
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
	t.Run("testAdminGetUsersWithWrongOrderByPara", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "GET", "/api/admin/users", request.AdminGetUsersRequest{
			Search:  "test_admin_get_users",
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
	t.Run("testAdminGetUsersOrderByNonExistingColumn", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "GET", "/api/admin/users", request.AdminGetUsersRequest{
			Search:  "test_admin_get_users",
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
	t.Run("testAdminGetUsersOrderByNonExistingOrder", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "GET", "/api/admin/users", request.AdminGetUsersRequest{
			Search:  "test_admin_get_users",
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
	t.Run("testAdminGetUsersDefaultLimit", func(t *testing.T) {
		t.Parallel()
		// DL: default limit
		users := make([]models.User, 25)
		for i := 0; i < 25; i++ {
			users[i] = models.User{
				Username: fmt.Sprintf("test_DL_admin_get_users_%d", i),
				Nickname: fmt.Sprintf("test_DL_admin_get_users_n_%d", i),
				Email:    fmt.Sprintf("test_DL_admin_get_users_%d@e.e", i),
				Password: fmt.Sprintf("test_DL_admin_get_users_pwd_%d", i),
			}
			base.DB.Create(&users[i])
		}
		httpResp := makeResp(makeReq(t, "GET", "/api/admin/users", request.AdminGetUsersRequest{
			Search: "test_DL_admin_get_users_",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		jsonEQ(t, response.AdminGetUsersResponse{
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
