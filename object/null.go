package object

// Null wraps nothing and implements our Object interface.
type Null struct{}

func (n *Null) Type() Type {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) InvokeMethod(method string, args ...Object) Object {
	return NewError("type error: %s object has no method %s", n.Type(), method)
}

func (n *Null) ToInterface() interface{} {
	return "<NULL>"
}

func (n *Null) Compare(other Object) (int, error) {
	return CompareTypes(n, other), nil
}
