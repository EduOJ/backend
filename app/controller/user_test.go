package controller_test

import (
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
		user1 := models.User{
			Username: "test_change_passwd_1",
			Nickname: "test_change_passwd_1_rand_str",
			Email:    "test_change_passwd_1@mail.com",
			Password: utils.HashPassword("test_change_passwd_old_passwd"),
		}
		assert.Nil(t, base.DB.Create(&user1).Error)
		user2 := models.User{
			Username: "test_change_passwd_2",
			Nickname: "test_change_passwd_2_rand_str",
			Email:    "test_change_passwd_2@mail.com",
			Password: utils.HashPassword("test_change_passwd_old_passwd"),
		}
		assert.Nil(t, base.DB.Create(&user2).Error)
		mainToken := models.Token{
			Token: utils.RandStr(32),
			User:  user1,
		}
		assert.Nil(t, base.DB.Create(&mainToken).Error)
		otherToken1 := models.Token{
			Token: utils.RandStr(32),
			User:  user1,
		}
		assert.Nil(t, base.DB.Create(&otherToken1).Error)
		otherToken2 := models.Token{
			Token: utils.RandStr(32),
			User:  user1,
		}
		assert.Nil(t, base.DB.Create(&otherToken2).Error)
		otherUserToken := models.Token{
			Token: utils.RandStr(32),
			User:  user2,
		}
		assert.Nil(t, base.DB.Create(&otherUserToken).Error)
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/user/change_password", request.ChangePasswordRequest{
			OldPassword: "test_change_passwd_old_passwd",
			NewPassword: "test_change_passwd_new_passwd",
		}, headerOption{
			"Authorization": {mainToken.Token},
		}))
		var tokens []models.Token
		var updatedUser models.User
		assert.Nil(t, base.DB.Preload("User").Where("user_id = ?", user1.ID).Find(&tokens).Error)
		mainToken, _ = utils.GetToken(mainToken.Token)
		assert.Nil(t, base.DB.First(&updatedUser, user1.ID).Error)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		jsonEQ(t, response.Response{
			Message: "SUCCESS",
			Error:   nil,
			Data:    nil,
		}, httpResp)
		assert.Equal(t, []models.Token{
			mainToken,
		}, tokens)
		assert.True(t, utils.VerifyPassword("test_change_passwd_new_passwd", updatedUser.Password))
	})

	t.Run("testChangePasswordWithWrongPassword", func(t *testing.T) {
		user3 := models.User{
			Username: "test_change_passwd_3",
			Nickname: "test_change_passwd_3_rand_str",
			Email:    "test_change_passwd_3@mail.com",
			Password: utils.HashPassword("test_change_passwd_old_passwd"),
		}
		assert.Nil(t, base.DB.Create(&user3).Error)
		mainToken := models.Token{
			Token: utils.RandStr(32),
			User:  user3,
		}
		assert.Nil(t, base.DB.Create(&mainToken).Error)
		t.Parallel()
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
