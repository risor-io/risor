package object

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/gofrs/uuid"
)

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
// for a single Go type. There may be a way to use generics here...
type TypeConverter interface {

	// To converts a Tamarin object to a Go object.
	To(Object) (interface{}, error)

	// From converts a Go object to a Tamarin object.
	From(interface{}) (Object, error)

	// Type that this TypeConverter is responsible for.
	Type() reflect.Type
}

var (
	intType         = reflect.TypeOf(int(0))
	int64Type       = reflect.TypeOf(int64(0))
	stringType      = reflect.TypeOf("")
	bytesType       = reflect.TypeOf([]byte{})
	float32Type     = reflect.TypeOf(float32(0))
	float64Type     = reflect.TypeOf(float64(0))
	booleanType     = reflect.TypeOf(false)
	timeType        = reflect.TypeOf(time.Time{})
	mapStrIfaceType = reflect.TypeOf(map[string]interface{}{})
	errType         = reflect.TypeOf(errors.New(""))
	contextType     = reflect.TypeOf(context.Background())
)

// Int64Converter converts between int64 and Integer.
type Int64Converter struct{}

func (c *Int64Converter) To(obj Object) (interface{}, error) {
	integer, ok := obj.(*Int)
	if !ok {
		return nil, fmt.Errorf("type error: expected an integer (%s given)", obj.Type())
	}
	return integer.value, nil
}

func (c *Int64Converter) From(obj interface{}) (Object, error) {
	return NewInt(obj.(int64)), nil
}

func (c *Int64Converter) Type() reflect.Type {
	return int64Type
}

// IntConverter converts between int and Integer.
type IntConverter struct{}

func (c *IntConverter) To(obj Object) (interface{}, error) {
	integer, ok := obj.(*Int)
	if !ok {
		return nil, fmt.Errorf("type error: expected an integer (%s given)", obj.Type())
	}
	return int(integer.value), nil
}

func (c *IntConverter) From(obj interface{}) (Object, error) {
	return NewInt(int64(obj.(int))), nil
}

func (c *IntConverter) Type() reflect.Type {
	return intType
}

// StringConverter converts between string and String.
type StringConverter struct{}

func (c *StringConverter) To(obj Object) (interface{}, error) {
	s, ok := obj.(*String)
	if !ok {
		return nil, fmt.Errorf("type error: expected a string (%s given)", obj.Type())
	}
	return s.value, nil
}

func (c *StringConverter) From(obj interface{}) (Object, error) {
	return NewString(obj.(string)), nil
}

func (c *StringConverter) Type() reflect.Type {
	return stringType
}

// BytesConverter converts between []byte and BSlice.
type BytesConverter struct{}

func (c *BytesConverter) To(obj Object) (interface{}, error) {
	switch obj := obj.(type) {
	case *BSlice:
		return obj.value, nil
	case *String:
		return []byte(obj.value), nil
	default:
		return nil, fmt.Errorf("type error: expected a string (%s given)", obj.Type())
	}
}

func (c *BytesConverter) From(obj interface{}) (Object, error) {
	return NewBSlice(obj.([]byte)), nil
}

func (c *BytesConverter) Type() reflect.Type {
	return bytesType
}

// Float64Converter converts between float64 and Float.
type Float64Converter struct{}

func (c *Float64Converter) To(obj Object) (interface{}, error) {
	f, ok := obj.(*Float)
	if !ok {
		return nil, fmt.Errorf("type error: expected a float (%s given)", obj.Type())
	}
	return f.value, nil
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
		return nil, fmt.Errorf("type error: expected a float (%s given)", obj.Type())
	}
	return float32(f.value), nil
}

func (c *Float32Converter) From(obj interface{}) (Object, error) {
	return NewFloat(float64(obj.(float32))), nil
}

func (c *Float32Converter) Type() reflect.Type {
	return float32Type
}

// BooleanConverter converts between bool and Bool.
type BooleanConverter struct{}

func (c *BooleanConverter) To(obj Object) (interface{}, error) {
	b, ok := obj.(*Bool)
	if !ok {
		return nil, fmt.Errorf("type error: expected a boolean (%s given)", obj.Type())
	}
	return b.value, nil
}

func (c *BooleanConverter) From(obj interface{}) (Object, error) {
	return NewBool(obj.(bool)), nil
}

func (c *BooleanConverter) Type() reflect.Type {
	return booleanType
}

// TimeConverter converts between time.Time and Time.
type TimeConverter struct{}

func (c *TimeConverter) To(obj Object) (interface{}, error) {
	t, ok := obj.(*Time)
	if !ok {
		return nil, fmt.Errorf("type error: expected a time (%s given)", obj.Type())
	}
	return t.value, nil
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
	mapObj, ok := obj.(*Map)
	if !ok {
		return nil, fmt.Errorf("type error: expected a map (%s given)", obj.Type())
	}
	return mapObj.Interface(), nil
}

func (c *MapStringIfaceConverter) From(obj interface{}) (Object, error) {
	m := obj.(map[string]interface{})
	objMap := make(map[string]Object, len(m))
	for k, v := range m {
		objMap[k] = FromGoType(v)
	}
	return NewMap(objMap), nil
}

func (c *MapStringIfaceConverter) Type() reflect.Type {
	return mapStrIfaceType
}

// StructConverter converts between a struct and a Map via JSON marshaling.
type StructConverter struct {
	Prototype interface{}
	AsPointer bool
}

func (c *StructConverter) To(obj Object) (interface{}, error) {
	m, ok := obj.(*Map)
	if !ok {
		return nil, fmt.Errorf("type error: expected a map (%s given)", obj.Type())
	}
	jsonBytes, err := json.Marshal(m.Interface())
	if err != nil {
		return nil, err
	}
	inst := reflect.New(reflect.TypeOf(c.Prototype))
	instValue := inst.Interface()
	if err := json.Unmarshal(jsonBytes, &instValue); err != nil {
		return nil, err
	}
	if c.AsPointer {
		return inst.Interface(), nil
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
		return nil, fmt.Errorf("type error: expected an error (%s given)", obj.Type())
	}
	return e.Message, nil
}

func (c *ErrorConverter) From(obj interface{}) (Object, error) {
	if obj == nil {
		return nil, nil
	}
	return NewError(obj.(error)), nil
}

func (c *ErrorConverter) Type() reflect.Type {
	return errType
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

func (c *ContextConverter) Type() reflect.Type {
	return contextType
}
