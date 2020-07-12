package main

import (
	"context"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"os"
	"os/signal"
	"syscall"
)

func serve() {
	readConfig()
	initLog()
	// TODO: init database
	initEcho()

	go base.Echo.StartServer(base.Echo.Server)
	log.Fatal("Server started at port ", base.Echo.Server.Addr)
	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	<-s
	log.Fatal("Server closing.")
	log.Fatal("Hit ctrl+C again to forceShutdown quit.")
	shutdownCtx, forceShutdown := context.WithCancel(context.Background())
	go func() {
		<-s
		forceShutdown()
	}()
	err := base.Echo.Shutdown(shutdownCtx)
	if err != nil {
		if err.Error() == "context canceled" {
			log.Fatal("Force quitting.")
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Server closed. Quitting.")
	}
}
