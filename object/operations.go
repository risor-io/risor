package object

import (
	"fmt"

	"github.com/cloudcmds/tamarin/v2/op"
)

// Compare two objects using the given comparison operator. An Error object is
// returned if either of the objects is not comparable.
func Compare(opType op.CompareOpType, a, b Object) Object {

	switch opType {
	case op.Equal:
		return a.Equals(b)
	case op.NotEqual:
		return Not(a.Equals(b).(*Bool))
	}

	comparable, ok := a.(Comparable)
	if !ok {
		return NewError(fmt.Errorf("type error: expected a comparable object (got %s)", a.Type()))
	}
	value, err := comparable.Compare(b)
	if err != nil {
		return NewError(err)
	}

	switch opType {
	case op.LessThan:
		return NewBool(value < 0)
	case op.LessThanOrEqual:
		return NewBool(value <= 0)
	case op.GreaterThan:
		return NewBool(value > 0)
	case op.GreaterThanOrEqual:
		return NewBool(value >= 0)
	default:
		panic(fmt.Errorf("unknown object comparison operator: %d", opType))
	}
}

// BinaryOp performs a binary operation on two objects, given an operator.
func BinaryOp(opType op.BinaryOpType, a, b Object) Object {
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
