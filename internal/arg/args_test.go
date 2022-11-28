package arg_test

import (
	"testing"

	"github.com/cloudcmds/tamarin/internal/arg"
	"github.com/cloudcmds/tamarin/object"
	"github.com/stretchr/testify/require"
)

func TestRequire(t *testing.T) {
	var err *object.Error

	err = arg.Require(
		"foo",
		1,
		[]object.Object{
			&object.Integer{Value: 1},
			&object.Integer{Value: 1},
			&object.Integer{Value: 1},
		},
	)
	require.NotNil(t, err)
	require.Equal(t, "type error: foo() takes exactly 1 argument (3 given)",
		err.Message)

	err = arg.Require(
		"bar",
		2,
		[]object.Object{
			&object.Integer{Value: 1},
		},
	)
	require.NotNil(t, err)
	require.Equal(t, "type error: bar() takes exactly 2 arguments (1 given)",
		err.Message)
}
