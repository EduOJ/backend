package config

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMapNode_String(t *testing.T) {
	m := &MapNode{
		M: map[string]Node{
			"intNode":    IntNode(123),
			"stringNode": StringNode("123\"\""),
			"boolNode":   BoolNode(false),
			"boolNode2":  BoolNode(true),
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
	}
	str := m.String()
	assert.Equal(t, `{"boolNode":false,"boolNode2":true,"intNode":123,"mapNode":{"intNode":123},"sliceNode":[123,"123\"\"",{"intNode":123}],"stringNode":"123\"\""}`, str)
}

func TestMapNode_Child(t *testing.T) {
	m := &MapNode{
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
	}
	tests := []struct {
		Key   string
		Value Node
		Err   error
	}{
		{
			"intNode", IntNode(123), nil,
		},
		{
			"boolNode", BoolNode(false), nil,
		},
		{
			"stringNode", StringNode("123\"\""), nil,
		},
		{
			"mapNode", &MapNode{
				M: map[string]Node{
					"intNode": IntNode(123),
				},
			}, nil,
		},
		{
			"mapNode.intNode", IntNode(123), nil,
		},
		{
			"sliceNode", &SliceNode{
				S: []Node{
					IntNode(123),
					StringNode("123\"\""),
					&MapNode{
						M: map[string]Node{
							"intNode": IntNode(123),
						},
					},
				},
			}, nil,
		},
		{
			"sliceNode.0", IntNode(123), nil,
		},
		{
			"sliceNode.1", StringNode("123\"\""), nil,
		},
		{
			"sliceNode.2", &MapNode{
				M: map[string]Node{
					"intNode": IntNode(123),
				},
			}, nil,
		},
		{
			"sliceNode.2.intNode", IntNode(123), nil,
		},
		{
			"sliceNode.2.intNode", IntNode(123), nil,
		},
		{
			".", m, nil,
		},
		{
			"sliceNode...2...intNode.", IntNode(123), nil,
		},
		{
			"123", nil, ErrKeyNotFound,
		},
		{
			"sliceNode.123", nil, ErrKeyNotFound,
		},
		{
			"sliceNode.2.intNode.123", nil, ErrNodeNotIndexable,
		},
		{
			"sliceNode.1.123", nil, ErrNodeNotIndexable,
		},
		{
			"boolNode.123", nil, ErrNodeNotIndexable,
		},
		{
			"sliceNode.intNode", nil, errors.New("strconv.Atoi: parsing \"intNode\": invalid syntax"),
		},
	}
	for _, test := range tests {
		t.Run("getWithKey:"+test.Key, func(t *testing.T) {
			v, err := m.Get(test.Key)
			assert.Equal(t, test.Value, v)
			if err == nil || test.Err == nil {
				assert.Equal(t, test.Err, err)
			} else {
				assert.Equal(t, test.Err.Error(), err.Error())
			}
		})
		t.Run("mustGetWithKey:"+test.Key, func(t *testing.T) {
			if test.Err == nil {
				assert.NotPanics(t, func() {
					v := m.MustGet(test.Key, "")
					assert.Equal(t, test.Value, v)
				})
			} else {
				assert.NotPanics(t, func() {
					v := m.MustGet(test.Key, "")
					assert.Equal(t, StringNode(""), v)
				})
			}
		})
	}
	t.Run("mustGetWithInvalidDefault", func(t *testing.T) {
		assert.Panics(t, func() {
			m.MustGet("asdasdasdasd", struct{}{})
		})
	})
	t.Run("mustGet", func(t *testing.T) {
		assert.NotPanics(t, func() {
			v := m.MustGet("mapNode", "").MustGet(".", "")
			assert.Equal(t, v, m.MustGet("mapNode", ""))
			v = m.MustGet("mapNode", "").MustGet(".123", "")
			assert.Equal(t, v, StringNode(""))
		})
		assert.Panics(t, func() {
			m.MustGet("mapNode", "").MustGet(".123", struct{}{})
		})
		assert.NotPanics(t, func() {
			v := m.MustGet("sliceNode", "").MustGet(".", "")
			assert.Equal(t, v, m.MustGet("sliceNode", ""))
			v = m.MustGet("sliceNode", "").MustGet(".123", "")
			assert.Equal(t, v, StringNode(""))
		})
		assert.Panics(t, func() {
			m.MustGet("sliceNode", "").MustGet(".123", struct{}{})
		})
		assert.NotPanics(t, func() {
			v := m.MustGet("intNode", "").MustGet(".", "")
			assert.Equal(t, v, m.MustGet("intNode", ""))
			v = m.MustGet("intNode", "").MustGet(".123", "")
			assert.Equal(t, v, StringNode(""))
		})
		assert.Panics(t, func() {
			m.MustGet("intNode", "").MustGet(".123", struct{}{})
		})
		assert.NotPanics(t, func() {
			v := m.MustGet("stringNode", "").MustGet(".", "")
			assert.Equal(t, v, m.MustGet("stringNode", ""))
			v = m.MustGet("stringNode", "").MustGet(".123", "")
			assert.Equal(t, v, StringNode(""))
		})
		assert.Panics(t, func() {
			m.MustGet("stringNode", "").MustGet(".123", struct{}{})
		})
		assert.NotPanics(t, func() {
			v := m.MustGet("boolNode", "").MustGet(".", "")
			assert.Equal(t, v, m.MustGet("boolNode", ""))
			v = m.MustGet("boolNode", "").MustGet(".123", "")
			assert.Equal(t, v, StringNode(""))
		})
		assert.Panics(t, func() {
			m.MustGet("boolNode", "").MustGet(".123", struct{}{})
		})
	})
}

func TestNodeValue(t *testing.T) {
	tests := []struct {
		Node  Node
		Value interface{}
	}{
		{
			&MapNode{M: map[string]Node{
				"123": IntNode(123),
			}}, map[string]Node{
				"123": IntNode(123),
			},
		}, {
			&SliceNode{[]Node{
				IntNode(123),
			}}, []Node{
				IntNode(123),
			},
		}, {
			IntNode(123),123,
		},{
			BoolNode(false),false,
		},{
			StringNode("123"),"123",
		},
	}
	for _, test := range tests {
		tt := reflect.TypeOf(test.Node)
		t.Run("test"+tt.Name()+"Value", func(t *testing.T) {
			assert.Equal(t, test.Value, test.Node.Value())
		})
	}
}
