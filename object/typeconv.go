package object

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"time"

	"github.com/gofrs/uuid"
)

var kindConverters = map[reflect.Kind]TypeConverter{
	reflect.Bool:    &BoolConverter{},
	reflect.Int:     &IntConverter{},
	reflect.Int8:    &Int8Converter{},
	reflect.Int16:   &Int16Converter{},
	reflect.Int32:   &Int32Converter{},
	reflect.Int64:   &Int64Converter{},
	reflect.Uint:    &UintConverter{},
	reflect.Uint8:   &Uint8Converter{},
	reflect.Uint16:  &Uint16Converter{},
	reflect.Uint32:  &Uint32Converter{},
	reflect.Uint64:  &Uint64Converter{},
	reflect.Float32: &Float32Converter{},
	reflect.Float64: &Float64Converter{},
	reflect.String:  &StringConverter{},
}

var typeConverters = map[reflect.Type]TypeConverter{}

// Outside of those, we want to handle:
// * Array
// * Func
// * Interface (Partial?)
// * Map (Partial?)

// Kinds do NOT intend to handle for now:
// * Chan
// * Complex64
// * Complex128
// * UnsafePointer

// *****************************************************************************
// Type assertion helpers
// *****************************************************************************

func AsBool(obj Object) (bool, *Error) {
	b, ok := obj.(*Bool)
	if !ok {
		return false, Errorf("type error: expected a bool (%s given)", obj.Type())
	}
	return b.value, nil
}

func AsString(obj Object) (string, *Error) {
	switch obj := obj.(type) {
	case *String:
		return obj.value, nil
	case *BSlice:
		return string(obj.value), nil
	case *Buffer:
		return obj.value.String(), nil
	default:
		return "", Errorf("type error: expected a string (%s given)", obj.Type())
	}
}

func AsInt(obj Object) (int64, *Error) {
	i, ok := obj.(*Int)
	if !ok {
		return 0, Errorf("type error: expected an integer (%s given)", obj.Type())
	}
	return i.value, nil
}

func AsFloat(obj Object) (float64, *Error) {
	switch obj := obj.(type) {
	case *Int:
		return float64(obj.value), nil
	case *Float:
		return obj.value, nil
	default:
		return 0.0, Errorf("type error: expected a number (%s given)", obj.Type())
	}
}

func AsList(obj Object) (*List, *Error) {
	list, ok := obj.(*List)
	if !ok {
		return nil, Errorf("type error: expected a list (%s given", obj.Type())
	}
	return list, nil
}

func AsMap(obj Object) (*Map, *Error) {
	m, ok := obj.(*Map)
	if !ok {
		return nil, Errorf("type error: expected a map (%s given", obj.Type())
	}
	return m, nil
}

func AsTime(obj Object) (result time.Time, err *Error) {
	s, ok := obj.(*Time)
	if !ok {
		return time.Time{}, Errorf("type error: expected a time (%s given", obj.Type())
	}
	return s.value, nil
}

func AsSet(obj Object) (*Set, *Error) {
	set, ok := obj.(*Set)
	if !ok {
		return nil, Errorf("type error: expected a set (%s given", obj.Type())
	}
	return set, nil
}

func AsBytes(obj Object) ([]byte, *Error) {
	switch obj := obj.(type) {
	case *BSlice:
		return obj.value, nil
	case *Buffer:
		return obj.value.Bytes(), nil
	case *String:
		return []byte(obj.value), nil
	default:
		return nil, Errorf("type error: expected bytes (%s given)", obj.Type())
	}
}

func AsReader(obj Object) (io.Reader, *Error) {
	switch obj := obj.(type) {
	case *BSlice:
		return bytes.NewBuffer(obj.value), nil
	case *String:
		return bytes.NewBufferString(obj.value), nil
	case *File:
		return obj.value, nil
	case *HttpResponse:
		return obj.resp.Body, nil
	case io.Reader:
		return obj, nil
	default:
		return nil, Errorf("type error: expected a readable object (%s given)", obj.Type())
	}
}

