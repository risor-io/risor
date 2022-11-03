package object

import (
	"fmt"
	"sort"
	"strings"
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
func (b *Boolean) InvokeMethod(method string, env Environment, args ...Object) Object {
	if method == "methods" {
		static := []string{"methods"}
		dynamic := env.Names("bool.")

		var names []string
		names = append(names, static...)
		for _, e := range dynamic {
			bits := strings.Split(e, ".")
			names = append(names, bits[1])
		}
		sort.Strings(names)

		result := make([]Object, len(names))
		for i, txt := range names {
			result[i] = &String{Value: txt}
		}
		return &Array{Elements: result}
	}
	return nil
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (b *Boolean) ToInterface() interface{} {
	return b.Value
}
