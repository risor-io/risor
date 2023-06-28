package object_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/cloudcmds/tamarin/v2/object"
	"github.com/stretchr/testify/require"
)

// Used to confirm we can proxy method calls that use complex types.
type ProxyTestOpts struct {
	A int
	B string
	C bool `json:"c"`
}

// We use this struct embedded in ProxyService to prove that methods provided by
// embedded structs are also proxied.
type ProxyServiceEmbedded struct{}

func (e ProxyServiceEmbedded) Flub(opts ProxyTestOpts) string {
	return fmt.Sprintf("flubbed:%d.%s.%v", opts.A, opts.B, opts.C)
}

func (e ProxyServiceEmbedded) Increment(ctx context.Context, i int64) int64 {
	return i + 1
}

// This represents a "service" provided by Go code that we want to call from
// Tamarin code using a proxy.
type ProxyService struct {
	ProxyServiceEmbedded
}

func (pt *ProxyService) ToUpper(s string) string {
	return strings.ToUpper(s)
}

func (pt *ProxyService) ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

type proxyTestType1 []string

func (p proxyTestType1) Len() int {
	return len(p)
}

func TestProxyNonStruct(t *testing.T) {
	proxy, err := object.NewProxy(proxyTestType1{"a", "b", "c"})
	require.Nil(t, err)
	fmt.Println(proxy)

	goType := proxy.GoType()
	fmt.Println("goType:", goType)

	require.Equal(t, []string{"Len"}, goType.AttributeNames())
	attr, ok := goType.GetAttribute("Len")
	require.True(t, ok)
	require.Equal(t, "Len", attr.Name())

	method, ok := attr.(*object.GoMethod)
	require.True(t, ok)
	require.Equal(t, "Len", method.Name())
	require.Equal(t, 1, method.NumIn())
	require.Equal(t, 1, method.NumOut())

	m, ok := proxy.GetAttr("Len")
	require.True(t, ok)
	lenBuiltin, ok := m.(*object.Builtin)
	require.True(t, ok)
	res := lenBuiltin.Call(context.Background())
	require.Equal(t, int64(3), res.(*object.Int).Value())
}

type proxyTestType2 struct {
	A    int
	B    map[string]int
	c    string
	Anon struct {
		X int
	}
	Nested proxyTestType1
}

func (p proxyTestType2) D(x int, y float32) (int, error) {
	return x + int(y), nil
}

func (p proxyTestType2) e() int {
	return 43
}

func TestProxyTestType2(t *testing.T) {
	proxy, err := object.NewProxy(&proxyTestType2{
		A: 99,
		B: map[string]int{
			"foo": 123,
			"bar": 456,
		},
		c:    "hello",
		Anon: struct{ X int }{99},
		Nested: proxyTestType1{
			"baz",
			"qux",
		},
	})
	require.Nil(t, err)
	fmt.Println(proxy)

	goType := proxy.GoType()
	require.Equal(t, "*object_test.proxyTestType2", goType.Name())
	fmt.Println("goType:", goType)

	require.Equal(t, []string{"A", "Anon", "B", "D", "Nested"},
		goType.AttributeNames())

	aAttr, ok := goType.GetAttribute("A")
	require.True(t, ok)
	require.Equal(t, "A", aAttr.Name())
	field, ok := aAttr.(*object.GoField)
	require.True(t, ok)
	require.Equal(t, "A", field.Name())
	require.Equal(t, "int", field.ReflectType().Name())

	anonAttr, ok := goType.GetAttribute("Anon")
	require.True(t, ok)
	require.Equal(t, "Anon", anonAttr.Name())
	field, ok = anonAttr.(*object.GoField)
	require.True(t, ok)
	require.Equal(t, "Anon", field.Name())
	require.Equal(t, "", field.ReflectType().Name())
	require.Equal(t, []string{"X"}, field.GoType().AttributeNames())

	attr, ok := goType.GetAttribute("D")
	require.True(t, ok)
	require.Equal(t, "D", attr.Name())

	method, ok := attr.(*object.GoMethod)
	require.True(t, ok)
	require.Equal(t, "D", method.Name())
	require.Equal(t, 3, method.NumIn())
	require.Equal(t, 2, method.NumOut())

	in0 := method.InType(0)
	require.Equal(t, "*object_test.proxyTestType2", in0.Name())
	in1 := method.InType(1)
	require.Equal(t, "int", in1.Name())
	in2 := method.InType(2)
	require.Equal(t, "float32", in2.Name())

	out0 := method.OutType(0)
	require.Equal(t, "int", out0.Name())
	out1 := method.OutType(1)
	require.Equal(t, "error", out1.Name())

	require.True(t, method.ProducesError())
	require.Equal(t, []int{1}, method.ErrorIndices())

	nestedAttr, ok := goType.GetAttribute("Nested")
	require.True(t, ok)
	require.Equal(t, "Nested", nestedAttr.Name())
	field, ok = nestedAttr.(*object.GoField)
	require.True(t, ok)
	require.Equal(t, "Nested", field.Name())
	require.Equal(t, "proxyTestType1", field.ReflectType().Name())
	require.Equal(t, []string{"Len"}, field.GoType().AttributeNames())

	ptt1, err := object.NewGoType(reflect.TypeOf(proxyTestType1{}))
	require.Nil(t, err)
	require.Equal(t, ptt1, field.GoType())

	aValue, getOk := proxy.GetAttr("A")
	require.True(t, getOk)
	require.Equal(t, object.NewInt(99), aValue)
}

