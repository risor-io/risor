package net

import (
	"context"
	"fmt"
	"net"

	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

var _ object.Object = (*IP)(nil)

const IPTYPE object.Type = "net.ip"

type IP struct {
	value net.IP
}

func (ip *IP) IsTruthy() bool {
	return true
}

func (ip *IP) Type() object.Type {
	return IPTYPE
}

func (ip *IP) Value() net.IP {
	return ip.value
}

func (ip *IP) Inspect() string {
	return fmt.Sprintf("%s(%s)", IPTYPE, ip.value.String())
}

func (ip *IP) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: cannot set %q on %s object", name, IPTYPE)
}

func (ip *IP) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "default_mask":
		return object.NewBuiltin("default_mask", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("net.ip.default_mask", 0, args); err != nil {
				return err
			}
			return object.NewByteSlice(ip.value.DefaultMask())
		}), true
	case "equal":
		return object.NewBuiltin("equal", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("net.ip.equal", 1, args); err != nil {
				return err
			}
			other, ok := args[0].(*IP)
			if !ok {
				return object.Errorf("eval error: expected %s, got %s", IPTYPE, args[0].Type())
			}
			return object.NewBool(ip.value.Equal(other.value))
		}), true
	case "is_global_unicast":
		return object.NewBuiltin("is_global_unicast", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("net.ip.is_global_unicast", 0, args); err != nil {
				return err
			}
			return object.NewBool(ip.value.IsGlobalUnicast())
		}), true
	case "is_loopback":
		return object.NewBuiltin("is_loopback", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("net.ip.is_loopback", 0, args); err != nil {
				return err
			}
			return object.NewBool(ip.value.IsLoopback())
		}), true
	case "is_multicast":
		return object.NewBuiltin("is_multicast", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("net.ip.is_multicast", 0, args); err != nil {
				return err
			}
			return object.NewBool(ip.value.IsMulticast())
		}), true
	case "is_unspecified":
		return object.NewBuiltin("is_unspecified", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("net.ip.is_unspecified", 0, args); err != nil {
				return err
			}
			return object.NewBool(ip.value.IsUnspecified())
		}), true
	case "mask":
		return object.NewBuiltin("mask", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("net.ip.mask", 1, args); err != nil {
				return err
			}
			mask, err := object.AsBytes(args[0])
			if err != nil {
				return err
			}
			return NewIP(ip.value.Mask(mask))
		}), true
	case "to16":
		return object.NewBuiltin("to16", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("net.ip.to16", 0, args); err != nil {
				return err
			}
			return NewIP(ip.value.To16())
		}), true
	case "to4":
		return object.NewBuiltin("to4", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("net.ip.to4", 0, args); err != nil {
				return err
			}
			return NewIP(ip.value.To4())
		}), true
	case "string":
		return object.NewBuiltin("string", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("net.ip.string", 0, args); err != nil {
				return err
			}
			return object.NewString(ip.value.String())
		}), true
	default:
		return nil, false
	}
}

func (ip *IP) String() string {
	return ip.value.String()
}

func (ip *IP) Interface() interface{} {
	return ip.value
}

func (ip *IP) Equals(other object.Object) object.Object {
	if ip == other {
		return object.True
	}
	return object.False
}

func (ip *IP) Cost() int {
	return 0
}

func (ip *IP) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.EvalErrorf("eval error: unsupported operation for %s: %v", IPTYPE, opType)
}

func NewIP(v net.IP) *IP {
	return &IP{value: v}
}
