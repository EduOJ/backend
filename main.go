package main

import (
	"flag"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"os"
)

func main() {
	flag.Parse()
	file, err := os.Open(configFileFlag)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	config.ReadConfig(file)

}
