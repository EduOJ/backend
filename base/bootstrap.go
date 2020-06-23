package base

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"sync"
)

// Flag to show if redis is ready, for
// those who's init needs redis.
var redisReady sync.WaitGroup

func init() {
	E = echo.New()

	// REDIS
	redisReady.Add(1)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisReady.Done()

	// e.g. Connect to postgresql server
	// Read configuration file
	// Parse command line arguments.
}
