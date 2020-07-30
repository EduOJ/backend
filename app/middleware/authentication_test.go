package middleware_test

import (
	"github.com/go-playground/assert/v2"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/middleware"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"time"

	"github.com/leoleoasd/EduOJBackend/base"
	"net/http"
	"testing"
)

func testController(context echo.Context) error {
	user := context.Get("user")
	if user == nil {
		user = models.User{}
	}
	return context.JSON(http.StatusOK, response.Response{
		Code:    0,
		Message: "success",
		Error:   nil,
		Data:    user,
	})
}

func TestAuthentication(t *testing.T) {
	oldEcho := base.Echo
	base.Echo = echo.New()
	t.Cleanup(func() {
		base.Echo = oldEcho
	})
	httpSuccessResponse := response.Response{
		Code:    0,
		Message: "success",
		Error:   nil,
		Data:    nil,
	}

	base.Echo.Use(middleware.Authentication)
	base.Echo.POST("/test_authentication", testController)

	testUser := models.User{
		Username: "testAuthenticationMiddle",
		Nickname: "testAuthenticationMiddle",
		Email:    "testAuthenticationMiddle@e.com",
		Password: "",
	}
	assert.Equal(t, nil, base.DB.Save(&testUser).Error)

	effectiveToken := models.Token{
		Token: utils.RandStr(32),
		User:  testUser,
	}
	expiredToken := models.Token{
		Token:     utils.RandStr(32),
		User:      testUser,
		UpdatedAt: time.Now().Add(-1 * time.Second * time.Duration(2000)),
	}
	effectiveRememberMeToken := models.Token{
		Token:      utils.RandStr(32),
		User:       testUser,
		RememberMe: true,
	}
	expiredRememberMeToken := models.Token{
		Token:      utils.RandStr(32),
		User:       testUser,
		UpdatedAt:  time.Now().Add(-1 * time.Second * time.Duration(720000)),
		RememberMe: true,
	}
	assert.Equal(t, nil, base.DB.Save(&effectiveToken).Error)
	assert.Equal(t, nil, base.DB.Save(&expiredToken).Error)
	assert.Equal(t, nil, base.DB.Save(&effectiveRememberMeToken).Error)
	assert.Equal(t, nil, base.DB.Save(&expiredRememberMeToken).Error)

	failTests := []struct {
		name        string
		tokenString string
		statusCode  int
		resp        response.Response
	}{
		{
			name:        "testNon-existingToken",
			tokenString: "Non-existingToken",
			statusCode:  http.StatusUnauthorized,
			resp:        response.ErrorResp(1, "Token not found", nil),
		},
		{
			name:        "testExpiredToken",
			tokenString: expiredToken.Token,
			statusCode:  http.StatusRequestTimeout,
			resp:        response.ErrorResp(1, "session expired", nil),
		},
		{
			name:        "testExpiredRememberMeToken",
			tokenString: expiredRememberMeToken.Token,
			statusCode:  http.StatusRequestTimeout,
			resp:        response.ErrorResp(1, "session expired", nil),
		},
	}
	for _, test := range failTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			req := MakeReq(t, "POST", "/test_authentication", nil)
			req.Header.Set("Authorization", test.tokenString)
			httpResp := MakeResp(req)
			resp := response.Response{}
			MustJsonDecode(httpResp, &resp)
			assert.Equal(t, test.statusCode, httpResp.StatusCode)
			assert.Equal(t, test.resp, resp)
		})
	}

	successTests := []struct {
		name        string
		tokenString string
		user        models.User
		statusCode  int
		resp        response.Response
	}{
		{
			name:        "testEmptyToken",
			tokenString: "",
			user:        models.User{},
			statusCode:  http.StatusOK,
			resp:        httpSuccessResponse,
		},
		{
			name:        "testEffectiveToken",
			tokenString: effectiveToken.Token,
			user:        effectiveToken.User,
			statusCode:  http.StatusOK,
			resp:        httpSuccessResponse,
		},
		{
			name:        "effectiveRememberMeToken",
			tokenString: effectiveRememberMeToken.Token,
			user:        effectiveRememberMeToken.User,
			statusCode:  http.StatusOK,
			resp:        httpSuccessResponse,
		},
	}

	for _, test := range successTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			req := MakeReq(t, "POST", "/test_authentication", nil)
			req.Header.Set("Authorization", test.tokenString)
			httpResp := MakeResp(req)
			resp := response.Response{}
			MustJsonDecode(httpResp, &resp)
			assert.Equal(t, test.statusCode, httpResp.StatusCode)
			JsonEQ(t, response.Response{
				Code:    0,
				Message: "success",
				Error:   nil,
				Data:    test.user,
			}, resp)
		})
	}

	base.Echo.POST("/test_loginCheck", testController, middleware.LoginCheck)

	t.Run("testLoginCheckFail", func(t *testing.T) {
		t.Parallel()
		LoginCheckReq := MakeReq(t, "POST", "/test_loginCheck", nil)
		LoginCheckReq.Header.Set("Authorization", "")
		httpResp := MakeResp(LoginCheckReq)
		resp := response.Response{}
		MustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusUnauthorized, httpResp.StatusCode)
		assert.Equal(t, response.ErrorResp(1, "Unauthorized", nil), resp)
	})

	t.Run("testLoginCheckSuccess", func(t *testing.T) {
		t.Parallel()
		LoginCheckReq := MakeReq(t, "POST", "/test_loginCheck", nil)
		LoginCheckReq.Header.Set("Authorization", effectiveToken.Token)
		httpResp := MakeResp(LoginCheckReq)
		resp := response.Response{}
		MustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		JsonEQ(t, response.Response{
			Code:    0,
			Message: "success",
			Error:   nil,
			Data:    effectiveToken.User,
		}, resp)
	})
}
