package object

import (
	"fmt"

	"github.com/cloudcmds/tamarin/v2/op"
)

type NilType struct{}

func (n *NilType) Type() Type {
	return NIL
}

func (n *NilType) Inspect() string {
	return "nil"
}

func (n *NilType) String() string {
	return "nil()"
}

func (n *NilType) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (n *NilType) Interface() interface{} {
	return nil
}

func (n *NilType) HashKey() HashKey {
	return HashKey{Type: n.Type()}
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

func (n *NilType) IsTruthy() bool {
	return false
}

func (n *NilType) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for nil: %v", opType))
}
