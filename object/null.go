package object

type NilType struct{}

func (n *NilType) Type() Type {
	return NIL
}

func (n *NilType) Inspect() string {
	return "nil"
}

func (n *NilType) InvokeMethod(method string, args ...Object) Object {
	return NewError("type error: %s object has no method %s", n.Type(), method)
}

func (n *NilType) ToInterface() interface{} {
	return nil
}

func (n *NilType) Compare(other Object) (int, error) {
	return CompareTypes(n, other), nil
}
