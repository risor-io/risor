package vm

import (
	"fmt"

	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/object"
)

func compare(opType op.CompareOpType, a, b object.Object) object.Object {

	switch opType {
	case op.Equal:
		return a.Equals(b)
	case op.NotEqual:
		return object.Not(a.Equals(b).(*object.Bool))
	}

	comparable, ok := a.(object.Comparable)
	if !ok {
		return object.NewError(fmt.Errorf("object is not comparable: %T", a))
	}
	value, err := comparable.Compare(b)
	if err != nil {
		return object.NewError(err)
	}

	switch opType {
	case op.LessThan:
		return object.NewBool(value < 0)
	case op.LessThanOrEqual:
		return object.NewBool(value <= 0)
	case op.GreaterThan:
		return object.NewBool(value > 0)
	case op.GreaterThanOrEqual:
		return object.NewBool(value >= 0)
	default:
		panic(fmt.Errorf("unknown comparison operator: %q", opType))
	}
}

func binaryOp(opType op.BinaryOpType, a, b object.Object) object.Object {
	switch opType {
	case op.And:
		if a.IsTruthy() && b.IsTruthy() {
			return b
		}
		return b
	case op.Or:
		if a.IsTruthy() {
			return a
		}
		return b
	}
	return a.RunOperation(opType, b)
}
