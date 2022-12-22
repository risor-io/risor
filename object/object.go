// Package object defines all available object types in Tamarin.
//
// For external users of Tamarin, most often an object.Object interface
// will be type asserted to a specific object type, such as *object.Float.
//
// For example:
//
//	switch obj := obj.(type) {
//	case *object.String:
//		// do something with obj.Value
//	case *object.Float:
//		// do something with obj.Value
//	}
//
// The Type() method of each object may also be used to get a string
// name of the object type, such as "STRING" or "FLOAT".
package object

// Type defines the type of an object.
type Type string

// Type constants
const (
	INT           Type = "int"
	FLOAT         Type = "float"
	BOOL          Type = "bool"
	NIL           Type = "nil"
	ERROR         Type = "error"
	FUNCTION      Type = "function"
	STRING        Type = "string"
	BUILTIN       Type = "builtin"
	LIST          Type = "list"
	MAP           Type = "map"
	FILE          Type = "file"
	REGEXP        Type = "regexp"
	SET           Type = "set"
	MODULE        Type = "module"
	RESULT        Type = "result"
	HTTP_RESPONSE Type = "http_response"
	DB_CONNECTION Type = "db_connection"
	TIME          Type = "time"
	PROXY         Type = "proxy"
	RETURN_VALUE  Type = "return_value"
	BREAK_VALUE   Type = "break_value"
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
}

// Slice is used to specify a range or slice of items in a container.
type Slice struct {
	Start Object
	Stop  Object
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
