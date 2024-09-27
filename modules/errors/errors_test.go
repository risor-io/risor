package errors

import (
	"context"
	"errors"
	"testing"

	"github.com/risor-io/risor/errz"
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
	require.False(t, errObj.IsRaised())
	require.Equal(t, "error \"foo bar\" 42", errObj.Value().Error())
}

func TestEmptyError(t *testing.T) {
	e := New(context.Background())
	require.IsType(t, &object.Error{}, e)
	errObj, ok := e.(*object.Error)
	require.True(t, ok)
	require.False(t, errObj.IsRaised())
	require.Equal(t, "", errObj.Value().Error())
}

func TestErrorTypes(t *testing.T) {
	e1 := EvalError(context.Background(), object.NewString("e1")).(*object.Error)
	e2 := TypeError(context.Background(), object.NewString("e2")).(*object.Error)
	e3 := ArgsError(context.Background(), object.NewString("e3")).(*object.Error)

	var evalErr *errz.EvalError
	var typeErr *errz.TypeError
	var argsErr *errz.ArgsError

	require.True(t, errors.As(e1.Value(), &evalErr))
	require.False(t, errors.As(e1.Value(), &typeErr))
	require.False(t, errors.As(e1.Value(), &argsErr))

	require.False(t, errors.As(e2.Value(), &evalErr))
	require.True(t, errors.As(e2.Value(), &typeErr))
	require.False(t, errors.As(e2.Value(), &argsErr))

	require.False(t, errors.As(e3.Value(), &evalErr))
	require.False(t, errors.As(e3.Value(), &typeErr))
	require.True(t, errors.As(e3.Value(), &argsErr))
}

func TestErrorsAs(t *testing.T) {
	e1 := EvalError(context.Background(), object.NewString("e1")).(*object.Error)
	evalErr := EvalError(context.Background()).(*object.Error)
	typeErr := TypeError(context.Background()).(*object.Error)
	genericErr := New(context.Background()).(*object.Error)

	require.Equal(t, object.True, As(context.Background(), e1, evalErr))
	require.Equal(t, object.False, As(context.Background(), e1, typeErr))
	require.Equal(t, object.True, As(context.Background(), e1, genericErr))
}
