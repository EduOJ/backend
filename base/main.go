package base

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/logging"
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
		Echo.Logger = &logging.EchoLogger{}
		go func() {
			logging.Fatal(Echo.Start(fmt.Sprintf(":%d", port)))
		}()
		return nil
	} else {
		return errors.New("web server configuration should be a map")
	}
}
