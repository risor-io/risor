package compiler

import (
	"testing"

	"github.com/cloudcmds/tamarin/v2/ast"
	"github.com/cloudcmds/tamarin/v2/op"
	"github.com/stretchr/testify/require"
)

func TestNil(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	scope, err := c.Compile(&ast.Nil{})
	require.Nil(t, err)
	require.Len(t, scope.Instructions, 1)
	instr := scope.Instructions[0]
	require.Equal(t, op.Nil, op.Code(instr))
}

// func TestAdd(t *testing.T) {
// 	program, err := parser.Parse(`
// 	x := 1
// 	y := 2
// 	x + y
// 	`)
// 	require.Nil(t, err)

// 	c := New(Options{})
// 	err = c.compile(program)
// 	require.Nil(t, err)

// 	instrs := c.instructions
// 	for {
// 		var opCode op.Code
// 		var operands []int
// 		opCode, operands, instrs = ReadInstruction(instrs)
// 		opInfo := op.GetInfo(opCode)
// 		fmt.Println(opInfo.Name, operands, instrs)
// 		if len(instrs) == 0 {
// 			break
// 		}
// 	}
// 	require.True(t, false)

// 	// vm.New()
// }
