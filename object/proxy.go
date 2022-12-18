package object

import (
	"context"
	"fmt"
	"reflect"
)

var (
	errorInterface   = reflect.TypeOf((*error)(nil)).Elem()
	contextInterface = reflect.TypeOf((*context.Context)(nil)).Elem()
)

// GoType represents a single Go type whose methods and fields can be proxied.
type GoType struct {
	attrs      []GoAttr
	methods    []*GoMethod
	fields     []*GoField
	rType      reflect.Type
	structType reflect.Type
}

func (gt *GoType) Name() string {
	return gt.structType.Name()
}

func (gt *GoType) Attrs() []GoAttr {
	return gt.attrs
}

func (gt *GoType) AttrByName(name string) (GoAttr, bool) {
	for _, a := range gt.attrs {
		if a.Name() == name {
			return a, true
		}
	}
	return nil, false
}

func (gt *GoType) Methods() []*GoMethod {
	return gt.methods
}

func (gt *GoType) MethodByName(name string) (*GoMethod, bool) {
	for _, m := range gt.methods {
		if m.name == name {
			return m, true
		}
	}
	return nil, false
}

func (gt *GoType) Fields() []*GoField {
	return gt.fields
}

func (gt *GoType) FieldByName(name string) (*GoField, bool) {
	for _, f := range gt.fields {
		if f.name == name {
			return f, true
		}
	}
	return nil, false
}

func (gt *GoType) Type() reflect.Type {
	return gt.rType
}

func (gt *GoType) StructType() reflect.Type {
	return gt.structType
}

// GoAttrType is used to indicate whether a GoAttr is a field or a method.
type GoAttrType string

const (
	// GoAttrTypeMethod indicates that the GoAttr is a method.
	GoAttrTypeMethod GoAttrType = "method"

	// GoAttrTypeField indicates that the GoAttr is a field.
	GoAttrTypeField GoAttrType = "field"
)

// GoAttr is an interface to represent an attribute on a Go type. This could
// be either a field or a method.
type GoAttr interface {

	// Name of the attribute.
	Name() string

	// Type indicates whether the attribute is a method or a field.
	AttrType() GoAttrType
}

// GoMethod represents a single method on a Go type that can be proxied.
type GoMethod struct {
	name             string
	method           reflect.Method
	numIn            int
	numOut           int
	inputConverters  []TypeConverter
	outputConverters []TypeConverter
	outputHasErr     bool
	outputErrIndex   int
}

func (m *GoMethod) Name() string {
	return m.name
}

func (m *GoMethod) AttrType() GoAttrType {
	return GoAttrTypeMethod
}

// GoField represents a single field on a Go type that can be read or written.
type GoField struct {
	name      string
	rType     reflect.Type
	tag       reflect.StructTag
	field     reflect.StructField
	converter TypeConverter
}

func (f *GoField) Name() string {
	return f.name
}

func (f *GoField) AttrType() GoAttrType {
	return GoAttrTypeField
}

func (f *GoField) Type() reflect.Type {
	return f.rType
}

func (f *GoField) Tag() reflect.StructTag {
	return f.tag
}

func (f *GoField) Converter() TypeConverter {
	return f.converter
}

// GoTypeRegistry is an interface that defines a way to register Go types and call
// methods on instances of those types.
type GoTypeRegistry interface {

	// Register determines type and method information for the provided
	// object and saves that information for use in later method call proxying.
	Register(obj interface{}) (*GoType, error)

	// GetType returns the GoType for the given object and a boolean that
	// indicates whether the type was found. Only types that were previously
	// registered will be found.
	GetType(obj interface{}) (*GoType, bool)

	// GetAttr returns the GoAttr for the given object and name.
	GetAttr(obj interface{}, name string) (GoAttr, bool)
}

// DefaultTypeRegistry implements the GoTypeRegistry interface.
type DefaultTypeRegistry struct {
	types      map[reflect.Type]*GoType
	converters map[reflect.Type]TypeConverter
}

// TypeRegistryOpts contains options used to create a GoTypeRegistry.
type TypeRegistryOpts struct {
	// Converters is a list of TypeConverters that will be used to convert
	// input and output types for method calls.
	Converters []TypeConverter

	// NoDefaults indicates that the default TypeConverters should not be
	// automatically used by the registry. If this is set, the caller should
	// provide their own TypeConverters.
	NoDefaults bool
}

