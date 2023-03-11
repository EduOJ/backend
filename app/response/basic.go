package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ValidationError struct {
	Field       string `json:"field"`
	Reason      string `json:"reason"`
	Translation string `json:"translation"`
}

type Response struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    interface{} `json:"data"`
}

func InternalErrorResp(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, MakeInternalErrorResp())
}

func MakeInternalErrorResp() Response {
	return Response{
		Message: "INTERNAL_ERROR",
		Error:   nil,
		Data:    nil,
	}
}

func ErrorResp(message string, error interface{}) Response {
	return Response{
		Message: message,
		Error:   error,
		Data:    nil,
	}
}
