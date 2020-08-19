package app

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

type Validator struct {
	v *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
	return cv.v.Struct(i)
}

var UsernameRegex = regexp.MustCompile("^[a-zA-Z0-9_]+$")

func ValidateUsername(fl validator.FieldLevel) bool {
	return UsernameRegex.MatchString(fl.Field().String())
}
