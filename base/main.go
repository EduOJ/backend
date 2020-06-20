package base

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

var E *echo.Echo
var redisClient *redis.Client

func Redis() *redis.Client {
	return redisClient
}
