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

type testClass struct {
	ID uint `gorm:"primary_key" json:"id"`
}

func (c testClass) TypeName() string {
	return "test_class"
}

func (c testClass) GetID() uint {
	return c.ID
}

func TestLogin(t *testing.T) {
	t.Parallel()
	// strip monotonic time
	t.Run("loginWithoutParams", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: "",
			Password:        "",
		}))
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)
		assert.Equal(t, response.Response{
			Message: "VALIDATION_ERROR",
			Error: []interface{}{
				map[string]interface{}{
					"field":       "UsernameOrEmail",
					"reason":      "required",
					"translation": "用户名为必填字段",
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
	t.Run("loginNotFound", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: "test_login_1_not_found",
			Password:        "test_login_password",
		}))
		resp := response.LoginResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
		assert.Equal(t, "NOT_FOUND", resp.Message)
		assert.Equal(t, nil, resp.Error)
	})
	t.Run("loginWithUsernameSuccess", func(t *testing.T) {
		user := models.User{
			Username: "test_login_1",
			Nickname: "test_login_1_rand_str",
			Email:    "test_login_1@mail.com",
			Password: utils.HashPassword("test_login_password"),
			Roles:    []models.UserHasRole{},
		}
		base.DB.Create(&user)
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: user.Username,
			Password:        "test_login_password",
			RememberMe:      false,
		}))
		resp := response.LoginResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.Equal(t, "SUCCESS", resp.Message)
		assert.Equal(t, nil, resp.Error)
		jsonEQ(t, user, resp.Data.User)
		token, err := utils.GetToken(resp.Data.Token)
		base.DB.Where("id = ?", user.ID).First(&user)
		assert.Equal(t, nil, err)
		token.User.LoadRoles()
		assert.True(t, user.UpdatedAt.Equal(token.User.UpdatedAt))
		jsonEQ(t, user, token.User)
		assert.False(t, token.RememberMe)
	})
	t.Run("loginWithUsernameAndRememberMeSuccess", func(t *testing.T) {
		user := models.User{
			Username: "test_login_2",
			Nickname: "test_login_2_rand_str",
			Email:    "test_login_2@mail.com",
			Password: utils.HashPassword("test_login_password"),
			Roles:    []models.UserHasRole{},
		}
		base.DB.Create(&user)
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: user.Username,
			Password:        "test_login_password",
			RememberMe:      true,
		}))
		resp := response.LoginResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.Equal(t, "SUCCESS", resp.Message)
		assert.Equal(t, nil, resp.Error)
		jsonEQ(t, user, resp.Data.User)
		token, err := utils.GetToken(resp.Data.Token)
		base.DB.Where("id = ?", user.ID).First(&user)
		assert.Equal(t, nil, err)
		token.User.LoadRoles()
		assert.True(t, user.UpdatedAt.Equal(token.User.UpdatedAt))
		jsonEQ(t, user, token.User)
		assert.True(t, token.RememberMe)
	})
	t.Run("loginWithEmailSuccess", func(t *testing.T) {
		user := models.User{
			Username: "test_login_3",
			Nickname: "test_login_3_rand_str",
			Email:    "test_login_3@mail.com",
			Password: utils.HashPassword("test_login_password"),
			Roles:    []models.UserHasRole{},
		}
		base.DB.Create(&user)
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: user.Email,
			Password:        "test_login_password",
		}))
		resp := response.LoginResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.Equal(t, "SUCCESS", resp.Message)
		assert.Equal(t, nil, resp.Error)
		jsonEQ(t, user, resp.Data.User)
		token, err := utils.GetToken(resp.Data.Token)
		base.DB.Where("id = ?", user.ID).First(&user)
		assert.Equal(t, nil, err)
		token.User.LoadRoles()
		jsonEQ(t, user, token.User)
	})
	t.Run("loginWrongPassword", func(t *testing.T) {
		user := models.User{
			Username: "test_login_4",
			Nickname: "test_login_4_rand_str",
			Email:    "test_login_4@mail.com",
			Password: utils.HashPassword("test_login_password"),
		}
		base.DB.Create(&user)
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: user.Email,
			Password:        "WRONG_PASSWORD",
		}))
		resp := response.LoginResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusForbidden, httpResp.StatusCode)
		assert.Equal(t, "WRONG_PASSWORD", resp.Message)
		assert.Equal(t, nil, resp.Error)
	})
	t.Run("loginWithUsernameAndRolesSuccess", func(t *testing.T) {
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
			Username: "test_login_5",
			Nickname: "test_login_5_rand_str",
			Email:    "test_login_5@mail.com",
			Password: utils.HashPassword("test_login_password"),
			Roles:    []models.UserHasRole{},
		}
		base.DB.Create(&user)
		user.GrantRole(adminRole, classA)

		httpResp := makeResp(makeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: user.Username,
			Password:        "test_login_password",
			RememberMe:      false,
		}))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)

		tokens := make([]models.Token, 0)
		base.DB.Model(&user).Related(&tokens)
		assert.Equal(t, 1, len(tokens))

		respUser := models.User{}
		base.DB.Model(&tokens[0]).Preload("Roles.Role.Permissions").Related(&respUser)
		base.DB.First(&user, user.ID)
		assert.Equal(t, user, respUser)
		jsonEQ(t, response.LoginResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				models.User `json:"user"`
				Token       string `json:"token"`
			}{
				respUser,
				tokens[0].Token,
			},
		}, httpResp)
		assert.False(t, tokens[0].RememberMe)
	})
}

