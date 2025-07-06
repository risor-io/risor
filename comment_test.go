package risor

import (
	"context"
	"testing"

	"github.com/risor-io/risor/ast"
	"github.com/risor-io/risor/parser"
	"github.com/stretchr/testify/require"
)

// TestCommentParsing verifies that comments are properly parsed into the AST
func TestCommentParsing(t *testing.T) {
	input := `// Single-line comment with //
x := 42
# Single-line comment with #
y := "hello"
/*
Multi-line comment
spanning multiple lines
*/
z := true`

	program, err := parser.Parse(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, program)

	statements := program.Statements()
	require.Len(t, statements, 6) // 3 comments + 3 variable declarations

	// Check first comment
	comment1, ok := statements[0].(*ast.Comment)
	require.True(t, ok)
	require.Equal(t, "// Single-line comment with //", comment1.Text())

	// Check first variable declaration
	var1, ok := statements[1].(*ast.Var)
	require.True(t, ok)
	name, value := var1.Value()
	require.Equal(t, "x", name)
	require.NotNil(t, value)

	// Check second comment
	comment2, ok := statements[2].(*ast.Comment)
	require.True(t, ok)
	require.Equal(t, "# Single-line comment with #", comment2.Text())

	// Check second variable declaration
	var2, ok := statements[3].(*ast.Var)
	require.True(t, ok)
	name, value = var2.Value()
	require.Equal(t, "y", name)
	require.NotNil(t, value)

	// Check multi-line comment
	comment3, ok := statements[4].(*ast.Comment)
	require.True(t, ok)
	expectedMultiLine := "/*\nMulti-line comment\nspanning multiple lines\n*/"
	require.Equal(t, expectedMultiLine, comment3.Text())

	// Check third variable declaration
	var3, ok := statements[5].(*ast.Var)
	require.True(t, ok)
	name, value = var3.Value()
	require.Equal(t, "z", name)
	require.NotNil(t, value)
}

// TestCommentInExpression verifies comments work alongside expressions
func TestCommentInExpression(t *testing.T) {
	input := `x := 1 + 2 // inline comment`

	program, err := parser.Parse(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, program)

	statements := program.Statements()
	require.Len(t, statements, 2) // 1 declaration + 1 comment

	// Check variable declaration
	varDecl, ok := statements[0].(*ast.Var)
	require.True(t, ok)
	name, _ := varDecl.Value()
	require.Equal(t, "x", name)

	// Check inline comment
	comment, ok := statements[1].(*ast.Comment)
	require.True(t, ok)
	require.Equal(t, "// inline comment", comment.Text())
}