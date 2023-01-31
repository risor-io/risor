package vm

import (
	"testing"

	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/object"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	constants := []object.Object{
		object.NewInt(3),
		object.NewInt(4),
	}
	code := []op.Code{
		op.LoadConst,
		0,
		0,
		op.LoadConst,
		1,
		0,
		op.BinaryOp,
	}
	vm := New(constants, code)
	err := vm.Run()
	require.Nil(t, err)

	tos, ok := vm.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewInt(7), tos)
}

// https://opensource.com/article/18/4/introduction-python-bytecode

func TestCall(t *testing.T) {
	constants := []object.Object{
		object.NewInt(3),
		object.NewInt(4),
		object.NewFunctionWithCode(3, 4),
	}
	code := []op.Code{
		op.JumpForward,
		// Main addr
		10,
		0,
		// Function
		op.LoadFast,
		0,
		op.LoadFast,
		1,
		op.BinaryOp,
		op.Code(op.Add),
		op.ReturnValue,
		1,
		// Main
		op.LoadConst,
		0,
		0,
		op.LoadConst,
		1,
		0,
		op.LoadConst,
		2,
		0,
		op.Call,
		// Arg count
		2,
	}
	vm := New(constants, code)
	err := vm.Run()
	require.Nil(t, err)

	tos, ok := vm.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewInt(7), tos)
}
