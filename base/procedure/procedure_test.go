package procedure

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestRegisterProcedure(t *testing.T) {
	t.Cleanup(func() {
		handlers = make(map[string]interface{})
	})
	tests := []struct {
		handler interface{}
		Args    []interface{}
		Results []interface{}
	}{
		{
			handler: func(a int, b string, c time.Time) (int, string, time.Time) {
				return a, b, c
			},
			Args:    append(make([]interface{}, 0), 1, "123", time.Date(2001, 2, 3, 4, 5, 6, 7, time.UTC)),
			Results: append(make([]interface{}, 0), 1, "123", time.Date(2001, 2, 3, 4, 5, 6, 7, time.UTC)),
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint("procedure_test_", i), func(t *testing.T) {
			t.Cleanup(func() {
				handlers = make(map[string]interface{})
			})
			RegisterProcedure("procedure_test", test.handler)
			result, err := CallProcedure("procedure_test", test.Args...)
			if err != nil {
				t.Error("Errors when calling procedure: ", err)
			}
			for i, v := range result {
				vv := reflect.ValueOf(v)
				assert.Equal(t, vv.Type(), reflect.TypeOf(test.Results[i]), "Type of result should be same.")
				assert.Equal(t, v, test.Results[i], "Value of result should be same.")
			}
		})
	}
}

func TestCallProcedure(t *testing.T) {
	t.Cleanup(func() {
		handlers = make(map[string]interface{})
	})
	RegisterProcedure("test_call_procedure", func() int {
		return 123
	})
	result, err := CallProcedure("test_call_procedure")
	assert.Equal(t, err, nil, "Should not have error.")
	assert.Equal(t, result[0], 123, "Should be the same.")

	RegisterProcedure("test_call_procedure_1", func(int) int {
		return 123
	})
	result, err = CallProcedure("test_call_procedure_1")
	assert.NotEqual(t, err, nil, "Should have error.")
	assert.Equal(t, err.Error(), "reflect: Call with too few input arguments", "Error should be too few arguments.")
	assert.Equal(t, result, []interface{}(nil), "Should not have result on error.")

	assert.PanicsWithValue(t, "Trying to register a non-func handler!", func() {
		RegisterProcedure("test_call_procedure_2", 123)
	}, "Should panic on non-function handlers.")
	assert.Panics(t, func() {
		RegisterProcedure("test_call_procedure_1", func(int) int {
			return 123
		})
	})
}