func TestProxyCall(t *testing.T) {
	proxy, err := object.NewProxy(&proxyTestType2{})
	require.Nil(t, err)

	m, ok := proxy.GetAttr("D")
	require.True(t, ok)

	b, ok := m.(*object.Builtin)
	require.True(t, ok)

	result := b.Call(context.Background(),
		object.NewInt(1),
		object.NewFloat(2.0))

	require.Equal(t, object.NewInt(3), result)
}

func TestProxySetGetAttr(t *testing.T) {

	proxy, err := object.NewProxy(&proxyTestType2{})
	require.Nil(t, err)

	// A starts at 0
	value, ok := proxy.GetAttr("A")
	require.True(t, ok)
	require.Equal(t, object.NewInt(0), value)

	// Set to 42
	require.Nil(t, proxy.SetAttr("A", object.NewInt(42)))

	// Confirm 42
	value, ok = proxy.GetAttr("A")
	require.True(t, ok)
	require.Equal(t, object.NewInt(42), value)

	// Set to -3
	require.Nil(t, proxy.SetAttr("A", object.NewInt(-3)))

	// Confirm -3
	value, ok = proxy.GetAttr("A")
	require.True(t, ok)
	require.Equal(t, object.NewInt(-3), value)

}

func TestAttemptProxyOnStructValue(t *testing.T) {
	// Cannot create a proxy on a struct value. It has to be a pointer.
	_, err := object.NewProxy(proxyTestType2{})
	require.NotNil(t, err)
	require.Equal(t, "type error: unable to proxy type (object_test.proxyTestType2 given)", err.Error())
}

func TestProxyBytesBuffer(t *testing.T) {

	ctx := context.Background()
	buf := bytes.NewBuffer([]byte("abc"))
	var reader io.Reader = buf

	// Creating a proxy on an interface really means creating a proxy on the
	// underlying concrete type.
	proxy, err := object.NewProxy(reader)
	require.Nil(t, err)

	// Confirm the GoType is actually *bytes.Buffer
	goType := proxy.GoType()
	require.Equal(t, "*bytes.Buffer", goType.Name())

	// The proxy should have attributes available for all public attributes
	// on *bytes.Buffer
	method, ok := proxy.GetAttr("Len")
	require.True(t, ok)

	// Confirm we can call a method
	lenMethod, ok := method.(*object.Builtin)
	require.True(t, ok)
	require.Equal(t, object.NewInt(3), lenMethod.Call(ctx))

	// Write to the buffer and confirm the length changes
	buf.WriteString("defg")
	require.Equal(t, object.NewInt(7), lenMethod.Call(ctx))

	// Confirm we can call Bytes() and get a byte_slice back
	getBytes, ok := proxy.GetAttr("Bytes")
	require.True(t, ok)
	bytes := getBytes.(*object.Builtin).Call(ctx)
	require.Equal(t, object.NewByteSlice([]byte("abcdefg")), bytes)
}

func TestProxyMethodError(t *testing.T) {

	// Using the ReadByte method as an example, call it in a situation that will
	// have it return an error, then confirm a Tamarin *Error is returned.

	// func (b *Buffer) ReadByte() (byte, error)
	// If no byte is available, it returns error io.EOF.

	ctx := context.Background()
	buf := bytes.NewBuffer(nil) // empty buffer!
	proxy, err := object.NewProxy(buf)
	require.Nil(t, err)

	method, ok := proxy.GetAttr("ReadByte")
	require.True(t, ok)

	readByte, ok := method.(*object.Builtin)
	require.True(t, ok)

	result := readByte.Call(ctx)
	errObj, ok := result.(*object.Error)
	require.True(t, ok)
	require.Equal(t, object.Errorf("EOF"), errObj)
}

func TestProxyHasher(t *testing.T) {
	ctx := context.Background()
	h := sha256.New()

	proxy, err := object.NewProxy(h)
	require.Nil(t, err)

	method, ok := proxy.GetAttr("Write")
	require.True(t, ok)
	write, ok := method.(*object.Builtin)
	require.True(t, ok)

	method, ok = proxy.GetAttr("Sum")
	require.True(t, ok)
	sum, ok := method.(*object.Builtin)
	require.True(t, ok)

	result := write.Call(ctx, object.NewByteSlice([]byte("abc")))
	require.Equal(t, object.NewInt(3), result)

	result = write.Call(ctx, object.NewByteSlice([]byte("de")))
	require.Equal(t, object.NewInt(2), result)

	result = sum.Call(ctx, object.NewByteSlice(nil))
	byte_slice, ok := result.(*object.ByteSlice)
	require.True(t, ok)

	other := sha256.New()
	other.Write([]byte("abcde"))
	expected := other.Sum(nil)

	require.Equal(t, expected, byte_slice.Value())
}