// *****************************************************************************
// Converting types from Go to Tamarin
// *****************************************************************************

func FromGoType(obj interface{}) Object {
	switch obj := obj.(type) {
	case nil:
		return Nil
	case int:
		return NewInt(int64(obj))
	case int16:
		return NewInt(int64(obj))
	case int32:
		return NewInt(int64(obj))
	case int64:
		return NewInt(obj)
	case uint:
		return NewInt(int64(obj))
	case uint16:
		return NewInt(int64(obj))
	case uint32:
		return NewInt(int64(obj))
	case uint64:
		return NewInt(int64(obj))
	case float32:
		return NewFloat(float64(obj))
	case float64:
		return NewFloat(obj)
	case string:
		return NewString(obj)
	case byte:
		return NewInt(int64(obj))
	case []byte:
		return NewBSlice(obj)
	case *bytes.Buffer:
		return NewBuffer(obj.Bytes())
	case bool:
		if obj {
			return True
		}
		return False
	case [16]uint8:
		return NewString(uuid.UUID(obj).String())
	case time.Time:
		return NewTime(obj)
	case []interface{}:
		items := make([]Object, 0, len(obj))
		for _, item := range obj {
			listItem := FromGoType(item)
			if IsError(listItem) {
				return listItem
			}
			items = append(items, listItem)
		}
		return NewList(items)
	case map[string]interface{}:
		m := make(map[string]Object, len(obj))
		for k, v := range obj {
			valueObj := FromGoType(v)
			if IsError(valueObj) {
				return valueObj
			}
			m[k] = valueObj
		}
		return NewMap(m)
	default:
		return Errorf("type error: unmarshaling %v (%v)",
			obj, reflect.TypeOf(obj))
	}
}

// *****************************************************************************
// TypeConverter interface and implementations
//   - These are applicable when the Go type(s) are known in advance
// *****************************************************************************

// TypeConverter is an interface used to convert between Go and Tamarin objects
// for a single Go type.
type TypeConverter interface {

	// To converts to a Go object from a Tamarin object.
	To(Object) (interface{}, error)

	// From converts a Go object to a Tamarin object.
	From(interface{}) (Object, error)
}

// GetConverter returns a TypeConverter for the given Go kind and type.
// Converters are cached internally for reuse.
func GetConverter(kind reflect.Kind, typ reflect.Type) (TypeConverter, error) {
	goTypeMutex.Lock()
	defer goTypeMutex.Unlock()

	return getConverter(kind, typ)
}

func getConverter(kind reflect.Kind, typ reflect.Type) (TypeConverter, error) {
	if conv, ok := kindConverters[kind]; ok {
		return conv, nil
	}
	if conv, ok := typeConverters[typ]; ok {
		return conv, nil
	}
	var err error
	var converter TypeConverter
	switch kind {
	case reflect.Struct:
		converter, err = NewStructConverter(typ)
		if err != nil {
			return nil, err
		}
	case reflect.Pointer:
		converter, err = NewPointerConverter(typ.Elem())
		if err != nil {
			return nil, err
		}
	case reflect.Slice:
		converter, err = NewSliceConverter(typ.Elem())
		if err != nil {
			return nil, err
		}
	case reflect.Map:
		if typ.Key().Kind() == reflect.String {
			converter, err = NewMapConverter(typ.Elem())
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("type error: unsupported map key type in %s", typ)
		}
	case reflect.Interface:
		if typ.Implements(errorInterface) {
			converter = &ErrorConverter{}
		} else if typ.Implements(contextInterface) {
			converter = &ContextConverter{}
		} else { // TODO: io.*?
			return nil, fmt.Errorf("type error: unsupported interface type %s", typ)
		}
	default:
		return nil, fmt.Errorf("type error: unsupported type %s", typ)
	}
	typeConverters[typ] = converter
	return converter, nil
}

// PointerConverter converts between *T and the Tamarin equivalent of T.
type PointerConverter struct {
	valueConverter TypeConverter
}

