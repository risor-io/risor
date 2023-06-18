package httputil

import (
	"bytes"
	"context"
	"encoding/json"
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
	var isJSON bool

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

	// Set the request body from the "body" or "data" parameters
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
	} else if dataObj := params.GetWithDefault("data", nil); dataObj != nil {
		data, err := json.Marshal(dataObj.Interface())
		if err != nil {
			return nil, 0, object.NewError(err)
		}
		body = bytes.NewBuffer(data)
		isJSON = true
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, 0, object.NewError(err)
	}

	// Add headers to the request
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

	// Automatically set content type if JSON data was supplied
	if isJSON && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add query parameters
	if paramsObj := params.GetWithDefault("params", nil); paramsObj != nil {
		paramsMap, err := object.AsMap(paramsObj)
		if err != nil {
			return nil, 0, err
		}
		q := req.URL.Query()
		for _, k := range paramsMap.StringKeys() {
			value := paramsMap.Get(k)
			switch value := value.(type) {
			case *object.String:
				q.Add(k, value.Value())
			default:
				q.Add(k, value.Inspect())
			}
		}
		req.URL.RawQuery = q.Encode()
	}

	return req, timeout, nil
}
