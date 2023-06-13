package http

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudcmds/tamarin/v2/internal/httputil"
	"github.com/cloudcmds/tamarin/v2/object"
)

func Fetch(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 2 {
		return object.NewArgsRangeError("fetch", 1, 2, len(args))
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
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, timeout, errObj := httputil.NewRequestFromParams(ctx, urlArg, params)
	if errObj != nil {
		return errObj
	}
	if timeout != 0 {
		client.Timeout = timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return object.NewError(err)
	}
	return object.NewHttpResponse(resp)
}

func Builtins() map[string]object.Object {
	return map[string]object.Object{
		"fetch": object.NewBuiltin("fetch", Fetch),
	}
}
