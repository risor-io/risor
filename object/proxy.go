package object

import (
	"context"
	"fmt"
	"reflect"
)

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

// ProxyType represents a single Go type whose methods can be proxied.
type ProxyType struct {
	NumMethod int
	Methods   []*ProxyMethod
	Value     interface{}
	Type      reflect.Type
}

// ProxyMethod represents a single method on a Go type that can be proxied.
type ProxyMethod struct {
	Name             string
	Method           reflect.Method
	NumIn            int
	NumOut           int
	InputConverters  []TypeConverter
	OutputConverters []TypeConverter
	OutputHasErr     bool
}

// ProxyManager is an interface that defines a way to register Go types and call
// methods on instances of those types.
type ProxyManager interface {

	// RegisterType determines type and method information for the provided
	// object and saves that information for use in later method call proxying.
	RegisterType(obj interface{}) (*ProxyType, error)

	// HasType returns true if the type of the objects has been registered.
	HasType(obj interface{}) bool

	// GetMethod returns a wrapped method from the given object and a boolean
	// that indicates whether the method was found.
	GetMethod(obj interface{}, method string) (*ProxyMethod, bool)

	// Call the named method on the object with the given arguments.
	// The type of the object must have been previously registered, otherwise
	// a Tamarin error object is returned.
	Call(ctx context.Context, obj interface{}, method *ProxyMethod, args ...Object) Object
}

// DefaultProxyManager implements the ProxyManager interface.
type DefaultProxyManager struct {
	types      map[reflect.Type]*ProxyType
	converters map[reflect.Type]TypeConverter
}

type ProxyManagerOpts struct {
	Types      []any
	Converters []TypeConverter
	NoDefaults bool
}

// NewProxyManager creates a ProxyManager that can be used to proxy method calls
// to various struct types. The provided type conversion functions are used to
// translate between Go and Tamarin types.
func NewProxyManager(opts ProxyManagerOpts) (*DefaultProxyManager, error) {
	mgr := &DefaultProxyManager{
		types:      map[reflect.Type]*ProxyType{},
		converters: map[reflect.Type]TypeConverter{},
	}
	if !opts.NoDefaults {
		defaultConverters := []TypeConverter{
			&IntConverter{},
			&Int64Converter{},
			&Float32Converter{},
			&Float64Converter{},
			&TimeConverter{},
			&StringConverter{},
			&BooleanConverter{},
			&MapStringIfaceConverter{},
		}
		for _, tc := range defaultConverters {
			mgr.converters[tc.Type()] = tc
		}
		for _, obj := range opts.Types {
			sc := &StructConverter{Prototype: obj}
			mgr.converters[sc.Type()] = sc
		}
	}
	for _, tc := range opts.Converters {
		mgr.converters[tc.Type()] = tc
	}
	for _, obj := range opts.Types {
		if _, err := mgr.RegisterType(obj); err != nil {
			return nil, err
		}
	}
	return mgr, nil
}

func (p *DefaultProxyManager) RegisterType(obj interface{}) (*ProxyType, error) {
	typ := reflect.TypeOf(obj)
	if proxyType, found := p.types[typ]; found {
		return proxyType, nil
	}
	proxyType := &ProxyType{Type: typ, NumMethod: typ.NumMethod()}
	// Collect information about each public method on the type
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		numIn := m.Type.NumIn()
		numOut := m.Type.NumOut()
		proxyMethod := &ProxyMethod{
			Name:   m.Name,
			NumIn:  numIn,
			NumOut: numOut,
			Method: m,
		}
		// Determine the converter for each input
		for j := 1; j < numIn; j++ { // Skip the "self" parameter at j=0
			inType := m.Type.In(j)
			converter, found := p.converters[inType]
			if !found {
				if inType.Implements(errorInterface) {
					converter = &ErrorConverter{}
				} else {
					return nil, fmt.Errorf("type error: no type conversion function found for %s (%s.%s([%d]))",
						inType, typ, m.Name, j-1)
				}
			}
			proxyMethod.InputConverters = append(proxyMethod.InputConverters, converter)
		}
		// Determine the converter for each output
		for j := 0; j < numOut; j++ {
			outType := m.Type.Out(j)
			converter, found := p.converters[outType]
			if !found {
				if outType.Implements(errorInterface) {
					converter = &ErrorConverter{}
					proxyMethod.OutputHasErr = true
				} else {
					return nil, fmt.Errorf("type error: no type conversion function found for %s (%s.%s()[%d])",
						outType, typ, m.Name, j)
				}
			}
			proxyMethod.OutputConverters = append(proxyMethod.OutputConverters, converter)
		}
		proxyType.Methods = append(proxyType.Methods, proxyMethod)
	}
	p.types[typ] = proxyType
	return proxyType, nil
}

