package object

import (
	"bytes"
	"fmt"
	"strings"
)

// List of objects
type List struct {
	// Items holds the list of objects
	Items []Object
}

// Type returns the type of this object.
func (ls *List) Type() Type {
	return LIST
}

// Inspect returns a string-representation of the given object.
func (ls *List) Inspect() string {
	var out bytes.Buffer
	elements := make([]string, 0)
	for _, e := range ls.Items {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

func (ls *List) InvokeMethod(method string, args ...Object) Object {
	if method == "len" {
		return &Int{Value: int64(len(ls.Items))}
	}
	return NewError("type error: %s object has no method %s", ls.Type(), method)
}

func (ls *List) ToInterface() interface{} {
	return "<ARRAY>"
}

func (ls *List) String() string {
	items := make([]string, 0, len(ls.Items))
	for _, item := range ls.Items {
		items = append(items, fmt.Sprintf("%s", item))
	}
	return fmt.Sprintf("List([%s])", strings.Join(items, ", "))
}

func (ls *List) Compare(other Object) (int, error) {
	typeComp := CompareTypes(ls, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherArr := other.(*List)
	if len(ls.Items) > len(otherArr.Items) {
		return 1, nil
	} else if len(ls.Items) < len(otherArr.Items) {
		return -1, nil
	}
	for i := 0; i < len(ls.Items); i++ {
		comparable, ok := ls.Items[i].(Comparable)
		if !ok {
			return 0, fmt.Errorf("type error: %s object is not comparable",
				ls.Items[i].Type())
		}
		comp, err := comparable.Compare(otherArr.Items[i])
		if err != nil {
			return 0, err
		}
		if comp != 0 {
			return comp, nil
		}
	}
	return 0, nil
}

func (ls *List) Reversed() *List {
	result := &List{Items: make([]Object, 0, len(ls.Items))}
	size := len(ls.Items)
	for i := 0; i < size; i++ {
		result.Items = append(result.Items, ls.Items[size-1-i])
	}
	return result
}

func NewList(items []Object) *List {
	return &List{Items: items}
}

func NewStringList(s []string) *List {
	array := &List{Items: make([]Object, 0, len(s))}
	for _, item := range s {
		array.Items = append(array.Items, &String{Value: item})
	}
	return array
}