// NewTypeRegistry creates a GoTypeRegistry that can be used to proxy method
// calls to various struct types. The provided type conversion functions are
// used to translate between Go and Tamarin types.
func NewTypeRegistry(opts ...TypeRegistryOpts) (*DefaultTypeRegistry, error) {
	var providedOpts TypeRegistryOpts
	if len(opts) > 0 {
		providedOpts = opts[0]
	}
	mgr := &DefaultTypeRegistry{
		types:      map[reflect.Type]*GoType{},
		converters: map[reflect.Type]TypeConverter{},
	}
	if !providedOpts.NoDefaults {
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
	}
	for _, tc := range providedOpts.Converters {
		mgr.converters[tc.Type()] = tc
	}
	return mgr, nil
}

func (p *DefaultTypeRegistry) getConverter(t reflect.Type) (TypeConverter, bool) {
	if converter, ok := p.converters[t]; ok {
		return converter, true
	}
	switch t.Kind() {
	case reflect.Struct:
		return &StructConverter{Prototype: reflect.New(t).Elem().Interface()}, true
	case reflect.Interface:
		if t.Implements(errorInterface) {
			return &ErrorConverter{}, true
		}
		if t.Implements(contextInterface) {
			return &ContextConverter{}, true
		}
	case reflect.Pointer:
		converter, found := p.getConverter(t.Elem())
		if !found {
			return nil, false
		}
		if structConv, ok := converter.(*StructConverter); ok {
			structConv.AsPointer = true
			return structConv, true
		}
		return nil, false
	}
	return nil, false
}

func (p *DefaultTypeRegistry) processMethod(m reflect.Method) (*GoMethod, error) {
	goMethod := &GoMethod{
		name:   m.Name,
		method: m,
		numIn:  m.Type.NumIn(),
		numOut: m.Type.NumOut(),
	}
	// Choose a converter for each input, skipping the "self" param at i=0
	for i := 1; i < goMethod.numIn; i++ {
		inType := m.Type.In(i)
		converter, found := p.getConverter(inType)
		if !found {
			return nil, fmt.Errorf("type error: no type conversion function found for %s", inType)
		}
		goMethod.inputConverters = append(goMethod.inputConverters, converter)
	}
	// Choose a converter for each output
	for i := 0; i < goMethod.numOut; i++ {
		outType := m.Type.Out(i)
		converter, found := p.getConverter(outType)
		if !found {
			return nil, fmt.Errorf("type error: no type conversion function found for %s", outType)
		}
		if _, ok := converter.(*ErrorConverter); ok {
			goMethod.outputHasErr = true
			goMethod.outputErrIndex = i
		}
		goMethod.outputConverters = append(goMethod.outputConverters, converter)
	}
	return goMethod, nil
}

func (p *DefaultTypeRegistry) processField(f reflect.StructField) (*GoField, error) {
	converter, found := p.converters[f.Type]
	if !found {
		return nil, fmt.Errorf("type error: no type conversion function found for %s", f.Type)
	}
	return &GoField{
		name:      f.Name,
		rType:     f.Type,
		tag:       f.Tag,
		field:     f,
		converter: converter,
	}, nil
}

func (p *DefaultTypeRegistry) Register(obj interface{}) (*GoType, error) {
	typ := reflect.TypeOf(obj)
	if goType, found := p.types[typ]; found {
		return goType, nil
	}
	var structType reflect.Type
	if typ.Kind() == reflect.Ptr {
		structType = typ.Elem()
	} else {
		structType = typ
	}
	goType := &GoType{rType: typ, structType: structType}
	// Discover type methods
	for i := 0; i < typ.NumMethod(); i++ {
		m, err := p.processMethod(typ.Method(i))
		if err != nil {
			// return nil, err
			continue // Log warning?
		}
		goType.attrs = append(goType.attrs, m)
		goType.methods = append(goType.methods, m)
	}
	// Discover type fields
	count := structType.NumField()
	for i := 0; i < count; i++ {
		f := structType.Field(i)
		if !f.IsExported() {
			continue
		}
		goField, err := p.processField(f)
		if err != nil {
			// return nil, err
			continue // Log warning?
		}
		goType.attrs = append(goType.attrs, goField)
		goType.fields = append(goType.fields, goField)
	}
	p.types[typ] = goType
	return goType, nil
}

func (p *DefaultTypeRegistry) GetType(obj interface{}) (*GoType, bool) {
	goType, found := p.types[reflect.TypeOf(obj)]
	return goType, found
}

