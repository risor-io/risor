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
	INTEGER_OBJ       = "INTEGER"
	FLOAT_OBJ         = "FLOAT"
	BOOLEAN_OBJ       = "BOOLEAN"
	NULL_OBJ          = "NULL"
	RETURN_VALUE_OBJ  = "RETURN_VALUE"
	ERROR_OBJ         = "ERROR"
	FUNCTION_OBJ      = "FUNCTION"
	STRING_OBJ        = "STRING"
	BUILTIN_OBJ       = "BUILTIN"
	ARRAY_OBJ         = "ARRAY"
	HASH_OBJ          = "HASH"
	FILE_OBJ          = "FILE"
	REGEXP_OBJ        = "REGEXP"
	SET_OBJ           = "SET"
	MODULE_OBJ        = "MODULE"
	RESULT_OBJ        = "RESULT"
	HTTP_RESPONSE_OBJ = "HTTP_RESPONSE"
	DB_CONNECTION_OBJ = "DB_CONNECTION"
	TIME_OBJ          = "TIME"
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
}

// Hashable type can be hashed.
type Hashable interface {

	// HashKey returns a hash key for the given object.
	HashKey() HashKey
}
