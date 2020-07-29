package main

import (
	"github.com/leoleoasd/EduOJBackend/base/exit"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"os"
)

func clean() {
	readConfig()
	initGorm()
	initLog()
	err := utils.CleanUpExpiredTokens()
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	exit.Close()
	exit.QuitWG.Wait()
	log.Fatal("Clean succeed!")
}
