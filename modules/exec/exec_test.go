package exec

import (
	"context"
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
