package exec_test

import (
	"context"
	"testing"

	"github.com/cloudcmds/tamarin/core/exec"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	ctx := context.Background()
	result, err := exec.Execute(ctx, exec.Opts{
		Input: `len([1,2,3])`,
	})
	require.Nil(t, err)
	require.Equal(t, "3", result.Inspect())
}

func TestExecError(t *testing.T) {
	ctx := context.Background()
	_, err := exec.Execute(ctx, exec.Opts{Input: `bogus()`})
	require.NotNil(t, err)
	require.Equal(t, "name error: \"bogus\" is not defined", err.Error())
}
