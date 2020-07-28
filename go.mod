module github.com/leoleoasd/EduOJBackend

go 1.14

require (
	github.com/go-playground/assert/v2 v2.0.1
	github.com/go-playground/validator/v10 v10.3.0
	github.com/go-redis/redis/v8 v8.0.0-beta.6
	github.com/jessevdk/go-flags v1.4.0
	github.com/jinzhu/gorm v1.9.15
	github.com/kami-zh/go-capturer v0.0.0-20171211120116-e492ea43421d
	github.com/labstack/echo/v4 v4.1.16
	github.com/labstack/gommon v0.3.0
	github.com/lib/pq v1.7.1 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	github.com/valyala/fasttemplate v1.2.0 // indirect
	go.opentelemetry.io/otel v0.9.0 // indirect
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sys v0.0.0-20200722175500-76b94024e4b6 // indirect
	golang.org/x/text v0.3.3 // indirect
	gopkg.in/gormigrate.v1 v1.6.0
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/stretchr/testify v1.6.1 => github.com/leoleoasd/testify v1.6.2-0.20200728091548-dbfc7ee10e01