func (p *DefaultProxyManager) HasType(obj interface{}) bool {
	_, found := p.types[reflect.TypeOf(obj)]
	return found
}

func (p *DefaultProxyManager) GetType(obj interface{}) (*ProxyType, bool) {
	typ := reflect.TypeOf(obj)
	if proxyType, found := p.types[typ]; found {
		return proxyType, true
	}
	return nil, false
}

func (p *DefaultProxyManager) GetMethod(obj interface{}, method string) (*ProxyMethod, bool) {
	typ := reflect.TypeOf(obj)
	proxyType, found := p.types[typ]
	if !found {
		return nil, false
	}
	for _, m := range proxyType.Methods {
		if m.Name == method {
			return m, true
		}
	}
	return nil, false
}

func (p *DefaultProxyManager) Call(ctx context.Context, obj interface{}, m *ProxyMethod, args ...Object) Object {
	if len(args) != m.NumIn-1 {
		return NewError("wrong number of arguments. got=%d, want=%d", len(args), m.NumIn-1)
	}
	inputs := make([]reflect.Value, m.NumIn)
	inputs[0] = reflect.ValueOf(obj)
	for i, arg := range args {
		input, err := m.InputConverters[i].To(arg)
		if err != nil {
			return NewError("error converting input %d: %s", i, err)
		}
		inputs[i+1] = reflect.ValueOf(input)
	}
	// TODO: pass through context
	// TODO: handle panic and translate to error
	outputs := m.Method.Func.Call(inputs)
	if len(outputs) == 0 {
		return Nil
	} else if len(outputs) == 1 {
		if m.OutputHasErr {
			if obj != nil {
				err := outputs[0].Interface().(error)
				return NewErrorResult(err.Error())
			}
			return NewOkResult(Nil)
		}
		obj, err := m.OutputConverters[0].From(outputs[0].Interface())
		if err != nil {
			return NewError("error converting output: %s", err)
		}
		return obj
	} else if len(outputs) == 2 {
		if !m.OutputHasErr {
			return NewError("too many outputs")
		}
		obj0, err := m.OutputConverters[0].From(outputs[0].Interface())
		if err != nil {
			return NewError("error converting output: %s", err)
		}
		obj1, err := m.OutputConverters[1].From(outputs[1].Interface())
		if err != nil {
			return NewError("error converting output: %s", err)
		}
		if obj1 != nil {
			errObj := obj1.(*Error)
			return NewErrorResult(errObj.Message)
		}
		return NewOkResult(obj0)
	}
	return NewError("too many outputs")
}

// Proxy is a Tamarin type that proxies method calls to a wrapped Go struct.
// Only the public methods of the Go type are proxied.
type Proxy struct {
	mgr ProxyManager
	obj interface{}
}

func (p *Proxy) Type() Type {
	return PROXY
}

func (p *Proxy) Inspect() string {
	return fmt.Sprintf("%v", p.obj)
}

func (p *Proxy) GetAttr(name string) (Object, bool) {
	method, found := p.mgr.GetMethod(p.obj, name)
	if !found {
		return nil, false
	}
	return &Builtin{
		Name: fmt.Sprintf("%v.%s", reflect.TypeOf(p.obj), name),
		Fn: func(ctx context.Context, args ...Object) Object {
			return p.mgr.Call(ctx, p.obj, method, args...)
		},
	}, true
}

func (p *Proxy) ToInterface() interface{} {
	return p.obj
}

func (p *Proxy) String() string {
	return fmt.Sprintf("%v", p.obj)
}

func (p *Proxy) Equals(other Object) Object {
	if other.Type() == PROXY && p.obj == other.(*Proxy).obj {
		return True
	}
	return False
}

// NewProxy returns a new Tamarin proxy object that wraps the given Go object.
// The Go type should previously been registered with the ProxyManager.
func NewProxy(mgr ProxyManager, obj interface{}) *Proxy {
	return &Proxy{mgr: mgr, obj: obj}
}
