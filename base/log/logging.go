package log

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
							writerConf.MustGet("level", "DEBUG").Value().(string))],
					})
				case "database":
					w := &databaseWriter{
						Level: stringToLevel[strings.ToUpper(
							writerConf.MustGet("level", "DEBUG").Value().(string))],
					}
					w.init()
					logger0.addWriter(w)
				case "event":
					// nothing to do.
				default:
					return errors.New("invalid writer name")
				}
			} else {
				return errors.New("invalid writer name")
			}
		}
		logger0.addWriter(&eventWriter{})
		logger0.setReady()
		return nil
	}
	return errors.New("log configuration should be an array")
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
