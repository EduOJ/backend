package event

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestEvent(t *testing.T) {
	t.Cleanup(func() {
		listeners = make(map[string][]interface{})
	})
	tests := []struct {
		Listener interface{}
		Args     []interface{}
		Results  []interface{}
	}{
		{
			Listener: func(a int, b string, c time.Time) (int, string, time.Time) {
				return a, b, c
			},
			Args:    append(make([]interface{}, 0), 1, "123", time.Date(2001, 2, 3, 4, 5, 6, 7, time.UTC)),
			Results: append(make([]interface{}, 0), 1, "123", time.Date(2001, 2, 3, 4, 5, 6, 7, time.UTC)),
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint("event_test_", i), func(t *testing.T) {
			t.Cleanup(func() {
				listeners = make(map[string][]interface{})
			})
			RegisterListener("event_test", test.Listener)
			result, err := FireEvent("event_test", test.Args...)
			if err != nil {
				t.Error("Errors when calling hook: ", err)
			}
			assert.Equal(t, len(result), 1, "Result length should be 1.")
			for i, v := range result[0] {
				vv := reflect.ValueOf(v)
				assert.Equal(t, vv.Type(), reflect.TypeOf(test.Results[i]), "Type of result should be same.")
				assert.Equal(t, v, test.Results[i], "Value of result should be same.")
			}
		})
	}
}

func TestFireEvent(t *testing.T) {
	t.Cleanup(func() {
		listeners = make(map[string][]interface{})
	})
	RegisterListener("test_fire_event", func() int {
		return 123
	})
	result, err := FireEvent("test_fire_event")
	assert.Equal(t, err, nil, "Should not have error.")
	assert.Equal(t, result[0][0], 123, "Should be the same.")

	RegisterListener("test_fire_event_1", func(int) int {
		return 123
	})
	result, err = FireEvent("test_fire_event_1")
	assert.NotEqual(t, err, nil, "Should have error.")
	assert.Equal(t, err.Error(), "reflect: Call with too few input arguments", "Error should be too few arguments.")
	assert.Equal(t, result, [][]interface{}(nil), "Should not have result on error.")

	assert.PanicsWithValue(t, "Trying to register a non-func listener!", func() {
		RegisterListener("test_fire_event_2", 123)
	}, "Should panic on non-function listeners.")
}
