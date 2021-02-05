package log

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
)

var logger0 _logger = &logger{}

type logConf []struct {
	Name  string
	Level string
}

func InitFromConfig() (err error) {
	if logger0.isReady() {
		return errors.New("already initialized")
	}
	var confSlice logConf
	err = viper.UnmarshalKey("log", &confSlice)
	if err != nil {
		return errors.Wrap(err, "Wrong log conf")
	}
	for _, c := range confSlice {
		switch c.Name {
		case "console":
			logger0.addWriter(&consoleWriter{
				Level: StringToLevel[strings.ToUpper(c.Level)],
			})
		case "database":
			w := &databaseWriter{
				Level: StringToLevel[strings.ToUpper(c.Level)],
			}
			w.init()
			logger0.addWriter(w)
		case "event":
			// nothing to do.
		default:
			return errors.New("invalid writer name")
		}
	}
	logger0.addWriter(&eventWriter{})
	logger0.setReady()
	return nil
}

func Debug(items ...interface{}) {
	logger0.Debug(items...)
}

func Info(items ...interface{}) {
	logger0.Info(items...)
}

func Warning(items ...interface{}) {
	logger0.Warning(items...)
}

func Error(items ...interface{}) {
	logger0.Error(items...)
}

func Fatal(items ...interface{}) {
	logger0.Fatal(items...)
}

func Debugf(format string, items ...interface{}) {
	logger0.Debugf(format, items...)
}

func Infof(format string, items ...interface{}) {
	logger0.Infof(format, items...)
}

func Warningf(format string, items ...interface{}) {
	logger0.Warningf(format, items...)
}

func Errorf(format string, items ...interface{}) {
	logger0.Errorf(format, items...)
}

func Fatalf(format string, items ...interface{}) {
	logger0.Fatalf(format, items...)
}

func Disable() {
	logger0.Disable()
}

func Disabled() bool {
	return logger0.Disabled()
}
