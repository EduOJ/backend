package utils

import (
	"github.com/spf13/viper"
)

var Origins []string

func InitOrigin() {
	var origin []string
	viper.UnmarshalKey("server.origin", &origin)
	for _, v := range origin {
		Origins = append(Origins, v)
	}
}
