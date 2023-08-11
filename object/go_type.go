package object

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"github.com/risor-io/risor/op"
)

// GoType wraps a single native Go type to make it easier to work with in Risor
// and also to be able to represent the type as a Risor object.
type GoType struct {
	*base
	typ            reflect.Type
	name           *String
	packagePath    *String
	attributes     map[string]GoAttribute
	attributeNames []string
	indirectType   *GoType
	converter      TypeConverter
	isPointerType  bool
	isDirectMethod map[string]bool
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
	case "is_pointer_type":
		return NewBool(t.isPointerType), true
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

func (t *GoType) IsPointerType() bool {
	return t.isPointerType
}

func (t *GoType) IndirectType() *GoType {
	return t.indirectType
}

func (t *GoType) ValueType() *GoType {
	if t.isPointerType {
		return t.indirectType
	}
	return t
}

func (t *GoType) PointerType() *GoType {
	if t.isPointerType {
		return t
	}
	return t.indirectType
}

func (t *GoType) New() reflect.Value {
	return reflect.New(t.ValueType().typ)
}

func (t *GoType) HasDirectMethod(name string) bool {
	return t.isDirectMethod[name]
}

func (t *GoType) GetConverter() (TypeConverter, error) {
	if t.converter != nil {
		return t.converter, nil
	}
	conv, err := getTypeConverter(t.typ)
	if err != nil {
		return nil, err
	}
	t.converter = conv
	return conv, nil
}

func (t *GoType) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name          string            `json:"name"`
		Package       string            `json:"package"`
		Attributes    map[string]Object `json:"attributes"`
		IsPointerType bool              `json:"is_pointer_type"`
	}{
		Name:          t.name.value,
		Package:       t.packagePath.value,
		Attributes:    t.attrMap(),
		IsPointerType: t.isPointerType,
	})
}

// newGoType creates and registers a new GoType for the type of the given object.
// This is NOT threadsafe. The caller must be holding goTypeMutex.
func newGoType(typ reflect.Type) (*GoType, error) {

	// Return the existing type if it's already registered
	if goType, ok := goTypeRegistry[typ]; ok {
		return goType, nil
	}

	// Just like Go does, we want to provide some equivalence between a type and
	// a pointer to that type. The "indirect type" is the opposite form from
	// what we're given.
	kind := typ.Kind()
	isPointer := kind == reflect.Ptr
	var indirectType reflect.Type
	var indirectKind reflect.Kind
	if isPointer {
		// Get the corresponding non-pointer type
		indirectType = typ.Elem()
		indirectKind = indirectType.Kind()
	} else {
		// Get the corresponding pointer type
		indirectType = reflect.PtrTo(typ)
		indirectKind = reflect.Ptr
	}

	name := typ.Name()
	if name == "" {
		name = typ.String()
	}

	// Add the new type to the registry
	goType := &GoType{
		attributes:     map[string]GoAttribute{},
		typ:            typ,
		name:           NewString(name),
		packagePath:    NewString(typ.PkgPath()),
		isPointerType:  isPointer,
		isDirectMethod: map[string]bool{},
	}

	// Add the new type to the registry before calling newGoType recursively
	goTypeRegistry[typ] = goType

	// Register the indirect type as well (recursive call!)
	indirectGoType, err := newGoType(indirectType)
	if err != nil {
		return nil, err
	}
	goType.indirectType = indirectGoType

	// If this is a struct, discover all its exported fields
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
				return nil, err
			}
			goType.attributes[field.Name] = goField
		}
	}

	// Discover methods on the indirect type
	indirectMethods, err := getMethods(indirectType)
	if err != nil {
		return nil, err
	}
	for name, method := range indirectMethods {
		goType.attributes[name] = method
	}

	// Discover methods on the direct type (higher precedence)
	directMethods, err := getMethods(typ)
	if err != nil {
		return nil, err
	}
	for name, method := range directMethods {
		goType.attributes[name] = method
		goType.isDirectMethod[name] = true
	}

	// Now that all attributes have been discovered, create a sorted list of
	// attribute names for use in the proxy.
	for attrName := range goType.attributes {
		goType.attributeNames = append(goType.attributeNames, attrName)
	}
	sort.Strings(goType.attributeNames)
	return goType, nil
}

// NewGoType registers and returns a Risor GoType for the type of the given
// native Go object. This is safe for concurrent use by multiple goroutines.
// A type registry is maintained behind the scenes to ensure that each type
// is only registered once.
func NewGoType(typ reflect.Type) (*GoType, error) {
	goTypeMutex.Lock()
	defer goTypeMutex.Unlock()

	return newGoType(typ)
}
