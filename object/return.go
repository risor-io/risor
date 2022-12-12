package object

// ReturnValue wraps Object and implements Object interface.
type ReturnValue struct {
	// Value is the object that is to be returned
	Value Object
}

// Type returns the type of this object.
func (rv *ReturnValue) Type() Type {
	return RETURN_VALUE
}

// Inspect returns a string-representation of the given object.
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// InvokeMethod invokes a method against the object.
// (Built-in methods only.)
func (rv *ReturnValue) InvokeMethod(method string, args ...Object) Object {
	return nil
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (rv *ReturnValue) ToInterface() interface{} {
	return "<RETURN_VALUE>"
}
