package app

import "github.com/go-playground/validator/v10"

type Validator struct {
	v *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
	return cv.v.Struct(i)
}
