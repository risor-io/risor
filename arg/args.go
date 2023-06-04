package arg

import (
	"fmt"

	"github.com/cloudcmds/tamarin/v2/object"
)

func Require(funcName string, count int, args []object.Object) *object.Error {
	nArgs := len(args)
	if nArgs != count {
		if count == 1 {
			return object.Errorf(
				fmt.Sprintf("type error: %s() takes exactly 1 argument (%d given)",
					funcName, nArgs))
		}
		return object.Errorf(
			fmt.Sprintf("type error: %s() takes exactly %d arguments (%d given)",
				funcName, count, nArgs))
	}
	return nil
}
