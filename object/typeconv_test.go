package object

import (
	"bytes"
	"reflect"
	"testing"
	"time"

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
	c, err := newMapConverter(reflect.TypeOf(""))
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
	c, err := newMapConverter(reflect.TypeOf(""))
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

	c, err := newPointerConverter(reflect.TypeOf(float64(0)))
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
	c, err := newSliceConverter(reflect.TypeOf(0.0))
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

	c, err := newStructConverter(reflect.TypeOf(foo{}))
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

func TestTimeConverter(t *testing.T) {

	now := time.Now()
	typ := reflect.TypeOf(now)

	c, err := NewTypeConverter(typ)
	require.Nil(t, err)

	tT, err := c.From(now)
	require.Nil(t, err)
	require.Equal(t, NewTime(now), tT)

	gT, err := c.To(NewTime(now))
	require.Nil(t, err)
	goTime, ok := gT.(time.Time)
	require.True(t, ok)
	require.Equal(t, now, goTime)
}

func TestBufferConverter(t *testing.T) {

	buf := bytes.NewBufferString("hello")
	typ := reflect.TypeOf(buf)

	c, err := NewTypeConverter(typ)
	require.Nil(t, err)

	tBuf, err := c.From(buf)
	require.Nil(t, err)
	require.Equal(t, NewBuffer(buf), tBuf)

	gBuf, err := c.To(NewBufferFromBytes([]byte("hello")))
	require.Nil(t, err)
	goBuf, ok := gBuf.(*bytes.Buffer)
	require.True(t, ok)
	require.Equal(t, buf, goBuf)
}

func TestBSliceConverter(t *testing.T) {

	buf := []byte("abc")
	typ := reflect.TypeOf(buf)

	c, err := NewTypeConverter(typ)
	require.Nil(t, err)

	tBuf, err := c.From(buf)
	require.Nil(t, err)
	require.Equal(t, NewBSlice([]byte("abc")), tBuf)

	gBuf, err := c.To(NewBSlice([]byte("abc")))
	require.Nil(t, err)
	goBuf, ok := gBuf.([]byte)
	require.True(t, ok)
	require.Equal(t, buf, goBuf)
}

func TestGenericMapConverter(t *testing.T) {

	m := map[string]interface{}{
		"foo": 1,
		"bar": "two",
	}
	typ := reflect.TypeOf(m)

	c, err := NewTypeConverter(typ)
	require.Nil(t, err)

	tMap, err := c.From(m)
	require.Nil(t, err)
	require.Equal(t, NewMap(map[string]Object{
		"foo": NewInt(1),
		"bar": NewString("two"),
	}), tMap)
}
