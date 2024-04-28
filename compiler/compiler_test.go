package compiler

import (
	"context"
	"testing"

	"github.com/risor-io/risor/ast"
	"github.com/risor-io/risor/op"
	"github.com/risor-io/risor/parser"
	"github.com/risor-io/risor/token"
	"github.com/stretchr/testify/require"
)

func TestNil(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	scope, err := c.Compile(&ast.Nil{})
	require.Nil(t, err)
	require.Equal(t, 1, scope.InstructionCount())
	instr := scope.Instruction(0)
	require.Equal(t, op.Nil, op.Code(instr))
}

func TestUndefinedVariable(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	_, err = c.Compile(ast.NewIdent(token.Token{
		Type:          token.IDENT,
		Literal:       "foo",
		StartPosition: token.Position{Line: 1, Column: 1},
	}))
	require.NotNil(t, err)
	require.Equal(t, "compile error: undefined variable \"foo\" (line 2)", err.Error())
}

func TestCompileErrors(t *testing.T) {
	testCase := []struct {
		name   string
		input  string
		errMsg string
	}{
		{
			name:   "undefined variable foo",
			input:  "foo",
			errMsg: "compile error: undefined variable \"foo\" (line 1)",
		},
		{
			name:   "undefined variable x",
			input:  "x = 1",
			errMsg: "compile error: undefined variable \"x\" (line 1)",
		},
		{
			name:   "undefined variable y",
			input:  "x := 1;\nx, y = [1, 2]",
			errMsg: "compile error: undefined variable \"y\" (line 2)",
		},
		{
			name:   "undefined variable z",
			input:  "\n\n z++;",
			errMsg: "compile error: undefined variable \"z\" (line 3)",
		},
	}
	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New()
			require.Nil(t, err)
			ast, err := parser.Parse(context.Background(), tt.input)
			require.Nil(t, err)
			_, err = c.Compile(ast)
			require.NotNil(t, err)
			require.Equal(t, tt.errMsg, err.Error())
		})
	}
}
