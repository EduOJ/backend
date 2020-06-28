package config

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
			assert.NotPanics(t, func() {
				buf := bytes.NewBufferString(test.string)
				node, err := ReadConfig(buf)
				assert.Equal(t, test.Node, node)
				assert.Equal(t, nil, err)
			})
		})
	}
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
			assert.NotPanics(t, func() {
				buf := bytes.NewBufferString(test.string)
				node, err := ReadConfig(buf)
				assert.Equal(t, nil, node)
				assert.NotEqual(t, nil, err)
				assert.Equal(t, test.error.Error(), err.Error())
			})
		})
	}
}
