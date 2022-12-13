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
	INT           = "int"
	FLOAT         = "float"
	BOOL          = "bool"
	NIL           = "nil"
	ERROR         = "error"
	FUNCTION      = "function"
	STRING        = "string"
	BUILTIN       = "builtin"
	LIST          = "list"
	MAP           = "map"
	FILE          = "file"
	REGEXP        = "regexp"
	SET           = "set"
	MODULE        = "module"
	RESULT        = "result"
	HTTP_RESPONSE = "http_response"
	DB_CONNECTION = "db_connection"
	TIME          = "time"
	PROXY         = "proxy"
	RETURN_VALUE  = "return_value"
	BREAK_VALUE   = "break_value"
)

var (
	Nil   = &NilType{}
	True  = &Bool{Value: true}
	False = &Bool{Value: false}
)

// Object is the interface that all object types in Tamarin must implement.
type Object interface {

	// Type of this object.
	Type() Type

	// Inspect returns a string-representation of the given object.
	Inspect() string

	// InvokeMethod invokes a method against the object.
	// (Built-in methods only.)
	InvokeMethod(method string, args ...Object) Object

	// ToInterface converts the given object to a "native" golang value,
	// which is required to ensure that we can use the object in our
	// `sprintf` or `printf` primitives.
	ToInterface() interface{}

	// Returns True if the given object is equal to this object.
	Equals(other Object) Object
}

// Hashable types can be hashed and consequently used in a set.
type Hashable interface {

	// HashKey returns a hash key for the given object.
	HashKey() Key
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
