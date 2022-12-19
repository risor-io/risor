package object

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

// List of objects
type List struct {
	// Items holds the list of objects
	Items []Object
}

func (ls *List) Type() Type {
	return LIST
}

func (ls *List) Inspect() string {
	var out bytes.Buffer
	items := make([]string, 0)
	for _, e := range ls.Items {
		items = append(items, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("]")
	return out.String()
}

func (ls *List) GetAttr(name string) (Object, bool) {
	switch name {
	case "append":
		return &Builtin{
			Name: "list.append",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("list.append", 1, len(args))
				}
				ls.Append(args[0])
				return ls
			},
		}, true
	case "clear":
		return &Builtin{
			Name: "list.clear",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("list.clear", 0, len(args))
				}
				ls.Clear()
				return ls
			},
		}, true
	case "copy":
		return &Builtin{
			Name: "list.copy",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("list.copy", 0, len(args))
				}
				return ls.Copy()
			},
		}, true
	case "count":
		return &Builtin{
			Name: "list.count",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("list.count", 1, len(args))
				}
				return NewInt(ls.Count(args[0]))
			},
		}, true
	case "extend":
		return &Builtin{
			Name: "list.extend",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("list.extend", 1, len(args))
				}
				other, err := AsList(args[0])
				if err != nil {
					return err
				}
				ls.Extend(other)
				return ls
			},
		}, true
	case "index":
		return &Builtin{
			Name: "list.index",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("list.index", 1, len(args))
				}
				return NewInt(ls.Index(args[0]))
			},
		}, true
	case "insert":
		return &Builtin{
			Name: "list.insert",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 2 {
					return NewArgsError("list.insert", 2, len(args))
				}
				index, err := AsInteger(args[0])
				if err != nil {
					return err
				}
				ls.Insert(index, args[1])
				return ls
			},
		}, true
	case "pop":
		return &Builtin{
			Name: "list.pop",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("list.pop", 1, len(args))
				}
				index, err := AsInteger(args[0])
				if err != nil {
					return err
				}
				return ls.Pop(index)
			},
		}, true
	case "remove":
		return &Builtin{
			Name: "list.remove",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("list.remove", 1, len(args))
				}
				ls.Remove(args[0])
				return ls
			},
		}, true
	case "reverse":
		return &Builtin{
			Name: "list.reverse",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("list.reverse", 0, len(args))
				}
				ls.Reverse()
				return ls
			},
		}, true
	case "sort":
		return &Builtin{
			Name: "list.sort",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("list.sort", 0, len(args))
				}
				if err := Sort(ls.Items); err != nil {
					return err
				}
				return ls
			},
		}, true
	case "map":
		return &Builtin{
			Name: "list.map",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("list.map", 1, len(args))
				}
				return ls.Map(ctx, args[0])
			},
		}, true
	case "filter":
		return &Builtin{
			Name: "list.filter",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("list.filter", 1, len(args))
				}
				return ls.Filter(ctx, args[0])
			},
		}, true
	case "each":
		return &Builtin{
			Name: "list.each",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("list.each", 1, len(args))
				}
				return ls.Each(ctx, args[0])
			},
		}, true
	}
	return nil, false
}

func (ls *List) Map(ctx context.Context, fn Object) Object {
	callFunc, found := GetCallFunc(ctx)
	if !found {
		return NewError("eval error: list.map() context did not contain a call function")
	}
	var numParameters int
	switch obj := fn.(type) {
	case *Builtin:
		numParameters = 1
	case *Function:
		numParameters = len(obj.Parameters)
	default:
		return NewError("type error: list.map() expected a function (%s given)", obj.Type())
	}
	if numParameters < 1 || numParameters > 2 {
		return NewError("type error: list.map() received an incompatible function")
	}
	var index Int
	mapArgs := make([]Object, 2)
	result := make([]Object, 0, len(ls.Items))
	for i, value := range ls.Items {
		index.Value = int64(i)
		mapArgs[0] = &index
		mapArgs[1] = value
		var outputValue Object
		if numParameters == 1 {
			outputValue = callFunc(ctx, nil, fn, mapArgs[1:])
		} else {
			outputValue = callFunc(ctx, nil, fn, mapArgs)
		}
		if IsError(outputValue) {
			return outputValue
		}
		result = append(result, outputValue)
	}
	return NewList(result)
}

