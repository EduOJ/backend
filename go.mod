module github.com/leoleoasd/EduOJBackend

go 1.14

require (
	github.com/fatih/color v1.9.0
	github.com/gabriel-vasile/mimetype v1.1.1
	github.com/go-gormigrate/gormigrate/v2 v2.0.0
	github.com/go-ini/ini v1.60.2 // indirect
	github.com/go-playground/assert/v2 v2.0.1
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.3.0
	github.com/go-redis/redis/v8 v8.0.0-beta.6
	github.com/jessevdk/go-flags v1.4.0
	github.com/johannesboyne/gofakes3 v0.0.0-20200716060623-6b2b4cb092cc
	github.com/kami-zh/go-capturer v0.0.0-20171211120116-e492ea43421d
	github.com/labstack/echo/v4 v4.1.16
	github.com/labstack/gommon v0.3.0
	github.com/lib/pq v1.7.1 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/minio/minio-go v6.0.14+incompatible
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/valyala/fasttemplate v1.2.0 // indirect
	github.com/xlab/treeprint v1.0.0
	go.opentelemetry.io/otel v0.9.0 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200831180312-196b9ba8737a // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	gorm.io/driver/mysql v1.0.1
	gorm.io/driver/postgres v1.0.0
	gorm.io/driver/sqlite v1.1.1
	gorm.io/gorm v1.20.11
)

replace github.com/stretchr/testify v1.6.1 => github.com/leoleoasd/testify v1.6.2-0.20200818074144-885db91dbfe9
