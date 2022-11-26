package object

import (
	"bytes"
	"fmt"
	"strings"
)

// HashKey is the structure used for hash-keys
type HashKey struct {
	// Type holds the type of the object.
	Type Type

	// Value holds the actual hash-key value.
	Value uint64
}

type Hash struct {
	Map map[string]Object
}

func (h *Hash) Type() Type {
	return HASH_OBJ
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := make([]string, 0)
	for k, v := range h.Map {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

func (h *Hash) InvokeMethod(method string, args ...Object) Object {
	switch method {
	case "keys":
		return h.Keys()
	case "values":
		return h.Values()
	}
	return nil
}

func (h *Hash) Keys() *Array {
	elements := make([]Object, 0, len(h.Map))
	for k := range h.Map {
		elements = append(elements, &String{Value: k})
	}
	return &Array{Elements: elements}
}

func (h *Hash) Values() *Array {
	elements := make([]Object, 0, len(h.Map))
	for _, v := range h.Map {
		elements = append(elements, v)
	}
	return &Array{Elements: elements}
}

func (h *Hash) GetWithObject(key *String) Object {
	value, found := h.Map[key.Value]
	if !found {
		return NULL
	}
	return value
}

func (h *Hash) Get(key string) Object {
	value, found := h.Map[key]
	if !found {
		return NULL
	}
	return value
}

func (h *Hash) Delete(key string) Object {
	delete(h.Map, key)
	return NULL
}

func (h *Hash) ToInterface() interface{} {
	result := make(map[string]any, len(h.Map))
	for k, v := range h.Map {
		result[k] = v.ToInterface()
	}
	return result
}

func NewHash(m map[string]interface{}) *Hash {
	result := &Hash{Map: make(map[string]Object, len(m))}
	for k, v := range m {
		value := FromGoType(v)
		if value == nil {
			panic(fmt.Sprintf("type error: cannot convert %v to a tamarin object", v))
		}
		result.Map[k] = value
	}
	return result
}
