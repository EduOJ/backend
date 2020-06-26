package config

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type Node interface {
	String() string
	Build(interface{}) error
	Child(index string) (Node, error)
	MustGet(index string, def interface{}) Node
}

type MapNode struct {
	M map[string]Node
}
type SliceNode struct {
	S []Node
}
type StringNode string
type IntNode int
type BoolNode bool
var ErrKeyNotFound = errors.New("key not found in node")
var ErrNodeNotIndexable = errors.New("node is not indexable")
var ErrTypeDontMatchError = errors.New("given type don't match with node type")
var ErrIllegalTypeError = errors.New("illegal type")


func (m *MapNode) String() string {
	rst, _ := json.Marshal(m)
	return string(rst)
}

func (m *MapNode) MarshalJSON() (data []byte, err error) {
	var buf bytes.Buffer
	buf.WriteString("{")
	for k, v := range m.M {
		data, err := json.Marshal(k)
		if err != nil {
			return nil, errors.Wrap(err, "could not marshal map key")
		}
		buf.Write(data)
		buf.WriteString(`:`)
		data, err = json.Marshal(v)
		if err != nil {
			return nil, errors.Wrap(err, "could not marshal map value")
		}
		buf.Write(data)
		buf.WriteString(",")
	}
	buf.Bytes()[len(buf.Bytes())-1] = '}'
	return buf.Bytes(), nil
}

func (m *MapNode) Child(index string) (Node, error) {
	if index == "" {
		return m, nil
	}
	strs := strings.SplitN(index, ".", 2)
	if c, ok := m.M[strs[0]]; ok {
		return c.Child(strs[1])
	}
	return nil, ErrKeyNotFound
}

func (m *MapNode) Build(data interface{}) error {
	m.M = map[string]Node{}
	if mapData, ok := data.(map[interface{}]interface{}); ok {
		for k,v := range mapData {
			if _, ok := k.(string); !ok {
				return ErrTypeDontMatchError
			}
			t, err := buildOne(v)
			if err != nil {
				return err
			}
			m.M[k.(string)] = t
		}
		return nil
	}
	return ErrTypeDontMatchError
}

func (m *MapNode) MustGet(index string, def interface{}) Node {
	if v, err := m.Child(index); err != nil {
		return v
	}
	return def.(Node)
}

func (s *SliceNode) String() string {
	rst, _ := json.Marshal(s)
	return string(rst)
}

func (s *SliceNode) MarshalJSON() (data []byte, err error) {
	var buf bytes.Buffer
	buf.WriteString("[")
	for _, v := range s.S {
		data, err = json.Marshal(v)
		if err != nil {
			return nil, errors.Wrap(err, "could not marshal map value")
		}
		buf.Write(data)
		buf.WriteString(",")
	}
	buf.Bytes()[len(buf.Bytes())-1] = ']'
	return buf.Bytes(), nil
}


func (s *SliceNode) Child(index string) (Node, error) {
	if index == "" {
		return s, nil
	}
	strs := strings.SplitN(index, ".", 2)
	intIndex, err := strconv.Atoi(strs[0])
	if err != nil {
		return nil, err
	}
	if intIndex >= len(strs) {
		return nil, ErrKeyNotFound
	}
	return s.S[intIndex], nil
}

func (s *SliceNode) Build(data interface{}) error {
	if sliceData, ok := data.([]interface{}); ok {
		s.S = make([]Node, len(sliceData))
		for k,v := range sliceData {
			vv, err := buildOne(v)
			if err != nil {
				return err
			}
			s.S[k] = vv
		}
		return nil
	}
	return ErrTypeDontMatchError
}

func (s *SliceNode) MustGet(index string, def interface{}) Node {
	if v, err := s.Child(index); err != nil {
		return v
	}
	return def.(Node)
}


func (s StringNode) String() string {
	return string(s)
}

func (s StringNode) Child(index string) (Node, error) {
	return nil, ErrNodeNotIndexable
}

func (s StringNode) Build(data interface{}) error {
	if stringData, ok := data.(string); ok {
		s = StringNode(stringData)
		return nil
	}
	return ErrTypeDontMatchError
}

func (s StringNode) MustGet(index_ string, def interface{}) Node {
	return def.(Node)
}

func (s IntNode) String() string {
	return string(s)
}

func (s IntNode) Child(index string) (Node, error) {
	return nil, ErrNodeNotIndexable
}

func (s IntNode) Build(data interface{}) error {
	if intData, ok := data.(int); ok {
		s = IntNode(intData)
		return nil
	}
	return ErrTypeDontMatchError
}

func (s IntNode) MustGet(index_ string, def interface{}) Node {
	return def.(Node)
}

func (s BoolNode) String() string {
	if s {
		return "true"
	}
	return "false"
}

func (s BoolNode) Child(index string) (Node, error) {
	return nil, ErrNodeNotIndexable
}

func (s BoolNode) Build(data interface{}) error {
	if boolData, ok := data.(bool); ok {
		s = BoolNode(boolData)
		return nil
	}
	return ErrTypeDontMatchError
}

func (s BoolNode) MustGet(index_ string, def interface{}) Node {
	return def.(Node)
}


func buildOne(data interface{}) (Node, error) {
	switch data.(type) {
	case map[interface{}]interface{}:
		v := &MapNode{}
		err := v.Build(data)
		if err != nil {
			return nil, err
		}
		return v, nil
	case []interface{}:
		v := &SliceNode{}
		err := v.Build(data)
		if err != nil {
			return nil, err
		}
		return v, nil
	case int:
		v := IntNode(data.(int))
		return v, nil
	case string:
		v := StringNode(data.(string))
		return v, nil
	case bool:
		v := BoolNode(data.(bool))
		return v, nil
	default:
		return nil, ErrIllegalTypeError
	}

}
