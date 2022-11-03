package object

import (
	"reflect"
	"time"

	"github.com/gofrs/uuid"
)

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
		return &Float{Value: float64(obj.Unix())} // TODO: improve
	case []interface{}:
		array := &Array{Elements: make([]Object, 0, len(obj))}
		for _, item := range obj {
			arrayItem := FromGoType(item)
			if arrayItem == nil {
				return nil
			}
			array.Elements = append(array.Elements, arrayItem)
		}
		return array
	case map[string]interface{}:
		hash := &Hash{
			Pairs: make(map[HashKey]HashPair, len(obj)),
		}
		for k, v := range obj {
			hashKey := FromGoType(k)
			if hashKey == nil {
				return nil
			}
			hashVal := FromGoType(v)
			if hashVal == nil {
				return nil
			}
			hashable, ok := hashKey.(Hashable)
			if !ok {
				return nil
			}
			hash.Pairs[hashable.HashKey()] = HashPair{
				Key:   hashKey,
				Value: hashVal,
			}
		}
		return hash
	default:
		return NewError("type error: unmarshaling %v (%v)",
			obj, reflect.TypeOf(obj))
	}
}

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
	case *Array:
		array := make([]interface{}, 0, len(obj.Elements))
		for _, item := range obj.Elements {
			array = append(array, ToGoType(item))
		}
		return array
	case *Set:
		array := make([]interface{}, 0, len(obj.Items))
		for _, item := range obj.Items {
			array = append(array, ToGoType(item))
		}
		return array
	case *Hash:
		m := make(map[interface{}]interface{}, len(obj.Pairs))
		for _, v := range obj.Pairs {
			key := ToGoType(v.Key)
			val := ToGoType(v.Value)
			m[key] = val
		}
		return m
	default:
		return NewError("type error: marshaling %v (%v)",
			obj, reflect.TypeOf(obj))
	}
}
