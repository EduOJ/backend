package response

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    interface{} `json:"data"`
}

func InternalErrorResp(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, Response{
		Code:    -1,
		Message: "Internal error",
		Error:   nil,
		Data:    nil,
	})
}

func ErrorResp(code int, message string, error interface{}) Response {
	return Response{
		Code:    code,
		Message: message,
		Error:   error,
		Data:    nil,
	}
}
