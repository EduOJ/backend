package procedure

import (
	"github.com/pkg/errors"
	"reflect"
)

var handlers = map[string]interface{}{}

// RegisterProcedure registers a procedure handlers with given name.
// The handler should be a function.
func RegisterProcedure(procedureName string, handler interface{}) {
	listenerValue := reflect.ValueOf(handler)
	if listenerValue.Kind() != reflect.Func {
		panic("Trying to register a non-func handler!")
	}
	if _, ok := handlers[procedureName]; ok {
		panic("Trying to re-register a handler!")
	}
	handlers[procedureName] = handler
}

// CallProcedure calls the procedure with given name, providing the procedure with given args and returns what the procedure returns as return value.
func CallProcedure(procedureName string, args ...interface{}) (result []interface{}, err error) {
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
	handlerFunc := reflect.ValueOf(handlers[procedureName])
	rst := handlerFunc.Call(argsValue)
	result = make([]interface{}, len(rst))
	for i, v := range rst {
		result[i] = v.Interface()
	}
	return
}
