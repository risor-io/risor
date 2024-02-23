package http

import (
	"context"
	"testing"
	"time"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestListenAndServeNoCloneCall(t *testing.T) {
	ctx := context.Background()
	h := object.NewFunction(compiler.NewFunction(compiler.FunctionOpts{}))
	res := ListenAndServe(ctx, object.NewString("localhost:8080"), h)
	require.NotNil(t, res)
	errObj, ok := res.(*object.Error)
	require.True(t, ok)
	require.Equal(t,
		errObj.Value().Error(),
		"http.listen_and_serve: no clone-call function found in context")
}

func TestListenAndServeTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5)
	defer cancel()
	callFunc := func(ctx context.Context, fn *object.Function, args []object.Object) (object.Object, error) {
		return object.Nil, nil
	}
	ctx = object.WithCloneCallFunc(ctx, object.CallFunc(callFunc))
	h := object.NewFunction(compiler.NewFunction(compiler.FunctionOpts{}))
	res := ListenAndServe(ctx, object.NewString("localhost:8080"), h)
	require.NotNil(t, res)
	errObj, ok := res.(*object.Error)
	require.True(t, ok)
	require.Equal(t, "http: Server closed", errObj.Value().Error())
}
