package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/risor-io/risor/limits"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const HTTP_REQUEST object.Type = "http_request"

type HttpRequest struct {
	req     *http.Request
	client  *http.Client
	timeout time.Duration
}

func (r *HttpRequest) IsTruthy() bool {
	return true
}

func (r *HttpRequest) Type() object.Type {
	return HTTP_REQUEST
}

func (r *HttpRequest) Inspect() string {
	return fmt.Sprintf("http.request(url: %s, method: %s)",
		r.req.URL.String(), r.req.Method)
}

func (r *HttpRequest) SetAttr(name string, value object.Object) error {
	switch name {
	case "timeout":
		tStr, objErr := object.AsString(value)
		if objErr != nil {
			return objErr.Value()
		}
		t, err := time.ParseDuration(tStr)
		if err != nil {
			return err
		}
		r.timeout = t
	case "header":
		headersMap, err := object.AsMap(value)
		if err != nil {
			return err.Value()
		}
		r.AddHeaders(headersMap)
	case "params":
		paramsMap, err := object.AsMap(value)
		if err != nil {
			return err.Value()
		}
		r.SetParams(paramsMap)
	case "body":
		if err := r.SetBody(value); err != nil {
			return err.Value()
		}
	case "data":
		if err := r.SetData(value); err != nil {
			return err.Value()
		}
	default:
		return object.TypeErrorf("type error: %s object has no attribute %q", HTTP_REQUEST, name)
	}

	return nil
}

func (r *HttpRequest) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "url":
		return r.URL(), true
	case "query":
		rm := make(map[string]object.Object)
		for k, v := range r.req.URL.Query() {
			rm[k] = object.NewString(v[0])
		}
		return object.NewMap(rm), true
	case "content_length":
		return r.ContentLength(), true
	case "header":
		return r.Header(), true
	case "path_value":
		return object.NewBuiltin("http.request.path_value", func(ctx context.Context, args ...object.Object) object.Object {
			if numArgs := len(args); numArgs != 1 {
				return object.NewArgsError("http.request.path_value", 1, numArgs)
			}
			key, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			return object.NewString(r.req.PathValue(key))
		}), true
	case "send":
		return object.NewBuiltin("http.request.send", func(ctx context.Context, args ...object.Object) object.Object {
			return r.Send(ctx)
		}), true
	case "add_header":
		return object.NewBuiltin("http.request.add_header", func(ctx context.Context, args ...object.Object) object.Object {
			numArgs := len(args)
			if numArgs != 2 {
				return object.NewArgsError("http.request.add_header", 2, numArgs)
			}
			name, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			r.AddHeader(name, args[1])
			return nil
		}), true
	case "add_cookie":
		return object.NewBuiltin("http.request.add_cookie", func(ctx context.Context, args ...object.Object) object.Object {
			if numArgs := len(args); numArgs != 2 {
				return object.NewArgsError("http.request.add_cookie", 2, numArgs)
			}
			name, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			value, err := object.AsMap(args[1])
			if err != nil {
				return err
			}
			return r.AddCookie(name, value)
		}), true
	case "set_body":
		return object.NewBuiltin("http.request.set_body", func(ctx context.Context, args ...object.Object) object.Object {
			if numArgs := len(args); numArgs != 1 {
				return object.NewArgsError("http.request.set_body", 1, numArgs)
			}
			return r.SetBody(args[0])
		}), true
	case "set_data":
		return object.NewBuiltin("http.request.set_data", func(ctx context.Context, args ...object.Object) object.Object {
			if numArgs := len(args); numArgs != 1 {
				return object.NewArgsError("http.request.set_data", 1, numArgs)
			}
			return r.SetData(args[0])
		}), true
	}
	return nil, false
}

func (r *HttpRequest) URL() *object.String {
	return object.NewString(r.req.URL.String())
}

func (r *HttpRequest) ContentLength() *object.Int {
	return object.NewInt(r.req.ContentLength)
}

func (r *HttpRequest) Header() *object.Map {
	hdr := r.req.Header
	m := make(map[string]object.Object, len(hdr))
	for k, v := range hdr {
		m[k] = object.NewStringList(v)
	}
	return object.NewMap(m)
}

