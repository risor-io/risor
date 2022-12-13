package object

import (
	"bytes"
	"fmt"
	"strings"
)

type Set struct {
	Items map[Key]Object
}

func (s *Set) Type() Type {
	return SET
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
	switch method {
	case "values":
		return s.List()
	case "contains":
		if len(args) == 0 {
			return NewError("type error: set.contains() expects at least one argument")
		}
		if s.Contains(args...) {
			return True
		}
		return False
	case "add":
		if len(args) == 0 {
			return NewError("type error: set.add() expects at least one argument")
		}
		if err := s.Add(args...); err != nil {
			return NewError(err.Error())
		}
		return Nil
	case "remove":
		if len(args) == 0 {
			return NewError("type error: set.remove() expects at least one argument")
		}
		if err := s.Remove(args...); err != nil {
			return NewError(err.Error())
		}
		return Nil
	case "union":
		if len(args) != 1 {
			return NewError("type error: set.union() expects one argument")
		}
		other, err := AsSet(args[0])
		if err != nil {
			return err
		}
		return s.Union(other)
	case "intersection":
		if len(args) != 1 {
			return NewError("type error: set.intersection() expects one argument")
		}
		other, err := AsSet(args[0])
		if err != nil {
			return err
		}
		return s.Intersection(other)
	case "difference":
		if len(args) != 1 {
			return NewError("type error: set.difference() expects one argument")
		}
		other, err := AsSet(args[0])
		if err != nil {
			return err
		}
		return s.Difference(other)
	default:
		return NewError("type error: %s object has no method %s", s.Type(), method)
	}
}

func (s *Set) ToInterface() interface{} {
	items := make([]interface{}, 0, len(s.Items))
	for _, v := range s.Items {
		items = append(items, v.ToInterface())
	}
	return items
}

func (s *Set) Size() int {
	return len(s.Items)
}

func (s *Set) Add(items ...Object) error {
	for _, item := range items {
		hashable, ok := item.(Hashable)
		if !ok {
			return fmt.Errorf("type error: %s object is unhashable", item.Type())
		}
		s.Items[hashable.HashKey()] = item
	}
	return nil
}

func (s *Set) Remove(items ...Object) error {
	for _, item := range items {
		hashable, ok := item.(Hashable)
		if !ok {
			return fmt.Errorf("type error: %s object is unhashable", item.Type())
		}
		delete(s.Items, hashable.HashKey())
	}
	return nil
}

func (s *Set) Contains(items ...Object) bool {
	for _, item := range items {
		hashable, ok := item.(Hashable)
		if !ok {
			return false
		}
		if _, ok = s.Items[hashable.HashKey()]; !ok {
			return false
		}
	}
	return true
}

// Union returns a new set that is the union of the two sets.
func (s *Set) Union(other *Set) *Set {
	union := &Set{Items: map[Key]Object{}}
	for k, v := range s.Items {
		union.Items[k] = v
	}
	for k, v := range other.Items {
		union.Items[k] = v
	}
	return union
}

// Intersection returns a new set that is the intersection of the two sets.
func (s *Set) Intersection(other *Set) *Set {
	intersection := &Set{Items: map[Key]Object{}}
	for k, v := range s.Items {
		if _, ok := other.Items[k]; ok {
			intersection.Items[k] = v
		}
	}
	return intersection
}

// Difference returns a new set that is the difference of the two sets.
func (s *Set) Difference(other *Set) *Set {
	difference := &Set{Items: map[Key]Object{}}
	for k, v := range s.Items {
		if _, ok := other.Items[k]; !ok {
			difference.Items[k] = v
		}
	}
	return difference
}

func (s *Set) List() *List {
	l := &List{Items: make([]Object, 0, len(s.Items))}
	for _, item := range s.Items {
		l.Items = append(l.Items, item)
	}
	return l
}

func (s *Set) Equals(other Object) Object {
	if other.Type() != SET {
		return False
	}
	otherSet := other.(*Set)
	if len(s.Items) != len(otherSet.Items) {
		return False
	}
	for k, v := range s.Items {
		if otherV, ok := otherSet.Items[k]; !ok || !v.Equals(otherV).(*Bool).Value {
			return False
		}
	}
	return True
}

func NewSetWithSize(size int) *Set {
	return &Set{Items: make(map[Key]Object, size)}
}
