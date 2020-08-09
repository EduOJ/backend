package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/pkg/errors"
	"net/http"
)

func PanicIfDBError(db *gorm.DB, message string) {
	if db.Error != nil {
		panic(errors.Wrap(db.Error, message))
	}
}

type ValidateFunc func(req interface{}) (bool, string)
type ExtraValidate struct {
	field         string
	validate      ValidateFunc
	fieldAppeared bool
	tag           string
}

func BindAndValidate(req interface{}, c *echo.Context, extraValidates ...ExtraValidate) (err error, ok bool) {
	if err := (*c).Bind(req); err != nil {
		panic(err)
	}
	var extraErrors []ExtraValidate
	for _, ev := range extraValidates {
		ok, tag := ev.validate(req)
		if !ok {
			ev.tag = tag
			extraErrors = append(extraErrors, ev)
		}
	}
	var validationErrors []response.ValidationError
	if err := (*c).Validate(req); err != nil {
		e, ok := err.(validator.ValidationErrors)
		if !ok {
			log.Error(errors.Wrap(err, "validate failed"), *c)
			return response.InternalErrorResp(*c), false
		}
		validationErrors = make([]response.ValidationError, len(e))
		for i, v := range e {
			field := v.Field()
			tag := v.Tag()
			for _, ev := range extraErrors {
				if field == ev.field {
					tag += ", " + ev.tag
					ev.fieldAppeared = true
				}
			}
			validationErrors[i] = response.ValidationError{
				Field:  field,
				Reason: tag,
			}
		}
	} else if len(extraErrors) == 0 {
		return nil, true
	}
	for _, ev := range extraErrors {
		if !ev.fieldAppeared {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:  ev.field,
				Reason: ev.tag,
			})
		}
	}
	return (*c).JSON(http.StatusBadRequest, response.ErrorResp("VALIDATION_ERROR", validationErrors)), false
}
