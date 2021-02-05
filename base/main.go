package base

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

var Echo *echo.Echo
var Redis *redis.Client
var DB *gorm.DB
var Storage *minio.Client
