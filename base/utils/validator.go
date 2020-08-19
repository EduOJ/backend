package utils

import (
	zhLocal "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/pkg/errors"
	"regexp"
)

type Validator struct {
	V *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
	return cv.V.Struct(i)
}

var Validate *validator.Validate
var Trans ut.Translator

func init() { // TODO: add to init list?
	zh := zhLocal.New()
	uni := ut.New(zh, zh)
	var found bool
	Trans, found = uni.GetTranslator("zh")
	if !found {
		log.Warning("could not found zh translator")
	}
	Validate = validator.New()
	// add custom translation here
	if err := zhTranslations.RegisterDefaultTranslations(Validate, Trans); err != nil {
		log.Error(errors.Wrap(err, "could not register default translations"))
	}
}

var UsernameRegex = regexp.MustCompile("^[a-zA-Z0-9_]+$")

func ValidateUsername(fl validator.FieldLevel) bool {
	return UsernameRegex.MatchString(fl.Field().String())
}
