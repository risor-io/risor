package evaluator

import "github.com/cloudcmds/tamarin/object"

func isTruthy(obj object.Object) bool {
	switch obj {
	case object.NULL:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		switch obj := obj.(type) {
		case *object.Integer:
			return obj.Value != 0
		case *object.Float:
			return obj.Value != 0.0
		case *object.String:
			return obj.Value != ""
		case *object.Array:
			return len(obj.Elements) > 0
		case *object.Hash:
			return len(obj.Map) > 0
		case *object.Set:
			return len(obj.Items) > 0
		case *object.Boolean:
			return obj.Value
		}
		return true
	}
}

func objectToNativeBoolean(o object.Object) bool {
	if r, ok := o.(*object.ReturnValue); ok {
		o = r.Value
	}
	switch obj := o.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.String:
		return obj.Value != ""
	case *object.Regexp:
		return obj.Value != ""
	case *object.Null:
		return false
	case *object.Integer:
		return obj.Value != 0
	case *object.Float:
		return obj.Value != 0.0
	case *object.Array:
		return len(obj.Elements) != 0
	case *object.Hash:
		return len(obj.Map) != 0
	case *object.Set:
		return len(obj.Items) != 0
	default:
		return true
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return object.TRUE
	}
	return object.FALSE
}
