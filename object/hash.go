package object

import (
	"bytes"
	"fmt"
	"strings"
)

type Map struct {
	Items map[string]Object
}

func (m *Map) Type() Type {
	return HASH
}

func (m *Map) Inspect() string {
	var out bytes.Buffer
	pairs := make([]string, 0)
	for k, v := range m.Items {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

func (m *Map) InvokeMethod(method string, args ...Object) Object {
	switch method {
	case "keys":
		return m.Keys()
	case "values":
		return m.Values()
	}
	return NewError("type error: %s object has no method %s", m.Type(), method)
}

func (m *Map) Keys() *List {
	items := make([]Object, 0, len(m.Items))
	for k := range m.Items {
		items = append(items, &String{Value: k})
	}
	return &List{Items: items}
}

func (m *Map) Values() *List {
	items := make([]Object, 0, len(m.Items))
	for _, v := range m.Items {
		items = append(items, v)
	}
	return &List{Items: items}
}

func (m *Map) GetWithObject(key *String) Object {
	value, found := m.Items[key.Value]
	if !found {
		return Null
	}
	return value
}

func (m *Map) Get(key string) Object {
	value, found := m.Items[key]
	if !found {
		return Null
	}
	return value
}

func (m *Map) Delete(key string) Object {
	delete(m.Items, key)
	return Null
}

func (m *Map) Set(key string, value Object) {
	m.Items[key] = value
}

func (m *Map) Size() int {
	return len(m.Items)
}

func (m *Map) ToInterface() interface{} {
	result := make(map[string]any, len(m.Items))
	for k, v := range m.Items {
		result[k] = v.ToInterface()
	}
	return result
}

func NewMap(m map[string]interface{}) *Map {
	result := &Map{Items: make(map[string]Object, len(m))}
	for k, v := range m {
		value := FromGoType(v)
		if value == nil {
			panic(fmt.Sprintf("type error: cannot convert %v to a tamarin object", v))
		}
		result.Items[k] = value
	}
	return result
}
