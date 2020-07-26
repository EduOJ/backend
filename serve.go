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
	startEcho()
	initAuth()
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-s
	log.Fatal("Server closing.")
	log.Fatal("Hit ctrl+C again to force quit.")
	exit.Close()
	go func() {
		<-s
		os.Exit(-1)
	}()
	exit.QuitWG.Wait()
}
