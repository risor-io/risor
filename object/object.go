// Package object defines the standard set of object types in Tamarin.
//
// For external users of Tamarin, often an object.Object interface
// will be type asserted to a specific object type, such as *object.Float.
//
// For example:
//
//	switch obj := obj.(type) {
//	case *object.String:
//		// do something with obj.Value()
//	case *object.Float:
//		// do something with obj.Value()
//	}
//
// The Type() method of each object may also be used to get a string
// name of the object type, such as "string" or "float".
package object

import "github.com/cloudcmds/tamarin/v2/op"

// Type of an object as a string.
type Type string

// Type constants
const (
	BOOL          Type = "bool"
	BUFFER        Type = "buffer"
	BUILTIN       Type = "builtin"
	BSLICE        Type = "bslice"
	BSLICE_ITER   Type = "bslice_iter"
	CELL          Type = "cell"
	CODE          Type = "code"
	ERROR         Type = "error"
	FILE          Type = "file"
	FLOAT         Type = "float"
	FUNCTION      Type = "function"
	GO_TYPE       Type = "go_type"
	GO_FIELD      Type = "go_field"
	GO_METHOD     Type = "go_method"
	HTTP_RESPONSE Type = "http_response"
	INT           Type = "int"
	ITER_ENTRY    Type = "iter_entry"
	LIST          Type = "list"
	LIST_ITER     Type = "list_iter"
	MAP           Type = "map"
	MAP_ITER      Type = "map_iter"
	MODULE        Type = "module"
	NIL           Type = "nil"
	PARTIAL       Type = "partial"
	PROXY         Type = "proxy"
	REGEXP        Type = "regexp"
	RESULT        Type = "result"
	SET           Type = "set"
	SET_ITER      Type = "set_iter"
	STRING        Type = "string"
	STRING_ITER   Type = "string_iter"
	TIME          Type = "time"
)

var (
	Nil   = &NilType{}
	True  = &Bool{value: true}
	False = &Bool{value: false}
)

// Object is the interface that all object types in Tamarin must implement.
type Object interface {

	// Type of the object.
	Type() Type

	// Inspect returns a string representation of the given object.
	Inspect() string

	// Interface converts the given object to a native Go value.
	Interface() interface{}

	// Returns True if the given object is equal to this object.
	Equals(other Object) Object

	// GetAttr returns the attribute with the given name from this object.
	GetAttr(name string) (Object, bool)

	// IsTruthy returns true if the object is considered "truthy".
	IsTruthy() bool

	// RunOperation runs an operation on this object with the given
	// right-hand side object.
	RunOperation(opType op.BinaryOpType, right Object) Object

	// Cost returns the incremental processing cost of this object.
	Cost() int
}

// Slice is used to specify a range or slice of items in a container.
type Slice struct {
	Start Object
	Stop  Object
}

// IteratorEntry is a single item returned by an iterator.
type IteratorEntry interface {
	Object
	Key() Object
	Value() Object
	Primary() Object
}

// Iterator is an interface used to iterate over a container.
type Iterator interface {
	Object

	// Next advances the iterator and then returns the current object and a
	// bool indicating whether the returned item is valid. Once Next() has been
	// called, the Entry() method can be used to get an IteratorEntry.
	Next() (Object, bool)

	// Entry returns the current entry in the iterator and a bool indicating
	// whether the returned item is valid.
	Entry() (IteratorEntry, bool)
}

type Container interface {

	// GetItem implements the [key] operator for a container type.
	GetItem(key Object) (Object, *Error)

	// GetSlice implements the [start:stop] operator for a container type.
	GetSlice(s Slice) (Object, *Error)

	// SetItem implements the [key] = value operator for a container type.
	SetItem(key, value Object) *Error

	// DelItem implements the del [key] operator for a container type.
	DelItem(key Object) *Error

	// Contains returns true if the given item is found in this container.
	Contains(item Object) *Bool

	// Len returns the number of items in this container.
	Len() *Int

	// Iter returns an iterator for this container.
	Iter() Iterator
}

// Hashable types can be hashed and consequently used in a set.
type Hashable interface {

	// Hash returns a hash key for the given object.
	HashKey() HashKey
}

// Comparable is an interface used to compare two objects.
//
//	-1 if this < other
//	 0 if this == other
//	 1 if this > other
type Comparable interface {
	Compare(other Object) (int, error)
}

func CompareTypes(a, b Object) int {
	aType := a.Type()
	bType := b.Type()
	if aType != bType {
		if aType < bType {
			return -1
		}
		return 1
	}
	return 0
}

// HashKey is used to identify unique values in a set.
type HashKey struct {
	// Type of the object being referenced.
	Type Type
	// FltValue is used as the key for floats.
	FltValue float64
	// IntValue is used as the key for integers.
	IntValue int64
	// StrValue is used as the key for strings.
	StrValue string
}
