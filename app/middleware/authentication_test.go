package middleware_test

import (
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/middleware"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
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
		Message: "SUCCESS",
		Error:   nil,
		Data:    user,
	})
}

func TestAuthenticationAndLoginCheck(t *testing.T) {
	t.Parallel()
	e := echo.New()
	httpSuccessResponse := response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	}

	e.Use(middleware.Authentication)
	e.POST("/test_authentication", testController)
	e.POST("/test_loginCheck", testController, middleware.LoginCheck)

	testUser := models.User{
		Username: "testAuthenticationMiddle",
		Nickname: "testAuthenticationMiddle",
		Email:    "testAuthenticationMiddle@e.com",
		Password: "",
	}
	assert.Equal(t, nil, base.DB.Save(&testUser).Error)

	activeToken := models.Token{
		Token: utils.RandStr(32),
		User:  testUser,
	}
	expiredToken := models.Token{
		Token:     utils.RandStr(32),
		User:      testUser,
		UpdatedAt: time.Now().Add(-1 * time.Second * time.Duration(2000)),
	}
	activeRememberMeToken := models.Token{
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
	assert.Equal(t, nil, base.DB.Save(&activeToken).Error)
	assert.Equal(t, nil, base.DB.Save(&expiredToken).Error)
	assert.Equal(t, nil, base.DB.Save(&activeRememberMeToken).Error)
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
			resp:        response.ErrorResp("AUTH_TOKEN_NOT_FOUND", nil),
		},
		{
			name:        "testExpiredToken",
			tokenString: expiredToken.Token,
			statusCode:  http.StatusRequestTimeout,
			resp:        response.ErrorResp("AUTH_SESSION_EXPIRED", nil),
		},
		{
			name:        "testExpiredRememberMeToken",
			tokenString: expiredRememberMeToken.Token,
			statusCode:  http.StatusRequestTimeout,
			resp:        response.ErrorResp("AUTH_SESSION_EXPIRED", nil),
		},
	}
	for _, test := range failTests {
		test := test
		t.Run("failTests"+test.name, func(t *testing.T) {
			t.Parallel()
			req := makeReq(t, "POST", "/test_authentication", nil)
			req.Header.Set("Authorization", test.tokenString)
			httpResp := makeResp(req, e)
			resp := response.Response{}
			mustJsonDecode(httpResp, &resp)
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
			tokenString: activeToken.Token,
			user:        activeToken.User,
			statusCode:  http.StatusOK,
			resp:        httpSuccessResponse,
		},
		{
			name:        "activeRememberMeToken",
			tokenString: activeRememberMeToken.Token,
			user:        activeRememberMeToken.User,
			statusCode:  http.StatusOK,
			resp:        httpSuccessResponse,
		},
	}

	for _, test := range successTests {
		test := test
		t.Run("successTests"+test.name, func(t *testing.T) {
			t.Parallel()
			req := makeReq(t, "POST", "/test_authentication", nil)
			req.Header.Set("Authorization", test.tokenString)
			httpResp := makeResp(req, e)
			resp := response.Response{}
			mustJsonDecode(httpResp, &resp)
			assert.Equal(t, test.statusCode, httpResp.StatusCode)
			jsonEQ(t, response.Response{
				Message: "SUCCESS",
				Error:   nil,
				Data:    test.user,
			}, resp)
		})
	}

	t.Run("testLoginCheckFail", func(t *testing.T) {
		t.Parallel()
		LoginCheckReq := makeReq(t, "POST", "/test_loginCheck", nil)
		LoginCheckReq.Header.Set("Authorization", "")
		httpResp := makeResp(LoginCheckReq, e)
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusUnauthorized, httpResp.StatusCode)
		assert.Equal(t, response.ErrorResp("AUTH_NEED_TOKEN", nil), resp)
	})

	t.Run("testLoginCheckSuccess", func(t *testing.T) {
		t.Parallel()
		LoginCheckReq := makeReq(t, "POST", "/test_loginCheck", nil)
		LoginCheckReq.Header.Set("Authorization", activeToken.Token)
		httpResp := makeResp(LoginCheckReq, e)
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		jsonEQ(t, response.Response{
			Message: "SUCCESS",
			Error:   nil,
			Data:    activeToken.User,
		}, resp)
	})
}
