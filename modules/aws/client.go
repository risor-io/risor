//go:build aws
// +build aws

package aws

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

type Client struct {
	client  interface{}
	service string
	methods map[string]*GoMethod
}

func (c *Client) Inspect() string {
	return fmt.Sprintf("aws.client(service=%q)", c.service)
}

func (c *Client) Type() object.Type {
	return "aws.client"
}

func (c *Client) Value() interface{} {
	return c.client
}

func (c *Client) GetAttr(name string) (object.Object, bool) {
	method, ok := c.methods[name]
	if !ok {
		return nil, false
	}
	methodName := fmt.Sprintf("aws.%s.%s", c.service, method.Name)
	return NewMethod(methodName, c.client, method), true
}

func (c *Client) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: aws.client object has no attribute %q", name)
}

func (c *Client) Interface() interface{} {
	return c.client
}

func (c *Client) String() string {
	return c.Inspect()
}

func (c *Client) Compare(other object.Object) (int, error) {
	return 0, errors.New("type error: unable to compare aws.client")
}

func (c *Client) Equals(other object.Object) object.Object {
	if c == other {
		return object.True
	}
	return object.False
}

func (c *Client) IsTruthy() bool {
	return c.service != "" && c.client != nil
}

func (c *Client) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("eval error: unsupported operation for aws.client: %v ", opType))
}

func (c *Client) Cost() int {
	return 0
}

func NewClient(service string, client interface{}) *Client {
	return &Client{service: service, client: client, methods: loadMethods(client)}
}

func loadMethods(obj interface{}) map[string]*GoMethod {
	typ := reflect.TypeOf(obj)
	methods := make(map[string]*GoMethod, typ.NumMethod())
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		goMethod := &GoMethod{
			Method:     m,
			Name:       m.Name,
			NumIn:      m.Type.NumIn(),
			NumOut:     m.Type.NumOut(),
			IsVariadic: m.Type.IsVariadic(),
		}
		for i := 0; i < goMethod.NumIn; i++ {
			goMethod.InTypes = append(goMethod.InTypes, m.Type.In(i))
		}
		for i := 0; i < goMethod.NumOut; i++ {
			goMethod.OutTypes = append(goMethod.OutTypes, m.Type.Out(i))
		}
		methods[m.Name] = goMethod
	}
	return methods
}

type GoMethod struct {
	Name       string
	Method     reflect.Method
	NumIn      int
	NumOut     int
	InTypes    []reflect.Type
	OutTypes   []reflect.Type
	IsVariadic bool
}
