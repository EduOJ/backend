package validator

import (
	zhLocal "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/labstack/echo/v4"
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

var Trans ut.Translator

func InitValidator(e *echo.Echo) {
	zh := zhLocal.New()
	uni := ut.New(zh, zh)
	var found bool
	Trans, found = uni.GetTranslator("zh")
	if !found {
		log.Fatal("could not found zh translator")
		panic("could not found zh translator")
	}
	v := validator.New()
	// add custom translation here
	if err := zhTranslations.RegisterDefaultTranslations(v, Trans); err != nil {
		log.Fatal(errors.Wrap(err, "could not register default translations"))
		panic(errors.Wrap(err, "could not register default translations"))
	}
	err := v.RegisterValidation("username", ValidateUsername)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not register validation"))
		panic(errors.Wrap(err, "could not register validation"))
	}
	e.Validator = &Validator{
		V: v,
	}
}

var UsernameRegex = regexp.MustCompile("^[a-zA-Z0-9_]+$")

func ValidateUsername(fl validator.FieldLevel) bool {
	return UsernameRegex.MatchString(fl.Field().String())
}
