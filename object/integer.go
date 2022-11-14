package object

import (
	"fmt"
)

// Integer wraps int64 and implements Object and Hashable interfaces.
type Integer struct {
	// Value holds the integer value this object wraps
	Value int64
}

// Inspect returns a string-representation of the given object.
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Type returns the type of this object.
func (i *Integer) Type() Type {
	return INTEGER_OBJ
}

// HashKey returns a hash key for the given object.
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// InvokeMethod invokes a method against the object.
// (Built-in methods only.)
func (i *Integer) InvokeMethod(method string, args ...Object) Object {
	if method == "chr" {
		return &String{Value: string(rune(i.Value))}
	}
	return nil
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (i *Integer) ToInterface() interface{} {
	return i.Value
}

func (i *Integer) String() string {
	return fmt.Sprintf("Integer(%v)", i.Value)
}

func NewInteger(value int64) *Integer {
	return &Integer{Value: value}
}
