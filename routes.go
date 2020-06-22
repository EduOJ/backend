package main

import (
	"github.com/leoleoasd/EduOJBackend/app/controllers"
	"github.com/leoleoasd/EduOJBackend/base"
)

func init() {
	base.E.GET("/", controllers.Recv)
	// TODO: routes.
}
