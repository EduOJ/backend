package utils

import (
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/log"
)

var Origins []string

func InitOrigin() {
	if n, err := config.Get("server.origin"); err == nil {
		for _, v := range n.(*config.SliceNode).S {
			if vv, ok := v.Value().(string); ok {
				Origins = append(Origins, vv)
			} else {
				log.Fatal("Illegal origin name" + v.String())
				panic("Illegal origin name" + v.String())
			}
		}
	} else {
		log.Fatal("Illegal origin config", err)
		panic("Illegal origin config")
	}
}
