package risor

import (
	"context"
	"errors"
	"testing"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/object"
	ros "github.com/risor-io/risor/os"
	"github.com/risor-io/risor/parser"
	"github.com/risor-io/risor/vm"
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
			expectedErr: "compile error: undefined variable \"keys\"\n\nlocation: unknown:1:1 (line 1, column 1)",
		},
		{
			input:       "any([0, 0, 1])",
			expectedErr: "compile error: undefined variable \"any\"\n\nlocation: unknown:1:1 (line 1, column 1)",
		},
		{
			input:       "string(42)",
			expectedErr: "compile error: undefined variable \"string\"\n\nlocation: unknown:1:1 (line 1, column 1)",
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

func TestWithDenyList(t *testing.T) {
	type testCase struct {
		input       string
		expectedErr error
	}
	testCases := []testCase{
		{
			input:       "keys({foo: 1})",
			expectedErr: nil,
		},
		{
			input:       "any([0, 0, 1])",
			expectedErr: errors.New("compile error: undefined variable \"any\"\n\nlocation: unknown:1:1 (line 1, column 1)"),
		},
		{
			input:       "json.marshal(42)",
			expectedErr: errors.New("compile error: undefined variable \"json\"\n\nlocation: unknown:1:1 (line 1, column 1)"),
		},
		{
			input:       `os.getenv("USER")`,
			expectedErr: nil,
		},
		{
			input:       "os.exit(1)",
			expectedErr: errors.New(`type error: attribute "exit" not found on module object`),
		},
		{
			input: `from os import exit
exit(1)`,
			expectedErr: errors.New(`import error: cannot import name "exit" from "os"`),
		},
		{
			input:       "cat /etc/issue",
			expectedErr: errors.New("compile error: undefined variable \"cat\"\n\nlocation: unknown:1:1 (line 1, column 1)"),
		},
	}
	for _, tc := range testCases {
		_, err := Eval(context.Background(), tc.input, WithoutGlobals("any", "os.exit", "json", "cat"))
		// t.Logf("want: %q; got: %v", tc.expectedErr, err)
		if tc.expectedErr != nil {
			require.NotNil(t, err)
			require.Equal(t, tc.expectedErr.Error(), err.Error())
			continue
		}
		require.Nil(t, err)
	}
}

func TestWithoutDefaultGlobals(t *testing.T) {
	_, err := Eval(context.Background(), "json.marshal(42)", WithoutDefaultGlobals())
	require.NotNil(t, err)
	require.Equal(t, errors.New("compile error: undefined variable \"json\"\n\nlocation: unknown:1:1 (line 1, column 1)"), err)
}

func TestWithoutDefaultGlobal(t *testing.T) {
	_, err := Eval(context.Background(), "json.marshal(42)", WithoutGlobal("json"))
	require.NotNil(t, err)
	require.Equal(t, errors.New("compile error: undefined variable \"json\"\n\nlocation: unknown:1:1 (line 1, column 1)"), err)
}

func TestWithVirtualOSStdinBuffer(t *testing.T) {
	ctx := context.Background()
	stdinBuf := ros.NewBufferFile([]byte("hello"))
	vos := ros.NewVirtualOS(ctx, ros.WithStdin(stdinBuf))

	result, err := Eval(ctx, "os.stdin.read()", WithOS(vos))
	require.Nil(t, err)
	require.Equal(t, object.NewByteSlice([]byte("hello")), result)
}

func TestWithVirtualOSStdinBufferViaContext(t *testing.T) {
	ctx := context.Background()
	stdinBuf := ros.NewBufferFile([]byte("hello"))
	vos := ros.NewVirtualOS(ctx, ros.WithStdin(stdinBuf))
	ctx = ros.WithOS(ctx, vos)

	result, err := Eval(ctx, "os.stdin.read()")
	require.Nil(t, err)
	require.Equal(t, object.NewByteSlice([]byte("hello")), result)
}

func TestWithVirtualOSStdoutBuffer(t *testing.T) {
	ctx := context.Background()
	stdoutBuf := ros.NewBufferFile(nil)
	vos := ros.NewVirtualOS(ctx, ros.WithStdout(stdoutBuf))

	result, err := Eval(ctx, "os.stdout.write('foo')", WithOS(vos))
	require.Nil(t, err)
	require.Equal(t, object.NewInt(int64(len("foo"))), result)

	require.Equal(t, "foo", string(stdoutBuf.Bytes()))
}

func TestStdinListBuffer(t *testing.T) {
	ctx := context.Background()
	stdinBuf := ros.NewBufferFile([]byte("foo\nbar!"))
	vos := ros.NewVirtualOS(ctx, ros.WithStdin(stdinBuf))

	result, err := Eval(ctx, "list(os.stdin)", WithOS(vos))
	require.Nil(t, err)
	require.Equal(t, object.NewList([]object.Object{
		object.NewString("foo"),
		object.NewString("bar!"),
	}), result)
}

func TestWithVirtualOSStdin(t *testing.T) {
	ctx := context.Background()
	stdinBuf := ros.NewInMemoryFile([]byte("hello"))
	vos := ros.NewVirtualOS(ctx, ros.WithStdin(stdinBuf))

	result, err := Eval(ctx, "os.stdin.read()", WithOS(vos))
	require.Nil(t, err)
	require.Equal(t, object.NewByteSlice([]byte("hello")), result)
}

func TestWithVirtualOSStdout(t *testing.T) {
	ctx := context.Background()
	stdoutBuf := ros.NewInMemoryFile(nil)
	vos := ros.NewVirtualOS(ctx, ros.WithStdout(stdoutBuf))

	result, err := Eval(ctx, "os.stdout.write('foo')", WithOS(vos))
	require.Nil(t, err)
	require.Equal(t, object.NewInt(int64(len("foo"))), result)

	stdoutBuf.Rewind()
	require.Equal(t, "foo", string(stdoutBuf.Bytes()))
}

func TestStdinList(t *testing.T) {
	ctx := context.Background()
	stdinBuf := ros.NewInMemoryFile([]byte("foo\nbar!"))
	vos := ros.NewVirtualOS(ctx, ros.WithStdin(stdinBuf))

	result, err := Eval(ctx, "list(os.stdin)", WithOS(vos))
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

func TestWithoutGlobal1(t *testing.T) {
	cfg := NewConfig(
		WithoutDefaultGlobals(),
		WithGlobals(map[string]any{
			"foo": object.NewInt(1),
			"bar": object.NewInt(2),
		}),
		WithoutGlobal("bar"),
	)
	require.Equal(t, map[string]any{
		"foo": object.NewInt(1),
	}, cfg.Globals())

	require.Equal(t, map[string]any{
		"foo": object.NewInt(1),
	}, cfg.CombinedGlobals())
}

func TestWithoutGlobal2(t *testing.T) {
	cfg := NewConfig(
		WithoutDefaultGlobals(),
		WithGlobals(map[string]any{
			"foo": object.NewInt(1),
			"bar": object.NewInt(2),
		}),
		WithoutGlobal("xyz"),
	)
	require.Equal(t, map[string]any{
		"foo": object.NewInt(1),
		"bar": object.NewInt(2),
	}, cfg.Globals())
}

func TestWithoutGlobal3(t *testing.T) {
	cfg := NewConfig(
		WithGlobal("foo", object.NewInt(1)),
		WithoutGlobal("foo"),
	)
	_, hasFoo := cfg.Globals()["foo"]
	require.False(t, hasFoo)
}

func TestWithoutGlobal4(t *testing.T) {
	cfg := NewConfig(
		WithoutDefaultGlobals(),
		WithGlobals(map[string]any{
			"foo": object.NewBuiltinsModule("foo", map[string]object.Object{
				"bar": object.NewBuiltinsModule("bar", map[string]object.Object{
					"baz": object.NewInt(1),
					"qux": object.NewInt(2),
				}),
			}),
		}),
		WithoutGlobal("foo.bar.baz"),
	)
	require.Equal(t, map[string]any{
		"foo": object.NewBuiltinsModule("foo", map[string]object.Object{
			"bar": object.NewBuiltinsModule("bar", map[string]object.Object{
				"qux": object.NewInt(2),
			}),
		}),
	}, cfg.Globals())
}

func TestWithGlobalOverride(t *testing.T) {
	cfg := NewConfig(
		WithoutDefaultGlobals(),
		WithGlobals(map[string]any{
			"foo": object.NewInt(1),
			"bar": object.NewBuiltinsModule("bar", map[string]object.Object{
				"baz": object.NewInt(1),
				"qux": object.NewInt(2),
			}),
		}),
		WithGlobalOverride("foo", object.NewString("FOO")),
		WithGlobalOverride("bar.baz", object.NewString("BAZ")),
	)

	require.Equal(t, map[string]any{
		"foo": object.NewString("FOO"),
		"bar": object.NewBuiltinsModule("bar", map[string]object.Object{
			"baz": object.NewString("BAZ"),
			"qux": object.NewInt(2),
		}),
	}, cfg.Globals())
}

func TestWithLocalImporter(t *testing.T) {
	result, err := Eval(context.Background(),
		`import data; data.mydata.count`,
		WithLocalImporter("./vm/fixtures"))
	require.Nil(t, err)
	require.Equal(t, object.NewInt(1), result)
}

func TestWithConcurrency(t *testing.T) {
	script := `c := chan(); go func() { c <- 33 }(); <-c`

	result, err := Eval(context.Background(), script, WithConcurrency())
	require.Nil(t, err)
	require.Equal(t, object.NewInt(33), result)

	_, err = Eval(context.Background(), script)
	require.NotNil(t, err)
	require.Equal(t, "eval error: context did not contain a spawn function", err.Error())
}

func TestStructFieldModification(t *testing.T) {
	type Object struct {
		A int
	}

	testCases := []struct {
		script   string
		expected int64
	}{
		{"Object.A = 9; Object.A *= 3; Object.A", 27},
		{"Object.A = 10; Object.A += 5; Object.A", 15},
		{"Object.A = 10; Object.A -= 3; Object.A", 7},
		{"Object.A = 20; Object.A /= 4; Object.A", 5},
	}

	for _, tc := range testCases {
		result, err := Eval(context.Background(),
			tc.script,
			WithGlobal("Object", &Object{}))

		require.Nil(t, err, "script: %s", tc.script)
		require.Equal(t, object.NewInt(tc.expected), result, "script: %s", tc.script)
	}
}

func TestWithExistingVM(t *testing.T) {
	ctx := context.Background()

	vm, err := vm.NewEmpty()
	require.Nil(t, err)

	result, err := Eval(ctx,
		"foo",
		WithVM(vm),
		WithGlobals(map[string]any{"foo": object.NewInt(3)}),
	)
	require.Nil(t, err)
	require.Equal(t, int64(3), result.Interface())

	result, err = Eval(ctx,
		"bar",
		WithVM(vm),
		WithGlobals(map[string]any{"bar": object.NewInt(4)}),
	)
	require.Nil(t, err)
	require.Equal(t, int64(4), result.Interface())
}

func TestDefaultGlobalsFunc(t *testing.T) {
	globals := DefaultGlobals(DefaultGlobalsOpts{
		ListenersAllowed: true,
	})
	expectedNames := []string{ // non-exhaustive
		"base64",
		"bytes",
		"errors",
		"exec",
		"filepath",
		"fmt",
		"http",
		"cat",
		"ls",
		"print",
	}
	for _, name := range expectedNames {
		_, ok := globals[name]
		require.True(t, ok, "expected global %s", name)
	}
}
