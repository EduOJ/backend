package base

import (
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

var Echo *echo.Echo
var Redis *redis.Client
var Gorm *gorm.DB

// Postgresql Client