func (c *PointerConverter) To(obj Object) (interface{}, error) {
	v, err := c.valueConverter.To(obj)
	if err != nil {
		return nil, err
	}
	vp := reflect.New(reflect.TypeOf(v))
	vp.Elem().Set(reflect.ValueOf(v))
	return vp.Interface(), nil
}

func (c *PointerConverter) From(obj interface{}) (Object, error) {
	v := reflect.ValueOf(obj).Elem().Interface()
	return c.valueConverter.From(v)
}

// NewPointerConverter creates a TypeConverter for pointers that point to
// items of the given type.
func NewPointerConverter(indirectType reflect.Type) (*PointerConverter, error) {
	indirectKind := indirectType.Kind()
	indirectConv, err := getConverter(indirectKind, indirectType)
	if err != nil {
		return nil, err
	}
	return &PointerConverter{valueConverter: indirectConv}, nil
}

// SliceConverter converts between []T and the Tamarin equivalent of []T.
type SliceConverter struct {
	valueConverter TypeConverter
	valueType      reflect.Type
}

func (c *SliceConverter) To(obj Object) (interface{}, error) {
	list := obj.(*List)
	slice := reflect.MakeSlice(reflect.SliceOf(c.valueType), 0, len(list.items))
	for _, v := range list.items {
		item, err := c.valueConverter.To(v)
		if err != nil {
			return nil, fmt.Errorf("type error: failed to convert slice element: %v", err)
		}
		slice = reflect.Append(slice, reflect.ValueOf(item))
	}
	return slice.Interface(), nil
}

func (c *SliceConverter) From(iface interface{}) (Object, error) {
	v := reflect.ValueOf(iface)
	count := v.Len()
	items := make([]Object, 0, count)
	for i := 0; i < count; i++ {
		item, err := c.valueConverter.From(v.Index(i).Interface())
		if err != nil {
			return nil, fmt.Errorf("type error: failed to convert slice element: %v", err)
		}
		items = append(items, item)
	}
	return NewList(items), nil
}

// NewSliceConverter creates a TypeConverter for slices containing the given
// value type, where the items can be converted using the given TypeConverter.
func NewSliceConverter(indirectType reflect.Type) (*SliceConverter, error) {
	indirectKind := indirectType.Kind()
	indirectConv, err := getConverter(indirectKind, indirectType)
	if err != nil {
		return nil, err
	}
	return &SliceConverter{
		valueType:      indirectType,
		valueConverter: indirectConv,
	}, nil
}

// BoolConverter converts between bool and *Bool.
type BoolConverter struct{}

func (c *BoolConverter) To(obj Object) (interface{}, error) {
	return obj.(*Bool).value, nil
}

func (c *BoolConverter) From(obj interface{}) (Object, error) {
	return NewBool(obj.(bool)), nil
}

// IntConverter converts between int and *Int.
type IntConverter struct{}

func (c *IntConverter) To(obj Object) (interface{}, error) {
	return int(obj.(*Int).value), nil
}

func (c *IntConverter) From(obj interface{}) (Object, error) {
	return NewInt(int64(obj.(int))), nil
}

// Int8Converter converts between int8 and *Int.
type Int8Converter struct{}

func (c *Int8Converter) To(obj Object) (interface{}, error) {
	return int8(obj.(*Int).value), nil
}

func (c *Int8Converter) From(obj interface{}) (Object, error) {
	return NewInt(int64(obj.(int8))), nil
}

// Int16Converter converts between int16 and *Int.
type Int16Converter struct{}

func (c *Int16Converter) To(obj Object) (interface{}, error) {
	return int16(obj.(*Int).value), nil
}

func (c *Int16Converter) From(obj interface{}) (Object, error) {
	return NewInt(int64(obj.(int16))), nil
}

// Int32Converter converts between int32 and *Int.
type Int32Converter struct{}

func (c *Int32Converter) To(obj Object) (interface{}, error) {
	return int32(obj.(*Int).value), nil
}

func (c *Int32Converter) From(obj interface{}) (Object, error) {
	return NewInt(int64(obj.(int32))), nil
}

// Int64Converter converts between int64 and *Int.
type Int64Converter struct{}

