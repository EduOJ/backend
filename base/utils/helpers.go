package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/database/models"
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

func BindAndValidate(req interface{}, c *echo.Context) (err error, ok bool) {
	if err := (*c).Bind(req); err != nil {
		panic(err)
	}
	if err := (*c).Validate(req); err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			validationErrors := make([]response.ValidationError, len(e))
			for i, v := range e {
				validationErrors[i] = response.ValidationError{
					Field:  v.Field(),
					Reason: v.Tag(),
				}
			}
			return (*c).JSON(http.StatusBadRequest, response.ErrorResp("VALIDATION_ERROR", validationErrors)), false
		}
		log.Error(errors.Wrap(err, "validate failed"), c)
		return response.InternalErrorResp(*c), false
	}
	return nil, true
}

func FindUser(id string) (*models.User, error) {
	user := models.User{}
	err := base.DB.Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		err = base.DB.Where("username = ?", id).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, err
			} else {
				panic(errors.Wrap(err, "could not query username"))
			}
		}
	} else if err != nil {
		panic(errors.Wrap(err, "could not query id"))
	}
	return &user, nil
}
