package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/logging"
	"github.com/pkg/errors"
	"os"
)

type _opt struct {
	Config string `short:"c" long:"verbose" description:"Config file name" default:"config.yml"`
}

var parser *flags.Parser
var opt _opt
var open = os.Open

func init() {
	parser = flags.NewNamedParser("eduOJ server", flags.HelpFlag|flags.PassDoubleDash)
	_, _ = parser.AddGroup("Application", "Application Options", &opt)
}

func parse() {
	// TODO: remove useless logs for parser debugging.
	logging.Debug("Parsing command-line arguments.")
	args, err := parser.Parse()
	if err != nil {
		logging.Fatal(errors.Wrap(err, "could not parse argument "))
		os.Exit(-1)
	}
	logging.Debug(args, err)
	logging.Debug(opt)
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
}
