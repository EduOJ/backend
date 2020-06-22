package base

import (
	"errors"
	"reflect"
)

var listeners = map[string][]interface{}{}

func RegisterListener(eventName string, listener interface{}) {
	listenerValue := reflect.ValueOf(listener)
	if listenerValue.Kind() != reflect.Func {
		panic("Trying to register a non-func listener!")
	}
	listeners[eventName] = append(listeners[eventName], listener)
}

func FireEvent(hookName string, args ...interface{}) (result [][]interface{}, err error) {
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
	result = make([][]interface{}, len(listeners[hookName]))
	for i, hook := range listeners[hookName] {
		hookFunc := reflect.ValueOf(hook)
		rst := hookFunc.Call(argsValue)
		result[i] = make([]interface{}, len(rst))
		for j, v := range rst {
			result[i][j] = v.Interface()
		}
	}
	return
}
