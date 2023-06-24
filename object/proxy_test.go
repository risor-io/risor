package object_test

import (
	"context"
	"fmt"
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

func TestProxy(t *testing.T) {

	// ctx := context.Background()

	// reg, err := object.NewTypeRegistry()
	// require.Nil(t, err)

	// _, err = reg.Register(&ProxyService{})
	// require.Nil(t, err)

	// proxyType, found := reg.GetType(&ProxyService{})
	// require.True(t, found)
	// methods := proxyType.Methods()
	// require.Len(t, methods, 4)

	// sort.Slice(methods, func(i, j int) bool {
	// 	return methods[i].Name() < methods[j].Name()
	// })

	// require.Equal(t, "Flub", methods[0].Name())
	// require.Equal(t, "Increment", methods[1].Name())
	// require.Equal(t, "ParseInt", methods[2].Name())
	// require.Equal(t, "ToUpper", methods[3].Name())

	// // Create a proxy around an instance of ProxyService
	// proxy, err := object.NewProxy(reg, &ProxyService{})
	// require.Nil(t, err)

	// getMethod := func(name string) *object.Builtin {
	// 	method, ok := proxy.GetAttr(name)
	// 	require.True(t, ok)
	// 	return method.(*object.Builtin)
	// }

	// flub := getMethod("Flub")
	// inc := getMethod("Increment")
	// toUpper := getMethod("ToUpper")
	// parseInt := getMethod("ParseInt")

	// // Call Flub and check the result
	// res := flub.Call(ctx, object.NewMap(map[string]object.Object{
	// 	"A": object.NewInt(99),
	// 	"B": object.NewString("B"),
	// 	"C": object.NewBool(true),
	// }))
	// require.Equal(t, "flubbed:99.B.true", res.(*object.String).Value())

	// // Try calling Increment
	// res = inc.Call(ctx, object.NewInt(123))
	// require.Equal(t, int64(124), res.(*object.Int).Value())

	// // Try calling ToUpper
	// res = toUpper.Call(ctx, object.NewString("hello"))
	// require.Equal(t, "HELLO", res.(*object.String).Value())

	// Call ParseInt and check that an Ok result is returned
	// res = parseInt.Call(ctx, object.NewString("234"))
	// result, ok := res.(*object.Result)
	// require.True(t, ok)
	// require.True(t, result.IsOk())
	// require.Equal(t, int64(234), result.Unwrap().(*object.Int).Value())

	// // Call ParseInt with an invalid input and check that an Err result is returned
	// res = parseInt.Call(ctx, object.NewString("not-an-int"))
	// result, ok = res.(*object.Result)
	// require.True(t, ok)
	// require.True(t, result.IsErr())
	// require.Equal(t, "strconv.Atoi: parsing \"not-an-int\": invalid syntax",
	// 	result.UnwrapErr().Message().Value())
}

type proxyTestType1 []string

func (p proxyTestType1) Len() int {
	return len(p)
}

func TestProxyNonStruct(t *testing.T) {
	proxy, err := object.NewProxy(proxyTestType1{})
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

	ptt1, err := object.NewGoType(proxyTestType1{})
	require.Nil(t, err)
	require.Equal(t, ptt1, field.GoType())
}
