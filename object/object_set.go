package object

import (
	"bytes"
	"strings"
)

type Set struct {
	Items  map[HashKey]Object
	offset int
}

func (s *Set) Type() Type {
	return SET_OBJ
}

func (s *Set) Inspect() string {
	var out bytes.Buffer
	items := make([]string, 0, len(s.Items))
	for _, item := range s.Items {
		items = append(items, item.Inspect())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("}")
	return out.String()
}

func (s *Set) InvokeMethod(method string, args ...Object) Object {
	return nil
}

func (s *Set) Reset() {
	s.offset = 0
}

func (s *Set) Next() (Object, bool) {
	if s.offset < len(s.Items) {
		idx := 0
		for _, item := range s.Items {
			if s.offset == idx {
				s.offset++
				return item, true
			}
			idx++
		}
	}
	return nil, false
}

func (s *Set) ToInterface() interface{} {
	return "<SET>"
}

func (s *Set) Array() []Object {
	array := make([]Object, 0, len(s.Items))
	for _, item := range s.Items {
		array = append(array, item)
	}
	return array
}
