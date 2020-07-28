package middleware_test

import (
	"bytes"
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
	base.Echo.POST("/test_authentication", func(context echo.Context) error {
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
	})
	base.Echo.POST("/test_loginCheck", func(context echo.Context) error {
		return context.JSON(http.StatusOK, httpSuccessResponse)
	}, middleware.LoginCheck)

	req := MakeReq(t, "POST", "/test_authentication", &bytes.Buffer{})
	testUser := models.User{
		Username: "testUser",
		Nickname: "testUserNickname",
		Email:    "testUser@e.com",
		Password: "",
	}
	base.DB.Save(&testUser)
	effectiveToken := models.Token{
		Token: utils.RandStr(32),
		User:  testUser,
	}
	base.DB.Save(&effectiveToken)
	expiredToken := models.Token{
		Token:     utils.RandStr(32),
		User:      testUser,
		UpdatedAt: time.Now().Add(-1 * time.Second * time.Duration(200*3600)),
	}
	base.DB.Save(&expiredToken)

	type testType struct {
		name        string
		tokenString string
		user        models.User
		statusCode  int
		resp        response.Response
	}

	failTests := []testType{
		{
			name:        "testNon-existingToken",
			tokenString: "Non-existingToken",
			statusCode:  http.StatusUnauthorized,
			resp:        response.ErrorResp(1, "Unauthorized", nil),
		},
		{
			name:        "testExpiredToken",
			tokenString: expiredToken.Token,
			statusCode:  http.StatusRequestTimeout,
			resp:        response.ErrorResp(1, "session expired", nil),
		},
	}
	for _, test := range failTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			req.Header.Set("token", test.tokenString)
			httpResp := MakeResp(req)
			resp := response.Response{}
			MustJsonDecode(httpResp, &resp)
			assert.Equal(t, test.statusCode, httpResp.StatusCode)
			assert.Equal(t, test.resp, resp)
		})
	}

	successTests := []testType{
		{
			name:        "testEmptyToken",
			tokenString: "",
			user:        models.User{},
			statusCode:  http.StatusOK,
			resp:        httpSuccessResponse,
		}, {
			name:        "testEffectiveToken",
			tokenString: effectiveToken.Token,
			user:        effectiveToken.User,
			statusCode:  http.StatusOK,
			resp:        httpSuccessResponse,
		},
	}

	for _, test := range successTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			req.Header.Set("token", test.tokenString)
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

}
