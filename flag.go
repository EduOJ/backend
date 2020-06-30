package main

import (
	"github.com/jessevdk/go-flags"
)

type _opt struct {
	Config string `short:"c" long:"verbose" description:"Config file name" default:"config.yml"`
}

var parser *flags.Parser
var opt _opt

func init() {
	parser = flags.NewNamedParser("eduOJ server", flags.HelpFlag | flags.PassDoubleDash)
	parser.AddGroup("Application", "Application Options", &opt)
}
