package main

import (
	"github.com/leoleoasd/EduOJBackend/base/log"
	"os"
)

func main() {
	parse()
	if len(args) < 1 {
		log.Fatal("Please specific a command to run.")
		// TODO: output usage
		os.Exit(-1)
	}
	switch args[0] {
	case "test-config":
		testConfig()
	case "serve", "server", "http", "run":
		serve()
	}
}
