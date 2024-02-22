package object

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockCallable struct {
	returnValue Object
}

func (m *mockCallable) Call(ctx context.Context, args ...Object) Object {
	return m.returnValue
}

func TestThread(t *testing.T) {
	c := &mockCallable{
		returnValue: NewInt(42),
	}
	thr := NewThread(context.Background(), c, []Object{})
	require.Equal(t, Type("thread"), thr.Type())
	waitFnObj, ok := thr.GetAttr("wait")
	require.True(t, ok)
	require.NotNil(t, waitFnObj)

	waitFn, ok := waitFnObj.(*Builtin)
	require.True(t, ok)
	require.NotNil(t, waitFn)

	result := waitFn.Call(context.Background())
	require.Equal(t, NewInt(42), result)
}
