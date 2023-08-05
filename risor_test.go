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
			expectedErr: "undefined variable: keys",
		},
		{
			input:       "any([0, 0, 1])",
			expectedErr: "undefined variable: any",
		},
		{
			input:       "string(42)",
			expectedErr: "undefined variable: string",
		},
	}
	for _, tc := range testCases {
		_, err := Eval(context.Background(), tc.input)
		require.NotNil(t, err)
		require.Equal(t, tc.expectedErr, err.Error())
	}
}

func TestWithBuiltins(t *testing.T) {
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
	}
	for _, tc := range testCases {
		result, err := Eval(context.Background(), tc.input, WithDefaultBuiltins())
		require.Nil(t, err)
		require.Equal(t, tc.expected, result)
	}
}

func TestConfirmNoModules(t *testing.T) {
	_, err := Eval(context.Background(), "json.marshal(42)")
	require.NotNil(t, err)
	require.Equal(t, errors.New("undefined variable: json"), err)
}

func TestWithModules(t *testing.T) {
	result, err := Eval(context.Background(), "json.marshal(42)", WithDefaultModules())
	require.Nil(t, err)
	require.Equal(t, object.NewString("42"), result)
}

func TestWithCode(t *testing.T) {
	ast, err := parser.Parse(context.Background(), "x := 3")
	require.Nil(t, err)

	main, err := compiler.Compile(ast)
	require.Nil(t, err)

	result, err := Eval(context.Background(), "x + 1", WithCode(main))
	require.Nil(t, err)
	require.Equal(t, object.NewInt(4), result)
}

func TestWithVirtualOSStdin(t *testing.T) {

	ctx := context.Background()
	stdinBuf := ros.NewInMemoryFile([]byte("hello"))
	vos := ros.NewVirtualOS(ctx, ros.WithStdin(stdinBuf))
	ctx = ros.WithOS(ctx, vos)

	result, err := Eval(ctx, "os.stdin.read()", WithDefaultModules())
	require.Nil(t, err)
	require.Equal(t, object.NewByteSlice([]byte("hello")), result)
}

func TestWithVirtualOSStdout(t *testing.T) {

	ctx := context.Background()
	stdoutBuf := ros.NewInMemoryFile(nil)
	vos := ros.NewVirtualOS(ctx, ros.WithStdout(stdoutBuf))
	ctx = ros.WithOS(ctx, vos)

	result, err := Eval(ctx, "os.stdout.write('foo')", WithDefaultModules())
	require.Nil(t, err)
	require.Equal(t, object.NewInt(int64(len("foo"))), result)

	require.Equal(t, "foo", string(stdoutBuf.Bytes()))
}

func TestStdinList(t *testing.T) {

	ctx := context.Background()
	stdinBuf := ros.NewInMemoryFile([]byte("foo\nbar!"))
	vos := ros.NewVirtualOS(ctx, ros.WithStdin(stdinBuf))
	ctx = ros.WithOS(ctx, vos)

	result, err := Eval(ctx, "list(os.stdin)", WithDefaultModules(), WithDefaultBuiltins())
	require.Nil(t, err)
	require.Equal(t, object.NewList([]object.Object{
		object.NewString("foo"),
		object.NewString("bar!"),
	}), result)
}