func (c *Int64Converter) To(obj Object) (interface{}, error) {
	return obj.(*Int).value, nil
}

func (c *Int64Converter) From(obj interface{}) (Object, error) {
	return NewInt(obj.(int64)), nil
}

// UintConverter converts between uint and *Int.
type UintConverter struct{}

func (c *UintConverter) To(obj Object) (interface{}, error) {
	return uint(obj.(*Int).value), nil
}

func (c *UintConverter) From(obj interface{}) (Object, error) {
	return NewInt(int64(obj.(uint))), nil
}

// Uint8Converter converts between uint8 and *Int.
type Uint8Converter struct{}

func (c *Uint8Converter) To(obj Object) (interface{}, error) {
	return uint8(obj.(*Int).value), nil
}

func (c *Uint8Converter) From(obj interface{}) (Object, error) {
	return NewInt(int64(obj.(uint8))), nil
}

// Uint16Converter converts between uint16 and *Int.
type Uint16Converter struct{}

func (c *Uint16Converter) To(obj Object) (interface{}, error) {
	return uint16(obj.(*Int).value), nil
}

func (c *Uint16Converter) From(obj interface{}) (Object, error) {
	return NewInt(int64(obj.(uint16))), nil
}

// Uint32Converter converts between uint32 and *Int.
type Uint32Converter struct{}

func (c *Uint32Converter) To(obj Object) (interface{}, error) {
	return uint32(obj.(*Int).value), nil
}

func (c *Uint32Converter) From(obj interface{}) (Object, error) {
	return NewInt(int64(obj.(uint32))), nil
}

// Uint64Converter converts between uint64 and *Int.
type Uint64Converter struct{}

func (c *Uint64Converter) To(obj Object) (interface{}, error) {
	v := obj.(*Int).value
	if v < 0 {
		return nil, fmt.Errorf("value error: %d is out of range for uint64", v)
	}
	return uint64(obj.(*Int).value), nil
}

func (c *Uint64Converter) From(obj interface{}) (Object, error) {
	v := obj.(uint64)
	if v > math.MaxInt64 {
		return nil, fmt.Errorf("value error: %d is out of range for int64", v)
	}
	return NewInt(int64(v)), nil
}

// Float32Converter converts between float32 and *Float.
type Float32Converter struct{}

func (c *Float32Converter) To(obj Object) (interface{}, error) {
	return float32(obj.(*Float).value), nil
}

func (c *Float32Converter) From(obj interface{}) (Object, error) {
	return NewFloat(float64(obj.(float32))), nil
}

// Float64Converter converts between float64 and *Float.
type Float64Converter struct{}

func (c *Float64Converter) To(obj Object) (interface{}, error) {
	return obj.(*Float).value, nil
}

func (c *Float64Converter) From(obj interface{}) (Object, error) {
	return NewFloat(obj.(float64)), nil
}

// StringConverter converts between string and *String.
type StringConverter struct{}

func (c *StringConverter) To(obj Object) (interface{}, error) {
	return obj.(*String).value, nil
}

func (c *StringConverter) From(obj interface{}) (Object, error) {
	return NewString(obj.(string)), nil
}

// BSliceConverter converts between []byte and BSlice.
type BSliceConverter struct{}

func (c *BSliceConverter) To(obj Object) (interface{}, error) {
	return obj.(*BSlice).value, nil
}

func (c *BSliceConverter) From(obj interface{}) (Object, error) {
	return NewBSlice(obj.([]byte)), nil
}

// TimeConverter converts between time.Time and Time.
type TimeConverter struct{}

func (c *TimeConverter) To(obj Object) (interface{}, error) {
	return obj.(*Time).value, nil
}

func (c *TimeConverter) From(obj interface{}) (Object, error) {
	return NewTime(obj.(time.Time)), nil
}

// MapConverter converts between map[string]interface{} and *Map.
type MapConverter struct {
	valueConverter TypeConverter
	valueType      reflect.Type
}

