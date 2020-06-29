package main

import "flag"

var configFileFlag string

func init() {
	flag.StringVar(&configFileFlag, "c", "config.yml", "Config file name")
	flag.StringVar(&configFileFlag, "config", "config.yml", "Config file name")

}
