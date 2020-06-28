package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

func ReadConfig(file io.Reader) (Node, error) {
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	config := make(map[interface{}]interface{})
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}
	ret := &MapNode{}
	err = ret.Build(config)
	if err != nil {
		return nil, errors.Wrap(err, "could not build map node")
	}
	return ret, nil
}
