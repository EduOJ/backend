// @title                       EduOJ Backend
// @version                     0.1.0
// @description                 The backend module for the EduOJ project.
// @BasePath                    /api
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
package main

import (
	"os"

	"github.com/EduOJ/backend/base/log"
	_ "github.com/EduOJ/backend/docs"
)

func main() {
	parse()
	if len(args) < 1 {
		log.Fatal("Please specific a command to run.")
		// TODO: output usage aa
		os.Exit(-1)
	}
	switch args[0] {
	case "test-config":
		testConfig()
	case "serve", "server", "http", "run":
		serve()
	case "migrate", "migration":
		doMigrate()
	case "clean", "clean-up", "clean-db":
		clean()
	case "permission", "perm":
		permission()
	}
}
