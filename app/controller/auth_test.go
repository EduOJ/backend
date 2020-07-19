package controller_test

import (
	"encoding/json"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestRegister(t *testing.T) {
	t.Run("registerUserWithoutParams", func(t *testing.T) {
		t.Parallel()
		resp := MakeResp(MakeReq(t, "POST", "/api/auth/register", request.RegisterRequest{
			Username: "",
			Nickname: "",
			Email:    "",
			Password: "",
		}))
		JsonEQ(t, response.Response{
			Code:    1,
			Message: "validation error",
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
		JsonEQ(t, response.ErrorResp(2, "duplicate email", nil), respResponse)
		respResponse = MakeResp(MakeReq(t, "POST", "/api/auth/register", request.RegisterRequest{
			Username: "test_registerUserSuccess_0",
			Nickname: "test_registerUserSuccess_0",
			Email:    "test_registerUserSuccess_1@mail.com",
			Password: "test_registerUserSuccess_0",
		}))
		JsonEQ(t, response.ErrorResp(3, "duplicate username", nil), respResponse)
	})
}
