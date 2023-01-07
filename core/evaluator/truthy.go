package evaluator

import "github.com/cloudcmds/tamarin/core/object"

func nativeBoolToBooleanObject(input bool) *object.Bool {
	if input {
		return object.True
	}
	return object.False
}
