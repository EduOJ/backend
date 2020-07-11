package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

var conf Node

// Get gets a value from the config object
// Example:
// value, err := conf.get("logging.0.level")
var Get func(index string) (Node, error)

// MustGet warps the Get function, when error accounts, returns the default value.
var MustGet func(index string, def interface{}) Node

func ReadConfig(file io.Reader) error {
	if conf != nil {
		return errors.New("could not read config: already read!")
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	config := make(map[interface{}]interface{})
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return err
	}
	ret := &MapNode{}
	err = ret.Build(config)
	if err != nil {
		return errors.Wrap(err, "could not build map node")
	}
	conf = ret
	Get = conf.Get
	MustGet = conf.MustGet
	return nil
}
