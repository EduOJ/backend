package config

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test read error")
}

func TestReadConfig(t *testing.T) {
	tests := []struct {
		string
		Node
	}{
		{
			``, &MapNode{M: map[string]Node{}},
		},
		{
			`intNode: 123
boolNode: false
stringNode: "123\"\""
mapNode:
  intNode: 123
sliceNode:
  - 123
  - "123\"\""
  - intNode: 123`, &MapNode{
				M: map[string]Node{
					"intNode":    IntNode(123),
					"boolNode":   BoolNode(false),
					"stringNode": StringNode("123\"\""),
					"mapNode": &MapNode{
						M: map[string]Node{
							"intNode": IntNode(123),
						},
					},
					"sliceNode": &SliceNode{
						S: []Node{
							IntNode(123), StringNode("123\"\""), &MapNode{
								M: map[string]Node{
									"intNode": IntNode(123),
								},
							},
						},
					},
				},
			},
		},
	}
	for id, test := range tests {
		t.Run(fmt.Sprint("testReadConfig", id), func(t *testing.T) {
			t.Cleanup(func() {
				conf = nil
			})
			assert.NotPanics(t, func() {
				buf := bytes.NewBufferString(test.string)
				err := ReadConfig(buf)
				assert.Equal(t, test.Node, conf)
				assert.Equal(t, nil, err)
			})
		})
	}
	t.Run("testReadConfigTwice", func(t *testing.T) {
		t.Cleanup(func() {
			conf = nil
		})
		assert.NotPanics(t, func() {
			buf := bytes.NewBufferString(tests[0].string)
			err := ReadConfig(buf)
			assert.Equal(t, tests[0].Node, conf)
			assert.Equal(t, nil, err)
		})
		assert.NotPanics(t, func() {
			buf := bytes.NewBufferString(tests[0].string)
			err := ReadConfig(buf)
			assert.Equal(t, tests[0].Node, conf)
			assert.NotEqual(t, nil, err)
			assert.Equal(t, "could not read config: already read!", err.Error())
		})
	})
}

func TestReadConfigFail(t *testing.T) {
	tests := []struct {
		string
		error
	}{
		{
			`\a\a\a\a`, errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `\\a\\a\\a\\a` into map[interface {}]interface {}"),
		},
		{
			`true: 123
false: 233`, errors.Wrap(ErrTypeDontMatchError, "could not build map node"),
		},
	}
	for id, test := range tests {
		t.Run(fmt.Sprint("testReadConfigFail", id), func(t *testing.T) {
			t.Cleanup(func() {
				conf = nil
			})
			assert.NotPanics(t, func() {
				buf := bytes.NewBufferString(test.string)
				err := ReadConfig(buf)
				assert.Equal(t, nil, conf)
				assert.NotEqual(t, nil, err)
				assert.Equal(t, test.error.Error(), err.Error())
			})
		})
	}
	t.Run("testReadError", func(t *testing.T) {
		assert.NotPanics(t, func() {
			err := ReadConfig(errReader(0))
			assert.Error(t, err)
			assert.Equal(t, "test read error", err.Error())
		})
	})
}
