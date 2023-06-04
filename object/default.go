package object

import (
	"fmt"

	"github.com/cloudcmds/tamarin/v2/op"
)

type DefaultImpl struct{}

func (d *DefaultImpl) Type() Type {
	panic("not implemented")
}

func (d *DefaultImpl) Inspect() string {
	panic("not implemented")
}

func (d *DefaultImpl) Interface() interface{} {
	return nil
}

func (d *DefaultImpl) Equals(other Object) Object {
	if d == other {
		return True
	}
	return False
}

func (d *DefaultImpl) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (d *DefaultImpl) IsTruthy() bool {
	return true
}

func (d *DefaultImpl) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for default: %v", opType))
}
