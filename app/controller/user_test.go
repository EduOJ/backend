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
	"testing"
)

func getToken(t *testing.T) (token models.Token) {
	token = models.Token{
		Token: utils.RandStr(32),
	}
	assert.Nil(t, base.DB.Create(&token).Error)
	return
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	token := getToken(t)

	t.Run("getUserNonExistId", func(t *testing.T) {
		t.Parallel()
		resp := makeResp(makeReq(t, "GET", "/api/user/10004", request.GetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		jsonEQ(t, response.ErrorResp("USER_NOT_FOUND", nil), resp)
	})
	t.Run("getUserNonExistUsername", func(t *testing.T) {
		t.Parallel()
		resp := makeResp(makeReq(t, "GET", "/api/user/test_get_non_existing_user", request.GetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		jsonEQ(t, response.ErrorResp("USER_NOT_FOUND", nil), resp)
	})
	t.Run("getUserSuccessWithId", func(t *testing.T) {
		user4 := models.User{
			Username: "test_get_user_4",
			Nickname: "test_get_user_4_rand_str",
			Email:    "test_get_user_4@mail.com",
			Password: utils.HashPassword("test_get_user_4_password"),
		}
		assert.Nil(t, base.DB.Create(&user4).Error)
		t.Parallel()
		resp := makeResp(makeReq(t, "GET", fmt.Sprintf("/api/user/%d", user4.ID), request.GetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		jsonEQ(t, response.GetUserResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*models.User `json:"user"`
			}{
				&user4,
			},
		}, resp)
	})
	t.Run("getUserSuccessWithUsername", func(t *testing.T) {
		user5 := models.User{
			Username: "test_get_user_5",
			Nickname: "test_get_user_5_rand_str",
			Email:    "test_get_user_5@mail.com",
			Password: utils.HashPassword("test_get_user_5_password"),
		}
		assert.Nil(t, base.DB.Create(&user5).Error)
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
				&user5,
			},
		}, resp)
	})
}

func TestGetUserMe(t *testing.T) {
	t.Parallel()

	t.Run("getUserSuccess", func(t *testing.T) {
		user1 := models.User{
			Username:   "test_get_user_me_1",
			Nickname:   "test_get_user_me_1_rand_str",
			Email:      "test_get_user_me_1@mail.com",
			Password:   utils.HashPassword("test_get_user_me_1_password"),
			RoleLoaded: true,
			Roles:      []models.UserHasRole{},
		}
		assert.Nil(t, base.DB.Create(&user1).Error)
		token := models.Token{
			Token: utils.RandStr(32),
			User:  user1,
		}
		assert.Nil(t, base.DB.Create(&token).Error)
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
				&user1,
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
		assert.Nil(t, base.DB.Create(&adminRole).Error)
		adminRole.AddPermission("all")
		user2 := models.User{
			Username:   "test_get_user_me_2",
			Nickname:   "test_get_user_me_2_rand_str",
			Email:      "test_get_user_me_2@mail.com",
			Password:   utils.HashPassword("test_get_user_me_2_password"),
			RoleLoaded: true,
			Roles:      []models.UserHasRole{},
		}
		assert.Nil(t, base.DB.Create(&user2).Error)
		user2.GrantRole(adminRole, classA)
		token := models.Token{
			Token: utils.RandStr(32),
			User:  user2,
		}
		assert.Nil(t, base.DB.Create(&token).Error)
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
				&user2,
			},
		}, resp)
	})
}


func TestChangePassword(t *testing.T) {
	t.Parallel()

	token := getToken(t)

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
