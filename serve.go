package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/EduOJ/backend/base/exit"
	"github.com/EduOJ/backend/base/log"
)

func serve() {
	readConfig()
	initGorm()
	initLog()
	initRedis()
	initStorage()
	initWebAuthn()
	initMail()
	initEvent()
	startEcho()
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-s

	go func() {
		<-s
		log.Fatal("Force quitting")
		os.Exit(-1)
	}()

	log.Fatal("Server closing.")
	log.Fatal("Hit ctrl+C again to force quit.")
	exit.Close()
	exit.QuitWG.Wait()
}
