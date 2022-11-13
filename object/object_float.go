package object

import (
	"fmt"
	"hash/fnv"
	"strconv"
)

// Float wraps float64 and implements Object and Hashable interfaces.
type Float struct {
	// Value holds the float-value this object wraps.
	Value float64
}

// Inspect returns a string-representation of the given object.
func (f *Float) Inspect() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

// Type returns the type of this object.
func (f *Float) Type() Type {
	return FLOAT_OBJ
}

// HashKey returns a hash key for the given object.
func (f *Float) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(f.Inspect()))
	return HashKey{Type: f.Type(), Value: h.Sum64()}
}

// InvokeMethod invokes a method against the object.
// (Built-in methods only.)
func (f *Float) InvokeMethod(method string, args ...Object) Object {
	return nil
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (f *Float) ToInterface() interface{} {
	return f.Value
}

func (f *Float) String() string {
	return fmt.Sprintf("Float(%v)", f.Value)
}

func NewFloat(value float64) *Float {
	return &Float{Value: value}
}
