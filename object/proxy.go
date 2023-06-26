package object

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/cloudcmds/tamarin/v2/op"
)

var (
	goTypeMutex      = &sync.RWMutex{}
	errorInterface   = reflect.TypeOf((*error)(nil)).Elem()
	contextInterface = reflect.TypeOf((*context.Context)(nil)).Elem()
	goTypeRegistry   = map[reflect.Type]*GoType{}
)

func IsProxyableKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Invalid, reflect.UnsafePointer, reflect.Chan, reflect.Array:
		return false
	default:
		return true
	}
}

func IsProxyableType(typ reflect.Type) bool {
	kind := typ.Kind()
	if kind == reflect.Ptr {
		// Indirection is allowed only for structs, and only one level
		return typ.Elem().Kind() == reflect.Struct
	}
	return IsProxyableKind(kind)
}

// GoAttribute is an interface to represent an attribute on a Go type. This could
// be either a field or a method.
type GoAttribute interface {
	// Name of the attribute.
	Name() string
}

func LookupGoType(obj interface{}) (*GoType, bool) {
	goType, found := goTypeRegistry[reflect.TypeOf(obj)]
	return goType, found
}

// Proxy is a Tamarin type that proxies method calls to a wrapped Go struct.
// Only the public methods of the Go type are proxied.
type Proxy struct {
	typ *GoType
	obj interface{}
}

func (p *Proxy) Type() Type {
	return PROXY
}

func (p *Proxy) Interface() interface{} {
	return p.obj
}

func (p *Proxy) Inspect() string {
	return fmt.Sprintf("%v", p.obj)
}

func (p *Proxy) String() string {
	return fmt.Sprintf("proxy(%v)", p.obj)
}

func (p *Proxy) GoType() *GoType {
	return p.typ
}

func (p *Proxy) GetAttr(name string) (Object, bool) {
	attr, found := p.typ.GetAttribute(name)
	if !found {
		return nil, false
	}
	switch attr := attr.(type) {
	case *GoField:
		conv, ok := attr.Converter()
		if !ok {
			return Errorf("type error: no converter for field %s", name), true
		}
		var value interface{}
		if _, ok := p.typ.IndirectType(); ok {
			value = reflect.ValueOf(p.obj).Elem().FieldByName(name).Interface()
		} else {
			value = reflect.ValueOf(p.obj).FieldByName(name).Interface()
		}
		result, err := conv.From(value)
		if err != nil {
			return NewError(err), true
		}
		return result, true
	case *GoMethod:
		return &Builtin{
			name: fmt.Sprintf("%s.%s", p.typ.Name(), name),
			fn: func(ctx context.Context, args ...Object) Object {
				return p.call(ctx, attr, args...)
			},
		}, true
	}
	return nil, false
}

func (p *Proxy) Equals(other Object) Object {
	if p == other {
		return True
	}
	return False
}

func (p *Proxy) IsTruthy() bool {
	return true
}

func (p *Proxy) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for proxy: %v", opType))
}

func (p *Proxy) Cost() int {
	return 8
}

func (p *Proxy) call(ctx context.Context, m *GoMethod, args ...Object) Object {
	// methodName := fmt.Sprintf("%s.%s", p.typ.Name(), m.name)
	// isVariadic := m.method.Type.IsVariadic()
	// var argIndex int
	// inputs := make([]reflect.Value, 1, m.numIn)
	// inputs[0] = reflect.ValueOf(p.obj)
	// minArgs := m.numIn
	// if isVariadic {
	// 	minArgs--
	// }
	// for i := 1; i < m.numIn; i++ {
	// 	converter := m.inputConverters[i-1]
	// 	if _, ok := converter.(*ContextConverter); ok {
	// 		inputs = append(inputs, reflect.ValueOf(ctx))
	// 		continue
	// 	}
	// 	if argIndex >= len(args) {
	// 		break
	// 	}
	// 	input, err := m.inputConverters[i-1].To(args[argIndex])
	// 	if err != nil {
	// 		return Errorf("type error: failed to convert argument %d in %s() call: %s", i, methodName, err)
	// 	}
	// 	inputs = append(inputs, reflect.ValueOf(input))
	// 	argIndex++
	// }
	// if len(inputs) < minArgs {
	// 	return Errorf("type error: %s() requires %d arguments, but %d were given",
	// 		methodName, minArgs, len(inputs))
	// }
	// outputs := m.method.Func.Call(inputs)
	// if len(outputs) == 0 {
	// 	return Nil
	// } else if len(outputs) == 1 {
	// 	if m.outputHasErr {
	// 		errObj := outputs[0].Interface()
	// 		if errObj != nil {
	// 			return NewErrResult(NewError(errObj.(error)))
	// 		}
	// 		return NewOkResult(Nil)
	// 	}
	// 	obj, err := m.outputConverters[0].From(outputs[0].Interface())
	// 	if err != nil {
	// 		return Errorf("call error: failed to convert output from %s() call: %s", methodName, err)
	// 	}
	// 	return obj
	// } else if len(outputs) == 2 {
	// 	if !m.outputHasErr {
	// 		return Errorf("call error: too many outputs from %s() call", methodName)
	// 	}
	// 	obj0, err := m.outputConverters[0].From(outputs[0].Interface())
	// 	if err != nil {
	// 		return Errorf("call error: failed to convert output from %s() call: %s", methodName, err)
	// 	}
	// 	obj1, err := m.outputConverters[1].From(outputs[1].Interface())
	// 	if err != nil {
	// 		return Errorf("call error: failed to convert output from %s() call: %s", methodName, err)
	// 	}
	// 	var resObj, errObj Object
	// 	if m.outputErrIndex == 0 {
	// 		errObj = obj0
	// 		resObj = obj1
	// 	} else {
	// 		errObj = obj1
	// 		resObj = obj0
	// 	}
	// 	if errObj != nil {
	// 		return NewErrResult(errObj.(*Error))
	// 	}
	// 	return NewOkResult(resObj)
	// }
	// return Errorf("call error: method %s has too many outputs", methodName)
	return NewInt(42)
}

// NewProxy returns a new Tamarin proxy object that wraps the given Go object.
// This operation may fail if the Go type has attributes whose types cannot be
// converted to Tamarin types.
func NewProxy(obj interface{}) (*Proxy, error) {

	typ := reflect.TypeOf(obj)

	// Is this type proxyable?
	if !IsProxyableType(typ) {
		return nil, fmt.Errorf("type error: unsupported argument for go_type (%t given)", typ)
	}

	goType, err := NewGoType(typ)
	if err != nil {
		return nil, err
	}

	return &Proxy{typ: goType, obj: obj}, nil
}
