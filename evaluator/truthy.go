package evaluator

import "github.com/cloudcmds/tamarin/object"

func isTruthy(obj object.Object) bool {
	switch obj {
	case object.Nil:
		return false
	case object.True:
		return true
	case object.False:
		return false
	default:
		switch obj := obj.(type) {
		case *object.Int:
			return obj.Value != 0
		case *object.Float:
			return obj.Value != 0.0
		case *object.String:
			return obj.Value != ""
		case *object.List:
			return len(obj.Items) > 0
		case *object.Map:
			return len(obj.Items) > 0
		case *object.Set:
			return len(obj.Items) > 0
		case *object.Bool:
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
	case *object.Bool:
		return obj.Value
	case *object.String:
		return obj.Value != ""
	case *object.Regexp:
		return obj.Value != ""
	case *object.NilType:
		return false
	case *object.Int:
		return obj.Value != 0
	case *object.Float:
		return obj.Value != 0.0
	case *object.List:
		return len(obj.Items) != 0
	case *object.Map:
		return len(obj.Items) != 0
	case *object.Set:
		return len(obj.Items) != 0
	default:
		return true
	}
}

func nativeBoolToBooleanObject(input bool) *object.Bool {
	if input {
		return object.True
	}
	return object.False
}
