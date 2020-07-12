package main

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/logging"
	"github.com/pkg/errors"
	"os"
)

func readConfig() {
	logging.Debug("Reading config.")
	configFile, err := open(opt.Config)
	if err != nil {
		logging.Fatal(errors.Wrap(err, "could not open config file "+opt.Config))
		os.Exit(-1)
	}
	err = config.ReadConfig(configFile)
	if err != nil {
		logging.Fatal(errors.Wrap(err, "could not read config file "+opt.Config))
		os.Exit(-1)
	}
	logging.Debug("Config read.")
}

func initLog() {
	logging.Debug("Initializing log.")
	loggingConf, err := config.Get("log")
	if err != nil {
		logging.Fatal(errors.Wrap(err, "could not read log config"))
		os.Exit(-1)
	}
	err = logging.InitFromConfig(loggingConf)
	if err != nil {
		logging.Fatal(errors.Wrap(err, "could not init log with config "+loggingConf.String()))
		os.Exit(-1)
	}
	logging.Debug("Logging initialized.")
}

func initEcho() {
	logging.Debug("Starting echo server.")
	echoConf, err := config.Get("server")
	if err != nil {
		logging.Fatal(errors.Wrap(err, "could not read http server config"))
		os.Exit(-1)
	}
	err = base.InitEchoFromConfig(echoConf)
	if err != nil {
		logging.Fatal(errors.Wrap(err, "could not init http server with config "+echoConf.String()))
		os.Exit(-1)
	}
}
