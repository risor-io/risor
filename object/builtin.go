package object

import "context"

// BuiltinFunction holds the type of a built-in function.
type BuiltinFunction func(ctx context.Context, args ...Object) Object

// Builtin wraps func and implements Object interface.
type Builtin struct {
	// Value holds the function we wrap.
	Fn BuiltinFunction
}

// Type returns the type of this object.
func (b *Builtin) Type() Type {
	return BUILTIN_OBJ
}

// Inspect returns a string-representation of the given object.
func (b *Builtin) Inspect() string {
	return "builtin function"
}

// InvokeMethod invokes a method against the object.
// (Built-in methods only.)
func (b *Builtin) InvokeMethod(method string, args ...Object) Object {
	return nil
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (b *Builtin) ToInterface() interface{} {
	return "<BUILTIN>"
}
