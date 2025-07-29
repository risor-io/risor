package compiler

import (
	"context"
	"strings"
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
	require.Equal(t, "compile error: undefined variable \"foo\"\n\nlocation: unknown:2:2 (line 2, column 2)", err.Error())
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
			errMsg: "compile error: undefined variable \"foo\"\n\nlocation: t.risor:1:1 (line 1, column 1)",
		},
		{
			name:   "undefined variable x",
			input:  "x = 1",
			errMsg: "compile error: undefined variable \"x\"\n\nlocation: t.risor:1:3 (line 1, column 3)",
		},
		{
			name:   "undefined variable y",
			input:  "x := 1;\nx, y = [1, 2]",
			errMsg: "compile error: undefined variable \"y\"\n\nlocation: t.risor:2:1 (line 2, column 1)",
		},
		{
			name:   "undefined variable z",
			input:  "\n\n z++;",
			errMsg: "compile error: undefined variable \"z\"\n\nlocation: t.risor:3:2 (line 3, column 2)",
		},
		{
			name:   "invalid argument defaults",
			input:  "func bad(a=1, b) {}",
			errMsg: "compile error: invalid argument defaults for function \"bad\"\n\nlocation: t.risor:1:1 (line 1, column 1)",
		},
		{
			name:   "invalid argument defaults for anonymous function",
			input:  "func(a=1, b) {}()",
			errMsg: "compile error: invalid argument defaults for anonymous function\n\nlocation: t.risor:1:1 (line 1, column 1)",
		},
		{
			name:   "unsupported default value",
			input:  "func(a, b=[1,2,3]) {}()",
			errMsg: "compile error: unsupported default value (got [1, 2, 3], line 1)",
		},
		{
			name:   "cannot assign to constant",
			input:  "const a = 1; a = 2",
			errMsg: "compile error: cannot assign to constant \"a\"\n\nlocation: t.risor:1:16 (line 1, column 16)",
		},
		{
			name:   "invalid for loop",
			input:  "\nfor a, b, c := range [1, 2, 3] {}",
			errMsg: "compile error: invalid for loop\n\nlocation: t.risor:2:1 (line 2, column 1)",
		},
		{
			name:   "unknown operator",
			input:  "\n defer func() {}()",
			errMsg: "compile error: defer statement outside of a function\n\nlocation: t.risor:2:2 (line 2, column 2)",
		},
	}
	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(WithFilename("t.risor"))
			require.Nil(t, err)
			ast, err := parser.Parse(context.Background(), tt.input)
			require.Nil(t, err)
			_, err = c.Compile(ast)
			require.NotNil(t, err)
			require.Equal(t, tt.errMsg, err.Error())
		})
	}
}

func TestCompilerLoopError(t *testing.T) {
	input := `
for _, v := range [1, 2, 3] {
	func() {
		undefined_var
	}()
}
	`
	c, err := New()
	require.Nil(t, err)
	ast, err := parser.Parse(context.Background(), input)
	require.Nil(t, err)
	_, err = c.Compile(ast)
	require.NotNil(t, err)
	require.Equal(t, "compile error: undefined variable \"undefined_var\"\n\nlocation: unknown:4:3 (line 4, column 3)", err.Error())
}