func (ls *List) Filter(ctx context.Context, fn Object) Object {
	callFunc, found := GetCallFunc(ctx)
	if !found {
		return NewError("eval error: list.filter() context did not contain a call function")
	}
	switch obj := fn.(type) {
	case *Function, *Builtin:
		// Nothing do do here
	default:
		return NewError("type error: list.filter() expected a function (%s given)", obj.Type())
	}
	filterArgs := make([]Object, 1)
	var result []Object
	for _, value := range ls.Items {
		filterArgs[0] = value
		decision := callFunc(ctx, nil, fn, filterArgs)
		if IsError(decision) {
			return decision
		}
		if IsTruthy(decision) {
			result = append(result, value)
		}
	}
	return NewList(result)
}

func (ls *List) Each(ctx context.Context, fn Object) Object {
	callFunc, found := GetCallFunc(ctx)
	if !found {
		return NewError("eval error: list.each() context did not contain a call function")
	}
	switch obj := fn.(type) {
	case *Function, *Builtin:
		// Nothing do do here
	default:
		return NewError("type error: list.each() expected a function (%s given)", obj.Type())
	}
	eachArgs := make([]Object, 1)
	for _, value := range ls.Items {
		eachArgs[0] = value
		result := callFunc(ctx, nil, fn, eachArgs)
		if IsError(result) {
			return result
		}
	}
	return Nil
}

// Append adds an item at the end of the list.
func (ls *List) Append(obj Object) {
	ls.Items = append(ls.Items, obj)
}

// Clear removes all the items from the list.
func (ls *List) Clear() {
	ls.Items = []Object{}
}

// Copy returns a shallow copy of the list.
func (ls *List) Copy() *List {
	result := &List{Items: make([]Object, len(ls.Items))}
	copy(result.Items, ls.Items)
	return result
}

// Count returns the number of items with the specified value.
func (ls *List) Count(obj Object) int64 {
	count := int64(0)
	for _, item := range ls.Items {
		if Equals(obj, item) {
			count++
		}
	}
	return count
}

// Extend adds the items of a list to the end of the current list.
func (ls *List) Extend(other *List) {
	ls.Items = append(ls.Items, other.Items...)
}

// Index returns the index of the first item with the specified value.
func (ls *List) Index(obj Object) int64 {
	for i, item := range ls.Items {
		if Equals(obj, item) {
			return int64(i)
		}
	}
	return int64(-1)
}

// Insert adds an item at the specified position.
func (ls *List) Insert(index int64, obj Object) {
	// Negative index is relative to the end of the list
	if index < 0 {
		index = int64(len(ls.Items)) + index
		if index < 0 {
			index = 0
		}
	}
	if index == 0 {
		ls.Items = append([]Object{obj}, ls.Items...)
		return
	}
	if index >= int64(len(ls.Items)) {
		ls.Items = append(ls.Items, obj)
		return
	}
	ls.Items = append(ls.Items, nil)
	copy(ls.Items[index+1:], ls.Items[index:])
	ls.Items[index] = obj
}

// Pop removes the item at the specified position.
func (ls *List) Pop(index int64) Object {
	idx, err := ResolveIndex(index, int64(len(ls.Items)))
	if err != nil {
		return NewError(err.Error())
	}
	result := ls.Items[idx]
	ls.Items = append(ls.Items[:idx], ls.Items[idx+1:]...)
	return result
}

// Remove removes the first item with the specified value.
func (ls *List) Remove(obj Object) {
	index := ls.Index(obj)
	if index == -1 {
		return
	}
	ls.Items = append(ls.Items[:index], ls.Items[index+1:]...)
}

// Reverse reverses the order of the list.
func (ls *List) Reverse() {
	for i, j := 0, len(ls.Items)-1; i < j; i, j = i+1, j-1 {
		ls.Items[i], ls.Items[j] = ls.Items[j], ls.Items[i]
	}
}

