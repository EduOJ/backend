package main

import (
	"github.com/leoleoasd/EduOJBackend/base/exit"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"os"
	"os/signal"
	"syscall"
)

func serve() {
	readConfig()
	initGorm()
	initLog()
	initRedis()
	initStorage()
	initWebAuthn()
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
