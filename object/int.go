package object

import (
	"fmt"
	"math"

	"github.com/cloudcmds/tamarin/v2/op"
)

// Int wraps int64 and implements Object and Hashable interfaces.
type Int struct {
	// value holds the int64 wrapped by this object.
	value int64
}

func (i *Int) Inspect() string {
	return fmt.Sprintf("%d", i.value)
}

func (i *Int) Type() Type {
	return INT
}

func (i *Int) Value() int64 {
	return i.value
}

func (i *Int) HashKey() HashKey {
	return HashKey{Type: i.Type(), IntValue: i.value}
}

func (i *Int) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (i *Int) Interface() interface{} {
	return i.value
}

func (i *Int) String() string {
	return fmt.Sprintf("int(%v)", i.value)
}

func (i *Int) Compare(other Object) (int, error) {
	switch other := other.(type) {
	case *Float:
		thisFloat := float64(i.value)
		if thisFloat == other.value {
			return 0, nil
		}
		if thisFloat > other.value {
			return 1, nil
		}
		return -1, nil
	case *Int:
		if i.value == other.value {
			return 0, nil
		}
		if i.value > other.value {
			return 1, nil
		}
		return -1, nil
	default:
		return CompareTypes(i, other), nil
	}
}

func (i *Int) Equals(other Object) Object {
	switch other := other.(type) {
	case *Int:
		if i.value == other.value {
			return True
		}
	case *Float:
		if float64(i.value) == other.value {
			return True
		}
	}
	return False
}

func (i *Int) IsTruthy() bool {
	return i.value != 0
}

func (i *Int) RunOperation(opType op.BinaryOpType, right Object) Object {
	switch right := right.(type) {
	case *Int:
		return i.runOperationInt(opType, right.value)
	case *Float:
		return i.runOperationFloat(opType, right.value)
	default:
		return NewError(fmt.Errorf("eval error: unsupported operation for int: %v on type %s", opType, right.Type()))
	}
}

func (i *Int) runOperationInt(opType op.BinaryOpType, right int64) Object {
	switch opType {
	case op.Add:
		return NewInt(i.value + right)
	case op.Subtract:
		return NewInt(i.value - right)
	case op.Multiply:
		return NewInt(i.value * right)
	case op.Divide:
		return NewInt(i.value / right)
	case op.Modulo:
		return NewInt(i.value % right)
	case op.Xor:
		return NewInt(i.value ^ right)
	case op.Power:
		return NewInt(int64(math.Pow(float64(i.value), float64(right))))
	case op.LShift:
		return NewInt(i.value << uint(right))
	case op.RShift:
		return NewInt(i.value >> uint(right))
	case op.BitwiseAnd:
		return NewInt(i.value & right)
	case op.BitwiseOr:
		return NewInt(i.value | right)
	default:
		return NewError(fmt.Errorf("eval error: unsupported operation for int: %v on type int", opType))
	}
}

func (i *Int) runOperationFloat(opType op.BinaryOpType, right float64) Object {
	iValue := float64(i.value)
	switch opType {
	case op.Add:
		return NewFloat(iValue + right)
	case op.Subtract:
		return NewFloat(iValue - right)
	case op.Multiply:
		return NewFloat(iValue * right)
	case op.Divide:
		return NewFloat(iValue / right)
	case op.Power:
		return NewInt(int64(math.Pow(float64(i.value), float64(right))))
	default:
		return NewError(fmt.Errorf("eval error: unsupported operation for int: %v on type float", opType))
	}
}

func (i *Int) Cost() int {
	return 8
}

func NewInt(value int64) *Int {
	if value >= 0 && value < tableSize {
		return intCache[value]
	}
	return &Int{value: value}
}

const tableSize = 256

var intCache = []*Int{}

func init() {
	intCache = make([]*Int, tableSize)
	for i := 0; i < tableSize; i++ {
		intCache[i] = &Int{value: int64(i)}
	}
}
