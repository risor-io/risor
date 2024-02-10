package cli

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	ucli "github.com/urfave/cli/v2"
)

const FLAG object.Type = "cli.flag"

type Flag struct {
	value ucli.Flag
}

func (f *Flag) Type() object.Type {
	return "cli.flag"
}

func (f *Flag) Inspect() string {
	return "cli.flag"
}

func (f *Flag) Interface() interface{} {
	return f.value
}

func (f *Flag) IsTruthy() bool {
	return true
}

func (f *Flag) Cost() int {
	return 0
}

func (f *Flag) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("type error: unable to marshal %s", FLAG)
}

func (f *Flag) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("eval error: unsupported operation for %s: %v", FLAG, opType)
}

func (f *Flag) Equals(other object.Object) object.Object {
	if other.Type() != "cli.flag" {
		return object.False
	}
	return object.NewBool(f.value == other.(*Flag).value)
}

func (f *Flag) SetAttr(name string, value object.Object) error {
	if err := setNamedField(f.value, name, value); err != nil {
		return err
	}
	return nil
}

func (f *Flag) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "names":
		return object.NewStringList(f.value.Names()), true
	case "is_set":
		return object.NewBool(f.value.IsSet()), true
	case "string":
		return object.NewString(f.value.String()), true
	default:
		field, err := getNamedField(f.value, name)
		if err != nil {
			return nil, false
		}
		return field, true
	}
}

func NewFlag(f ucli.Flag) *Flag {
	return &Flag{value: f}
}

func getNamedField(i interface{}, name string) (object.Object, error) {
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %v", val.Kind())
	}
	name = snakeToCap(name)
	fieldVal := val.FieldByName(name)
	if !fieldVal.IsValid() {
		return nil, fmt.Errorf("no such field: %s in obj", name)
	}
	obj := object.FromGoType(fieldVal.Interface())
	if errObj, ok := obj.(*object.Error); ok {
		return nil, errObj.Value()
	}
	return obj, nil
}

func setNamedField(i interface{}, name string, value object.Object) error {
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got %v", val.Kind())
	}
	name = snakeToCap(name)
	fieldVal := val.FieldByName(name)
	if !fieldVal.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}
	if !fieldVal.CanSet() {
		return fmt.Errorf("cannot set field: %s in obj", name)
	}
	fieldVal.Set(reflect.ValueOf(value.Interface()))
	return nil
}

func snakeToCap(s string) string {
	var result string
	words := strings.Split(s, "_")
	for _, word := range words {
		result += strings.Title(word)
	}
	return result
}
