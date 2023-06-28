package fetch

import (
	"context"
	"net/http"

	"github.com/risor-io/risor/internal/httputil"
	"github.com/risor-io/risor/limits"
	"github.com/risor-io/risor/object"
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
	lim, ok := limits.GetLimits(ctx)
	if !ok {
		return object.NewError(limits.LimitsNotFound)
	}
	client := &http.Client{Timeout: lim.IOTimeout()}
	req, timeout, errObj := httputil.NewRequestFromParams(ctx, urlArg, params)
	if errObj != nil {
		return errObj
	}
	if timeout != 0 {
		if timeout < client.Timeout {
			client.Timeout = timeout
		}
	}
	if err := lim.TrackHTTPRequest(req); err != nil {
		return object.NewError(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return object.NewError(err)
	}
	if err := lim.TrackHTTPResponse(resp); err != nil {
		return object.NewError(err)
	}
	return object.NewHttpResponse(resp, client.Timeout, lim.MaxBufferSize())
}

func Builtins() map[string]object.Object {
	return map[string]object.Object{
		"fetch": object.NewBuiltin("fetch", Fetch),
	}
}
