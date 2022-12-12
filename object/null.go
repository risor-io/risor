package object

type NullType struct{}

func (n *NullType) Type() Type {
	return NULL
}

func (n *NullType) Inspect() string {
	return "null"
}

func (n *NullType) InvokeMethod(method string, args ...Object) Object {
	return NewError("type error: %s object has no method %s", n.Type(), method)
}

func (n *NullType) ToInterface() interface{} {
	return nil
}

func (n *NullType) Compare(other Object) (int, error) {
	return CompareTypes(n, other), nil
}
