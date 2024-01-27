package object

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type fooStruct struct{}

func (f *fooStruct) Inc(x int) error {
	return nil
}

func TestGoMethod(t *testing.T) {
	f := &fooStruct{}
	typ := reflect.TypeOf(f)

	m, ok := typ.MethodByName("Inc")
	require.True(t, ok)

	method, err := newGoMethod(m)
	require.Nil(t, err)

	require.Equal(t, "Inc", method.Name())
	require.Equal(t, 2, method.NumIn())
	require.Equal(t, 1, method.NumOut())

	in1 := method.InType(1)
	require.Equal(t, "int", in1.Name())

	out0 := method.OutType(0)
	require.Equal(t, "error", out0.Name())
}
