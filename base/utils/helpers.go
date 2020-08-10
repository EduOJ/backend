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

// The function for custom validations
type ValidateFunc func(req interface{}) (bool, string)

// The structure to save custom validations
type CustomValidation struct {
	requestField  string       // The field in the request
	validate      ValidateFunc // The custom validating function
	fieldAppeared bool         // If this field has standard validation error (this variable is useless now)
	tag           string       // The custom validation tag
}

func BindAndValidate(req interface{}, c *echo.Context, customValidations ...CustomValidation) (err error, ok bool) {
	if err := (*c).Bind(req); err != nil {
		panic(err)
	}
	// Execute custom validation, collect custom validation errors in "customErrors"
	// The final used length of "customErrors" isn't greater than len(customValidations)
	customErrors := make([]CustomValidation, 0, len(customValidations))
	for _, ev := range customValidations {
		ok, tag := ev.validate(req)
		if !ok {
			ev.tag = tag
			customErrors = append(customErrors, ev)
		}
	}
	// Combine the custom validation errors and standard validation errors
	var validationErrors []response.ValidationError
	if err := (*c).Validate(req); err != nil {
		// There are standard validation errors.
		// For those fields with both custom and standard validation errors, we insert
		// the custom errors' tags into standard errors' tags;
		// For those fields with only standard validation errors, the errors would be
		// record in "response.ValidationError" structures.
		e, ok := err.(validator.ValidationErrors)
		if !ok {
			log.Error(errors.Wrap(err, "validate failed"), *c)
			return response.InternalErrorResp(*c), false
		}
		// The final used length of "validationErrors" isn't greater than len(e)+len(customErrors)
		validationErrors = make([]response.ValidationError, 0, len(e)+len(customErrors))
		for i, v := range e {
			field := v.Field()
			tag := v.Tag()
			for _, ev := range customErrors {
				if field == ev.requestField {
					tag += ", " + ev.tag
					ev.fieldAppeared = true
				}
			}
			validationErrors[i] = response.ValidationError{
				Field:  field,
				Reason: tag,
			}
		}
	} else if len(customErrors) == 0 {
		//There are neither standard validation errors or custom validation errors.
		//So we just return
		return nil, true
	}
	// The fields with standard errors are processed by the if at 42, so the fields
	// only have custom errors here. We create "response.ValidationError" structures
	// to record them.
	for _, ev := range customErrors {
		if !ev.fieldAppeared {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:  ev.requestField,
				Reason: ev.tag,
			})
		}
	}
	return (*c).JSON(http.StatusBadRequest, response.ErrorResp("VALIDATION_ERROR", validationErrors)), false
}
