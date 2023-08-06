//go:build google
// +build google

package google

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"unicode"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

type Client struct {
	client  interface{}
	service string
	methods map[string]*GoMethod
}

func (c *Client) Inspect() string {
	return fmt.Sprintf("google.client(service=%s)", c.service)
}

func (c *Client) Type() object.Type {
	return "google.client"
}

func (c *Client) Value() interface{} {
	return c.client
}

func (c *Client) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "__service__":
		return object.NewString(c.service), true
	}
	method, ok := c.methods[name]
	if !ok {
		return nil, false
	}
	methodName := fmt.Sprintf("google.%s.%s", c.service, method.Name)
	return NewMethod(methodName, c.client, method), true
}

func (c *Client) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: google.client object has no attribute %q", name)
}

func (c *Client) Interface() interface{} {
	return c.client
}

func (c *Client) String() string {
	return c.Inspect()
}

func (c *Client) Compare(other object.Object) (int, error) {
	return 0, errors.New("type error: unable to compare google.client")
}

func (c *Client) Equals(other object.Object) object.Object {
	if c == other {
		return object.True
	}
	return object.False
}

func (c *Client) IsTruthy() bool {
	return true
}

func (c *Client) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("eval error: unsupported operation for google.client: %v ", opType))
}

func (c *Client) Cost() int {
	return 0
}

func (c *Client) MarshalJSON() ([]byte, error) {
	return nil, errors.New("type error: unable to marshal google.client")
}

func NewClient(service string, client interface{}) *Client {
	return &Client{
		service: service,
		client:  client,
		methods: loadMethods(client),
	}
}

func getClient(ctx context.Context, service, resource string) object.Object {
	switch service {
	case "compute":
		switch resource {
		case "instances":
			c, err := compute.NewInstancesRESTClient(ctx)
			if err != nil {
				return object.NewError(err)
			}
			return NewClient(service, c)
		default:
			return object.Errorf("unknown google compute resource: %s", resource)
		}
	default:
		return object.Errorf("unknown google service: %s", service)
	}
}

func loadMethods(obj interface{}) map[string]*GoMethod {
	typ := reflect.TypeOf(obj)
	methods := make(map[string]*GoMethod, typ.NumMethod())
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		name := toSnakeCase(m.Name)
		goMethod := &GoMethod{
			Method:     m,
			Name:       name,
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
		methods[name] = goMethod
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

func toSnakeCase(str string) string {
	var lastUpper bool
	var result string
	for i, v := range str {
		if unicode.IsUpper(v) {
			if i != 0 && !lastUpper {
				result += "_"
			}
			result += string(unicode.ToLower(v))
			lastUpper = true
		} else {
			result += string(v)
			lastUpper = false
		}
	}
	return result
}
