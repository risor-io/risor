package evaluator

import "github.com/cloudcmds/tamarin/object"

func nativeBoolToBooleanObject(input bool) *object.Bool {
	if input {
		return object.True
	}
	return object.False
}
