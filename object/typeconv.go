package object

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gofrs/uuid"
)

// *****************************************************************************
// Type assertion helpers
// *****************************************************************************

func AsString(obj Object) (result string, err *Error) {
	s, ok := obj.(*String)
	if !ok {
		return "", NewError("type error: expected a string (got %v)", obj.Type())
	}
	return s.Value, nil
}

func AsInteger(obj Object) (int64, *Error) {
	i, ok := obj.(*Integer)
	if !ok {
		return 0, NewError("type error: expected an integer (got %v)", obj.Type())
	}
	return i.Value, nil
}

func AsFloat(obj Object) (float64, *Error) {
	switch obj := obj.(type) {
	case *Integer:
		return float64(obj.Value), nil
	case *Float:
		return obj.Value, nil
	default:
		return 0.0, NewError("type error: expected a number (got %v)", obj.Type())
	}
}

func AsList(obj Object) (*List, *Error) {
	arr, ok := obj.(*List)
	if !ok {
		return nil, NewError("type error: expected an array (got %v)", obj.Type())
	}
	return arr, nil
}

func AsHash(obj Object) (*Hash, *Error) {
	hash, ok := obj.(*Hash)
	if !ok {
		return nil, NewError("type error: expected a hash (got %v)", obj.Type())
	}
	return hash, nil
}

func AsTime(obj Object) (result time.Time, err *Error) {
	s, ok := obj.(*Time)
	if !ok {
		return time.Time{}, NewError("type error: expected a time (got %v)", obj.Type())
	}
	return s.Value, nil
}

func AsSet(obj Object) (*Set, *Error) {
	set, ok := obj.(*Set)
	if !ok {
		return nil, NewError("type error: expected a set (got %v)", obj.Type())
	}
	return set, nil
}

// *****************************************************************************
// Converting types from Go to Tamarin
// *****************************************************************************

func FromGoType(obj interface{}) Object {
	switch obj := obj.(type) {
	case nil:
		return NULL
	case int:
		return &Integer{Value: int64(obj)}
	case int32:
		return &Integer{Value: int64(obj)}
	case int64:
		return &Integer{Value: obj}
	case float32:
		return &Float{Value: float64(obj)}
	case float64:
		return &Float{Value: obj}
	case string:
		return &String{Value: obj}
	case bool:
		if obj {
			return TRUE
		}
		return FALSE
	case [16]uint8:
		return &String{Value: uuid.UUID(obj).String()}
	case time.Time:
		return &Time{Value: obj}
	case []interface{}:
		ls := &List{Items: make([]Object, 0, len(obj))}
		for _, item := range obj {
			listItem := FromGoType(item)
			if listItem == nil {
				return nil
			}
			ls.Items = append(ls.Items, listItem)
		}
		return ls
	case map[string]interface{}:
		hash := &Hash{Map: make(map[string]Object, len(obj))}
		for k, v := range obj {
			hashVal := FromGoType(v)
			if hashVal == nil {
				return nil
			}
			hash.Map[k] = hashVal
		}
		return hash
	default:
		return NewError("type error: unmarshaling %v (%v)",
			obj, reflect.TypeOf(obj))
	}
}

// *****************************************************************************
// Converting types Tamarin to Go
// *****************************************************************************

func ToGoType(obj Object) interface{} {
	switch obj := obj.(type) {
	case *Null:
		return nil
	case *Integer:
		return obj.Value
	case *Float:
		return obj.Value
	case *String:
		return obj.Value
	case *Boolean:
		return obj.Value
	case *Time:
		return obj.Value
	case *List:
		ls := make([]interface{}, 0, len(obj.Items))
		for _, item := range obj.Items {
			ls = append(ls, ToGoType(item))
		}
		return ls
	case *Set:
		array := make([]interface{}, 0, len(obj.Items))
		for _, item := range obj.Items {
			array = append(array, ToGoType(item))
		}
		return array
	case *Hash:
		m := make(map[string]interface{}, len(obj.Map))
		for k, v := range obj.Map {
			m[k] = ToGoType(v)
		}
		return m
	default:
		return NewError("type error: marshaling %v (%v)",
			obj, reflect.TypeOf(obj))
	}
}

// *****************************************************************************
// TypeConverter interface and implementations
//   - These are applicable when the Go type(s) are known in advance
// *****************************************************************************

// TypeConverter is an interface used to convert between Go and Tamarin objects
// for a single Go type. There may be a way to use generics here...
type TypeConverter interface {

	// To converts a Tamarin object to a Go object.
	To(Object) (interface{}, error)

	// From converts a Go object to a Tamarin object.
	From(interface{}) (Object, error)

	// Type that this TypeConverter is responsbile for.
	Type() reflect.Type
}

var (
	intType         = reflect.TypeOf(int(0))
	int64Type       = reflect.TypeOf(int64(0))
	stringType      = reflect.TypeOf("")
	float32Type     = reflect.TypeOf(float32(0))
	float64Type     = reflect.TypeOf(float64(0))
	booleanType     = reflect.TypeOf(false)
	timeType        = reflect.TypeOf(time.Time{})
	mapStrIfaceType = reflect.TypeOf(map[string]interface{}{})
	errType         = reflect.TypeOf(errors.New(""))
)

// Int64Converter converts between int64 and Integer.
type Int64Converter struct{}

func (c *Int64Converter) To(obj Object) (interface{}, error) {
	integer, ok := obj.(*Integer)
	if !ok {
		return nil, fmt.Errorf("type error: expected an integer (got %v)", obj.Type())
	}
	return integer.Value, nil
}

func (c *Int64Converter) From(obj interface{}) (Object, error) {
	return NewInteger(obj.(int64)), nil
}

