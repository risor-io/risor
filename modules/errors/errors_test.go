package errors

import (
	"context"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	e := New(context.Background(),
		object.NewString("error %q %d"),
		object.NewString("foo bar"),
		object.NewInt(42),
	)
	require.IsType(t, &object.Error{}, e)
	errObj, ok := e.(*object.Error)
	require.True(t, ok)
	require.False(t, errObj.Raised())
	require.Equal(t, "error \"foo bar\" 42", errObj.Value().Error())
}

func TestBadErrorsCall(t *testing.T) {
	e := New(context.Background())
	require.IsType(t, &object.Error{}, e)
	errObj, ok := e.(*object.Error)
	require.True(t, ok)
	require.True(t, errObj.Raised())
	require.Equal(t, "type error: errors.new() takes 1 or more arguments (0 given)", errObj.Value().Error())
}
