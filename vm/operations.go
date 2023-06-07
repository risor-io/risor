package vm

import (
	"fmt"

	"github.com/cloudcmds/tamarin/v2/object"
	"github.com/cloudcmds/tamarin/v2/op"
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
		panic(fmt.Errorf("unknown comparison operator: %d", opType))
	}
}

func binaryOp(opType op.BinaryOpType, a, b object.Object) object.Object {
	switch opType {
	case op.And:
		aTruthy := a.IsTruthy()
		bTruthy := b.IsTruthy()
		if aTruthy && bTruthy {
			return b
		} else if aTruthy {
			return b // return b because it's falsy
		} else {
			return a // return a because it's falsy
		}
	case op.Or:
		if a.IsTruthy() {
			return a
		}
		return b
	}
	return a.RunOperation(opType, b)
}
