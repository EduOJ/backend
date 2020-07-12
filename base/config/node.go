package config

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"sort"
	"strconv"
	"strings"
)

type Node interface {
	String() string
	Get(index string) (Node, error)
	MustGet(index string, def interface{}) Node
	Value() interface{}
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
	// To sort the map by key.
	keys := make([]string, len(m.M))
	i := 0
	for k := range m.M {
		keys[i] = k
		i = i + 1
	}
	sort.Strings(keys)
	for _, k := range keys {
		data, _ := json.Marshal(k)
		buf.Write(data)
		buf.WriteString(`:`)
		datas := m.M[k].String()
		buf.WriteString(datas)
		buf.WriteString(",")
	}
	buf.Bytes()[len(buf.Bytes())-1] = '}'
	return buf.Bytes(), nil
}

func (m *MapNode) Get(index string) (Node, error) {
	if index == "" {
		return m, nil
	}
	strs := strings.SplitN(index, ".", 2)
	if strs[0] == "" {
		// index != "" and strs[0] == ""
		// which means that len(strs) == 2.
		return m.Get(strs[1])
	}
	if c, ok := m.M[strs[0]]; ok {
		if len(strs) == 2 {
			return c.Get(strs[1])
		}
		return c.Get("")
	}
	return nil, ErrKeyNotFound
}

func (m *MapNode) Build(data interface{}) error {
	m.M = map[string]Node{}
	if mapData, ok := data.(map[interface{}]interface{}); ok {
		for k, v := range mapData {
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
	if v, err := m.Get(index); err == nil {
		return v
	}
	v, err := buildOne(def)
	if err != nil {
		panic(err)
	}
	return v
}

func (m *MapNode) Value() interface{} {
	return m.M
}

func (s *SliceNode) String() string {
	rst, _ := json.Marshal(s)
	return string(rst)
}

func (s *SliceNode) MarshalJSON() (data []byte, err error) {
	var buf bytes.Buffer
	buf.WriteString("[")
	for _, v := range s.S {
		data := v.String()
		buf.WriteString(data)
		buf.WriteString(",")
	}
	buf.Bytes()[len(buf.Bytes())-1] = ']'
	return buf.Bytes(), nil
}

func (s *SliceNode) Get(index string) (Node, error) {
	if index == "" {
		return s, nil
	}
	strs := strings.SplitN(index, ".", 2)
	if strs[0] == "" {
		// index != "" and strs[0] == ""
		// which means that len(strs) == 2.
		return s.Get(strs[1])
	}
	intIndex, err := strconv.Atoi(strs[0])
	if err != nil {
		return nil, err
	}
	if intIndex >= len(s.S) {
		return nil, ErrKeyNotFound
	}
	if len(strs) >= 2 {
		return s.S[intIndex].Get(strs[1])
	}
	return s.S[intIndex], nil
}

func (s *SliceNode) Build(data interface{}) error {
	if sliceData, ok := data.([]interface{}); ok {
		s.S = make([]Node, len(sliceData))
		for k, v := range sliceData {
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
	if v, err := s.Get(index); err == nil {
		return v
	}
	v, err := buildOne(def)
	if err != nil {
		panic(err)
	}
	return v
}

func (s *SliceNode) Value() interface{} {
	return s.S
}

func (s StringNode) String() string {
	data, _ := json.Marshal(s)
	return string(data)
}

func (s StringNode) Get(index string) (Node, error) {
	if index != "" {
		return nil, ErrNodeNotIndexable
	}
	return s, nil
}

func (s StringNode) MustGet(index_ string, def interface{}) Node {
	if index_ == "." || index_ == "" {
		return s
	}
	v, err := buildOne(def)
	if err != nil {
		panic(err)
	}
	return v
}

func (s StringNode) Value() interface{} {
	return string(s)
}

func (s IntNode) String() string {
	return strconv.Itoa(int(s))
}

func (s IntNode) Get(index string) (Node, error) {
	if index != "" {
		return nil, ErrNodeNotIndexable
	}
	return s, nil
}

func (s IntNode) MustGet(index_ string, def interface{}) Node {
	if index_ == "." || index_ == "" {
		return s
	}
	v, err := buildOne(def)
	if err != nil {
		panic(err)
	}
	return v
}

func (s IntNode) Value() interface{} {
	return int(s)
}

func (s BoolNode) String() string {
	if s {
		return "true"
	}
	return "false"
}

func (s BoolNode) Get(index string) (Node, error) {
	if index != "" {
		return nil, ErrNodeNotIndexable
	}
	return s, nil
}

func (s BoolNode) MustGet(index_ string, def interface{}) Node {
	if index_ == "." || index_ == "" {
		return s
	}
	v, err := buildOne(def)
	if err != nil {
		panic(err)
	}
	return v
}

func (s BoolNode) Value() interface{} {
	return bool(s)
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