func (c *MapConverter) To(obj Object) (interface{}, error) {
	tMap := obj.(*Map)
	keyType := reflect.TypeOf("")
	mapType := reflect.MapOf(keyType, c.valueType)
	gMap := reflect.MakeMapWithSize(mapType, tMap.Size())
	for k, v := range tMap.items {
		conv, err := c.valueConverter.To(v)
		if err != nil {
			return nil, err
		}
		gMap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(conv))
	}
	return gMap.Interface(), nil
}

func (c *MapConverter) From(obj interface{}) (Object, error) {
	m := reflect.ValueOf(obj)
	o := make(map[string]Object, m.Len())
	for _, key := range m.MapKeys() {
		v := m.MapIndex(key)
		conv, err := c.valueConverter.From(v.Interface())
		if err != nil {
			return nil, err
		}
		o[key.Interface().(string)] = conv
	}
	return NewMap(o), nil
}

func NewMapConverter(valueType reflect.Type) (*MapConverter, error) {
	valueConverter, err := getConverter(valueType.Kind(), valueType)
	if err != nil {
		return nil, fmt.Errorf("type error: unsupported map value type %s", valueType)
	}
	return &MapConverter{
		valueConverter: valueConverter,
		valueType:      valueType,
	}, nil
}

// StructConverter converts between a Go struct and a Tamarin Map.
type StructConverter struct {
	typ             reflect.Type
	fieldConverters []TypeConverter
	fieldNames      []string
}

func (c *StructConverter) To(obj Object) (interface{}, error) {
	m, ok := obj.(*Map)
	if !ok {
		return nil, fmt.Errorf("type error: expected a map (%s given)", obj.Type())
	}
	v := reflect.New(c.typ).Elem()
	for i := 0; i < v.NumField(); i++ {
		attrName := c.fieldNames[i]
		if item, ok := m.items[attrName]; ok {
			if f := v.Field(i); f.CanSet() {
				conv := c.fieldConverters[i]
				attrValue, err := conv.To(item)
				if err != nil {
					return nil, err
				}
				f.Set(reflect.ValueOf(attrValue))
			}
		}
	}
	return v.Interface(), nil
}

func (c *StructConverter) From(obj interface{}) (Object, error) {
	items := map[string]Object{}
	objValue := reflect.ValueOf(obj)
	for i, conv := range c.fieldConverters {
		f := objValue.Field(i)
		if !f.CanInterface() {
			continue
		}
		item, err := conv.From(f.Interface())
		if err != nil {
			return nil, err
		}
		items[c.fieldNames[i]] = item
	}
	return NewMap(items), nil
}

// NewStructConverter creates a TypeConverter for a given type of struct.
func NewStructConverter(typ reflect.Type) (*StructConverter, error) {
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type error: expected a struct (%s given)", typ)
	}
	numField := typ.NumField()
	fieldConverters := make([]TypeConverter, numField)
	fieldNames := make([]string, numField)
	for i := 0; i < numField; i++ {
		field := typ.Field(i)
		fieldType := field.Type
		fieldNames[i] = field.Name
		conv, err := getConverter(fieldType.Kind(), fieldType)
		if err != nil {
			return nil, err
		}
		fieldConverters[i] = conv
	}
	return &StructConverter{
		typ:             typ,
		fieldConverters: fieldConverters,
		fieldNames:      fieldNames,
	}, nil
}

// ErrorConverter converts between error and *Error.
type ErrorConverter struct{}

func (c *ErrorConverter) To(obj Object) (interface{}, error) {
	return obj.(*Error).err, nil
}

func (c *ErrorConverter) From(obj interface{}) (Object, error) {
	return NewError(obj.(error)), nil
}

// ContextConverter converts between context.Context and Context.
type ContextConverter struct{}

func (c *ContextConverter) To(obj Object) (interface{}, error) {
	// Not actually called, but needed to satisfy the Converter interface.
	return nil, errors.New("not implemented")
}

func (c *ContextConverter) From(obj interface{}) (Object, error) {
	// Not actually called, but needed to satisfy the Converter interface.
	return nil, errors.New("not implemented")
}
