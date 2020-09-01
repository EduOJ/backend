package main

import (
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/pkg/errors"
	"os"
)

func testConfig() {
	//TODO: test config using config.Get
	readConfig()
	initGorm(false)
	initLog()
	initRedis()
	initStorage()
	c, err := config.Get(".")
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not get the root of config file "+opt.Config))
		os.Exit(-1)
	}
	log.Debug("config: ", c.String())
}
