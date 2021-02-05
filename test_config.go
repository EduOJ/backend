package main

import (
	"github.com/leoleoasd/EduOJBackend/base/log"
)

func testConfig() {
	// TODO: test config using config.Get
	readConfig()
	initGorm(false)
	initLog()
	initRedis()
	initStorage()
	log.Fatalf("should success.")
}
