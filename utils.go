package main

import (
	"context"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/database"
	"github.com/pkg/errors"
	"os"
)

var baseContext context.Context
var cancel context.CancelFunc

func init() {
	baseContext, cancel = context.WithCancel(context.Background())
}

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

func initEcho() {
	log.Debug("Starting echo server.")
	echoConf, err := config.Get("server")
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not read http server config"))
		os.Exit(-1)
	}
	err = base.InitEchoFromConfig(echoConf)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not init http server with config "+echoConf.String()))
		os.Exit(-1)
	}
}

func initRedis() {
	log.Debug("Starting redis client.")
	redisConf, err := config.Get("redis")
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not read redis config"))
		os.Exit(-1)
	}
	err = base.InitRedisFromConfig(redisConf)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not init redis with config "+redisConf.String()))
		os.Exit(-1)
	}
}

func initGorm(toMigrate ...bool) {
	log.Debug("Starting database client.")
	databaseConf, err := config.Get("database")
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not read database config"))
		os.Exit(-1)
	}
	err = base.InitGormFromConfig(databaseConf)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not init database with config "+databaseConf.String()))
		os.Exit(-1)
	}
	if len(toMigrate) == 0 || toMigrate[0] {
		database.Migrate()
	}
}