func (p *DefaultTypeRegistry) GetAttr(obj interface{}, name string) (GoAttr, bool) {
	goType, found := p.types[reflect.TypeOf(obj)]
	if !found {
		return nil, false
	}
	goAttr, found := goType.AttrByName(name)
	return goAttr, found
}

// Proxy is a Tamarin type that proxies method calls to a wrapped Go struct.
// Only the public methods of the Go type are proxied.
type Proxy struct {
	reg GoTypeRegistry
	typ *GoType
	obj interface{}
}

func (p *Proxy) Type() Type {
	return PROXY
}

func (p *Proxy) Inspect() string {
	return fmt.Sprintf("%v", p.obj)
}

func (p *Proxy) GetAttr(name string) (Object, bool) {
	attr, found := p.reg.GetAttr(p.obj, name)
	if !found {
		return nil, false
	}
	switch attr := attr.(type) {
	case *GoField:
		v := reflect.ValueOf(p.obj).Elem().FieldByName(name).Interface()
		obj, err := attr.converter.From(v)
		if err != nil {
			return NewError(err.Error()), true
		}
		return obj, true
	case *GoMethod:
		return &Builtin{
			Name: fmt.Sprintf("%s.%s", p.typ.Name(), name),
			Fn: func(ctx context.Context, args ...Object) Object {
				return p.call(ctx, attr, args...)
			},
		}, true
	}
	return nil, false
}

func (p *Proxy) ToInterface() interface{} {
	return p.obj
}

func (p *Proxy) String() string {
	return fmt.Sprintf("proxy(%v)", p.obj)
}

func (p *Proxy) Equals(other Object) Object {
	if other.Type() == PROXY && p.obj == other.(*Proxy).obj {
		return True
	}
	return False
}

func (p *Proxy) call(ctx context.Context, m *GoMethod, args ...Object) Object {
	methodName := fmt.Sprintf("%s.%s", p.typ.Name(), m.name)
	var argIndex int
	inputs := make([]reflect.Value, m.numIn)
	inputs[0] = reflect.ValueOf(p.obj)
	for i := 1; i < m.numIn; i++ {
		converter := m.inputConverters[i-1]
		if _, ok := converter.(*ContextConverter); ok {
			inputs[i] = reflect.ValueOf(ctx)
			continue
		}
		input, err := m.inputConverters[i-1].To(args[argIndex])
		if err != nil {
			return NewError("type error: failed to convert argument %d in %s() call: %s", i, methodName, err)
		}
		inputs[i] = reflect.ValueOf(input)
		argIndex++
	}
	outputs := m.method.Func.Call(inputs)
	if len(outputs) == 0 {
		return Nil
	} else if len(outputs) == 1 {
		if m.outputHasErr {
			errObj := outputs[0].Interface()
			if errObj != nil {
				return NewErrorResult(errObj.(error).Error())
			}
			return NewOkResult(Nil)
		}
		obj, err := m.outputConverters[0].From(outputs[0].Interface())
		if err != nil {
			return NewError("call error: failed to convert output from %s() call: %s", methodName, err)
		}
		return obj
	} else if len(outputs) == 2 {
		if !m.outputHasErr {
			return NewError("call error: too many outputs from %s() call", methodName)
		}
		obj0, err := m.outputConverters[0].From(outputs[0].Interface())
		if err != nil {
			return NewError("call error: failed to convert output from %s() call: %s", methodName, err)
		}
		obj1, err := m.outputConverters[1].From(outputs[1].Interface())
		if err != nil {
			return NewError("call error: failed to convert output from %s() call: %s", methodName, err)
		}
		var resObj, errObj Object
		if m.outputErrIndex == 0 {
			errObj = obj0
			resObj = obj1
		} else {
			errObj = obj1
			resObj = obj0
		}
		if errObj != nil {
			return NewErrorResult(errObj.(*Error).Message)
		}
		return NewOkResult(resObj)
	}
	return NewError("call error: method %s has too many outputs", methodName)
}

// NewProxy returns a new Tamarin proxy object that wraps the given Go object.
// The Go type is registered with the type registry, which has no effect if the
// type is already registered. This operation may fail if the Go type has
// attributes whose types cannot be converted to Tamarin types.
func NewProxy(reg GoTypeRegistry, obj interface{}) (*Proxy, error) {
	goType, err := reg.Register(obj)
	if err != nil {
		return nil, err
	}
	return &Proxy{reg: reg, typ: goType, obj: obj}, nil
}
