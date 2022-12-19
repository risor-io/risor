package object

type NilType struct{}

func (n *NilType) Type() Type {
	return NIL
}

func (n *NilType) Inspect() string {
	return "nil"
}

func (n *NilType) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (n *NilType) Interface() interface{} {
	return nil
}

func (n *NilType) Compare(other Object) (int, error) {
	return CompareTypes(n, other), nil
}

func (n *NilType) Equals(other Object) Object {
	if other.Type() == NIL {
		return True
	}
	return False
}
