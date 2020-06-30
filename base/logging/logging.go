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

var Debug = logger0.Debug
var Info = logger0.Info
var Warning = logger0.Warning
var Error = logger0.Error
var Fatal = logger0.Fatal
var Debugf = logger0.Debugf
var Infof = logger0.Infof
var Warningf = logger0.Warningf
var Errorf = logger0.Errorf
var Fatalf = logger0.Fatalf