func (c *Int64Converter) Type() reflect.Type {
	return int64Type
}

// IntConverter converts between int and Integer.
type IntConverter struct{}

func (c *IntConverter) To(obj Object) (interface{}, error) {
	integer, ok := obj.(*Integer)
	if !ok {
		return nil, fmt.Errorf("type error: expected an integer (got %v)", obj.Type())
	}
	return int(integer.Value), nil
}

func (c *IntConverter) From(obj interface{}) (Object, error) {
	return NewInteger(int64(obj.(int))), nil
}

func (c *IntConverter) Type() reflect.Type {
	return intType
}

// StringConverter converts between string and String.
type StringConverter struct{}

func (c *StringConverter) To(obj Object) (interface{}, error) {
	s, ok := obj.(*String)
	if !ok {
		return nil, fmt.Errorf("type error: expected a string (got %v)", obj.Type())
	}
	return s.Value, nil
}

func (c *StringConverter) From(obj interface{}) (Object, error) {
	return NewString(obj.(string)), nil
}

func (c *StringConverter) Type() reflect.Type {
	return stringType
}

// Float64Converter converts between float64 and Float.
type Float64Converter struct{}

func (c *Float64Converter) To(obj Object) (interface{}, error) {
	f, ok := obj.(*Float)
	if !ok {
		return nil, fmt.Errorf("type error: expected a float (got %v)", obj.Type())
	}
	return f.Value, nil
}

func (c *Float64Converter) From(obj interface{}) (Object, error) {
	return NewFloat(obj.(float64)), nil
}

func (c *Float64Converter) Type() reflect.Type {
	return float64Type
}

// Float32Converter converts between float32 and Float.
type Float32Converter struct{}

func (c *Float32Converter) To(obj Object) (interface{}, error) {
	f, ok := obj.(*Float)
	if !ok {
		return nil, fmt.Errorf("type error: expected a float (got %v)", obj.Type())
	}
	return float32(f.Value), nil
}

func (c *Float32Converter) From(obj interface{}) (Object, error) {
	return NewFloat(float64(obj.(float32))), nil
}

func (c *Float32Converter) Type() reflect.Type {
	return float32Type
}

// BooleanConverter converts between bool and Boolean.
type BooleanConverter struct{}

func (c *BooleanConverter) To(obj Object) (interface{}, error) {
	b, ok := obj.(*Boolean)
	if !ok {
		return nil, fmt.Errorf("type error: expected a boolean (got %v)", obj.Type())
	}
	return b.Value, nil
}

func (c *BooleanConverter) From(obj interface{}) (Object, error) {
	return NewBoolean(obj.(bool)), nil
}

func (c *BooleanConverter) Type() reflect.Type {
	return booleanType
}

// TimeConverter converts between time.Time and Time.
type TimeConverter struct{}

func (c *TimeConverter) To(obj Object) (interface{}, error) {
	t, ok := obj.(*Time)
	if !ok {
		return nil, fmt.Errorf("type error: expected a time (got %v)", obj.Type())
	}
	return t.Value, nil
}

func (c *TimeConverter) From(obj interface{}) (Object, error) {
	return NewTime(obj.(time.Time)), nil
}

func (c *TimeConverter) Type() reflect.Type {
	return timeType
}

// MapStringIfaceConverter converts between map[string]interface{} and Hash.
type MapStringIfaceConverter struct{}

func (c *MapStringIfaceConverter) To(obj Object) (interface{}, error) {
	hash, ok := obj.(*Hash)
	if !ok {
		return nil, fmt.Errorf("type error: expected a hash (got %v)", obj.Type())
	}
	m := make(map[string]interface{}, len(hash.Map))
	for k, v := range hash.Map {
		m[k] = ToGoType(v)
	}
	return m, nil
}

func (c *MapStringIfaceConverter) From(obj interface{}) (Object, error) {
	m := obj.(map[string]interface{})
	hash := NewHash(make(map[string]interface{}, len(m)))
	for k, v := range m {
		hash.Map[k] = FromGoType(v)
	}
	return hash, nil
}

func (c *MapStringIfaceConverter) Type() reflect.Type {
	return mapStrIfaceType
}

// StructConverter converts between a struct and a Hash via JSON marshaling.
type StructConverter struct {
	Prototype interface{}
}

func (c *StructConverter) To(obj Object) (interface{}, error) {
	hash, ok := obj.(*Hash)
	if !ok {
		return nil, fmt.Errorf("type error: expected a hash (got %v)", obj.Type())
	}
	goMap := ToGoType(hash)
	jsonBytes, err := json.Marshal(goMap)
	if err != nil {
		return nil, err
	}
	inst := reflect.New(reflect.TypeOf(c.Prototype))
	instValue := inst.Interface()
	if err := json.Unmarshal(jsonBytes, &instValue); err != nil {
		return nil, err
	}
	return inst.Elem().Interface(), nil
}

func (c *StructConverter) From(obj interface{}) (Object, error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	var goMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &goMap); err != nil {
		return nil, err
	}
	return FromGoType(goMap), nil
}

func (c *StructConverter) Type() reflect.Type {
	return reflect.TypeOf(c.Prototype)
}

// ErrorConverter converts between error and Error.

type ErrorConverter struct{}

func (c *ErrorConverter) To(obj Object) (interface{}, error) {
	e, ok := obj.(*Error)
	if !ok {
		return nil, fmt.Errorf("type error: expected an error (got %v)", obj.Type())
	}
	return e.Message, nil
}

func (c *ErrorConverter) From(obj interface{}) (Object, error) {
	if obj == nil {
		return nil, nil
	}
	return NewError(obj.(error).Error()), nil
}

func (c *ErrorConverter) Type() reflect.Type {
	return errType
}
