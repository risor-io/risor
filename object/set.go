package object

import (
	"bytes"
	"context"
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
	for _, item := range s.SortedItems() {
		items = append(items, item.Inspect())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("}")
	return out.String()
}

func (s *Set) GetAttr(name string) (Object, bool) {
	switch name {
	case "contains":
		return &Builtin{
			Name: "set.contains",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("set.contains", 1, len(args))
				}
				if s.Contains(args[0]) {
					return True
				}
				return False
			},
		}, true
	case "add":
		return &Builtin{
			Name: "set.add",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("set.add", 1, len(args))
				}
				if err := s.Add(args[0]); err != nil {
					return NewError(err.Error())
				}
				return s
			},
		}, true
	case "remove":
		return &Builtin{
			Name: "set.remove",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("set.remove", 1, len(args))
				}
				if err := s.Remove(args[0]); err != nil {
					return NewError(err.Error())
				}
				return s
			},
		}, true
	case "union":
		return &Builtin{
			Name: "set.union",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("set.union", 1, len(args))
				}
				other, err := AsSet(args[0])
				if err != nil {
					return err
				}
				return s.Union(other)
			},
		}, true
	case "intersection":
		return &Builtin{
			Name: "set.intersection",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("set.intersection", 1, len(args))
				}
				other, err := AsSet(args[0])
				if err != nil {
					return err
				}
				return s.Intersection(other)
			},
		}, true
	}
	return nil, false
}

func (s *Set) InvokeMethod(method string, args ...Object) Object {
	return NewError("type error: %s object has no method %s", s.Type(), method)
}

func (s *Set) ToInterface() interface{} {
	items := make([]interface{}, 0, len(s.Items))
	for _, item := range s.SortedItems() {
		items = append(items, item.ToInterface())
	}
	return items
}

func (s *Set) Size() int {
	return len(s.Items)
}

func (s *Set) SortedItems() []Object {
	items := make([]Object, 0, len(s.Items))
	for _, v := range s.Items {
		items = append(items, v)
	}
	Sort(items)
	return items
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

func (s *Set) Contains(item Object) bool {
	hashable, ok := item.(Hashable)
	if !ok {
		return false
	}
	_, ok = s.Items[hashable.HashKey()]
	return ok
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
