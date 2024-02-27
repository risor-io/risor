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

	method, err := newGoMethod(typ, m)
	require.Nil(t, err)

	require.Equal(t, "Inc", method.Name())
	require.Equal(t, 2, method.NumIn())
	require.Equal(t, 1, method.NumOut())

	in1 := method.InType(1)
	require.Equal(t, "int", in1.Name())

	out0 := method.OutType(0)
	require.Equal(t, "error", out0.Name())
}

type reflectService struct{}

func (svc *reflectService) Test() *reflect.Value {
	return nil
}

func TestGoMethodError(t *testing.T) {
	svc := &reflectService{}
	typ := reflect.TypeOf(svc)

	m, ok := typ.MethodByName("Test")
	require.True(t, ok)

	_, err := newGoMethod(typ, m)
	require.NotNil(t, err)

	expectedErr := `type error: (*object.reflectService).Test has input parameter of type *object.reflectService; 
(*object.reflectService).Test has output parameter of type *reflect.Value; 
(*reflect.Value).CanConvert has input parameter of type reflect.Type; 
(reflect.Type).Field has output parameter of type reflect.StructField; 
unsupported kind: uintptr`

	require.Equal(t, expectedErr, err.Error())
}
