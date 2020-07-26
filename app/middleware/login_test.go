package middleware_test

import (
	"github.com/go-playground/assert/v2"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/middleware"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"net/http"
	"testing"
)

func LoginMiddlewareTest(c echo.Context) error {
	return c.JSON(http.StatusOK, response.Response{
		Code:    0,
		Message: "login_middleware_test",
		Error:   nil,
		Data:	nil,
	})
}

func TestLogin(t *testing.T) {
	login := base.Echo.Group("/api").Group("/login",middleware.Login)
	login.POST("/login_middleware_test", LoginMiddlewareTest).Name = "login.login_middleware_test"

	t.Parallel()
	user1 := models.User{
		Username: "test_login_middleware_1",
		Nickname: "test_login_middleware_rand",
		Email:    "test_login_middleware@e.com",
		Password: utils.HashPassword("test_login_middleware_password"),
	}
	base.DB.Create(&user1)
	// strip monotonic time
	user1.CreatedAt = user1.CreatedAt.Round(0)
	user1.UpdatedAt = user1.UpdatedAt.Round(0)
	httpResp := MakeResp(MakeReq(t, "POST", "/api/auth/login", request.LoginRequest{
		UsernameOrEmail: user1.Email,
		Password:        "test_login_middleware_password",
	}))
	resp := response.LoginResponse{}
	MustJsonDecode(httpResp, &resp)
	user1Token := resp.Data.Token




	t.Run("requestWithoutToken", func(t *testing.T) {
		t.Parallel()
		httpResp := MakeResp(MakeReq(t, "POST", "/api/login/login_middleware_test", request.LoggedRequest{}))
		assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)
		JsonEQ(t, response.Response{
			Code:    1,
			Message: "validation error",
			Error: []map[string]string{
				{
					"field":  "Token",
					"reason": "required",
				},
			},
			Data: nil,
		}, httpResp)
	})
	t.Run("requestWithUnqualifiedToken", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			Token   string
			Reason	string
		}{
			{
				"tooShort", "len",
			},
			{
				"number23333", "alpha",
			},
			{
				"tooLongTooLongTooLongTooLongTooLong", "len",
			},
		}
		for _, test := range tests {
			httpResp := MakeResp(MakeReq(t, "POST", "/api/login/login_middleware_test", request.LoggedRequest{
				Token: test.Token,
			}))
			assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)
			JsonEQ(t, response.Response{
				Code:    1,
				Message: "validation error",
				Error: []map[string]string{
					{
						"field":  "Token",
						"reason": test.Reason,
					},
				},
				Data: nil,
			}, httpResp)
		}
	})
	t.Run("requestWithWrongToken", func(t *testing.T) {
		t.Parallel()
		httpResp := MakeResp(MakeReq(t, "POST", "/api/login/login_middleware_test", request.LoggedRequest{
			Token: utils.RandStr(32),
		}))
		assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)
		JsonEQ(t, response.Response{
			Code:    2,
			Message: "invalid token",
			Error: nil,
			Data: nil,
		}, httpResp)
	})
	t.Run("requestWithCorrectToken", func(t *testing.T) {
		t.Parallel()
		httpResp := MakeResp(MakeReq(t, "POST", "/api/login/login_middleware_test", request.LoggedRequest{
			Token: user1Token,
		}))
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		JsonEQ(t, response.Response{
			Code:    0,
			Message: "login_middleware_test",
			Error: nil,
			Data: nil,
		}, httpResp)
	})
}
