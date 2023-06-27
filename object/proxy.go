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

func IsProxyableType(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Interface, reflect.Slice:
		return true
	case reflect.Ptr:
		return typ.Elem().Kind() == reflect.Struct
	default:
		return false
	}
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
	*base
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
	return p.String()
}

func (p *Proxy) String() string {
	return fmt.Sprintf("proxy(%s(%v))", reflect.TypeOf(p.obj), p.obj)
}

func (p *Proxy) GoType() *GoType {
	return p.typ
}

func (p *Proxy) GetAttr(name string) (Object, bool) {
	if name == "__type__" {
		return p.typ, true
	}
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

func (p *Proxy) SetAttr(name string, value Object) error {
	attr, found := p.typ.GetAttribute(name)
	if !found {
		return fmt.Errorf("attribute error: %s has no attribute %s", p.typ.Name(), name)
	}
	switch attr := attr.(type) {
	case *GoField:
		conv, ok := attr.Converter()
		if !ok {
			return fmt.Errorf("type error: no converter for field %s", name)
		}
		var field reflect.Value
		if _, ok := p.typ.IndirectType(); ok {
			field = reflect.ValueOf(p.obj).Elem().FieldByName(name)
		} else {
			field = reflect.ValueOf(p.obj).FieldByName(name)
		}
		result, err := conv.To(value)
		if err != nil {
			return err
		}
		if field.CanSet() {
			field.Set(reflect.ValueOf(result))
			return nil
		} else {
			return fmt.Errorf("type error: cannot set field %s", name)
		}
	case *GoMethod:
		return fmt.Errorf("attribute error: cannot set method %s", name)
	}
	return fmt.Errorf("attribute error: unknown attribute type")
}

func (p *Proxy) Equals(other Object) Object {
	if p == other {
		return True
	}
	return False
}

func (p *Proxy) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for proxy: %v", opType))
}

func (p *Proxy) call(ctx context.Context, m *GoMethod, args ...Object) Object {
	methodName := fmt.Sprintf("%s.%s", p.typ.Name(), m.name)
	isVariadic := m.method.Type.IsVariadic()
	var argIndex int
	numIn := m.NumIn()
	inputs := make([]reflect.Value, 1, numIn)
	inputs[0] = reflect.ValueOf(p.obj)
	minArgs := numIn
	if isVariadic {
		minArgs--
	}
	for i := 1; i < numIn; i++ {
		inType := m.inputTypes[i]
		converter := inType.converter
		if _, ok := converter.(*ContextConverter); ok {
			inputs = append(inputs, reflect.ValueOf(ctx))
			continue
		}
		if argIndex >= len(args) {
			break
		}
		input, err := converter.To(args[argIndex])
		if err != nil {
			return Errorf("type error: failed to convert argument %d in %s() call: %s", i, methodName, err)
		}
		inputs = append(inputs, reflect.ValueOf(input))
		argIndex++
	}
	if len(inputs) < minArgs {
		return Errorf("type error: %s() requires %d arguments, but %d were given",
			methodName, minArgs, len(inputs))
	}
	outputs := m.method.Func.Call(inputs)
	if len(outputs) == 0 {
		return Nil
	}
	for _, errIndex := range m.errorIndices {
		errObj := outputs[errIndex].Interface()
		if errObj != nil {
			return NewError(errObj.(error))
		}
	}
	outputCount := len(outputs) - len(m.errorIndices)
	if outputCount <= 1 {
		for i, output := range outputs {
			if !m.IsOutputError(i) {
				outType := m.outputTypes[i]
				result, err := outType.converter.From(output.Interface())
				if err != nil {
					return Errorf("call error: failed to convert output from %s() call: %s", methodName, err)
				}
				return result
			}
		}
		return Nil
	}
	var results []Object
	for i, output := range outputs {
		if !m.IsOutputError(i) {
			outType := m.outputTypes[i]
			result, err := outType.converter.From(output.Interface())
			if err != nil {
				return Errorf("call error: failed to convert output from %s() call: %s", methodName, err)
			}
			results = append(results, result)
		}
	}
	return NewList(results)
}

// NewProxy returns a new Tamarin proxy object that wraps the given Go object.
// This operation may fail if the Go type has attributes whose types cannot be
// converted to Tamarin types.
func NewProxy(obj interface{}) (*Proxy, error) {

	typ := reflect.TypeOf(obj)

	// Is this type proxyable?
	if !IsProxyableType(typ) {
		return nil, fmt.Errorf("type error: unable to proxy type (%T given)", obj)
	}

	goType, err := NewGoType(typ)
	if err != nil {
		return nil, err
	}

	return &Proxy{typ: goType, obj: obj}, nil
}
