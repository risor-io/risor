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

	"github.com/risor-io/risor/object"
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
// Risor code using a proxy.
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

type proxyTestType3 struct {
	A int
	P *string
	I io.Reader
	M map[string]int
	S []string
}

func TestProxySetGetAttrNil(t *testing.T) {
	proxy, err := object.NewProxy(&proxyTestType3{})
	require.Nil(t, err)

	// A is not nillable
	err = proxy.SetAttr("A", object.Nil)
	require.Error(t, err)
	require.Equal(t, "type error: expected int (nil given)", err.Error())

	// P starts at nil
	value, ok := proxy.GetAttr("P")
	require.True(t, ok)
	require.Equal(t, object.Nil, value)

	// Set to "abc"
	require.Nil(t, proxy.SetAttr("P", object.NewString("abc")))

	// Confirm "abc"
	value, ok = proxy.GetAttr("P")
	require.True(t, ok)
	require.Equal(t, object.NewString("abc"), value)

	// Set to nil
	require.Nil(t, proxy.SetAttr("P", object.Nil))

	// Confirm nil
	value, ok = proxy.GetAttr("P")
	require.True(t, ok)
	require.Equal(t, object.Nil, value)

	// I starts at nil
	value, ok = proxy.GetAttr("I")
	require.True(t, ok)
	require.Equal(t, object.Nil, value)

	// Set to "abc"
	require.Nil(t, proxy.SetAttr("I", object.NewBuffer(bytes.NewBufferString("abc"))))

	// Confirm "abc"
	value, ok = proxy.GetAttr("I")
	require.True(t, ok)
	require.Equal(t, object.NewBuffer(bytes.NewBufferString("abc")), value)

	// Set to nil
	require.Nil(t, proxy.SetAttr("I", object.Nil))

	// Confirm nil
	value, ok = proxy.GetAttr("I")
	require.True(t, ok)
	require.Equal(t, object.Nil, value)

	// M starts at nil
	value, ok = proxy.GetAttr("M")
	require.True(t, ok)
	require.Equal(t, object.NewMap(map[string]object.Object{}), value)

	// Set to {"a": 1, "b": 2, "c": 3}
	require.Nil(t, proxy.SetAttr("M", object.NewMap(map[string]object.Object{
		"a": object.NewInt(1),
		"b": object.NewInt(2),
		"c": object.NewInt(3),
	})))

	// Confirm {"a": 1, "b": 2, "c": 3}
	value, ok = proxy.GetAttr("M")
	require.True(t, ok)
	require.Equal(t, object.NewMap(map[string]object.Object{
		"a": object.NewInt(1),
		"b": object.NewInt(2),
		"c": object.NewInt(3),
	}), value)

	// Set to nil
	require.Nil(t, proxy.SetAttr("M", object.Nil))

	// Confirm nil
	value, ok = proxy.GetAttr("M")
	require.True(t, ok)
	require.Equal(t, object.NewMap(map[string]object.Object{}), value)

	// S starts at nil
	value, ok = proxy.GetAttr("S")
	require.True(t, ok)
	require.Equal(t, object.NewList([]object.Object{}), value)

	// Set to ["a", "b", "c"]
	require.Nil(t, proxy.SetAttr("S", object.NewStringList([]string{"a", "b", "c"})))

	// Confirm ["a", "b", "c"]
	value, ok = proxy.GetAttr("S")
	require.True(t, ok)
	require.Equal(t, object.NewStringList([]string{"a", "b", "c"}), value)

	// Set to nil
	require.Nil(t, proxy.SetAttr("S", object.Nil))

	// Confirm nil
	value, ok = proxy.GetAttr("S")
	require.True(t, ok)
	require.Equal(t, object.NewList([]object.Object{}), value)
}

