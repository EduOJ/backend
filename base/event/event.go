package event

import (
	"errors"
	"reflect"
)

var listeners = map[string][]interface{}{}

// RegisterListener
func RegisterListener(eventName string, listener interface{}) {
	listenerValue := reflect.ValueOf(listener)
	if listenerValue.Kind() != reflect.Func {
		panic("Trying to register a non-func listener!")
	}
	listeners[eventName] = append(listeners[eventName], listener)
}

func FireEvent(eventName string, args ...interface{}) (result [][]interface{}, err error) {
	defer func() {
		if p := recover(); p != nil {
			result = nil
			err = errors.New(p.(string))
			return
		}
	}()
	argsValue := make([]reflect.Value, len(args))
	for i, arg := range args {
		argsValue[i] = reflect.ValueOf(arg)
	}
	result = make([][]interface{}, len(listeners[eventName]))
	for i, listener := range listeners[eventName] {
		listenerFunc := reflect.ValueOf(listener)
		rst := listenerFunc.Call(argsValue)
		result[i] = make([]interface{}, len(rst))
		for j, v := range rst {
			result[i][j] = v.Interface()
		}
	}
	return
}
