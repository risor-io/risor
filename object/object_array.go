package object

import (
	"bytes"
	"sort"
	"strings"
)

// Array wraps Object array and implements Object interface.
type Array struct {
	// Elements holds the individual members of the array we're wrapping.
	Elements []Object

	// offset holds our iteration-offset.
	offset int
}

// Type returns the type of this object.
func (ao *Array) Type() Type {
	return ARRAY_OBJ
}

// Inspect returns a string-representation of the given object.
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := make([]string, 0)
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// InvokeMethod invokes a method against the object.
// (Built-in methods only.)
func (ao *Array) InvokeMethod(method string, env Environment, args ...Object) Object {
	if method == "len" {
		return &Integer{Value: int64(len(ao.Elements))}
	}
	if method == "methods" {
		static := []string{"len", "methods"}
		dynamic := env.Names("array.")

		var names []string
		names = append(names, static...)
		for _, e := range dynamic {
			bits := strings.Split(e, ".")
			names = append(names, bits[1])
		}
		sort.Strings(names)

		result := make([]Object, len(names))
		for i, txt := range names {
			result[i] = &String{Value: txt}
		}
		return &Array{Elements: result}
	}
	return nil
}

// Reset implements the Iterable interface, and allows the contents
// of the array to be reset to allow re-iteration.
func (ao *Array) Reset() {
	ao.offset = 0
}

// Next implements the Iterable interface, and allows the contents
// of our array to be iterated over.
func (ao *Array) Next() (Object, Object, bool) {
	if ao.offset < len(ao.Elements) {
		ao.offset++

		element := ao.Elements[ao.offset-1]
		return element, &Integer{Value: int64(ao.offset - 1)}, true
	}

	return nil, &Integer{Value: 0}, false
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (ao *Array) ToInterface() interface{} {
	return "<ARRAY>"
}
