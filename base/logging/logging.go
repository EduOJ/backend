package logging

import (
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/pkg/errors"
	"strings"
)

var logger0 _logger = &logger{}

func InitFromConfig(logConf config.Node) (err error) {
	if logger0.isReady() {
		return errors.New("already initialized")
	}
	if sliceNode, ok := logConf.(*config.SliceNode); ok {
		for _, _i := range sliceNode.S {
			writerConf, ok := _i.(*config.MapNode)
			if !ok {
				return errors.New("writer configuration should be a map")
			}
			_name, err := writerConf.Get("name")
			if err != nil {
				return errors.New("writer configuration should contain name")
			}
			if name, ok := _name.(config.StringNode); ok {
				switch name {
				case "console":
					logger0.addWriter(&consoleWriter{
						Level: stringToLevel[strings.ToUpper(
							string(writerConf.MustGet("level", "DEBUG").(config.StringNode)))],
					})
				case "database":
					// TODO
				case "event":
					// nothing to do.
				default:
					return errors.New("invalid writer name")
				}
			} else {
				return errors.New("invalid writer name")
			}
		}
		logger0.setReady()
		return nil
	}
	return errors.New("log configuration should be an array")
}

func Debug(args ...interface{}) {
	logger0.Debug(args...)
}

func Info(args ...interface{}) {
	logger0.Info(args...)
}

func Warning(args ...interface{}) {
	logger0.Warning(args...)
}

func Error(args ...interface{}) {
	logger0.Error(args...)
}

func Fatal(args ...interface{}) {
	logger0.Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	logger0.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger0.Infof(format, args...)
}

func Warningf(format string, args ...interface{}) {
	logger0.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger0.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger0.Fatalf(format, args...)
}
