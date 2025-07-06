package github

import (
	"fmt"
	
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

// base is a common struct for all GitHub objects
type base struct{}

func (b *base) SetAttr(name string, value object.Object) error {
	return object.NewError(fmt.Errorf("attribute error: %s", name))
}

func (b *base) IsTruthy() bool {
	return true
}

func (b *base) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("unsupported operation"))
}

func (b *base) Cost() int {
	return 0
}