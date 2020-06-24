package base

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

// E : The echo object.
var E *echo.Echo
var redisClient *redis.Client

// Redis returns the Redis object.
func Redis() *redis.Client {
	return redisClient
}

// Postgresql Client