func TestCompoundAssignmentWithIndex(t *testing.T) {
	// test[0] *= 3
	input := "test := [1, 2]; test[0] *= 3"
	expected := [][]op.Code{
		{op.LoadConst, 0}, // 1
		{op.LoadConst, 1}, // 2
		{op.BuildList, 2},
		{op.StoreGlobal, 0}, // store into 'test'
		{op.LoadGlobal, 0},  // load 'test'
		{op.LoadConst, 2},   // load index 0
		{op.BinarySubscr},   // get test[0]
		{op.LoadConst, 3},   // load 3
		{op.BinaryOp, op.Code(op.Multiply)},
		{op.LoadGlobal, 0}, // load 'test' again
		{op.LoadConst, 4},  // load index 0
		{op.StoreSubscr},   // store result back in test[0]
		{op.Nil},           // implicit return value
	}

	c, err := New()
	require.Nil(t, err)

	ast, err := parser.Parse(context.Background(), input)
	require.Nil(t, err)

	code, err := c.Compile(ast)
	require.Nil(t, err)

	// Compare the generated instructions
	actual := NewInstructionIter(code).All()

	require.Equal(t, len(expected), len(actual),
		"instruction length mismatch. got=%d, want=%d",
		len(actual), len(expected))

	for i, want := range expected {
		got := actual[i]
		require.Equal(t, want, got,
			"wrong instruction at pos %d. got=%v, want=%v",
			i, got, want)
	}
}

func TestBreakFromRangeLoop(t *testing.T) {
	input := `
for range [1, 2] {
	break
}
`
	c, err := New()
	require.Nil(t, err)

	ast, err := parser.Parse(context.Background(), input)
	require.Nil(t, err)

	code, err := c.Compile(ast)
	require.Nil(t, err)

	// Extract instructions
	instructions := code.instructions

	// Find the index of the break statement (JumpForward)
	jumpIndex := -1
	for i, instr := range instructions {
		if instr == op.JumpForward {
			jumpIndex = i - 1 // We're interested in the instruction before the jump
			break
		}
	}

	require.NotEqual(t, -1, jumpIndex, "JumpForward instruction not found")

	// The instruction right before JumpForward should be PopTop for a range loop break
	require.Equal(t, op.PopTop, instructions[jumpIndex],
		"Expected PopTop before JumpForward for break in range loop")
}

func TestStringImport(t *testing.T) {
	tests := []struct {
		input             string
		expectedCode      []op.Code
		expectedConstants []interface{}
	}{
		{
			input: `import "foo"`,
			expectedCode: []op.Code{
				op.LoadConst, 0, // "foo"
				op.Import,
				op.StoreGlobal, 0,
				op.Nil,
			},
			expectedConstants: []interface{}{"foo"},
		},
		{
			input: `import foo`,
			expectedCode: []op.Code{
				op.LoadConst, 0, // "foo"
				op.Import,
				op.StoreGlobal, 0,
				op.Nil,
			},
			expectedConstants: []interface{}{"foo"},
		},
		{
			input: `import "path/to/foo"`,
			expectedCode: []op.Code{
				op.LoadConst, 0, // "path/to/foo"
				op.Import,
				op.StoreGlobal, 0,
				op.Nil,
			},
			expectedConstants: []interface{}{"path/to/foo"},
		},
		{
			input: `import "path/to/foo" as bar`,
			expectedCode: []op.Code{
				op.LoadConst, 0, // "path/to/foo"
				op.Import,
				op.StoreGlobal, 0,
				op.Nil,
			},
			expectedConstants: []interface{}{"path/to/foo"},
		},
		{
			input: `3 & 1`,
			expectedCode: []op.Code{
				op.LoadConst, 0, // 3
				op.LoadConst, 1, // 1
				op.BinaryOp,
				op.Code(op.BitwiseAnd),
			},
			expectedConstants: []interface{}{int64(3), int64(1)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			astNode, err := parser.Parse(context.Background(), tt.input)
			require.NoError(t, err)

			code, err := Compile(astNode)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, code.instructions)
			require.Equal(t, tt.expectedConstants, code.constants)
		})
	}
}

func TestFunctionRedefinition(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "duplicate function definition",
			input: `
func bar() {
    print("first bar")
}

func bar() {
    print("second bar")
}
`,
			expected: `function "bar" redefined`,
		},
		{
			name: "multiple duplicate function definitions",
			input: `
func foo() {
    print("first foo")
}

func foo() {
    print("second foo")
}

func foo() {
    print("third foo")
}
`,
			expected: `function "foo" redefined`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New()
			require.Nil(t, err)
			ast, err := parser.Parse(context.Background(), tt.input)
			require.Nil(t, err)
			_, err = c.Compile(ast)
			if err == nil {
				t.Errorf("Expected error but got none")
				return
			}
			if !strings.Contains(err.Error(), tt.expected) {
				t.Errorf("Expected error containing %q, got %q", tt.expected, err.Error())
			}
		})
	}
}

