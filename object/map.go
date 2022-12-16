package object

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

type Map struct {
	Items map[string]Object
}

func (m *Map) Type() Type {
	return MAP
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

func (m *Map) GetAttr(name string) (Object, bool) {
	switch name {
	case "keys":
		return &Builtin{
			Name: "map.keys",
			Fn: func(ctx context.Context, args ...Object) Object {
				return m.Keys()
			},
		}, true
	case "values":
		return &Builtin{
			Name: "map.values",
			Fn: func(ctx context.Context, args ...Object) Object {
				return m.Values()
			},
		}, true
	case "get":
		return &Builtin{
			Name: "map.get",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) < 1 || len(args) > 2 {
					return NewArgsError("map.get", 1, len(args))
				}
				key, err := AsString(args[0])
				if err != nil {
					return err
				}
				value, found := m.Items[key]
				if !found {
					if len(args) == 2 {
						return args[1]
					}
					return Nil
				}
				return value
			},
		}, true
	}
	return nil, false
}

func (m *Map) InvokeMethod(method string, args ...Object) Object {
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
		return Nil
	}
	return value
}

func (m *Map) Get(key string) Object {
	value, found := m.Items[key]
	if !found {
		return Nil
	}
	return value
}

func (m *Map) Delete(key string) Object {
	delete(m.Items, key)
	return Nil
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

func (m *Map) Equals(other Object) Object {
	if other.Type() != MAP {
		return False
	}
	otherMap := other.(*Map)
	if len(m.Items) != len(otherMap.Items) {
		return False
	}
	for k, v := range m.Items {
		otherValue, found := otherMap.Items[k]
		if !found {
			return False
		}
		if !v.Equals(otherValue).(*Bool).Value {
			return False
		}
	}
	return True
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
