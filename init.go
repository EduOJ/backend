package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/exit"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/database"
	"github.com/pkg/errors"
	"os"
)

func readConfig() {
	log.Debug("Reading config.")
	configFile, err := open(opt.Config)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not open config file "+opt.Config))
		os.Exit(-1)
	}
	err = config.ReadConfig(configFile)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not read config file "+opt.Config))
		os.Exit(-1)
	}
	log.Debug("Config read.")
}

func initLog() {
	log.Debug("Initializing log.")
	loggingConf, err := config.Get("log")
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not read log config"))
		os.Exit(-1)
	}
	err = log.InitFromConfig(loggingConf)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not init log with config "+loggingConf.String()))
		os.Exit(-1)
	}
	log.Debug("Logging initialized.")
}

func startEcho() {
	log.Debug("Starting echo server.")
	echoConf, err := config.Get("server")
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not read http server config"))
		os.Exit(-1)
	}
	if _, ok := echoConf.(*config.MapNode); ok {
		log.Fatal(errors.Wrap(errors.New("web server configuration should be a map"), "could not init http server with config "+echoConf.String()))
		os.Exit(-1)
	}
	port := echoConf.MustGet("port", 8080).Value().(int)
	base.Echo = echo.New()
	base.Echo.Logger = &log.EchoLogger{}
	base.Echo.HideBanner = true
	base.Echo.HidePort = true
	base.Echo.Use(middleware.Recover())
	base.Echo.Server.Addr = fmt.Sprintf(":%d", port)
	// TODO: load routes and middleware.
	go func() {
		err := base.Echo.StartServer(base.Echo.Server)
		if err != nil {
			log.Fatal(errors.Wrap(err, "server closed"))
		}
	}()
	log.Fatal("Server started at ", base.Echo.Server.Addr)

	// When server closes, closes web server.
	go func() {
		<-exit.BaseContext.Done()
		err := base.Echo.Shutdown(context.Background())
		if err != nil {
			if err.Error() == "context canceled" {
				log.Fatal("Force quitting.")
			} else {
				log.Fatal(err)
			}
		}
	}()
}

func initRedis() {
	log.Debug("Starting redis client.")
	redisConf, err := config.Get("redis")
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not read redis config"))
		os.Exit(-1)
	}
	if _, ok := redisConf.(*config.MapNode); ok {
		log.Fatal(errors.Wrap(errors.New("redis configuration should be a map"), "could not init http server with config "+redisConf.String()))
		os.Exit(-1)
	}
	port := redisConf.MustGet("port", 6379).Value().(int)
	host := redisConf.MustGet("host", "localhost").Value().(string)
	username := redisConf.MustGet("username", "").Value().(string)
	password := redisConf.MustGet("password", "").Value().(string)
	base.Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprint(host, ":", port),
		Username: username,
		Password: password,
	})
	// Test connection.
	_, err = base.Redis.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not init redis with config "+redisConf.String()))
		os.Exit(-1)
	}
	log.Debug("Redis client started.")

	// When server closes, closes this client.
	exit.QuitWG.Add(1)
	go func() {
		<-exit.BaseContext.Done()
		_ = base.Redis.Close()
		exit.QuitWG.Done()
	}()
}

func initGorm(toMigrate ...bool) {
	log.Debug("Starting database client.")
	databaseConf, err := config.Get("database")
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not read database config"))
		os.Exit(-1)
	}
	if _, ok := databaseConf.(*config.MapNode); ok {
		log.Fatal(errors.Wrap(errors.New("database configuration should be a map"), "could not init http server with config "+databaseConf.String()))
		os.Exit(-1)
	}
	dialect := databaseConf.MustGet("dialect", "").Value().(string)
	uri := databaseConf.MustGet("uri", "").Value().(string)
	base.DB, err = gorm.Open(dialect, uri)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not init database with config "+databaseConf.String()))
		os.Exit(-1)
	}
	if len(toMigrate) == 0 || toMigrate[0] {
		database.Migrate()
	}
	log.Debug("Database client started.")

	// Cause we need to wait until all logs are wrote to the db
	// So we dont close db connection here.
}
