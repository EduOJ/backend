package base

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/pkg/errors"
)

import _ "github.com/jinzhu/gorm/dialects/mysql"
import _ "github.com/jinzhu/gorm/dialects/postgres"
import _ "github.com/jinzhu/gorm/dialects/sqlite"

var Echo *echo.Echo
var Redis *redis.Client
var DB *gorm.DB

func InitGormFromConfig(_conf config.Node) error {
	if conf, ok := _conf.(*config.MapNode); ok {
		dialect := conf.MustGet("dialect", "").Value().(string)
		uri := conf.MustGet("uri", "").Value().(string)
		var err error
		DB, err = gorm.Open(dialect, uri)
		if err != nil {
			return errors.Wrap(err, "could not connect to database")
		}
		return nil
	} else {
		return errors.New("database configuration should be a map")
	}
}

func InitRedisFromConfig(_conf config.Node) error {
	if conf, ok := _conf.(*config.MapNode); ok {
		port := conf.MustGet("port", 6379).Value().(int)
		host := conf.MustGet("host", "localhost").Value().(string)
		username := conf.MustGet("username", "").Value().(string)
		password := conf.MustGet("password", "").Value().(string)
		Redis = redis.NewClient(&redis.Options{
			Addr:               fmt.Sprint(host, ":", port),
			Username:           username,
			Password:           password,
		})
		// Test connection.
		_, err := Redis.Ping(context.Background()).Result()
		if err != nil {
			return errors.Wrap(err, "could not connect to the redis server")
		}
		return nil
	} else {
		return errors.New("web server configuration should be a map")
	}
}

// InitEchoFromConfig initialize the Echo object and starts the web server according to the config.
func InitEchoFromConfig(_conf config.Node) error {
	if conf, ok := _conf.(*config.MapNode); ok {
		port := conf.MustGet("port", 8080).Value().(int)
		Echo = echo.New()
		Echo.Logger = &log.EchoLogger{}
		Echo.HideBanner = true
		Echo.HidePort = true
		Echo.Use(middleware.Recover())
		Echo.Server.Addr = fmt.Sprintf(":%d", port)
		// TODO: load routes and middleware.
		return nil
	} else {
		return errors.New("web server configuration should be a map")
	}
}
