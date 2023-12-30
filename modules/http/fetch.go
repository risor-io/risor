package http

import (
	"context"

	"github.com/risor-io/risor/object"
)

func Fetch(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 2 {
		return object.NewArgsRangeError("fetch", 1, 2, numArgs)
	}
	urlArg, argErr := object.AsString(args[0])
	if argErr != nil {
		return argErr
	}
	var errObj *object.Error
	var params *object.Map
	if numArgs == 2 {
		params, errObj = object.AsMap(args[1])
		if errObj != nil {
			return errObj
		}
	}
	req, errObj := NewRequestFromParams(urlArg, params)
	if errObj != nil {
		return errObj
	}

	return req.Send(ctx)
}
