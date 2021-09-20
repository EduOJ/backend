package base

import (
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/go-mail/mail"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"html/template"
)

var Echo *echo.Echo
var Redis *redis.Client
var DB *gorm.DB
var Storage *minio.Client
var WebAuthn *webauthn.WebAuthn
var Mail mail.Dialer
var Template *template.Template
