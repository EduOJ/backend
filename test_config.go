package main

import (
	"fmt"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/pkg/errors"
	"os"
)

func testConfig() {
	readConfig()
	initGorm(false)
	initLog()
	initRedis()
	initAuth()
	c, err := config.Get(".")
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not get the root of config file "+opt.Config))
		os.Exit(-1)
	}
	fmt.Println("config: ", c.String())
}
