package object

import (
	"fmt"
)

// Boolean wraps bool and implements Object and Hashable interface.
type Boolean struct {
	// Value holds the boolean value we wrap.
	Value bool
}

// Type returns the type of this object.
func (b *Boolean) Type() Type {
	return BOOLEAN_OBJ
}

// Inspect returns a string-representation of the given object.
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// HashKey returns a hash key for the given object.
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

// InvokeMethod invokes a method against the object.
// (Built-in methods only.)
func (b *Boolean) InvokeMethod(method string, args ...Object) Object {
	return nil
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (b *Boolean) ToInterface() interface{} {
	return b.Value
}

func (b *Boolean) String() string {
	return fmt.Sprintf("Boolean(%v)", b.Value)
}

func NewBoolean(value bool) *Boolean {
	if value {
		return TRUE
	}
	return FALSE
}
