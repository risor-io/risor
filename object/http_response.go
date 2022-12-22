package object

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type HttpResponse struct {
	resp *http.Response
}

func (r *HttpResponse) Type() Type {
	return HTTP_RESPONSE
}

func (r *HttpResponse) Inspect() string {
	return "http_response()"
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

func (r *HttpResponse) JSON() *Result {
	defer r.resp.Body.Close()
	var target interface{}
	if err := json.NewDecoder(r.resp.Body).Decode(&target); err != nil {
		return NewErrResult(NewError(err))
	}
	scriptObj := FromGoType(target)
	if scriptObj == nil {
		return NewErrResult(Errorf("eval error: unmarshal failed"))
	}
	return NewOkResult(scriptObj)
}

func (r *HttpResponse) Text() *Result {
	defer r.resp.Body.Close()
	bytes, err := io.ReadAll(r.resp.Body)
	if err != nil {
		return NewErrResult(NewError(err))
	}
	if err != nil {
		return NewErrResult(NewError(err))
	}
	return NewOkResult(NewString(string(bytes)))
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

func NewHttpResponse(resp *http.Response) *HttpResponse {
	return &HttpResponse{resp: resp}
}
