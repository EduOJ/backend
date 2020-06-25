package base

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

var echoObject *echo.Echo
var redisClient *redis.Client

func Echo() *echo.Echo {
	return echoObject
}

// Redis returns the Redis object.
func Redis() *redis.Client {
	return redisClient
}

// Postgresql Client
