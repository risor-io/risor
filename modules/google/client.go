//go:build google
// +build google

package google

import (
	"context"
	"errors"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

type Client struct {
	client  interface{}
	service string
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
	return nil, false
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
	}
}

func getClient(ctx context.Context, service string) object.Object {
	switch service {
	case "compute":
		c, err := compute.NewInstancesRESTClient(ctx)
		if err != nil {
			return object.NewError(err)
		}
		return NewClient(service, c)
	default:
		return object.Errorf("unknown google service: %s", service)
	}
}
