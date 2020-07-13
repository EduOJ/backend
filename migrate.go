package main

import (
	"github.com/leoleoasd/EduOJBackend/base/log"
)

func doMigrate() {
	readConfig()
	initGorm()
	initLog()
	initRedis()
	log.Fatal("Migrate succeed!")
}
