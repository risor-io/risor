package object

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

type Set struct {
	Items map[HashKey]Object
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
				return s.Contains(args[0])
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

func (s *Set) Interface() interface{} {
	items := make([]interface{}, 0, len(s.Items))
	for _, item := range s.SortedItems() {
		items = append(items, item.Interface())
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

// Union returns a new set that is the union of the two sets.
func (s *Set) Union(other *Set) *Set {
	union := &Set{Items: map[HashKey]Object{}}
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
	intersection := &Set{Items: map[HashKey]Object{}}
	for k, v := range s.Items {
		if _, ok := other.Items[k]; ok {
			intersection.Items[k] = v
		}
	}
	return intersection
}

// Difference returns a new set that is the difference of the two sets.
func (s *Set) Difference(other *Set) *Set {
	difference := &Set{Items: map[HashKey]Object{}}
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

func (s *Set) GetItem(key Object) (Object, *Error) {
	hashable, ok := key.(Hashable)
	if !ok {
		return nil, NewError("type error: %s object is unhashable", key.Type())
	}
	if _, ok := s.Items[hashable.HashKey()]; ok {
		return True, nil
	}
	return False, nil
}

// GetSlice implements the [start:stop] operator for a container type.
func (s *Set) GetSlice(slice Slice) (Object, *Error) {
	return nil, NewError("eval error: set does not support get slice")
}

// SetItem assigns a value to the given key in the map.
func (s *Set) SetItem(key, value Object) *Error {
	return NewError("eval error: set does not support set item")
}

// DelItem deletes the item with the given key from the map.
func (s *Set) DelItem(key Object) *Error {
	hashable, ok := key.(Hashable)
	if !ok {
		return NewError("type error: %s object is unhashable", key.Type())
	}
	delete(s.Items, hashable.HashKey())
	return nil
}

// Contains returns true if the given item is found in this container.
func (s *Set) Contains(key Object) *Bool {
	hashable, ok := key.(Hashable)
	if !ok {
		return False
	}
	_, ok = s.Items[hashable.HashKey()]
	return NewBool(ok)
}

// Len returns the number of items in this container.
func (s *Set) Len() *Int {
	return NewInt(int64(len(s.Items)))
}

func NewSetWithSize(size int) *Set {
	return &Set{Items: make(map[HashKey]Object, size)}
}
