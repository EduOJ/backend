package base

import (
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go"
)

import _ "github.com/jinzhu/gorm/dialects/mysql"
import _ "github.com/jinzhu/gorm/dialects/postgres"
import _ "github.com/jinzhu/gorm/dialects/sqlite"

var Echo *echo.Echo
var Redis *redis.Client
var DB *gorm.DB
var Storage *minio.Client
