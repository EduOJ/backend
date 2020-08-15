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
	"time"
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
		assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)
		jsonEQ(t, response.Response{
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
		httpResp := makeResp(makeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: "test_login_1_not_found",
			Password:        "test_login_password",
		}))
		resp := response.LoginResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
		assert.Equal(t, "WRONG_USERNAME", resp.Message)
		assert.Equal(t, nil, resp.Error)
	})
	t.Run("loginWithUsernameSuccess", func(t *testing.T) {
		user1 := models.User{
			Username:   "test_login_1",
			Nickname:   "test_login_1_rand_str",
			Email:      "test_login_1@mail.com",
			Password:   utils.HashPassword("test_login_password"),
			RoleLoaded: true,
			Roles:      []models.UserHasRole{},
		}
		base.DB.Create(&user1)
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: user1.Username,
			Password:        "test_login_password",
			RememberMe:      false,
		}))
		resp := response.LoginResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.Equal(t, "SUCCESS", resp.Message)
		assert.Equal(t, nil, resp.Error)
		jsonEQ(t, user1, resp.Data.User)
		token, err := utils.GetToken(resp.Data.Token)
		base.DB.Where("id = ?", user1.ID).First(&user1)
		assert.Equal(t, nil, err)
		token.User.LoadRoles()
		assert.True(t, user1.UpdatedAt.Equal(token.User.UpdatedAt))
		assert.Equal(t, user1, token.User)
		assert.False(t, token.RememberMe)
	})
	t.Run("loginWithUsernameAndRememberMeSuccess", func(t *testing.T) {
		user2 := models.User{
			Username:   "test_login_2",
			Nickname:   "test_login_2_rand_str",
			Email:      "test_login_2@mail.com",
			Password:   utils.HashPassword("test_login_password"),
			RoleLoaded: true,
			Roles:      []models.UserHasRole{},
		}
		base.DB.Create(&user2)
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: user2.Username,
			Password:        "test_login_password",
			RememberMe:      true,
		}))
		resp := response.LoginResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.Equal(t, "SUCCESS", resp.Message)
		assert.Equal(t, nil, resp.Error)
		jsonEQ(t, user2, resp.Data.User)
		token, err := utils.GetToken(resp.Data.Token)
		base.DB.Where("id = ?", user2.ID).First(&user2)
		assert.Equal(t, nil, err)
		token.User.LoadRoles()
		assert.True(t, user2.UpdatedAt.Equal(token.User.UpdatedAt))
		assert.Equal(t, user2, token.User)
		assert.True(t, token.RememberMe)
	})
	t.Run("loginWithEmailSuccess", func(t *testing.T) {
		user3 := models.User{
			Username:   "test_login_3",
			Nickname:   "test_login_3_rand_str",
			Email:      "test_login_3@mail.com",
			Password:   utils.HashPassword("test_login_password"),
			RoleLoaded: true,
			Roles:      []models.UserHasRole{},
		}
		base.DB.Create(&user3)
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: user3.Email,
			Password:        "test_login_password",
		}))
		resp := response.LoginResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.Equal(t, "SUCCESS", resp.Message)
		assert.Equal(t, nil, resp.Error)
		jsonEQ(t, user3, resp.Data.User)
		token, err := utils.GetToken(resp.Data.Token)
		base.DB.Where("id = ?", user3.ID).First(&user3)
		assert.Equal(t, nil, err)
		token.User.LoadRoles()
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
		httpResp := makeResp(makeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: user4.Email,
			Password:        "WRONG_PASSWORD",
		}))
		resp := response.LoginResponse{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusForbidden, httpResp.StatusCode)
		assert.Equal(t, "WRONG_PASSWORD", resp.Message)
		assert.Equal(t, nil, resp.Error)
	})
	t.Run("loginWithUsernameAndRolesSuccess", func(t *testing.T) {

		classA := testClass{ID: 1}
		dummy := "test_class"
		adminRole := models.Role{
			Name:   "admin",
			Target: &dummy,
		}
		base.DB.Create(&adminRole)
		adminRole.AddPermission("all")
		user5 := models.User{
			Username:   "test_login_5",
			Nickname:   "test_login_5_rand_str",
			Email:      "test_login_5@mail.com",
			Password:   utils.HashPassword("test_login_password"),
			RoleLoaded: true,
			Roles:      []models.UserHasRole{},
		}
		base.DB.Create(&user5)
		user5.GrantRole(adminRole, classA)

		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/api/auth/login", request.LoginRequest{
			UsernameOrEmail: user5.Username,
			Password:        "test_login_password",
			RememberMe:      false,
		}))
		resp := struct {
			Message string      `json:"message"`
			Error   interface{} `json:"error"`
			Data    struct {
				User struct {
					ID       uint   `gorm:"primary_key" json:"id"`
					Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5,username"`
					Nickname string `gorm:"index:nickname" json:"nickname"`
					Email    string `gorm:"unique_index" json:"email"`
					Password string `json:"-"`

					Roles      []models.Role `json:"roles"`
					RoleLoaded bool          `gorm:"-"`

					CreatedAt time.Time  `json:"created_at"`
					UpdatedAt time.Time  `json:"-"`
					DeletedAt *time.Time `sql:"index" json:"deleted_at"`
					//TODO: bio
				}
				Token string `json:"token"`
			} `json:"data"`
		}{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.Equal(t, "SUCCESS", resp.Message)
		assert.Equal(t, nil, resp.Error)

		jsonEQ(t, user5, resp.Data.User)
		token, err := utils.GetToken(resp.Data.Token)
		base.DB.Where("id = ?", user5.ID).First(&user5)
		assert.Equal(t, nil, err)
		token.User.LoadRoles()
		assert.True(t, user5.UpdatedAt.Equal(token.User.UpdatedAt))
		assert.Equal(t, user5, token.User)
		assert.False(t, token.RememberMe)
	})
}

func TestRegister(t *testing.T) {
	t.Run("registerUserWithoutParams", func(t *testing.T) {
		t.Parallel()
		resp := makeResp(makeReq(t, "POST", "/api/auth/register", request.RegisterRequest{
			Username: "",
			Nickname: "",
			Email:    "",
			Password: "",
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
	t.Run("registerUserSuccess", func(t *testing.T) {
		t.Parallel()
		respResponse := makeResp(makeReq(t, "POST", "/api/auth/register", request.RegisterRequest{
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
		jsonEQ(t, resp.Data.User, user)

		respResponse = makeResp(makeReq(t, "POST", "/api/auth/register", request.RegisterRequest{
			Username: "test_registerUserSuccess_0",
			Nickname: "test_registerUserSuccess_0",
			Email:    "test_registerUserSuccess_0@mail.com",
			Password: "test_registerUserSuccess_0",
		}))
		jsonEQ(t, response.ErrorResp("DUPLICATE_EMAIL", nil), respResponse)
		assert.Equal(t, http.StatusBadRequest, respResponse.StatusCode)
		respResponse = makeResp(makeReq(t, "POST", "/api/auth/register", request.RegisterRequest{
			Username: "test_registerUserSuccess_0",
			Nickname: "test_registerUserSuccess_0",
			Email:    "test_registerUserSuccess_1@mail.com",
			Password: "test_registerUserSuccess_0",
		}))
		jsonEQ(t, response.ErrorResp("DUPLICATE_USERNAME", nil), respResponse)
		assert.Equal(t, http.StatusBadRequest, respResponse.StatusCode)
	})
}
