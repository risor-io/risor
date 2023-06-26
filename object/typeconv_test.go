package object

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat64Converter(t *testing.T) {
	c := Float64Converter{}

	f, err := c.From(2.0)
	require.Nil(t, err)
	require.Equal(t, NewFloat(2.0), f)

	v, err := c.To(NewFloat(3.0))
	require.Nil(t, err)
	require.Equal(t, 3.0, v)
}

func TestMapStringConverter(t *testing.T) {
	c, err := NewMapConverter(reflect.TypeOf(""))
	require.Nil(t, err)

	m := map[string]string{
		"a": "apple",
		"b": "banana",
	}

	tM, err := c.From(m)
	require.Nil(t, err)
	require.Equal(t, NewMap(map[string]Object{
		"a": NewString("apple"),
		"b": NewString("banana"),
	}), tM)

	gM, err := c.To(NewMap(map[string]Object{
		"c": NewString("cod"),
		"d": NewString("deer"),
	}))
	require.Nil(t, err)
	require.Equal(t, map[string]string{
		"c": "cod",
		"d": "deer",
	}, gM)
}

func TestMapStringInterfaceConverter(t *testing.T) {
	c, err := NewMapConverter(reflect.TypeOf(""))
	require.Nil(t, err)

	m := map[string]string{
		"a": "apple",
		"b": "banana",
	}

	tM, err := c.From(m)
	require.Nil(t, err)
	require.Equal(t, NewMap(map[string]Object{
		"a": NewString("apple"),
		"b": NewString("banana"),
	}), tM)

	gM, err := c.To(NewMap(map[string]Object{
		"c": NewString("cod"),
		"d": NewString("deer"),
	}))
	require.Nil(t, err)
	require.Equal(t, map[string]string{
		"c": "cod",
		"d": "deer",
	}, gM)
}

func TestPointerConverter(t *testing.T) {

	c, err := NewPointerConverter(reflect.TypeOf(float64(0)))
	require.Nil(t, err)

	v := 2.0
	vPtr := &v

	f, err := c.From(vPtr)
	require.Nil(t, err)
	require.Equal(t, NewFloat(2.0), f)

	// Convert one Tamarin Float to a *float64
	outVal1, err := c.To(NewFloat(3.0))
	require.Nil(t, err)
	outValPtr1, ok := outVal1.(*float64)
	require.True(t, ok)
	require.Equal(t, 3.0, *outValPtr1)

	// Convert a second Tamarin Float to a *float64
	outVal2, err := c.To(NewFloat(4.0))
	require.Nil(t, err)
	outValPtr2, ok := outVal2.(*float64)
	require.True(t, ok)
	require.Equal(t, 4.0, *outValPtr2)

	// Confirm the two pointers are different
	require.Equal(t, 3.0, *outValPtr1)
	require.Equal(t, 4.0, *outValPtr2)
}

func TestCreatingPointerViaReflect(t *testing.T) {
	v := 3.0
	var vInterface interface{} = v

	vPointer := reflect.New(reflect.TypeOf(vInterface))
	vPointer.Elem().Set(reflect.ValueOf(v))
	floatPointer := vPointer.Interface()

	result, ok := floatPointer.(*float64)
	require.True(t, ok)
	require.NotNil(t, result)
	require.Equal(t, 3.0, *result)
	require.Equal(t, &v, result)
}

func TestSliceConverter(t *testing.T) {
	c, err := NewSliceConverter(reflect.TypeOf(0.0))
	require.Nil(t, err)

	v := []float64{1.0, 2.0, 3.0}

	f, err := c.From(v)
	require.Nil(t, err)
	require.Equal(t, NewList([]Object{
		NewFloat(1.0),
		NewFloat(2.0),
		NewFloat(3.0),
	}), f)

	list := NewList([]Object{
		NewFloat(9.0),
		NewFloat(-8.0),
	})
	result, err := c.To(list)
	require.Nil(t, err)
	require.Equal(t, []float64{9.0, -8.0}, result)
}

func TestStructConverter(t *testing.T) {

	type foo struct {
		A int
		B string
	}

	c, err := NewStructConverter(reflect.TypeOf(foo{}))
	require.Nil(t, err)

	f := foo{A: 1, B: "two"}

	tF, err := c.From(f)
	require.Nil(t, err)
	require.Equal(t, NewMap(map[string]Object{
		"A": NewInt(1),
		"B": NewString("two"),
	}), tF)

	gF, err := c.To(NewMap(map[string]Object{
		"A": NewInt(3),
		"B": NewString("four"),
	}))
	require.Nil(t, err)

	gFStruct, ok := gF.(foo)
	require.True(t, ok)
	require.Equal(t, foo{A: 3, B: "four"}, gFStruct)
}
