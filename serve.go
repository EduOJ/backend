package main

import (
	"context"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/logging"
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
	logging.Fatal("Server started at port ", base.Echo.Server.Addr)
	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	<-s
	logging.Fatal("Server closing.")
	logging.Fatal("Hit ctrl+C again to forceShutdown quit.")
	shutdownCtx, forceShutdown := context.WithCancel(context.Background())
	go func() {
		<-s
		forceShutdown()
	}()
	err := base.Echo.Shutdown(shutdownCtx)
	if err != nil {
		if err.Error() == "context canceled" {
			logging.Fatal("Force quitting.")
		} else {
			logging.Fatal(err)
		}
	} else {
		logging.Fatal("Server closed. Quitting.")
	}
}
