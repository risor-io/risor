package object

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudcmds/tamarin/v2/op"
)

const TenMB = 1024 * 1024 * 10

type HttpResponse struct {
	resp    *http.Response
	body    []byte
	bodyErr error
}

func (r *HttpResponse) Type() Type {
	return HTTP_RESPONSE
}

func (r *HttpResponse) Inspect() string {
	return fmt.Sprintf("http_response(status: %q, content_length: %d)",
		r.resp.Status, r.resp.ContentLength)
}

func (r *HttpResponse) GetAttr(name string) (Object, bool) {
	switch name {
	case "status":
		return r.Status(), true
	case "status_code":
		return r.StatusCode(), true
	case "proto":
		return r.Proto(), true
	case "content_length":
		return r.ContentLength(), true
	case "header":
		return r.Header(), true
	case "json":
		return &Builtin{
			name: "http_response.json",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("json", 0, len(args))
				}
				return r.JSON()
			},
		}, true
	case "text":
		return &Builtin{
			name: "http_response.text",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("text", 0, len(args))
				}
				return r.Text()

			},
		}, true
	}
	return nil, false
}

func (r *HttpResponse) Interface() interface{} {
	return r.resp
}

func (r *HttpResponse) readBody(limit int64) error {
	defer r.resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(r.resp.Body, limit))
	if err != nil {
		r.bodyErr = err
		return err
	}
	r.body = body
	return nil
}

func (r *HttpResponse) JSON() Object {
	if r.bodyErr != nil {
		return NewError(r.bodyErr)
	}
	var target interface{}
	if err := json.Unmarshal(r.body, &target); err != nil {
		return NewError(err)
	}
	scriptObj := FromGoType(target)
	if scriptObj == nil {
		return Errorf("eval error: unmarshal failed")
	}
	return scriptObj
}

func (r *HttpResponse) Text() Object {
	if r.bodyErr != nil {
		return NewError(r.bodyErr)
	}
	return NewString(string(r.body))
}

func (r *HttpResponse) Status() *String {
	return NewString(r.resp.Status)
}

func (r *HttpResponse) StatusCode() *Int {
	return NewInt(int64(r.resp.StatusCode))
}

func (r *HttpResponse) Proto() *String {
	return NewString(r.resp.Proto)
}

func (r *HttpResponse) ContentLength() *Int {
	return NewInt(r.resp.ContentLength)
}

func (r *HttpResponse) Header() *Map {
	hdr := r.resp.Header
	m := make(map[string]Object, len(hdr))
	for k, v := range hdr {
		m[k] = NewStringList(v)
	}
	return NewMap(m)
}

func (r *HttpResponse) Equals(other Object) Object {
	if other.Type() != HTTP_RESPONSE {
		return False
	}
	return NewBool(r.resp == other.(*HttpResponse).resp)
}

func (r *HttpResponse) IsTruthy() bool {
	return true
}

func (r *HttpResponse) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for http_response: %v", opType))
}

func NewHttpResponse(resp *http.Response) *HttpResponse {
	obj := &HttpResponse{resp: resp}
	// We have to guarantee that we read and close the HTTP response body
	// in order to not leak memory. When/if we need to support different body
	// size limits or streaming, we can add a new function to help with that.
	obj.readBody(TenMB)
	return obj
}
