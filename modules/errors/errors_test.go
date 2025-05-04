package errors

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/object"
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

func TestErrorsIs(t *testing.T) {
	var (
		ErrNotExist = errors.New("file does not exist")
		ErrExist    = errors.New("file already exists")
	)
	err1 := object.NewError(ErrNotExist)
	err2 := object.NewError(ErrExist)
	err3 := New(context.Background(), object.NewString("test"))
	require.Equal(t, object.True, Is(context.Background(), err1, object.NewError(ErrNotExist)))
	require.Equal(t, object.True, Is(context.Background(), err2, object.NewError(ErrExist)))
	require.Equal(t, object.False, Is(context.Background(), err2, object.NewError(ErrNotExist)))
	require.Equal(t, object.True, Is(context.Background(), err3, err3))
}

func TestErrorsIsWithWrappedErrors(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := fmt.Errorf("wrapped: %w", baseErr)

	base := object.NewError(baseErr)
	wrapped := object.NewError(wrappedErr)

	// A wrapped error should match its base error
	require.Equal(t, object.True, Is(context.Background(), wrapped, base))

	// But base error should not match the wrapped error
	require.Equal(t, object.False, Is(context.Background(), base, wrapped))
}

func TestErrorsIsInvalidArgs(t *testing.T) {
	ctx := context.Background()
	err := object.NewError(errors.New("test error"))

	// First argument is not an error
	result := Is(ctx, object.NewString("not an error"), err)
	require.IsType(t, &object.Error{}, result)

	// Second argument is not an error
	result = Is(ctx, err, object.NewString("not an error"))
	require.IsType(t, &object.Error{}, result)

	// Wrong number of arguments
	result = Is(ctx, err)
	require.IsType(t, &object.Error{}, result)
	errObj, ok := result.(*object.Error)
	require.True(t, ok)
	require.Contains(t, errObj.Value().Error(), "takes exactly 2 arguments")

	result = Is(ctx, err, err, err)
	require.IsType(t, &object.Error{}, result)
	errObj, ok = result.(*object.Error)
	require.True(t, ok)
	require.Contains(t, errObj.Value().Error(), "takes exactly 2 arguments")
}

func TestErrorsIsSentinelErrors(t *testing.T) {
	// Define custom sentinel errors
	var (
		ErrPermissionDenied = errors.New("permission denied")
		ErrCustom           = errors.New("custom error")
	)

	// Create wrapper errors
	permErr := fmt.Errorf("cannot access file: %w", ErrPermissionDenied)
	deepWrapper := fmt.Errorf("operation failed: %w", permErr)

	// Convert to Risor error objects
	risorPermErr := object.NewError(permErr)
	risorDeepWrapper := object.NewError(deepWrapper)
	sentinelErr := object.NewError(ErrPermissionDenied)
	otherErr := object.NewError(ErrCustom)

	// Test error chain matching
	require.Equal(t, object.True, Is(context.Background(), risorPermErr, sentinelErr))
	require.Equal(t, object.True, Is(context.Background(), risorDeepWrapper, sentinelErr))
	require.Equal(t, object.False, Is(context.Background(), risorPermErr, otherErr))
	require.Equal(t, object.False, Is(context.Background(), risorDeepWrapper, otherErr))
}

func TestErrorsIsWithStdlibSentinels(t *testing.T) {
	ctx := context.Background()

	// Import standard sentinel errors from multiple stdlib packages
	var (
		// os package errors
		errNotExist   = object.NewError(os.ErrNotExist)
		errExist      = object.NewError(os.ErrExist)
		errPermission = object.NewError(os.ErrPermission)

		// io package errors
		errEOF           = object.NewError(io.EOF)
		errUnexpectedEOF = object.NewError(io.ErrUnexpectedEOF)

		// context package errors
		errCanceled         = object.NewError(context.Canceled)
		errDeadlineExceeded = object.NewError(context.DeadlineExceeded)
	)

	// Create wrapped errors with stdlib sentinel errors
	fileErr := fmt.Errorf("could not open file: %w", os.ErrNotExist)
	readErr := fmt.Errorf("failed to read file: %w", io.EOF)
	timeoutErr := fmt.Errorf("operation timed out: %w", context.DeadlineExceeded)

	// Convert wrapped errors to Risor error objects
	risorFileErr := object.NewError(fileErr)
	risorReadErr := object.NewError(readErr)
	risorTimeoutErr := object.NewError(timeoutErr)

	// Test matching with stdlib sentinel errors
	require.Equal(t, object.True, Is(ctx, risorFileErr, errNotExist))
	require.Equal(t, object.False, Is(ctx, risorFileErr, errExist))
	require.Equal(t, object.False, Is(ctx, risorFileErr, errPermission))

	require.Equal(t, object.True, Is(ctx, risorReadErr, errEOF))
	require.Equal(t, object.False, Is(ctx, risorReadErr, errUnexpectedEOF))

	require.Equal(t, object.True, Is(ctx, risorTimeoutErr, errDeadlineExceeded))
	require.Equal(t, object.False, Is(ctx, risorTimeoutErr, errCanceled))

	// Test double-wrapped errors
	deepFileErr := fmt.Errorf("data access error: %w", fileErr)
	risorDeepFileErr := object.NewError(deepFileErr)
	require.Equal(t, object.True, Is(ctx, risorDeepFileErr, errNotExist))

	// Test combining multiple errors (only matches the wrapped one)
	combinedErr := fmt.Errorf("multiple issues: %w and also EOF", os.ErrPermission)
	risorCombinedErr := object.NewError(combinedErr)
	require.Equal(t, object.True, Is(ctx, risorCombinedErr, errPermission))
	require.Equal(t, object.False, Is(ctx, risorCombinedErr, errEOF))
}
