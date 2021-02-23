package controller

import (
	"bytes"
	"encoding/json"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

var cac = cache.New(5*time.Minute, 10*time.Minute)

func BeginRegistration(c echo.Context) error {
	user := c.Get("user").(models.User)
	log.Debug(base.WebAuthn.Config)
	options, sessionData, err := base.WebAuthn.BeginRegistration(&user)
	if err != nil {
		panic(errors.Wrap(err, "could not begin register webauthn"))
	}
	cac.Set("register"+user.Username, sessionData, cache.DefaultExpiration)
	return c.JSON(http.StatusOK, response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    options,
	})
}

func FinishRegistration(c echo.Context) error {
	user := c.Get("user").(models.User)
	sessionData, found := cac.Get("register" + user.Username)
	if !found {
		panic(errors.New("not registered"))
	}
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		panic(errors.Wrap(err, "could not read body"))
	}
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(b))
	if err != nil {
		panic(errors.Wrap(err, "could not parse credential from body"))
	}
	credential, err := base.WebAuthn.CreateCredential(&user, *sessionData.(*webauthn.SessionData), parsedResponse)
	if err != nil {
		panic(errors.Wrap(err, "could not create credential"))
	}
	webauthnCredential := models.WebauthnCredential{}
	m, err := json.Marshal(credential)
	if err != nil {
		panic(errors.Wrap(err, "could not marshal credential to json"))
	}
	webauthnCredential.Content = string(m)
	if err := base.DB.Model(&user).Association("Credentials").Append(&webauthnCredential); err != nil {
		panic(errors.Wrap(err, "could not save credentials"))
	}
	return c.JSON(http.StatusOK, response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}