func (r *HttpRequest) Interface() interface{} {
	return r.req
}

func (r *HttpRequest) Equals(other object.Object) object.Object {
	if other.Type() != HTTP_REQUEST {
		return object.False
	}
	return object.NewBool(r.req == other.(*HttpRequest).req)
}

func (r *HttpRequest) Cost() int {
	return 8 + int(r.req.ContentLength)
}

func (r *HttpRequest) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for http.request: %v", opType)
}

func (r *HttpRequest) AddHeaders(headers *object.Map) {
	for k, v := range headers.Value() {
		r.AddHeader(k, v)
	}
}

func (r *HttpRequest) AddHeader(name string, v object.Object) {
	switch v := v.(type) {
	case *object.String:
		r.req.Header.Add(name, v.Value())
		if strings.EqualFold(name, "host") {
			r.req.Host = v.Value()
		}
	case *object.List:
		for _, v := range v.Value() {
			if vStr, ok := v.(*object.String); ok {
				r.req.Header.Add(name, vStr.Value())
			} else {
				r.req.Header.Add(name, v.Inspect())
			}
		}
	default:
		r.req.Header.Add(name, v.Inspect())
	}
}

func (r *HttpRequest) AddCookie(name string, cookie *object.Map) *object.Error {
	c := &http.Cookie{
		Name: name,
	}
	if valueObj := cookie.GetWithDefault("value", nil); valueObj != nil {
		value, err := object.AsString(valueObj)
		if err != nil {
			return err
		}
		c.Value = value
	}
	if pathObj := cookie.GetWithDefault("path", nil); pathObj != nil {
		path, err := object.AsString(pathObj)
		if err != nil {
			return err
		}
		c.Path = path
	}
	if maxAgeObj := cookie.GetWithDefault("max_age", nil); maxAgeObj != nil {
		maxAge, objErr := object.AsString(maxAgeObj)
		if objErr != nil {
			return objErr
		}
		d, err := time.ParseDuration(maxAge)
		if err != nil {
			return object.NewError(err)
		}
		c.MaxAge = int(d.Seconds())
	}
	if secureObj := cookie.GetWithDefault("secure", nil); secureObj != nil {
		secure, err := object.AsBool(secureObj)
		if err != nil {
			return err
		}
		c.Secure = secure
	}
	if httpOnlyObj := cookie.GetWithDefault("http_only", nil); httpOnlyObj != nil {
		httpOnly, err := object.AsBool(httpOnlyObj)
		if err != nil {
			return err
		}
		c.HttpOnly = httpOnly
	}
	r.req.AddCookie(c)
	return nil
}

func (r *HttpRequest) SetParams(params *object.Map) {
	q := r.req.URL.Query()
	for _, k := range params.StringKeys() {
		value := params.Get(k)
		switch value := value.(type) {
		case *object.String:
			q.Add(k, value.Value())
		default:
			q.Add(k, value.Inspect())
		}
	}
	r.req.URL.RawQuery = q.Encode()
}

func (r *HttpRequest) SetBody(bodyObj object.Object) *object.Error {
	if bodyObj == nil {
		r.req.Body = nil
		r.req.ContentLength = 0
		return nil
	}
	var body io.Reader
	if reader, ok := bodyObj.(io.Reader); ok {
		body = reader
		switch v := body.(type) {
		case *bytes.Buffer:
			r.req.ContentLength = int64(v.Len())
			buf := v.Bytes()
			r.req.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return io.NopCloser(r), nil
			}
		case *bytes.Reader:
			r.req.ContentLength = int64(v.Len())
			snapshot := *v
			r.req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(&r), nil
			}
		case *strings.Reader:
			r.req.ContentLength = int64(v.Len())
			snapshot := *v
			r.req.GetBody = func() (io.ReadCloser, error) {
				rs := snapshot
				return io.NopCloser(&rs), nil
			}
		}
	} else {
		data, errObj := object.AsBytes(bodyObj)
		if errObj != nil {
			return errObj
		}
		body = bytes.NewBuffer(data)
		r.req.ContentLength = int64(len(data))
	}
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = io.NopCloser(body)
	}
	r.req.Body = rc
	return nil
}

