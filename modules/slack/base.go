package slack

import (
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

type base struct {
	typeName       object.Type
	interfaceValue interface{}
}

func (b *base) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("type error: cannot set %q on %s object", name, b.typeName)
}

func (b *base) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("type error: unsupported operation for %s object", b.typeName)
}

func (b *base) IsTruthy() bool {
	return true
}

func (b *base) Cost() int {
	return 0
}

func (b *base) Type() object.Type {
	return b.typeName
}

func (b *base) Interface() interface{} {
	return b.interfaceValue
}
