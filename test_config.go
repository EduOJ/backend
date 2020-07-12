package main

import (
	"fmt"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/logging"
	"github.com/pkg/errors"
	"os"
)

func testConfig() {
	readConfig()
	c, err := config.Get(".")
	if err != nil {
		logging.Fatal(errors.Wrap(err, "could not get the root of config file "+opt.Config))
		os.Exit(-1)
	}
	fmt.Println("config: ", c.String())
	os.Exit(0)
}
