package main

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/log"
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
