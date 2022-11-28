package object

import (
	"bytes"
	"fmt"
	"strings"
)

// Array wraps Object array and implements Object interface.
type Array struct {
	// Elements holds the individual members of the array we're wrapping.
	Elements []Object
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
func (ao *Array) InvokeMethod(method string, args ...Object) Object {
	if method == "len" {
		return &Integer{Value: int64(len(ao.Elements))}
	}
	return NewError("type error: %s object has no method %s", ao.Type(), method)
}

func (ao *Array) ToInterface() interface{} {
	return "<ARRAY>"
}

func (ao *Array) String() string {
	items := make([]string, 0, len(ao.Elements))
	for _, item := range ao.Elements {
		items = append(items, fmt.Sprintf("%s", item))
	}
	return fmt.Sprintf("Array([%s])", strings.Join(items, ", "))
}

func (ao *Array) Compare(other Object) (int, error) {
	typeComp := CompareTypes(ao, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherArr := other.(*Array)
	if len(ao.Elements) > len(otherArr.Elements) {
		return 1, nil
	} else if len(ao.Elements) < len(otherArr.Elements) {
		return -1, nil
	}
	for i := 0; i < len(ao.Elements); i++ {
		comparable, ok := ao.Elements[i].(Comparable)
		if !ok {
			return 0, fmt.Errorf("type error: %s object is not comparable",
				ao.Elements[i].Type())
		}
		comp, err := comparable.Compare(otherArr.Elements[i])
		if err != nil {
			return 0, err
		}
		if comp != 0 {
			return comp, nil
		}
	}
	return 0, nil
}

func (ao *Array) Reversed() *Array {
	result := &Array{Elements: make([]Object, 0, len(ao.Elements))}
	size := len(ao.Elements)
	for i := 0; i < size; i++ {
		result.Elements = append(result.Elements, ao.Elements[size-1-i])
	}
	return result
}

func NewStringArray(s []string) *Array {
	array := &Array{Elements: make([]Object, 0, len(s))}
	for _, item := range s {
		array.Elements = append(array.Elements, &String{Value: item})
	}
	return array
}
