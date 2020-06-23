package event

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func cleanup() {
	listeners = make(map[string][]interface{})
}

func TestEvent(t *testing.T) {
	defer cleanup()

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
			defer cleanup()
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