func TestRegister(t *testing.T) {
	t.Run("registerUserWithoutParams", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/auth/register", request.RegisterRequest{
			Username: "",
			Nickname: "",
			Email:    "",
			Password: "",
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
	t.Run("registerUserSuccess", func(t *testing.T) {
		t.Parallel()
		httpResponse := makeResp(makeReq(t, "POST", "/api/auth/register", request.RegisterRequest{
			Username: "test_registerUserSuccess_0",
			Nickname: "test_registerUserSuccess_0",
			Email:    "test_registerUserSuccess_0@mail.com",
			Password: "test_registerUserSuccess_0",
		}))
		assert.Equal(t, http.StatusCreated, httpResponse.StatusCode)
		resp := response.RegisterResponse{}
		respBytes, err := ioutil.ReadAll(httpResponse.Body)
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
		jsonEQ(t, resp.Data.User, user)

		resp2 := response.Response{}
		httpResponse = makeResp(makeReq(t, "POST", "/api/auth/register", request.RegisterRequest{
			Username: "test_registerUserSuccess_0",
			Nickname: "test_registerUserSuccess_0",
			Email:    "test_registerUserSuccess_0@mail.com",
			Password: "test_registerUserSuccess_0",
		}))
		mustJsonDecode(httpResponse, &resp2)
		assert.Equal(t, http.StatusConflict, httpResponse.StatusCode)
		assert.Equal(t, response.ErrorResp("DUPLICATE_EMAIL", nil), resp2)
		httpResponse = makeResp(makeReq(t, "POST", "/api/auth/register", request.RegisterRequest{
			Username: "test_registerUserSuccess_0",
			Nickname: "test_registerUserSuccess_0",
			Email:    "test_registerUserSuccess_1@mail.com",
			Password: "test_registerUserSuccess_0",
		}))
		mustJsonDecode(httpResponse, &resp2)
		assert.Equal(t, http.StatusConflict, httpResponse.StatusCode)
		assert.Equal(t, response.ErrorResp("DUPLICATE_USERNAME", nil), resp2)

	})
}

func TestEmailRegistered(t *testing.T) {
	t.Parallel()
	user := models.User{
		Username: "test_email_registered_username",
		Nickname: "test_email_registered_nickname",
		Email:    "test_email_registered@e.com",
		Password: "test_email_registered_passwd",
	}
	base.DB.Create(&user)
	t.Run("testEmailRegisteredConflict", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "GET", "/api/auth/email_registered", request.EmailRegistered{
			Email: "test_email_registered@e.com",
		}))
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusConflict, httpResp.StatusCode)
		assert.Equal(t, response.Response{
			Message: "EMAIL_REGISTERED",
			Error:   nil,
			Data:    nil,
		}, resp)
	})
	t.Run("testEmailRegisteredSuccess", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "GET", "/api/auth/email_registered", request.EmailRegistered{
			Email: "test_email_registered_0@e.com",
		}))
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.Equal(t, response.Response{
			Message: "SUCCESS",
			Error:   nil,
			Data:    nil,
		}, resp)
	})
}
