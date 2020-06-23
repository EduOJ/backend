package main

import (
	"github.com/leoleoasd/EduOJBackend/app/controllers"
	"github.com/leoleoasd/EduOJBackend/base"
)

func init() {
	base.E.GET("/login", controllers.Recv)

	base.E.GET("/admin", controllers.Recv)

	// TODO: routes.
}
