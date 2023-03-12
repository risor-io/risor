package vm

import (
	"fmt"
	"testing"

	"github.com/cloudcmds/tamarin/internal/compiler"
	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/internal/symbol"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/parser"
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
		op.Code(op.Add),
	}
	vm := New(&compiler.Bytecode{
		Scopes: []*compiler.Scope{
			{
				Constants:    constants,
				Instructions: code,
				Symbols:      symbol.NewTable(),
			},
		},
	})
	err := vm.Run()
	require.Nil(t, err)

	tos, ok := vm.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewInt(7), tos)
}

// https://opensource.com/article/18/4/introduction-python-bytecode

func TestAddCompilationAndExecution(t *testing.T) {

	program, err := parser.Parse(`
	x := 11
	y := 12
	x + y
	`)
	require.Nil(t, err)

	c := compiler.New(compiler.Options{})
	bytecode, err := c.Compile(program)
	require.Nil(t, err)
	scope := bytecode.Scopes[0]

	consts := scope.Constants
	require.Len(t, consts, 2)

	c1, ok := consts[0].(*object.Int)
	require.True(t, ok)
	require.Equal(t, int64(11), c1.Value())

	c2, ok := consts[1].(*object.Int)
	require.True(t, ok)
	require.Equal(t, int64(12), c2.Value())

	vm := New(bytecode)
	require.Nil(t, vm.Run())

	tos, ok := vm.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewInt(23), tos)
}

func TestConditional(t *testing.T) {

	program, err := parser.Parse(`
	x := 20
	if x > 10 {
		x = 99
	}
	x
	`)
	require.Nil(t, err)

	c := compiler.New(compiler.Options{})
	bytecode, err := c.Compile(program)
	require.Nil(t, err)

	vm := New(bytecode)
	require.Nil(t, vm.Run())

	tos, ok := vm.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewInt(99), tos)
}

func TestConditional3(t *testing.T) {
	result, err := Run(`
	x := 5
	y := 10
	if x > 1 {
		y
	} else {
		99
	}
	`)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, object.NewInt(10), result)
}

func TestConditional4(t *testing.T) {
	result, err := Run(`
	x := 5
	y := 22
	z := 33
	if x < 1 {
		x = y
	} else {
		x = z
	}
	x
	`)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, object.NewInt(33), result)
}

func TestLoop(t *testing.T) {
	result, err := Run(`
	y := 0
	for {
		y = y + 1
		if y > 10 {
			break
		}
	}
	y
	`)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, object.NewInt(11), result)
}

func TestAssign(t *testing.T) {
	result, err := Run(`
	y := 99
	y  = 3
	y += 6
	y /= 9
	y *= 2
	y
	`)
	fmt.Println(result, err)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, object.NewInt(2), result)
}

func TestCall(t *testing.T) {
	result, err := Run(`
	f := func(x) { x + 42 }
	y := f(1)
	z := y + 1
	z
	`)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, object.NewInt(43), result)
}
