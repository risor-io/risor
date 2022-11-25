package object_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/cloudcmds/tamarin/object"
)

type FlubOpts struct {
	A int
	B string
}

type Embedded struct{}

func (e *Embedded) Flub(opts FlubOpts) string {
	return fmt.Sprintf("flubbed:%d.%s", opts.A, opts.B)
}

func (e *Embedded) Test(i int64) int {
	return int(i)
}

type Whatever struct {
	*Embedded
}

func (w *Whatever) Hello(response string, allCaps bool) string {
	if allCaps {
		return strings.ToUpper(response)
	}
	return response
}

func TestProxy(t *testing.T) {
	w := &Whatever{
		Embedded: &Embedded{},
	}
	var v interface{} = w
	type foo struct{}
	fmt.Println("TYPE:", reflect.TypeOf(struct{}{}), reflect.TypeOf(foo{}), reflect.TypeOf(foo{}).Kind())

	mgr := object.NewProxyManager([]object.TypeConverter{
		&object.IntConverter{},
		&object.Int64Converter{},
		&object.StringConverter{},
		&object.BooleanConverter{},
		&object.ErrorConverter{},
		&object.StructConverter{Prototype: FlubOpts{}},
	})
	_, err := mgr.RegisterType("whatev", v)
	if err != nil {
		t.Fatal(err)
	}

	proxy := object.NewProxy(mgr, v)

	hashKey := object.NewString("A")
	hashVal := object.NewInteger(99)
	B := object.NewString("B")
	hash := &object.Hash{Pairs: map[object.HashKey]object.HashPair{
		hashKey.HashKey(): {Key: hashKey, Value: hashVal},
		B.HashKey():       {Key: B, Value: B},
	}}
	res := proxy.InvokeMethod("Flub", hash)
	fmt.Println("RES", res)
}
