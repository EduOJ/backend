package admin_test

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	adminRequest "github.com/leoleoasd/EduOJBackend/app/request/admin"
	"github.com/leoleoasd/EduOJBackend/app/response"
	adminResponse "github.com/leoleoasd/EduOJBackend/app/response/admin"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

type testClass struct {
	ID uint `gorm:"primary_key" json:"id"`
}

func (c testClass) TypeName() string {
	return "test_class"
}

func (c testClass) GetID() uint {
	return c.ID
}

func getToken(t *testing.T) (token models.Token) {
	token = models.Token{
		Token: utils.RandStr(32),
	}
	assert.Nil(t, base.DB.Create(&token).Error)
	return
}

func TestPostUser(t *testing.T) {
	t.Parallel()
	token := getToken(t)

	t.Run("postUserWithoutParams", func(t *testing.T) {
		t.Parallel()
		resp := makeResp(makeReq(t, "POST", "/api/admin/user", adminRequest.PostUserRequest{
			Username: "",
			Nickname: "",
			Email:    "",
			Password: "",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		jsonEQ(t, response.Response{
			Message: "VALIDATION_ERROR",
			Error: []map[string]string{
				{
					"field":  "Username",
					"reason": "required",
				},
				{
					"field":  "Nickname",
					"reason": "required",
				},
				{
					"field":  "Email",
					"reason": "required",
				},
				{
					"field":  "Password",
					"reason": "required",
				},
			},
			Data: nil,
		}, resp)
	})
	t.Run("postUserSuccess", func(t *testing.T) {
		t.Parallel()
		respResponse := makeResp(makeReq(t, "POST", "/api/admin/user", adminRequest.PostUserRequest{
			Username: "test_post_user_success_0",
			Nickname: "test_post_user_success_0",
			Email:    "test_post_user_success_0@mail.com",
			Password: "test_post_user_success_0",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusCreated, respResponse.StatusCode)
		resp := adminResponse.PostUserResponse{}
		respBytes, err := ioutil.ReadAll(respResponse.Body)
		assert.Equal(t, nil, err)
		err = json.Unmarshal(respBytes, &resp)
		assert.Equal(t, nil, err)
		user := models.User{}
		err = base.DB.Where("email = ?", "test_post_user_success_0@mail.com").First(&user).Error
		assert.Equal(t, nil, err)
		jsonEQ(t, resp.Data.User, user)

		respResponse = makeResp(makeReq(t, "POST", "/api/admin/user", adminRequest.PostUserRequest{
			Username: "test_post_user_success_0",
			Nickname: "test_post_user_success_0",
			Email:    "test_post_user_success_0@mail.com",
			Password: "test_post_user_success_0",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusBadRequest, respResponse.StatusCode)
		jsonEQ(t, response.ErrorResp("USER_DUPLICATE_EMAIL", nil), respResponse)
		respResponse = makeResp(makeReq(t, "POST", "/api/admin/user", adminRequest.PostUserRequest{
			Username: "test_post_user_success_0",
			Nickname: "test_post_user_success_0",
			Email:    "test_post_user_success_1@mail.com",
			Password: "test_post_user_success_0",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusBadRequest, respResponse.StatusCode)
		jsonEQ(t, response.ErrorResp("USER_DUPLICATE_USERNAME", nil), respResponse)
	})
}

func TestPutUser(t *testing.T) {
	t.Parallel()
	token := getToken(t)

	t.Run("putUserWithoutParams", func(t *testing.T) {
		user1 := models.User{
			Username: "test_put_user_1",
			Nickname: "test_put_user_1_rand_str",
			Email:    "test_put_user_1@mail.com",
			Password: utils.HashPassword("test_put_user_1_password"),
		}
		assert.Nil(t, base.DB.Create(&user1).Error)
		t.Parallel()
		resp := makeResp(makeReq(t, "PUT", fmt.Sprintf("/api/admin/user/%d", user1.ID), adminRequest.PutUserRequest{
			Username: "",
			Nickname: "",
			Email:    "",
			Password: "",
		}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		jsonEQ(t, response.Response{
			Message: "VALIDATION_ERROR",
			Error: []map[string]string{
				{
					"field":  "Username",
					"reason": "required",
				},
				{
					"field":  "Nickname",
					"reason": "required",
				},
				{
					"field":  "Email",
					"reason": "required",
				},
				{
					"field":  "Password",
					"reason": "required",
				},
			},
			Data: nil,
		}, resp)
	})
	t.Run("putUserNonExist", func(t *testing.T) {
		t.Run("deleteUserNonExistId", func(t *testing.T) {
			t.Parallel()
			resp := makeResp(makeReq(t, "PUT", "/api/admin/user/10001", adminRequest.PutUserRequest{
				Username: "test_put_user_non_exist",
				Nickname: "test_put_user_non_exist_nick",
				Email:    "test_put_user_non_exist@e.com",
				Password: "test_put_user_non_exist_passwd",
			}, headerOption{
				"Authorization": {token.Token},
			}))
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
			jsonEQ(t, response.ErrorResp("USER_NOT_FOUND", nil), resp)
		})
		t.Run("deleteUserNonExistUsername", func(t *testing.T) {
			t.Parallel()
			resp := makeResp(makeReq(t, "PUT", "/api/admin/user/test_put_non_existing_user", adminRequest.PutUserRequest{
				Username: "test_put_user_non_exist",
				Nickname: "test_put_user_non_exist_nick",
				Email:    "test_put_user_non_exist@e.com",
				Password: "test_put_user_non_exist_passwd",
			}, headerOption{
				"Authorization": {token.Token},
			}))
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
			jsonEQ(t, response.ErrorResp("USER_NOT_FOUND", nil), resp)
		})
	})
	t.Run("putUserWithParams", func(t *testing.T) {
		user2 := models.User{
			Username: "test_put_user_2",
			Nickname: "test_put_user_2_rand_str",
			Email:    "test_put_user_2@mail.com",
			Password: utils.HashPassword("test_put_user_2_password"),
		}
		assert.Nil(t, base.DB.Create(&user2).Error)
		user3 := models.User{
			Username: "test_put_user_3",
			Nickname: "test_put_user_3_rand_str",
			Email:    "test_put_user_3@mail.com",
			Password: utils.HashPassword("test_put_user_3_password"),
		}
		assert.Nil(t, base.DB.Create(&user3).Error)
		t.Run("putUserSuccessWithId", func(t *testing.T) {
			user4 := models.User{
				Username: "test_put_user_4",
				Nickname: "test_put_user_4_rand_str",
				Email:    "test_put_user_4@mail.com",
				Password: utils.HashPassword("test_put_user_4_password"),
			}
			assert.Nil(t, base.DB.Create(&user4).Error)
			t.Parallel()
			respResponse := makeResp(makeReq(t, "PUT", fmt.Sprintf("/api/admin/user/%d", user4.ID), adminRequest.PutUserRequest{
				Username: "test_putUserSuccess_0",
				Nickname: "test_putUserSuccess_0",
				Email:    "test_putUserSuccess_0@mail.com",
				Password: "test_putUserSuccess_0",
			}, headerOption{
				"Authorization": {token.Token},
			}))
			assert.Equal(t, http.StatusOK, respResponse.StatusCode)
			resp := adminResponse.PutUserResponse{}
			respBytes, err := ioutil.ReadAll(respResponse.Body)
			assert.Equal(t, nil, err)
			err = json.Unmarshal(respBytes, &resp)
			assert.Equal(t, nil, err)
			user := models.User{}
			err = base.DB.Where("id = ?", user4.ID).First(&user).Error
			assert.Equal(t, nil, err)
			jsonEQ(t, resp.Data.User, user)
		})
		t.Run("putUserSuccessWithUsername", func(t *testing.T) {
			user5 := models.User{
				Username: "test_put_user_5",
				Nickname: "test_put_user_5_rand_str",
				Email:    "test_put_user_5@mail.com",
				Password: utils.HashPassword("test_put_user_5_password"),
			}
			assert.Nil(t, base.DB.Create(&user5).Error)
			t.Parallel()
			respResponse := makeResp(makeReq(t, "PUT", "/api/admin/user/test_put_user_5", adminRequest.PutUserRequest{
				Username: "test_putUserSuccess_1",
				Nickname: "test_putUserSuccess_1",
				Email:    "test_putUserSuccess_1@mail.com",
				Password: "test_putUserSuccess_1",
			}, headerOption{
				"Authorization": {token.Token},
			}))
			assert.Equal(t, http.StatusOK, respResponse.StatusCode)
			resp := adminResponse.PutUserResponse{}
			respBytes, err := ioutil.ReadAll(respResponse.Body)
			assert.Equal(t, nil, err)
			err = json.Unmarshal(respBytes, &resp)
			assert.Equal(t, nil, err)
			user := models.User{}
			err = base.DB.Where("id = ?", user5.ID).First(&user).Error
			assert.Equal(t, nil, err)
			jsonEQ(t, resp.Data.User, user)
		})
		t.Run("putUserDuplicateEmail", func(t *testing.T) {
			t.Parallel()
			resp := makeResp(makeReq(t, "PUT", fmt.Sprintf("/api/admin/user/%d", user2.ID), adminRequest.PutUserRequest{
				Username: "test_put_user_2",
				Nickname: "test_put_user_2_rand_str",
				Email:    "test_put_user_3@mail.com",
				Password: "test_put_user_2_password",
			}, headerOption{
				"Authorization": {token.Token},
			}))
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			jsonEQ(t, response.ErrorResp("USER_DUPLICATE_EMAIL", nil), resp)
		})
		t.Run("putUserDuplicateUsername", func(t *testing.T) {
			t.Parallel()
			resp := makeResp(makeReq(t, "PUT", fmt.Sprintf("/api/admin/user/%d", user2.ID), adminRequest.PutUserRequest{
				Username: "test_put_user_3",
				Nickname: "test_put_user_2_rand_str",
				Email:    "test_put_user_2@mail.com",
				Password: "test_put_user_2_password",
			}, headerOption{
				"Authorization": {token.Token},
			}))
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			jsonEQ(t, response.ErrorResp("USER_DUPLICATE_USERNAME", nil), resp)
		})
	})
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	token := getToken(t)
	t.Run("deleteUserNonExistId", func(t *testing.T) {
		t.Parallel()
		resp := makeResp(makeReq(t, "DELETE", "/api/admin/user/10002", adminRequest.DeleteUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		jsonEQ(t, response.ErrorResp("USER_NOT_FOUND", nil), resp)
	})
	t.Run("deleteUserNonExistUsername", func(t *testing.T) {
		t.Parallel()
		resp := makeResp(makeReq(t, "DELETE", "/api/admin/user/test_delete_non_existing_user", adminRequest.DeleteUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		jsonEQ(t, response.ErrorResp("USER_NOT_FOUND", nil), resp)
	})
	t.Run("deleteUserSuccessWithId", func(t *testing.T) {
		user1 := models.User{
			Username: "test_delete_user_1",
			Nickname: "test_delete_user_1_rand_str",
			Email:    "test_delete_user_1@mail.com",
			Password: utils.HashPassword("test_delete_user_1_password"),
		}
		assert.Nil(t, base.DB.Create(&user1).Error)
		t.Parallel()
		resp := makeResp(makeReq(t, "DELETE", fmt.Sprintf("/api/admin/user/%d", user1.ID), adminRequest.DeleteUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		jsonEQ(t, response.Response{
			Message: "SUCCESS",
			Error:   nil,
			Data:    nil,
		}, resp)
		user := models.User{}
		err := base.DB.Where("id = ?", user1.ID).First(&user).Error
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
	t.Run("deleteUserSuccessWithUsername", func(t *testing.T) {
		user2 := models.User{
			Username: "test_delete_user_2",
			Nickname: "test_delete_user_2_rand_str",
			Email:    "test_delete_user_2@mail.com",
			Password: utils.HashPassword("test_delete_user_2_password"),
		}
		assert.Nil(t, base.DB.Create(&user2).Error)
		t.Parallel()
		resp := makeResp(makeReq(t, "DELETE", "/api/admin/user/test_delete_user_2", adminRequest.DeleteUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		jsonEQ(t, response.Response{
			Message: "SUCCESS",
			Error:   nil,
			Data:    nil,
		}, resp)
		user := models.User{}
		err := base.DB.Where("id = ?", user2.ID).First(&user).Error
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
func TestGetUser(t *testing.T) {
	t.Parallel()
	token := getToken(t)
	t.Run("getUserNonExistId", func(t *testing.T) {
		t.Parallel()
		resp := makeResp(makeReq(t, "GET", "/api/admin/user/10003", adminRequest.GetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		jsonEQ(t, response.ErrorResp("USER_NOT_FOUND", nil), resp)
	})
	t.Run("getUserNonExistUsername", func(t *testing.T) {
		t.Parallel()
		resp := makeResp(makeReq(t, "GET", "/api/admin/user/test_get_non_existing_user", adminRequest.GetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		jsonEQ(t, response.ErrorResp("USER_NOT_FOUND", nil), resp)
	})
	t.Run("getUserSuccessWithId", func(t *testing.T) {
		user1 := models.User{
			Username:   "test_get_user_1",
			Nickname:   "test_get_user_1_rand_str",
			Email:      "test_get_user_1@mail.com",
			Password:   utils.HashPassword("test_get_user_1_password"),
			RoleLoaded: true,
			Roles:      []models.UserHasRole{},
		}
		assert.Nil(t, base.DB.Create(&user1).Error)
		t.Parallel()
		resp := makeResp(makeReq(t, "GET", fmt.Sprintf("/api/admin/user/%d", user1.ID), adminRequest.GetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		jsonEQ(t, adminResponse.GetUserResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*models.User `json:"user"`
			}{
				&user1,
			},
		}, resp)
	})
	t.Run("getUserSuccessWithUsername", func(t *testing.T) {
		user2 := models.User{
			Username:   "test_get_user_2",
			Nickname:   "test_get_user_2_rand_str",
			Email:      "test_get_user_2@mail.com",
			Password:   utils.HashPassword("test_get_user_2_password"),
			RoleLoaded: true,
			Roles:      []models.UserHasRole{},
		}
		assert.Nil(t, base.DB.Create(&user2).Error)
		t.Parallel()
		resp := makeResp(makeReq(t, "GET", "/api/admin/user/test_get_user_2", adminRequest.GetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		jsonEQ(t, adminResponse.GetUserResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*models.User `json:"user"`
			}{
				&user2,
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
		user3 := models.User{
			Username:   "test_get_user_3",
			Nickname:   "test_get_user_3_rand_str",
			Email:      "test_get_user_3@mail.com",
			Password:   utils.HashPassword("test_get_user_3_password"),
			RoleLoaded: true,
			Roles:      []models.UserHasRole{},
		}
		assert.Nil(t, base.DB.Create(&user3).Error)
		user3.GrantRole(adminRole, classA)
		t.Parallel()
		resp := makeResp(makeReq(t, "GET", fmt.Sprintf("/api/admin/user/%d", user3.ID), adminRequest.GetUserRequest{}, headerOption{
			"Authorization": {token.Token},
		}))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		jsonEQ(t, adminResponse.GetUserResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*models.User `json:"user"`
			}{
				&user3,
			},
		}, resp)
	})
}
