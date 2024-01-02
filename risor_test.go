package risor

import (
	"context"
	"errors"
	"testing"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/object"
	ros "github.com/risor-io/risor/os"
	"github.com/risor-io/risor/parser"
	"github.com/stretchr/testify/require"
)

func ptrTo[T any](v T) *T {
	return &v
}

func TestBasicUsage(t *testing.T) {
	result, err := Eval(context.Background(), "1 + 1")
	require.Nil(t, err)
	require.Equal(t, object.NewInt(2), result)
}

func TestConfirmNoBuiltins(t *testing.T) {
	type testCase struct {
		input       string
		expectedErr string
	}
	testCases := []testCase{
		{
			input:       "keys({foo: 1})",
			expectedErr: "compile error: undefined variable \"keys\"",
		},
		{
			input:       "any([0, 0, 1])",
			expectedErr: "compile error: undefined variable \"any\"",
		},
		{
			input:       "string(42)",
			expectedErr: "compile error: undefined variable \"string\"",
		},
	}
	for _, tc := range testCases {
		_, err := Eval(context.Background(), tc.input, WithoutDefaultGlobals())
		require.NotNil(t, err)
		require.Equal(t, tc.expectedErr, err.Error())
	}
}

func TestDefaultGlobals(t *testing.T) {
	type testCase struct {
		input    string
		expected object.Object
	}
	testCases := []testCase{
		{
			input:    "keys({foo: 1})",
			expected: object.NewList([]object.Object{object.NewString("foo")}),
		},
		{
			input:    "any([0, 0, 1])",
			expected: object.True,
		},
		{
			input:    "try(func() { error(\"boom\") }, 42)",
			expected: object.NewInt(42),
		},
		{
			input:    "string(42)",
			expected: object.NewString("42"),
		},
		{
			input:    "json.marshal(42)",
			expected: object.NewString("42"),
		},
	}
	for _, tc := range testCases {
		result, err := Eval(context.Background(), tc.input)
		require.Nil(t, err)
		require.Equal(t, tc.expected, result)
	}
}

func TestWithoutDefaultGlobals(t *testing.T) {
	_, err := Eval(context.Background(), "json.marshal(42)", WithoutDefaultGlobals())
	require.NotNil(t, err)
	require.Equal(t, errors.New("compile error: undefined variable \"json\""), err)
}

func TestWithoutDefaultGlobal(t *testing.T) {
	_, err := Eval(context.Background(), "json.marshal(42)", WithoutDefaultGlobal("json"))
	require.NotNil(t, err)
	require.Equal(t, errors.New("compile error: undefined variable \"json\""), err)
}

func TestWithVirtualOSStdin(t *testing.T) {

	ctx := context.Background()
	stdinBuf := ros.NewInMemoryFile([]byte("hello"))
	vos := ros.NewVirtualOS(ctx, ros.WithStdin(stdinBuf))
	ctx = ros.WithOS(ctx, vos)

	result, err := Eval(ctx, "os.stdin.read()")
	require.Nil(t, err)
	require.Equal(t, object.NewByteSlice([]byte("hello")), result)
}

func TestWithVirtualOSStdout(t *testing.T) {

	ctx := context.Background()
	stdoutBuf := ros.NewInMemoryFile(nil)
	vos := ros.NewVirtualOS(ctx, ros.WithStdout(stdoutBuf))
	ctx = ros.WithOS(ctx, vos)

	result, err := Eval(ctx, "os.stdout.write('foo')")
	require.Nil(t, err)
	require.Equal(t, object.NewInt(int64(len("foo"))), result)

	require.Equal(t, "foo", string(stdoutBuf.Bytes()))
}

func TestStdinList(t *testing.T) {

	ctx := context.Background()
	stdinBuf := ros.NewInMemoryFile([]byte("foo\nbar!"))
	vos := ros.NewVirtualOS(ctx, ros.WithStdin(stdinBuf))
	ctx = ros.WithOS(ctx, vos)

	result, err := Eval(ctx, "list(os.stdin)")
	require.Nil(t, err)
	require.Equal(t, object.NewList([]object.Object{
		object.NewString("foo"),
		object.NewString("bar!"),
	}), result)
}

func TestEvalCode(t *testing.T) {

	ctx := context.Background()

	source := `
	x := 2
	y := 3
	func add(a, b) { a + b }
	result := add(x, y)
	x = 99
	result
	`

	ast, err := parser.Parse(ctx, source)
	require.Nil(t, err)
	code, err := compiler.Compile(ast)
	require.Nil(t, err)

	// Should be able to evaluate the precompiled code any number of times
	for i := 0; i < 100; i++ {
		result, err := EvalCode(ctx, code)
		require.Nil(t, err)
		require.Equal(t, object.NewInt(5), result)
	}
}

func TestCall(t *testing.T) {
	ctx := context.Background()
	source := `
	func add(a, b) { a + b }
	`
	ast, err := parser.Parse(ctx, source)
	require.Nil(t, err)
	code, err := compiler.Compile(ast)
	require.Nil(t, err)

	result, err := Call(ctx, code, "add", []object.Object{
		object.NewInt(9),
		object.NewInt(1),
	})
	require.Nil(t, err)
	require.Equal(t, object.NewInt(10), result)
}
