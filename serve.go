package main

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"os"
	"os/signal"
	"syscall"
)

func serve() {
	readConfig()
	initLog()
	initRedis()
	initGorm()
	// TODO: init database
	initEcho()

	go base.Echo.StartServer(base.Echo.Server)
	log.Fatal("Server started at port ", base.Echo.Server.Addr)
	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	<-s
	log.Fatal("Server closing.")
	log.Fatal("Hit ctrl+C again to forceShutdown quit.")
	go func() {
		<-s
		cancel()
	}()
	err := base.Echo.Shutdown(baseContext)
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
