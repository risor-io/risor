package httputil

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/cloudcmds/tamarin/v2/object"
)

func NewRequestFromParams(
	ctx context.Context,
	url string,
	params *object.Map,
) (*http.Request, time.Duration, *object.Error) {

	method := "GET"
	var timeout time.Duration
	var body io.Reader
	var errObj *object.Error

	// Simple request configuration with no parameters
	if params == nil {
		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return nil, 0, object.NewError(err)
		}
		return req, 0, nil
	}

	if methodObj := params.GetWithDefault("method", nil); methodObj != nil {
		method, errObj = object.AsString(methodObj)
		if errObj != nil {
			return nil, 0, errObj
		}
	}

	if timeoutObj := params.GetWithDefault("timeout", nil); timeoutObj != nil {
		timeoutInt, errObj := object.AsInt(timeoutObj)
		if errObj != nil {
			return nil, 0, errObj
		}
		timeout = time.Duration(timeoutInt) * time.Millisecond
	}

	if bodyObj := params.GetWithDefault("body", nil); bodyObj != nil {
		if reader, ok := bodyObj.(io.Reader); ok {
			body = reader
		} else {
			bodyStr, errObj := object.AsBytes(bodyObj)
			if errObj != nil {
				return nil, 0, errObj
			}
			body = bytes.NewBuffer(bodyStr)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, 0, object.NewError(err)
	}

	if headersObj := params.GetWithDefault("headers", nil); headersObj != nil {
		headersMap, err := object.AsMap(headersObj)
		if err != nil {
			return nil, 0, err
		}
		for k, v := range headersMap.Value() {
			switch v := v.(type) {
			case *object.String:
				req.Header.Add(k, v.Value())
			case *object.List:
				for _, v := range v.Value() {
					if vStr, ok := v.(*object.String); ok {
						req.Header.Add(k, vStr.Value())
					} else {
						req.Header.Add(k, v.Inspect())
					}
				}
			default:
				req.Header.Add(k, v.Inspect())
			}
		}
	}
	return req, timeout, nil
}
