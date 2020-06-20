package main

import (
	"github.com/leoleoasd/eduoj/backend/app/controllers"
	"github.com/leoleoasd/eduoj/backend/base"
)

func init() {
	base.E.GET("/", controllers.Root)
	// TODO: routes.
}