func TestProxyOnStructValue(t *testing.T) {
	p, err := object.NewProxy(proxyTestType2{A: 99})
	require.NoError(t, err)
	require.Equal(t, "*object_test.proxyTestType2", p.GoType().Name())
	attr, ok := p.GetAttr("A")
	require.True(t, ok)
	require.Equal(t, object.NewInt(99), attr)
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
	// have it return an error, then confirm a Risor *Error is returned.

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

type nestedStructA struct {
	B string
}

type nestedStructConfig struct {
	A nestedStructA
}

func TestProxyNestedStruct(t *testing.T) {
	config := &nestedStructConfig{}
	proxy, err := object.NewProxy(config)
	require.Nil(t, err)

	// Get the A field
	aField, ok := proxy.GetAttr("A")
	require.True(t, ok)

	// Verify A is a proxy to nestedStructA
	aProxy, ok := aField.(*object.Proxy)
	require.True(t, ok)
	require.Equal(t, "*object_test.nestedStructA", aProxy.GoType().Name())

	// Set B field directly on the A proxy
	err = aProxy.SetAttr("B", object.NewString("hello"))
	require.Nil(t, err)

	// Verify the value was set correctly
	require.Equal(t, "hello", config.A.B)
}

type testNilArg struct{}

func (t *testNilArg) Test(arg any) {
	// Method implementation doesn't matter for this test
}

func (t *testNilArg) TestMultiple(a, b any) {
	// Method implementation doesn't matter for this test
}

func (t *testNilArg) TestMixed(a string, b any) {
	// Method implementation doesn't matter for this test
}

func (t *testNilArg) TestReturnNil() any {
	return nil
}

func (t *testNilArg) TestReturnValue() string {
	return "hello"
}

func (t *testNilArg) TestReturnMultiple() (string, any) {
	return "hello", nil
}

func (t *testNilArg) TestReturnPointer() *string {
	s := "hello"
	return &s
}

func TestProxyNilArg(t *testing.T) {
	proxy, err := object.NewProxy(&testNilArg{})
	require.Nil(t, err)

	method, ok := proxy.GetAttr("Test")
	require.True(t, ok)

	b, ok := method.(*object.Builtin)
	require.True(t, ok)

	// Call Test with nil argument
	result := b.Call(context.Background(), object.Nil)
	require.Equal(t, object.Nil, result)
}

func TestProxyMultipleNilArgs(t *testing.T) {
	proxy, err := object.NewProxy(&testNilArg{})
	require.Nil(t, err)

	method, ok := proxy.GetAttr("TestMultiple")
	require.True(t, ok)

	b, ok := method.(*object.Builtin)
	require.True(t, ok)

	// Call TestMultiple with two nil arguments
	result := b.Call(context.Background(), object.Nil, object.Nil)
	require.Equal(t, object.Nil, result)
}

func TestProxyMixedArgs(t *testing.T) {
	proxy, err := object.NewProxy(&testNilArg{})
	require.Nil(t, err)

	method, ok := proxy.GetAttr("TestMixed")
	require.True(t, ok)

	b, ok := method.(*object.Builtin)
	require.True(t, ok)

	// Call TestMixed with a string and nil
	result := b.Call(context.Background(), object.NewString("hello"), object.Nil)
	require.Equal(t, object.Nil, result)
}

func TestProxyReturnNil(t *testing.T) {
	proxy, err := object.NewProxy(&testNilArg{})
	require.Nil(t, err)

	method, ok := proxy.GetAttr("TestReturnNil")
	require.True(t, ok)

	b, ok := method.(*object.Builtin)
	require.True(t, ok)

	// Call TestReturnNil and verify it returns nil
	result := b.Call(context.Background())
	require.Equal(t, object.Nil, result)
}

func TestProxyReturnValue(t *testing.T) {
	proxy, err := object.NewProxy(&testNilArg{})
	require.Nil(t, err)

	method, ok := proxy.GetAttr("TestReturnValue")
	require.True(t, ok)

	b, ok := method.(*object.Builtin)
	require.True(t, ok)

	// Call TestReturnValue and verify it returns the string
	result := b.Call(context.Background())
	str, ok := result.(*object.String)
	require.True(t, ok)
	require.Equal(t, "hello", str.Value())
}

func TestProxyReturnMultiple(t *testing.T) {
	proxy, err := object.NewProxy(&testNilArg{})
	require.Nil(t, err)

	method, ok := proxy.GetAttr("TestReturnMultiple")
	require.True(t, ok)

	b, ok := method.(*object.Builtin)
	require.True(t, ok)

	// Call TestReturnMultiple and verify it returns both values
	result := b.Call(context.Background())
	list, ok := result.(*object.List)
	require.True(t, ok)
	require.Equal(t, int64(2), list.Len().Value())

	// Check first value (string)
	str, err := list.GetItem(object.NewInt(0))
	require.Nil(t, err)
	require.Equal(t, "hello", str.(*object.String).Value())

	// Check second value (nil)
	str, err = list.GetItem(object.NewInt(1))
	require.Nil(t, err)
	require.Equal(t, object.Nil, str)
}

func TestProxyReturnPointer(t *testing.T) {
	proxy, err := object.NewProxy(&testNilArg{})
	require.Nil(t, err)

	method, ok := proxy.GetAttr("TestReturnPointer")
	require.True(t, ok)

	b, ok := method.(*object.Builtin)
	require.True(t, ok)

	// Call TestReturnPointer and verify it returns the pointer
	result := b.Call(context.Background())
	str, ok := result.(*object.String)
	require.True(t, ok)
	require.Equal(t, "hello", str.Value())
}

// Vector3D is a simple 3D vector type for testing struct field setting
type Vector3D struct {
	X, Y, Z float64
}

func (v Vector3D) Add(other Vector3D) Vector3D {
	return Vector3D{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
}

func (v Vector3D) String() string {
	return fmt.Sprintf("Vector3D{X:%v, Y:%v, Z:%v}", v.X, v.Y, v.Z)
}

type VectorFactory struct{}

func (f *VectorFactory) NewVector(x, y, z float64) Vector3D {
	return Vector3D{X: x, Y: y, Z: z}
}

// TestProxyMethodReturnStructValue tests that struct values returned from methods
// can have their fields modified
func TestProxyMethodReturnStructValue(t *testing.T) {
	// Create a proxy for the vector factory
	factory, err := object.NewProxy(&VectorFactory{})
	require.Nil(t, err)

	// Get the NewVector method
	newVectorMethod, ok := factory.GetAttr("NewVector")
	require.True(t, ok)
	newVector, ok := newVectorMethod.(*object.Builtin)
	require.True(t, ok)

	// Create a vector using the factory method
	ctx := context.Background()
	vector1 := newVector.Call(ctx, object.NewFloat(1), object.NewFloat(2), object.NewFloat(3))

	// Verify it's a proxy
	vectorProxy1, ok := vector1.(*object.Proxy)
	require.True(t, ok)

	// Verify the proxy is to a *Vector3D, not a Vector3D
	require.Equal(t, "*object_test.Vector3D", vectorProxy1.GoType().Name())

	// Get the Add method from the vector
	addMethod, ok := vectorProxy1.GetAttr("Add")
	require.True(t, ok)
	add, ok := addMethod.(*object.Builtin)
	require.True(t, ok)

	// Create another vector and add them
	vector2 := newVector.Call(ctx, object.NewFloat(4), object.NewFloat(5), object.NewFloat(6))
	result := add.Call(ctx, vector2)

	// Verify result is a proxy
	resultProxy, ok := result.(*object.Proxy)
	require.True(t, ok)

	// Verify the result proxy is to a *Vector3D, not a Vector3D. Struct values
	// should be converted to pointers automatically
	require.Equal(t, "*object_test.Vector3D", resultProxy.GoType().Name())

	// Now test that we can modify fields on the result
	err = resultProxy.SetAttr("X", object.NewFloat(15))
	require.Nil(t, err, "Should be able to set field X on the result")

	// Verify the field was updated
	x, ok := resultProxy.GetAttr("X")
	require.True(t, ok)
	require.Equal(t, object.NewFloat(15), x)

	// Test other fields too
	err = resultProxy.SetAttr("Y", object.NewFloat(25))
	require.Nil(t, err)
	y, ok := resultProxy.GetAttr("Y")
	require.True(t, ok)
	require.Equal(t, object.NewFloat(25), y)

	err = resultProxy.SetAttr("Z", object.NewFloat(35))
	require.Nil(t, err)
	z, ok := resultProxy.GetAttr("Z")
	require.True(t, ok)
	require.Equal(t, object.NewFloat(35), z)
}

// TestProxyStructConverterRoundTrip tests that struct values are properly converted
// when going from Go to Risor and back to Go
func TestProxyStructConverterRoundTrip(t *testing.T) {
	// Create a Vector3D struct
	original := Vector3D{X: 1, Y: 2, Z: 3}

	// Create a proxy for the struct
	proxy, err := object.NewProxy(original)
	require.Nil(t, err)

	// Verify it's a pointer type in the proxy
	require.Equal(t, "*object_test.Vector3D", proxy.GoType().Name())

	// Modify a field
	err = proxy.SetAttr("X", object.NewFloat(10))
	require.Nil(t, err)

	// Convert back to Go type
	// This should extract the value, not the pointer
	converter, err := object.NewTypeConverter(reflect.TypeOf(original))
	require.Nil(t, err)

	result, err := converter.To(proxy)
	require.Nil(t, err)

	// Verify the result is a Vector3D, not a *Vector3D
	resultVector, ok := result.(Vector3D)
	require.True(t, ok)

	// Verify the field was modified
	require.Equal(t, 10.0, resultVector.X)
	require.Equal(t, 2.0, resultVector.Y)
	require.Equal(t, 3.0, resultVector.Z)
}

// TestProxyVectorModificationTracking tests that when a struct value is returned
// from a method call, modifications to the struct value are reflected when later
// accessed from Go code
func TestProxyVectorModificationTracking(t *testing.T) {
	// Create a vector factory
	factory := &VectorFactory{}

	// Create a proxy for the factory
	factoryProxy, err := object.NewProxy(factory)
	require.Nil(t, err)

	// Get the NewVector method
	newVectorMethod, ok := factoryProxy.GetAttr("NewVector")
	require.True(t, ok)
	newVector, ok := newVectorMethod.(*object.Builtin)
	require.True(t, ok)

	// Call the NewVector method to create a vector
	ctx := context.Background()
	vector := newVector.Call(ctx, object.NewFloat(1), object.NewFloat(2), object.NewFloat(3))

	// Get the Add method
	vectorProxy, ok := vector.(*object.Proxy)
	require.True(t, ok)
	addMethod, ok := vectorProxy.GetAttr("Add")
	require.True(t, ok)
	add, ok := addMethod.(*object.Builtin)
	require.True(t, ok)

	// Call Add to create a result vector
	otherVector, err := object.NewProxy(Vector3D{X: 4, Y: 5, Z: 6})
	require.Nil(t, err)
	resultVector := add.Call(ctx, otherVector)
	resultProxy, ok := resultVector.(*object.Proxy)
	require.True(t, ok)

	// Extract the actual Go value from the proxy
	// We should be able to get the underlying Go object
	underlyingObj := resultProxy.Interface()

	// Modify the vector through the proxy
	err = resultProxy.SetAttr("X", object.NewFloat(99))
	require.Nil(t, err)
	err = resultProxy.SetAttr("Y", object.NewFloat(88))
	require.Nil(t, err)
	err = resultProxy.SetAttr("Z", object.NewFloat(77))
	require.Nil(t, err)

	// Check that the modifications are reflected in the underlying Go value
	// This is important - the changes should be visible to Go code
	underlyingVector, ok := underlyingObj.(*Vector3D)
	require.True(t, ok)
	require.Equal(t, 99.0, underlyingVector.X)
	require.Equal(t, 88.0, underlyingVector.Y)
	require.Equal(t, 77.0, underlyingVector.Z)
}
