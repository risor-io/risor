package object

import (
	"fmt"
	"reflect"

	"github.com/cloudcmds/tamarin/v2/op"
)

// GoField represents a single field on a Go type that can be read or written.
type GoField struct {
	*base
	field     reflect.StructField
	fieldType *GoType
	name      *String
	tag       *String
	converter TypeConverter
}

func (f *GoField) Name() string {
	return f.field.Name
}

func (f *GoField) ReflectType() reflect.Type {
	return f.field.Type
}

func (f *GoField) GoType() *GoType {
	return f.fieldType
}

func (f *GoField) Tag() reflect.StructTag {
	return f.field.Tag
}

func (f *GoField) Type() Type {
	return GO_TYPE
}

func (f *GoField) Inspect() string {
	return fmt.Sprintf("go_field(%s)", f.Name())
}

func (f *GoField) Interface() interface{} {
	return f.field
}

func (f *GoField) Equals(other Object) Object {
	if f == other {
		return True
	}
	return False
}

func (f *GoField) GetAttr(name string) (Object, bool) {
	switch name {
	case "name":
		return f.name, true
	case "type":
		return f.fieldType, true
	case "tag":
		return f.tag, true
	}
	return nil, false
}

func (f *GoField) IsTruthy() bool {
	return true
}

func (f *GoField) RunOperation(opType op.BinaryOpType, right Object) Object {
	return Errorf("type error: unsupported operation on go_field (%s)", opType)
}

func (f *GoField) Converter() (TypeConverter, bool) {
	return f.converter, f.converter != nil
}

func newGoField(f reflect.StructField) (*GoField, error) {
	fieldGoType, err := newGoType(f.Type)
	if err != nil {
		return nil, err
	}
	conv := kindConverters[f.Type.Kind()]
	return &GoField{
		field:     f,
		fieldType: fieldGoType,
		name:      NewString(f.Name),
		tag:       NewString(string(f.Tag)),
		converter: conv,
	}, nil
}
