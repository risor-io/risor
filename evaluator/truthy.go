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

func nativeBoolToBooleanObject(input bool) *object.Bool {
	if input {
		return object.True
	}
	return object.False
}
