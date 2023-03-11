package middleware_test

import (
	"net/http"
	"testing"

	"github.com/EduOJ/backend/app/middleware"
	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/database/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestValidateParams(t *testing.T) {
	t.Parallel()
	e := echo.New()
	group1 := e.Group("/group1", middleware.ValidateParams(map[string]string{
		"test_param1": "PARAM1_NOT_FOUND",
	}))
	group2 := e.Group("/group2", middleware.ValidateParams(map[string]string{
		"test_param1": "PARAM1_NOT_FOUND",
		"test_param2": "PARAM2_NOT_FOUND",
	}))
	group3 := e.Group("/group3", middleware.ValidateParams(map[string]string{
		"test_param1": "PARAM1_NOT_FOUND",
		"test_param2": "PARAM2_NOT_FOUND",
	}))
	group1.GET("/:test_param1/:test_param2/test", testController).Name = "group1.test"
	group2.GET("/:test_param1/:test_param2/test", testController).Name = "group2.test"
	group3.GET("/:test_param1/test", testController).Name = "group3.test"

	t.Run("Pass", func(t *testing.T) {
		t.Run("Group1", func(t *testing.T) {
			httpResp := makeResp(makeReq(t, "GET", e.Reverse("group1.test", 1, "non_int_string"), nil), e)
			assert.Equal(t, http.StatusOK, httpResp.StatusCode)
			resp := response.Response{}
			mustJsonDecode(httpResp, &resp)
			jsonEQ(t, response.Response{
				Message: "SUCCESS",
				Error:   nil,
				Data:    models.User{},
			}, resp)
		})
		t.Run("Group2", func(t *testing.T) {
			httpResp := makeResp(makeReq(t, "GET", e.Reverse("group2.test", "", -1), nil), e)
			assert.Equal(t, http.StatusOK, httpResp.StatusCode)
			resp := response.Response{}
			mustJsonDecode(httpResp, &resp)
			jsonEQ(t, response.Response{
				Message: "SUCCESS",
				Error:   nil,
				Data:    models.User{},
			}, resp)
		})
		t.Run("Group3", func(t *testing.T) {
			httpResp := makeResp(makeReq(t, "GET", e.Reverse("group3.test", -2), nil), e)
			assert.Equal(t, http.StatusOK, httpResp.StatusCode)
			resp := response.Response{}
			mustJsonDecode(httpResp, &resp)
			jsonEQ(t, response.Response{
				Message: "SUCCESS",
				Error:   nil,
				Data:    models.User{},
			}, resp)
		})
	})
	t.Run("Fail", func(t *testing.T) {
		t.Run("Group1", func(t *testing.T) {
			httpResp := makeResp(makeReq(t, "GET", e.Reverse("group1.test", "non_int_string", 2), nil), e)
			assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
			resp := response.Response{}
			mustJsonDecode(httpResp, &resp)
			jsonEQ(t, response.ErrorResp("PARAM1_NOT_FOUND", nil), resp)
		})
		t.Run("Group2", func(t *testing.T) {
			httpResp := makeResp(makeReq(t, "GET", e.Reverse("group2.test", 0, "non_int_string"), nil), e)
			assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
			resp := response.Response{}
			mustJsonDecode(httpResp, &resp)
			jsonEQ(t, response.ErrorResp("PARAM2_NOT_FOUND", nil), resp)
		})
		t.Run("Group3", func(t *testing.T) {
			httpResp := makeResp(makeReq(t, "GET", e.Reverse("group3.test", "non_int_string"), nil), e)
			assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
			resp := response.Response{}
			mustJsonDecode(httpResp, &resp)
			jsonEQ(t, response.ErrorResp("PARAM1_NOT_FOUND", nil), resp)
		})
	})
}
