package main

import (
	"github.com/leoleoasd/EduOJBackend/app/controllers"
	"github.com/leoleoasd/EduOJBackend/base"
)

func registerRoutes() {
	base.Echo.GET("/login", controllers.Recv)

	base.Echo.GET("/admin", controllers.Recv)

	// TODO: routes.
}
