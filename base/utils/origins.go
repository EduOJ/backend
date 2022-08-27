package utils

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var Origins []string

func InitOrigin() {
	err := viper.UnmarshalKey("server.origin", &Origins)
	if err != nil {
		panic(errors.Wrap(err, "could not read server.origin from config"))
	}
}
