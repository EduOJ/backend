package main

import (
	"github.com/leoleoasd/EduOJBackend/base"
)

func main() {
	base.Echo().Logger.Fatal(base.Echo().Start(":1323"))
}
