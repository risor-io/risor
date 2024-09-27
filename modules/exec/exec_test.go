package exec

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func getStdoutLines(result *Result) []string {
	stdout := result.Stdout().(*object.String)
	return strings.Split(stdout.Value(), "\n")
}

func TestExecOldWay(t *testing.T) {
	ctx := context.Background()
	resultObj := Exec(ctx,
		object.NewString("ls"),
		object.NewList([]object.Object{object.NewString("-l")}))
	result, ok := resultObj.(*Result)
	require.True(t, ok)
	lines := getStdoutLines(result)
	require.True(t, len(lines) > 1)
	var found bool
	for _, line := range lines {
		if strings.HasSuffix(line, "exec_test.go") {
			found = true
			break
		}
	}
	require.True(t, found)
}

func TestExecNewWay(t *testing.T) {
	ctx := context.Background()
	resultObj := Exec(ctx,
		object.NewList([]object.Object{
			object.NewString("ls"),
			object.NewString("-l"),
		}))
	result, ok := resultObj.(*Result)
	require.True(t, ok)
	lines := getStdoutLines(result)
	require.True(t, len(lines) > 1)
	var found bool
	for _, line := range lines {
		if strings.HasSuffix(line, "exec_test.go") {
			found = true
			break
		}
	}
	require.True(t, found)
}

func TestConfigureCommand(t *testing.T) {
	var stderr, stdout bytes.Buffer
	cmd := &exec.Cmd{}
	err := configureCommand(cmd,
		object.NewMap(map[string]object.Object{
			"dir":    object.NewString("/tmp"),
			"stdin":  object.NewString("hello"),
			"stderr": object.NewBuffer(&stderr),
			"stdout": object.NewBuffer(&stdout),
			"env": object.NewMap(map[string]object.Object{
				"FOO": object.NewString("bar"),
			}),
		}))
	require.NoError(t, err)

	require.Equal(t, "/tmp", cmd.Dir)
	require.NotNil(t, cmd.Stdin)
	require.NotNil(t, cmd.Stdout)
	require.NotNil(t, cmd.Stderr)
	require.Equal(t, []string{"FOO=bar"}, cmd.Env)

	data, err := io.ReadAll(cmd.Stdin)
	require.NoError(t, err)
	require.Equal(t, "hello", string(data))
}

func TestCommandFunc(t *testing.T) {
	ctx := context.Background()
	cmdObj, ok := CommandFunc(ctx,
		object.NewString("ls"),
		object.NewString("-l")).(*Command)
	require.True(t, ok)
	cmd := cmdObj.Value()
	_, end := path.Split(cmd.Path)
	require.Equal(t, "ls", end)
	require.Equal(t, []string{"ls", "-l"}, cmd.Args)
}

func TestLookPath(t *testing.T) {
	result := LookPath(context.Background(), object.NewString("ls"))
	path, ok := result.(*object.String)
	require.True(t, ok)
	require.NotEmpty(t, path.Value())
}

func TestConfigureWithBadMaps(t *testing.T) {
	cases := []struct {
		name     string
		params   *object.Map
		expected string
	}{
		{
			name: "stdin",
			params: object.NewMap(map[string]object.Object{
				"stdin": object.NewList([]object.Object{}),
			}),
			expected: "exec expected io.Reader for stdin (got list)",
		},
		{
			name: "stdout",
			params: object.NewMap(map[string]object.Object{
				"stdout": object.NewList([]object.Object{}),
			}),

			expected: "exec expected io.Writer for stdout (got list)",
		},
		{
			name: "stderr",
			params: object.NewMap(map[string]object.Object{
				"stderr": object.NewList([]object.Object{}),
			}),
			expected: "exec expected io.Writer for stderr (got list)",
		},
		{
			name: "dir",
			params: object.NewMap(map[string]object.Object{
				"dir": object.NewList([]object.Object{}),
			}),
			expected: "exec expected string for dir (got list)",
		},
		{
			name: "env",
			params: object.NewMap(map[string]object.Object{
				"env": object.NewList([]object.Object{}),
			}),
			expected: "exec expected map for env (got list)",
		},
		{
			name: "env-value",
			params: object.NewMap(map[string]object.Object{
				"env": object.NewMap(map[string]object.Object{
					"FOO": object.NewList([]object.Object{}),
				}),
			}),
			expected: "exec expected string for env value (got list)",
		},
		{
			name: "unknown-key",
			params: object.NewMap(map[string]object.Object{
				"oops": object.NewInt(0),
			}),
			expected: "exec found unexpected key \"oops\"",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &exec.Cmd{}
			err := configureCommand(cmd, tt.params)
			require.Error(t, err)
			require.Equal(t, tt.expected, err.Error())
		})
	}
}
