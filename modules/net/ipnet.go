package net

import (
	"context"
	"fmt"
	"net"

	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

var _ object.Object = (*IPNet)(nil)

const IPNET object.Type = "net.ipnet"

type IPNet struct {
	value *net.IPNet
}

func (n *IPNet) IsTruthy() bool {
	return true
}

func (n *IPNet) Type() object.Type {
	return IPNET
}

func (n *IPNet) Value() *net.IPNet {
	return n.value
}

func (n *IPNet) Inspect() string {
	return fmt.Sprintf("%s(%s)", IPNET, n.value.String())
}

func (n *IPNet) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("type error: cannot set %q on %s object", name, IPNET)
}

func (n *IPNet) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "contains":
		return object.NewBuiltin("contains", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("net.ipnet.contains", 1, args); err != nil {
				return err
			}
			switch arg := args[0].(type) {
			case *IP:
				return object.NewBool(n.value.Contains(arg.value))
			case *object.String:
				ipAddr := net.ParseIP(arg.Value())
				if ipAddr == nil {
					return object.Errorf("value error: invalid ip address %q", arg.Value())
				}
				return object.NewBool(n.value.Contains(ipAddr))
			default:
				return object.TypeErrorf("type error: expected ip address (got %s)", args[0].Type())
			}
		}), true
	case "network":
		return object.NewBuiltin("network", func(ctx context.Context, args ...object.Object) object.Object {
			return object.NewString(n.value.Network())
		}), true
	case "string":
		return object.NewBuiltin("string", func(ctx context.Context, args ...object.Object) object.Object {
			return object.NewString(n.value.String())
		}), true
	default:
		return nil, false
	}
}

func (n *IPNet) String() string {
	return n.value.String()
}

func (n *IPNet) Interface() interface{} {
	return n.value
}

func (n *IPNet) Equals(other object.Object) object.Object {
	if n == other {
		return object.True
	}
	return object.False
}

func (n *IPNet) Cost() int {
	return 0
}

func (n *IPNet) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for %s: %v", IPNET, opType)
}

func NewIPNet(v *net.IPNet) *IPNet {
	return &IPNet{value: v}
}