func (ls *List) Interface() interface{} {
	items := make([]interface{}, 0, len(ls.Items))
	for _, item := range ls.Items {
		items = append(items, item.Interface())
	}
	return items
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

func (ls *List) Equals(other Object) Object {
	if other.Type() != LIST {
		return False
	}
	otherList := other.(*List)
	if len(ls.Items) != len(otherList.Items) {
		return False
	}
	for i, v := range ls.Items {
		otherV := otherList.Items[i]
		if !Equals(v, otherV) {
			return False
		}
	}
	return True
}

func (ls *List) Reversed() *List {
	result := &List{Items: make([]Object, 0, len(ls.Items))}
	size := len(ls.Items)
	for i := 0; i < size; i++ {
		result.Items = append(result.Items, ls.Items[size-1-i])
	}
	return result
}

func (ls *List) Keys() Object {
	items := make([]Object, 0, len(ls.Items))
	for i := 0; i < len(ls.Items); i++ {
		items = append(items, NewInt(int64(i)))
	}
	return NewList(items)
}

func (ls *List) GetItem(key Object) (Object, *Error) {
	indexObj, ok := key.(*Int)
	if !ok {
		return nil, NewError("type error: list index must be an int (got %s)", key.Type())
	}
	idx, err := ResolveIndex(indexObj.Value, int64(len(ls.Items)))
	if err != nil {
		return nil, NewError(err.Error())
	}
	return ls.Items[idx], nil
}

// GetSlice implements the [start:stop] operator for a container type.
func (ls *List) GetSlice(s Slice) (Object, *Error) {
	start, stop, err := ResolveIntSlice(s, int64(len(ls.Items)))
	if err != nil {
		return nil, NewError(err.Error())
	}
	items := ls.Items[start:stop]
	itemsCopy := make([]Object, len(items))
	copy(itemsCopy, items)
	return NewList(itemsCopy), nil
}

// SetItem implements the [key] = value operator for a container type.
func (ls *List) SetItem(key, value Object) *Error {
	indexObj, ok := key.(*Int)
	if !ok {
		return NewError("type error: list index must be an int (got %s)", key.Type())
	}
	idx, err := ResolveIndex(indexObj.Value, int64(len(ls.Items)))
	if err != nil {
		return NewError(err.Error())
	}
	ls.Items[idx] = value
	return nil
}

// DelItem implements the del [key] operator for a container type.
func (ls *List) DelItem(key Object) *Error {
	return NewError("list does not support del operator")
}

// Contains returns true if the given item is found in this container.
func (ls *List) Contains(item Object) *Bool {
	for _, v := range ls.Items {
		if Equals(v, item) {
			return True
		}
	}
	return False
}

// Len returns the number of items in this container.
func (ls *List) Len() *Int {
	return NewInt(int64(len(ls.Items)))
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

// ResolveIndex checks that the index is inbounds and transforms a negative
// index into the corresponding positive index. If the index is out of bounds,
// an error is returned.
func ResolveIndex(idx int64, size int64) (int64, error) {
	max := size - 1
	if idx > max {
		return 0, fmt.Errorf("index error: index out of range: %d", idx)
	}
	if idx >= 0 {
		return idx, nil
	}
	// Handle negative indices, where -1 is the last item in the array
	reversed := idx + size
	if reversed < 0 || reversed > max {
		return 0, fmt.Errorf("index error: index out of range: %d", idx)
	}
	return reversed, nil
}

// ResolveIntSlice checks that the slice start and stop indices are inbounds and
// transforms negative indices into the corresponding positive indices. If the
// slice is out of bounds, an error is returned.
func ResolveIntSlice(slice Slice, size int64) (start int64, stop int64, err error) {
	var startIndex, stopIndex int64
	if slice.Start != nil {
		startObj, ok := slice.Start.(*Int)
		if !ok {
			err = fmt.Errorf("type error: slice start index must be an int (got %s)", slice.Start.Type())
			return
		}
		startIndex = startObj.Value
	}
	if slice.Stop != nil {
		stopObj, ok := slice.Stop.(*Int)
		if !ok {
			err = fmt.Errorf("type error: slice stop index must be an int (got %s)", slice.Stop.Type())
			return
		}
		stopIndex = stopObj.Value
	} else {
		stopIndex = size
	}
	start, err = ResolveIndex(startIndex, size+1)
	if err != nil {
		return
	}
	stop, err = ResolveIndex(stopIndex, size+1)
	if err != nil {
		return
	}
	if start > stop {
		err = fmt.Errorf("slice error: start index is greater than stop index")
		return
	}
	return start, stop, nil
}
