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
