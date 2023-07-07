package object

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"github.com/risor-io/risor/op"
)

// GoType represents a single Go type whose methods and fields can be proxied.
type GoType struct {
	*base
	typ            reflect.Type
	name           *String
	packagePath    *String
	attributes     map[string]GoAttribute
	attributeNames []string
	indirectType   *GoType
	indirectKind   reflect.Kind
	warnings       []error
	converter      TypeConverter
}

func (t *GoType) Type() Type {
	return GO_TYPE
}

func (t *GoType) Inspect() string {
	return fmt.Sprintf("go_type(%s)", t.Name())
}

func (t *GoType) Interface() interface{} {
	return t.typ
}

func (t *GoType) Equals(other Object) Object {
	if t == other {
		return True
	}
	return False
}

func (t *GoType) GetAttr(name string) (Object, bool) {
	switch name {
	case "name":
		return t.name, true
	case "package_path":
		return t.packagePath, true
	case "attributes":
		return NewMap(t.attrMap()), true
	}
	return nil, false
}

func (t *GoType) attrMap() map[string]Object {
	attrs := make(map[string]Object, len(t.attributes))
	for name, attr := range t.attributes {
		switch attr := attr.(type) {
		case *GoMethod:
			attrs[name] = attr
		case *GoField:
			attrs[name] = attr
		}
	}
	return attrs
}

func (t *GoType) IsTruthy() bool {
	return true
}

func (t *GoType) RunOperation(opType op.BinaryOpType, right Object) Object {
	return Errorf("type error: unsupported operation on go_type (%s)", opType)
}

func (t *GoType) Name() string {
	return t.name.value
}

func (t *GoType) PackagePath() string {
	return t.packagePath.value
}

func (t *GoType) AttributeNames() []string {
	return t.attributeNames
}

func (t *GoType) GetAttribute(name string) (GoAttribute, bool) {
	attr, ok := t.attributes[name]
	return attr, ok
}

func (t *GoType) ReflectType() reflect.Type {
	return t.typ
}

func (t *GoType) Warnings() []error {
	return t.warnings
}

func (t *GoType) IndirectType() (*GoType, bool) {
	return t.indirectType, t.indirectType != nil
}

func (t *GoType) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name":       t.name,
		"package":    t.packagePath,
		"attributes": t.attrMap(),
	})
}

// newGoType creates and registers a new GoType for the type of the given object.
// This is NOT threadsafe. The caller must be holding goTypeMutex.
func newGoType(typ reflect.Type) (*GoType, error) {

	// Is this type already registered?
	kind := typ.Kind()
	if goType, ok := goTypeRegistry[typ]; ok {
		return goType, nil
	}

	isPointer := kind == reflect.Ptr

	var indirectType reflect.Type
	var indirectKind reflect.Kind
	if isPointer {
		indirectType = typ.Elem()
		indirectKind = indirectType.Kind()
	}

	name := typ.Name()
	if name == "" {
		name = typ.String()
	}

	conv, err := getTypeConverter(typ)
	if err != nil {
		return nil, err
	}

	// Add the new type to the registry
	goType := &GoType{
		attributes:   map[string]GoAttribute{},
		typ:          typ,
		indirectKind: indirectKind,
		name:         NewString(name),
		packagePath:  NewString(typ.PkgPath()),
		converter:    conv,
	}
	goTypeRegistry[typ] = goType

	// Create/lookup the indirect type if there is one
	if indirectType != nil {
		indirectGoType, err := newGoType(indirectType)
		if err != nil {
			return nil, err
		}
		goType.indirectType = indirectGoType
	}

	// Discover and register the type of each field for struct types
	if kind == reflect.Struct || indirectKind == reflect.Struct {
		structType := typ
		if isPointer {
			structType = typ.Elem()
		}
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			if !field.IsExported() {
				continue
			}
			goField, err := newGoField(field)
			if err != nil {
				goType.warnings = append(goType.warnings, err)
				continue
			}
			goType.attributes[field.Name] = goField
		}
	}

	// Discover methods and register the types of their inputs and outputs
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if !method.IsExported() {
			continue
		}
		methodGoType, err := newGoMethod(method)
		if err != nil {
			goType.warnings = append(goType.warnings, err)
			continue
		}
		goType.attributes[method.Name] = methodGoType
	}

	// Now that all attributes have been discovered, create a sorted list of
	// attribute names for use in the proxy.
	for attrName := range goType.attributes {
		goType.attributeNames = append(goType.attributeNames, attrName)
	}
	sort.Strings(goType.attributeNames)

	return goType, nil
}

// NewGoType creates and registers a new GoType for the type of the given Go object.
// This is safe for use by multiple goroutines. A type registry is maintained
// behind the scenes to ensure that each type is only registered once.
func NewGoType(typ reflect.Type) (*GoType, error) {
	goTypeMutex.Lock()
	defer goTypeMutex.Unlock()

	return newGoType(typ)
}
