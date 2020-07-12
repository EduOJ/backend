package base

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/pkg/errors"
)

var Echo *echo.Echo
var Redis *redis.Client
var Gorm *gorm.DB

// 以下内容不知道是否会用到
//var Context context.Context
//var Cancel context.CancelFunc
//
//func init() {
//	Context, Cancel = context.WithCancel(context.Background())
//}

// InitEchoFromConfig initialize the Echo object and starts the web server according to the config.
func InitEchoFromConfig(_conf config.Node) error {
	if conf, ok := _conf.(*config.MapNode); ok {
		port := int(conf.MustGet("port", config.IntNode(8080)).(config.IntNode))
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
