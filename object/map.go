package object

import (
	"bytes"
	"context"
	"fmt"
	"sort"
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
	for _, k := range m.SortedKeys() {
		v := m.Items[k]
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
				if len(args) != 0 {
					return NewArgsError("map.keys", 0, len(args))
				}
				return m.Keys()
			},
		}, true
	case "values":
		return &Builtin{
			Name: "map.values",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("map.values", 0, len(args))
				}
				return m.Values()
			},
		}, true
	case "get":
		return &Builtin{
			Name: "map.get",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) < 1 || len(args) > 2 {
					return NewArgsRangeError("map.get", 1, 2, len(args))
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
	case "clear":
		return &Builtin{
			Name: "map.clear",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("map.clear", 0, len(args))
				}
				m.Clear()
				return m
			},
		}, true
	case "copy":
		return &Builtin{
			Name: "map.copy",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("map.copy", 0, len(args))
				}
				return m.Copy()
			},
		}, true
	case "contains":
		return &Builtin{
			Name: "map.contains",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("map.contains", 1, len(args))
				}
				key, err := AsString(args[0])
				if err != nil {
					return err
				}
				return m.Contains(key)
			},
		}, true
	case "items":
		return &Builtin{
			Name: "map.items",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("map.items", 0, len(args))
				}
				return m.ListItems()
			},
		}, true
	case "pop":
		return &Builtin{
			Name: "map.pop",
			Fn: func(ctx context.Context, args ...Object) Object {
				nArgs := len(args)
				if nArgs < 1 || nArgs > 2 {
					return NewArgsRangeError("map.pop", 1, 2, len(args))
				}
				key, err := AsString(args[0])
				if err != nil {
					return err
				}
				var def Object
				if nArgs == 2 {
					def = args[1]
				}
				return m.Pop(key, def)
			},
		}, true
	case "setdefault":
		return &Builtin{
			Name: "map.setdefault",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 2 {
					return NewArgsError("map.setdefault", 2, len(args))
				}
				key, err := AsString(args[0])
				if err != nil {
					return err
				}
				return m.SetDefault(key, args[1])
			},
		}, true
	case "update":
		return &Builtin{
			Name: "map.update",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("map.update", 1, len(args))
				}
				other, err := AsMap(args[0])
				if err != nil {
					return err
				}
				m.Update(other)
				return m
			},
		}, true
	}
	return nil, false
}

func (m *Map) ListItems() *List {
	items := make([]Object, 0, len(m.Items))
	for _, k := range m.SortedKeys() {
		items = append(items, &List{
			Items: []Object{&String{Value: k}, m.Items[k]},
		})
	}
	return &List{Items: items}
}

func (m *Map) Clear() {
	m.Items = map[string]Object{}
}

func (m *Map) Copy() *Map {
	items := make(map[string]Object, len(m.Items))
	for k, v := range m.Items {
		items[k] = v
	}
	return &Map{Items: items}
}

func (m *Map) Contains(key string) Object {
	_, found := m.Items[key]
	return NewBool(found)
}

func (m *Map) Pop(key string, def Object) Object {
	value, found := m.Items[key]
	if found {
		delete(m.Items, key)
		return value
	}
	if def != nil {
		return def
	}
	return Nil
}

func (m *Map) SetDefault(key string, value Object) Object {
	if _, found := m.Items[key]; !found {
		m.Items[key] = value
	}
	return m.Items[key]
}

func (m *Map) Update(other *Map) {
	for k, v := range other.Items {
		m.Items[k] = v
	}
}

func (m *Map) InvokeMethod(method string, args ...Object) Object {
	return NewError("type error: %s object has no method %s", m.Type(), method)
}

func (m *Map) SortedKeys() []string {
	keys := make([]string, 0, len(m.Items))
	for k := range m.Items {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (m *Map) Keys() *List {
	items := make([]Object, 0, len(m.Items))
	for _, k := range m.SortedKeys() {
		items = append(items, &String{Value: k})
	}
	return &List{Items: items}
}

func (m *Map) Values() *List {
	items := make([]Object, 0, len(m.Items))
	for _, k := range m.SortedKeys() {
		items = append(items, m.Items[k])
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
