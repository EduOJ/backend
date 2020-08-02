package middleware_test

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/exit"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/database"
	"github.com/leoleoasd/EduOJBackend/database/models"
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

func MustJsonDecode(data interface{}, out interface{}) {
	var err error
	if dataResp, ok := data.(*http.Response); ok {
		data, err = ioutil.ReadAll(dataResp.Body)
		if err != nil {
			panic(err)
		}
	}
	if dataString, ok := data.(string); ok {
		data = []byte(dataString)
	}
	if dataBytes, ok := data.([]byte); ok {
		err = json.Unmarshal(dataBytes, out)
		if err != nil {
			panic(err)
		}
	}
}

func MakeReq(t *testing.T, method string, path string, data interface{}) *http.Request {
	j, err := json.Marshal(data)
	assert.Equal(t, nil, err)
	req := httptest.NewRequest(method, path, bytes.NewReader(j))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return req
}

func MakeResp(req *http.Request, e *echo.Echo) *http.Response {
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Result()
}

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

func TestMain(m *testing.M) {
	defer database.SetupDatabaseForTest()()
	defer exit.SetupExitForTest()()
	configFile := bytes.NewBufferString(`debug: false
auth:
  session_timeout: 1200
  remember_me_timeout: 604800
  session_count: 10`)
	//configFile := bytes.NewBufferString("debug: false" +
	//	"\nauth:\n  session_timeout: 1200\n  remember_me_timeout: 604800\n  session_count: 10")
	err := config.ReadConfig(configFile)
	if err != nil {
		panic(err)
	}
	log.Disable()

	base.Echo = echo.New()
	app.Register(base.Echo)
	os.Exit(m.Run())
}
