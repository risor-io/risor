package vm

import (
	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/object"
)

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
