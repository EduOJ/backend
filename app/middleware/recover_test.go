package middleware_test

import (
	"bytes"
	"github.com/kami-zh/go-capturer"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/middleware"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/pkg/errors"
	"net/http"
	"testing"
)

func TestRecover(t *testing.T) {
	oldEcho := base.Echo
	base.Echo = echo.New()
	t.Cleanup(func() {
		base.Echo = oldEcho
	})

	base.Echo.Use(middleware.Recover)
	base.Echo.POST("/panics_with_error", func(context echo.Context) error {
		panic(errors.New("123"))
	})
	base.Echo.POST("/panics_with_other", func(context echo.Context) error {
		panic("123")
	})

	req := MakeReq(t, "POST", "/panics_with_error", &bytes.Buffer{})
	resp := (*http.Response)(nil)
	_ = capturer.CaptureOutput(func() {
		resp = MakeResp(req)
	})
	JsonEQ(t, response.MakeInternalErrorResp(), resp)
	req = MakeReq(t, "POST", "/panics_with_other", &bytes.Buffer{})
	_ = capturer.CaptureOutput(func() {
		resp = MakeResp(req)
	})
	JsonEQ(t, response.MakeInternalErrorResp(), resp)

}
