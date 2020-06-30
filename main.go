package main

import (
	"github.com/leoleoasd/EduOJBackend/base/logging"
)

func main() {

	args, err := parser.Parse()
	logging.Debug(args, err)
	logging.Debug(opt)

}
