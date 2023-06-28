package object

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestSetAttributeViaReflect(t *testing.T) {
	type test struct {
		A int
	}
	tStruct := test{A: 99}
	var tInterface interface{} = tStruct

	if reflect.TypeOf(tInterface).Kind() != reflect.Ptr {
		// Create a pointer to the value
		tInterfacePointer := reflect.New(reflect.TypeOf(tInterface))
		tInterfacePointer.Elem().Set(reflect.ValueOf(tInterface))
		tInterface = tInterfacePointer.Interface()
	}

	// Set the field "A"
	value := reflect.ValueOf(tInterface)
	value.Elem().FieldByName("A").Set(reflect.ValueOf(100))

	// Confirm the field was set
	require.Equal(t, 100, value.Elem().FieldByName("A").Interface())
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

func TestByteSliceConverter(t *testing.T) {

	buf := []byte("abc")
	typ := reflect.TypeOf(buf)

	c, err := NewTypeConverter(typ)
	require.Nil(t, err)

	tBuf, err := c.From(buf)
	require.Nil(t, err)
	require.Equal(t, NewByteSlice([]byte("abc")), tBuf)

	gBuf, err := c.To(NewByteSlice([]byte("abc")))
	require.Nil(t, err)
	goBuf, ok := gBuf.([]byte)
	require.True(t, ok)
	require.Equal(t, buf, goBuf)
}

func TestArrayConverterInt(t *testing.T) {

	arr := [4]int{2, 3, 4, 5}
	c, err := NewTypeConverter(reflect.TypeOf(arr))
	require.Nil(t, err)

	tList, err := c.From(arr)
	require.Nil(t, err)
	require.Equal(t, NewList([]Object{
		NewInt(2),
		NewInt(3),
		NewInt(4),
		NewInt(5),
	}), tList)

	goValue, err := c.To(NewList([]Object{
		NewInt(-1),
		NewInt(-2),
	}))
	require.Nil(t, err)

	goArray, ok := goValue.([4]int)
	require.True(t, ok)
	require.Equal(t, [4]int{-1, -2}, goArray)
}

func TestArrayConverterFloat64(t *testing.T) {

	arr := [2]float64{100, 101}
	c, err := NewTypeConverter(reflect.TypeOf(arr))
	require.Nil(t, err)

	tList, err := c.From(arr)
	require.Nil(t, err)
	require.Equal(t, NewList([]Object{
		NewFloat(100),
		NewFloat(101),
	}), tList)

	goValue, err := c.To(NewList([]Object{
		NewFloat(-1),
		NewFloat(-2),
	}))
	require.Nil(t, err)

	goArray, ok := goValue.([2]float64)
	require.True(t, ok)
	require.Equal(t, [2]float64{-1, -2}, goArray)
}

func TestGenericMapConverter(t *testing.T) {

	m := map[string]interface{}{
		"foo": 1,
		"bar": "two",
		"baz": []interface{}{
			"three",
			4,
			false,
			map[string]interface{}{
				"five": 5,
			},
		},
	}
	typ := reflect.TypeOf(m)

	c, err := NewTypeConverter(typ)
	require.Nil(t, err)

	tMap, err := c.From(m)
	require.Nil(t, err)
	require.Equal(t, NewMap(map[string]Object{
		"foo": NewInt(1),
		"bar": NewString("two"),
		"baz": NewList([]Object{
			NewString("three"),
			NewInt(4),
			False,
			NewMap(map[string]Object{
				"five": NewInt(5),
			}),
		}),
	}), tMap)
}

func TestGenericMapConverterFromJSON(t *testing.T) {
	m := `{
		"foo": 1,
		"bar": "two",
		"baz": [
			"three",
			4,
			false,
			{ "five": 5 }
		]
	}`
	var v interface{}
	err := json.Unmarshal([]byte(m), &v)
	require.Nil(t, err)

	fmt.Println(v, reflect.TypeOf(v))
	typ := reflect.TypeOf(v)

	c, err := NewTypeConverter(typ)
	require.Nil(t, err)

	tMap, err := c.From(v)
	require.Nil(t, err)
	require.Equal(t, NewMap(map[string]Object{
		"foo": NewFloat(1),
		"bar": NewString("two"),
		"baz": NewList([]Object{
			NewString("three"),
			NewFloat(4),
			False,
			NewMap(map[string]Object{
				"five": NewFloat(5),
			}),
		}),
	}), tMap)
}
