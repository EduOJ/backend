package base

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

func init() {
	E = echo.New()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// e.g. Connect to postgresql server
	// Read configuration file
	// Parse command line arguments.
}
