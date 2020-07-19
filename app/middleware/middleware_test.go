package middleware_test

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/exit"
	"github.com/leoleoasd/EduOJBackend/database"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func JsonEQ(t *testing.T, expected, actual interface{}) {
	assert.JSONEq(t, MustJsonEncode(t, expected), MustJsonEncode(t, actual))
}

func MustJsonEncode(t *testing.T, data interface{}) string {
	var err error
	if dataResp, ok := data.(*http.Response); ok {
		data, err = ioutil.ReadAll(dataResp.Body)
		assert.Equal(t, nil, err)
	}
	if dataString, ok := data.(string); ok {
		data = []byte(dataString)
	}
	if dataBytes, ok := data.([]byte); ok {
		err := json.Unmarshal(dataBytes, &data)
		assert.Equal(t, nil, err)
	}
	j, err := json.Marshal(data)
	if err != nil {
		t.Fatal(data, err)
	}
	return string(j)
}

func MakeReq(t *testing.T, method string, path string, data interface{}) *http.Request {
	j, err := json.Marshal(data)
	assert.Equal(t, nil, err)
	req := httptest.NewRequest(method, path, bytes.NewReader(j))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return req
}

func MakeResp(req *http.Request) *http.Response {
	rec := httptest.NewRecorder()
	base.Echo.ServeHTTP(rec, req)
	return rec.Result()
}

func TestMain(m *testing.M) {
	defer database.SetupDatabaseForTest()()
	defer exit.SetupExitForTest()()
	println(base.DB.HasTable("users"))
	base.Echo = echo.New()
	app.Register(base.Echo)
	os.Exit(m.Run())
}
