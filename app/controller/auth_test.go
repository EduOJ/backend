package controller_test

import (
	"encoding/json"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Parallel()
	// strip monotonic time
	t.Run("loginWithoutParams", func(t *testing.T) {
		t.Parallel()
		httpResp := MakeResp(MakeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: "",
			Password:        "",
		}))
		assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)
		JsonEQ(t, response.Response{
			Message: "VALIDATION_ERROR",
			Error: []map[string]string{
				{
					"field":  "UsernameOrEmail",
					"reason": "required",
				},
				{
					"field":  "Password",
					"reason": "required",
				},
			},
			Data: nil,
		}, httpResp)
	})
	t.Run("loginNotFound", func(t *testing.T) {
		t.Parallel()
		httpResp := MakeResp(MakeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: "test_login_1_not_found",
			Password:        "test_login_password",
		}))
		resp := response.LoginResponse{}
		MustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
		assert.Equal(t, "WRONG_USERNAME", resp.Message)
		assert.Equal(t, nil, resp.Error)
	})
	t.Run("loginWithUsernameSuccess", func(t *testing.T) {
		user1 := models.User{
			Username: "test_login_1",
			Nickname: "test_login_1_rand_str",
			Email:    "test_login_1@mail.com",
			Password: utils.HashPassword("test_login_password"),
		}
		base.DB.Create(&user1)
		t.Parallel()
		httpResp := MakeResp(MakeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: user1.Username,
			Password:        "test_login_password",
			RememberMe:      false,
		}))
		resp := response.LoginResponse{}
		MustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.Equal(t, "SUCCESS", resp.Message)
		assert.Equal(t, nil, resp.Error)
		JsonEQ(t, user1, resp.Data.User)
		token, err := utils.GetToken(resp.Data.Token)
		base.DB.Where("id = ?", user1.ID).First(&user1)
		assert.Equal(t, nil, err)
		assert.True(t, user1.UpdatedAt.Equal(token.User.UpdatedAt))
		assert.Equal(t, user1, token.User)
		assert.False(t, token.RememberMe)
	})
	t.Run("loginWithUsernameAndRememberMeSuccess", func(t *testing.T) {
		user2 := models.User{
			Username: "test_login_2",
			Nickname: "test_login_2_rand_str",
			Email:    "test_login_2@mail.com",
			Password: utils.HashPassword("test_login_password"),
		}
		base.DB.Create(&user2)
		t.Parallel()
		httpResp := MakeResp(MakeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: user2.Username,
			Password:        "test_login_password",
			RememberMe:      true,
		}))
		resp := response.LoginResponse{}
		MustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.Equal(t, "SUCCESS", resp.Message)
		assert.Equal(t, nil, resp.Error)
		JsonEQ(t, user2, resp.Data.User)
		token, err := utils.GetToken(resp.Data.Token)
		base.DB.Where("id = ?", user2.ID).First(&user2)
		assert.Equal(t, nil, err)
		assert.True(t, user2.UpdatedAt.Equal(token.User.UpdatedAt))
		assert.Equal(t, user2, token.User)
		assert.True(t, token.RememberMe)
	})
	t.Run("loginWithEmailSuccess", func(t *testing.T) {
		user3 := models.User{
			Username: "test_login_3",
			Nickname: "test_login_3_rand_str",
			Email:    "test_login_3@mail.com",
			Password: utils.HashPassword("test_login_password"),
		}
		base.DB.Create(&user3)
		t.Parallel()
		httpResp := MakeResp(MakeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: user3.Email,
			Password:        "test_login_password",
		}))
		resp := response.LoginResponse{}
		MustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.Equal(t, "SUCCESS", resp.Message)
		assert.Equal(t, nil, resp.Error)
		JsonEQ(t, user3, resp.Data.User)
		token, err := utils.GetToken(resp.Data.Token)
		base.DB.Where("id = ?", user3.ID).First(&user3)
		assert.Equal(t, nil, err)
		assert.Equal(t, user3, token.User)
	})
	t.Run("loginWrongPassword", func(t *testing.T) {
		user4 := models.User{
			Username: "test_login_4",
			Nickname: "test_login_4_rand_str",
			Email:    "test_login_4@mail.com",
			Password: utils.HashPassword("test_login_password"),
		}
		base.DB.Create(&user4)
		t.Parallel()
		httpResp := MakeResp(MakeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: user4.Email,
			Password:        "wrong_password",
		}))
		resp := response.LoginResponse{}
		MustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusForbidden, httpResp.StatusCode)
		assert.Equal(t, "WRONG_PASSWORD", resp.Message)
		assert.Equal(t, nil, resp.Error)
	})
}

func TestRegister(t *testing.T) {
	t.Run("registerUserWithoutParams", func(t *testing.T) {
		t.Parallel()
		resp := MakeResp(MakeReq(t, "POST", "/api/auth/register", request.RegisterRequest{
			Username: "",
			Nickname: "",
			Email:    "",
			Password: "",
		}))
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		JsonEQ(t, response.Response{
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
	t.Run("registerUserSuccess", func(t *testing.T) {
		t.Parallel()
		respResponse := MakeResp(MakeReq(t, "POST", "/api/auth/register", request.RegisterRequest{
			Username: "test_registerUserSuccess_0",
			Nickname: "test_registerUserSuccess_0",
			Email:    "test_registerUserSuccess_0@mail.com",
			Password: "test_registerUserSuccess_0",
		}))
		assert.Equal(t, http.StatusCreated, respResponse.StatusCode)
		resp := response.RegisterResponse{}
		respBytes, err := ioutil.ReadAll(respResponse.Body)
		assert.Equal(t, nil, err)
		err = json.Unmarshal(respBytes, &resp)
		assert.Equal(t, nil, err)
		user := models.User{}
		err = base.DB.Where("email = ?", "test_registerUserSuccess_0@mail.com").First(&user).Error
		assert.Equal(t, nil, err)
		token := models.Token{}
		err = base.DB.Where("token = ?", resp.Data.Token).Last(&token).Error
		assert.Equal(t, nil, err)
		assert.Equal(t, user.ID, token.ID)
		JsonEQ(t, resp.Data.User, user)

		respResponse = MakeResp(MakeReq(t, "POST", "/api/auth/register", request.RegisterRequest{
			Username: "test_registerUserSuccess_0",
			Nickname: "test_registerUserSuccess_0",
			Email:    "test_registerUserSuccess_0@mail.com",
			Password: "test_registerUserSuccess_0",
		}))
		JsonEQ(t, response.ErrorResp("DUPLICATE_EMAIL", nil), respResponse)
		assert.Equal(t, http.StatusBadRequest, respResponse.StatusCode)
		respResponse = MakeResp(MakeReq(t, "POST", "/api/auth/register", request.RegisterRequest{
			Username: "test_registerUserSuccess_0",
			Nickname: "test_registerUserSuccess_0",
			Email:    "test_registerUserSuccess_1@mail.com",
			Password: "test_registerUserSuccess_0",
		}))
		JsonEQ(t, response.ErrorResp("DUPLICATE_USERNAME", nil), respResponse)
		assert.Equal(t, http.StatusBadRequest, respResponse.StatusCode)
	})
}