func (r *HttpRequest) SetData(dataObj object.Object) *object.Error {
	if dataObj == nil {
		r.req.Body = nil
		r.req.ContentLength = 0
		return nil
	}
	data, err := json.Marshal(dataObj.Interface())
	if err != nil {
		return object.NewError(err)
	}
	body := bytes.NewBuffer(data)

	r.req.Body = io.NopCloser(body)
	r.req.ContentLength = int64(len(data))

	// Automatically set content type if JSON data was supplied
	if r.req.Header.Get("Content-Type") == "" {
		r.req.Header.Set("Content-Type", "application/json")
	}

	return nil
}

func (r *HttpRequest) Send(ctx context.Context) object.Object {
	lim, _ := limits.GetLimits(ctx)
	if r.req == nil {
		return object.Errorf("bad request")
	}
	if r.client == nil {
		r.client = &http.Client{}
	}
	if lim != nil {
		r.client.Timeout = lim.IOTimeout()
	}
	if r.timeout != 0 {
		if r.client.Timeout == 0 || r.timeout < r.client.Timeout {
			r.client.Timeout = r.timeout
		}
	}
	req := r.req.WithContext(ctx)
	if lim != nil {
		if err := lim.TrackHTTPRequest(req); err != nil {
			return object.NewError(err)
		}
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return object.NewError(err)
	}
	readerLimit := int64(-1)
	if lim != nil {
		if err := lim.TrackHTTPResponse(resp); err != nil {
			return object.NewError(err)
		}
		readerLimit = lim.MaxBufferSize()
	}
	return NewHttpResponse(resp, r.client.Timeout, readerLimit)
}

func NewRequestFromParams(url string, params *object.Map) (*HttpRequest, *object.Error) {
	method := "GET"
	var errObj *object.Error
	var isJSON bool

	// Simple request configuration with no parameters
	if params == nil {
		req, err := http.NewRequest(method, url, http.NoBody)
		if err != nil {
			return nil, object.NewError(err)
		}
		return &HttpRequest{req: req}, nil
	}

	r := &HttpRequest{}

	if methodObj := params.GetWithDefault("method", nil); methodObj != nil {
		method, errObj = object.AsString(methodObj)
		if errObj != nil {
			return nil, errObj
		}
	}

	if timeoutObj := params.GetWithDefault("timeout", nil); timeoutObj != nil {
		timeoutInt, errObj := object.AsInt(timeoutObj)
		if errObj != nil {
			return nil, errObj
		}
		r.timeout = time.Duration(timeoutInt) * time.Millisecond
	}

	var body io.Reader
	// Set the request body from the "body" or "data" parameters
	if bodyObj := params.GetWithDefault("body", nil); bodyObj != nil {
		if reader, ok := bodyObj.(io.Reader); ok {
			body = reader
		} else {
			bodyStr, errObj := object.AsBytes(bodyObj)
			if errObj != nil {
				return nil, errObj
			}
			body = bytes.NewBuffer(bodyStr)
		}
	} else if dataObj := params.GetWithDefault("data", nil); dataObj != nil {
		data, err := json.Marshal(dataObj.Interface())
		if err != nil {
			return nil, object.NewError(err)
		}
		body = bytes.NewBuffer(data)

		isJSON = true
	}

	var err error
	r.req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, object.NewError(err)
	}

	// Add headers to the request
	if headersObj := params.GetWithDefault("headers", nil); headersObj != nil {
		headersMap, err := object.AsMap(headersObj)
		if err != nil {
			return nil, err
		}
		r.AddHeaders(headersMap)
	}

	// Automatically set content type if JSON data was supplied
	if isJSON && r.req.Header.Get("Content-Type") == "" {
		r.req.Header.Set("Content-Type", "application/json")
	}

	// Add query parameters
	if paramsObj := params.GetWithDefault("params", nil); paramsObj != nil {
		paramsMap, err := object.AsMap(paramsObj)
		if err != nil {
			return nil, err
		}
		r.SetParams(paramsMap)
	}

	// Build the HTTP client
	c, err := NewHTTPClientFromParams(params)
	if err != nil {
		return nil, object.NewError(err)
	}

	r.client = c

	return r, nil
}

func NewRequest(r *http.Request) *HttpRequest {
	return &HttpRequest{req: r}
}