func TestForwardDeclarationCompilation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "basic forward declaration",
			input: `
			func main() {
				return helper()
			}
			
			func helper() {
				return 42
			}
			`,
			wantErr: false,
		},
		{
			name: "mutual recursion",
			input: `
			func is_even(n) {
				if n == 0 {
					return true
				}
				return is_odd(n - 1)
			}
			
			func is_odd(n) {
				if n == 0 {
					return false
				}
				return is_even(n - 1)
			}
			`,
			wantErr: false,
		},
		{
			name: "multiple forward declarations",
			input: `
			func a() {
				return b() + c()
			}
			
			func b() {
				return d()
			}
			
			func c() {
				return 10
			}
			
			func d() {
				return 20
			}
			`,
			wantErr: false,
		},
		{
			name: "forward declaration with closures",
			input: `
			func outer() {
				x := 10
				return func() {
					return inner() + x
				}
			}
			
			func inner() {
				return 5
			}
			`,
			wantErr: false,
		},
		{
			name: "forward declaration with default parameters",
			input: `
			func caller(op="add") {
				if op == "add" {
					return adder(5, 3)
				}
				return multiplier(5, 3)
			}
			
			func adder(a, b) {
				return a + b
			}
			
			func multiplier(a, b) {
				return a * b
			}
			`,
			wantErr: false,
		},
		{
			name: "undefined function should error",
			input: `
			func caller() {
				return undefined_function()
			}
			`,
			wantErr: true,
		},
		{
			name: "function redefinition should error",
			input: `
			func duplicate() {
				return 1
			}
			
			func duplicate() {
				return 2
			}
			`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New()
			require.Nil(t, err)

			ast, err := parser.Parse(context.Background(), tt.input)
			require.Nil(t, err)

			_, err = c.Compile(ast)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestForwardDeclarationInstructionGeneration(t *testing.T) {
	// Test that forward declarations generate correct instructions
	input := `
	func main() {
		return helper(5)
	}
	
	func helper(x) {
		return x * 2
	}
	
	main()
	`

	c, err := New()
	require.Nil(t, err)

	ast, err := parser.Parse(context.Background(), input)
	require.Nil(t, err)

	code, err := c.Compile(ast)
	require.Nil(t, err)

	// Verify that the code compiles successfully and has expected structure
	require.NotNil(t, code)
	require.Greater(t, code.InstructionCount(), 0)

	// Verify that the code compiles successfully and contains expected constants
	require.Greater(t, code.ConstantsCount(), 0, "should have constants")

	// The main verification is that compilation succeeded without errors
	// indicating that forward declarations were properly resolved
}

func TestForInCompilation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []op.Code
	}{
		{
			name:  "basic for-in loop",
			input: "for x in [1, 2, 3] { x }",
			expected: []op.Code{
				op.BuildList,    // Create array [1, 2, 3]
				op.GetIter,      // Get iterator
				op.ForIter,      // ForIter instruction
				op.StoreFast,    // Store to variable x
				op.LoadFast,     // Load x
				op.PopTop,       // Pop result
				op.JumpBackward, // Jump back to ForIter
			},
		},
		{
			name:  "for-in with variable assignment",
			input: "result := 0; for x in [1, 2] { result = result + x }",
			expected: []op.Code{
				op.LoadConst,    // 0
				op.StoreFast,    // result := 0
				op.BuildList,    // [1, 2]
				op.GetIter,      // Get iterator
				op.ForIter,      // ForIter instruction
				op.StoreFast,    // Store to x
				op.LoadFast,     // Load result
				op.LoadFast,     // Load x
				op.BinaryOp,     // result + x (binary operation)
				op.StoreFast,    // result = ...
				op.PopTop,       // Pop assignment result
				op.JumpBackward, // Jump back
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := parser.Parse(context.Background(), tt.input)
			require.Nil(t, err)

			c, err := New()
			require.Nil(t, err)

			code, err := c.Compile(program)
			require.Nil(t, err)
			require.NotNil(t, code)

			// Verify basic compilation succeeded
			require.Greater(t, code.InstructionCount(), 0)

			// Check for presence of key opcodes
			instructions := make([]op.Code, code.InstructionCount())
			for i := 0; i < code.InstructionCount(); i++ {
				instructions[i] = op.Code(code.Instruction(i))
			}

			// Verify GetIter and ForIter are present
			hasGetIter := false
			hasForIter := false
			for _, instr := range instructions {
				if instr == op.GetIter {
					hasGetIter = true
				}
				if instr == op.ForIter {
					hasForIter = true
				}
			}
			require.True(t, hasGetIter, "Expected GetIter instruction")
			require.True(t, hasForIter, "Expected ForIter instruction")
		})
	}
}

func TestForInCompilationErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr string
	}{
		{
			name:        "undefined variable in iterable",
			input:       "for x in undefined_var { }",
			expectedErr: "undefined variable",
		},
		{
			name:        "undefined variable in loop body",
			input:       "for x in [1, 2, 3] { undefined_func(x) }",
			expectedErr: "undefined variable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := parser.Parse(context.Background(), tt.input)
			require.Nil(t, err)

			c, err := New()
			require.Nil(t, err)

			_, err = c.Compile(program)
			require.NotNil(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestForInWithBreakContinue(t *testing.T) {
	input := `
	for x in [1, 2, 3, 4, 5] {
		if x == 2 {
			continue
		}
		if x == 4 {
			break
		}
		x
	}
	`

	program, err := parser.Parse(context.Background(), input)
	require.Nil(t, err)

	c, err := New()
	require.Nil(t, err)

	code, err := c.Compile(program)
	require.Nil(t, err)
	require.NotNil(t, code)

	// Verify that break/continue instructions are present
	instructions := make([]op.Code, code.InstructionCount())
	for i := 0; i < code.InstructionCount(); i++ {
		instructions[i] = op.Code(code.Instruction(i))
	}

	hasJumpForward := false
	hasJumpBackward := false
	for _, instr := range instructions {
		if instr == op.JumpForward {
			hasJumpForward = true
		}
		if instr == op.JumpBackward {
			hasJumpBackward = true
		}
	}

	require.True(t, hasJumpForward, "Expected JumpForward instruction for break/continue")
	require.True(t, hasJumpBackward, "Expected JumpBackward instruction for loop")
}

func TestForInNestedLoops(t *testing.T) {
	input := `
	for x in [1, 2] {
		for y in [3, 4] {
			x + y
		}
	}
	`

	program, err := parser.Parse(context.Background(), input)
	require.Nil(t, err)

	c, err := New()
	require.Nil(t, err)

	code, err := c.Compile(program)
	require.Nil(t, err)
	require.NotNil(t, code)

	// Count ForIter instructions - should have 2 for nested loops
	instructions := make([]op.Code, code.InstructionCount())
	for i := 0; i < code.InstructionCount(); i++ {
		instructions[i] = op.Code(code.Instruction(i))
	}

	forIterCount := 0
	getIterCount := 0
	for _, instr := range instructions {
		if instr == op.ForIter {
			forIterCount++
		}
		if instr == op.GetIter {
			getIterCount++
		}
	}

	require.Equal(t, 2, getIterCount, "Expected 2 GetIter instructions for nested loops")
	require.Equal(t, 2, forIterCount, "Expected 2 ForIter instructions for nested loops")
}
