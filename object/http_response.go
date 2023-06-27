package object

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cloudcmds/tamarin/v2/limits"
	"github.com/cloudcmds/tamarin/v2/op"
)

type HttpResponse struct {
	resp        *http.Response
	readerLimit int64
	once        sync.Once
	closed      chan bool
}

func (r *HttpResponse) Type() Type {
	return HTTP_RESPONSE
}

func (r *HttpResponse) Inspect() string {
	return fmt.Sprintf("http.response(status: %q, content_length: %d)",
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
			name: "http.response.json",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("json", 0, len(args))
				}
				return r.JSON()
			},
		}, true
	case "text":
		return &Builtin{
			name: "http.response.text",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("text", 0, len(args))
				}
				return r.Text()

			},
		}, true
	case "close":
		return &Builtin{
			name: "http.response.close",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("close", 0, len(args))
				}
				r.Close()
				return Nil
			},
		}, true
	}
	return nil, false
}

func (r *HttpResponse) Interface() interface{} {
	return r.resp
}

func (r *HttpResponse) Close() {
	r.once.Do(func() {
		r.resp.Body.Close()
		close(r.closed)
	})
}

func (r *HttpResponse) readBody() ([]byte, error) {
	if r.readerLimit > 0 && r.resp.ContentLength > r.readerLimit {
		return nil, limits.NewLimitsError("limit error: content length exceeded limit of %d bytes (got %d)",
			r.readerLimit, r.resp.ContentLength)
	}
	return limits.ReadAll(r.resp.Body, r.readerLimit)
}

func (r *HttpResponse) JSON() Object {
	body, err := r.readBody()
	if err != nil {
		return NewError(err)
	}
	var target interface{}
	if err := json.Unmarshal(body, &target); err != nil {
		return NewError(err)
	}
	scriptObj := FromGoType(target)
	if scriptObj == nil {
		return Errorf("value error: unmarshal failed")
	}
	return scriptObj
}

func (r *HttpResponse) Text() Object {
	body, err := r.readBody()
	if err != nil {
		return NewError(err)
	}
	return NewString(string(body))
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
	return NewError(fmt.Errorf("eval error: unsupported operation for http.response: %v", opType))
}

func (r *HttpResponse) Cost() int {
	return 8 + int(r.resp.ContentLength)
}

func NewHttpResponse(
	resp *http.Response,
	timeout time.Duration,
	readerLimit int64,
) *HttpResponse {
	obj := &HttpResponse{
		resp:        resp,
		readerLimit: readerLimit,
		closed:      make(chan bool),
	}
	if timeout > 0 {
		// Guarantee that the response body is closed after the timeout
		// elapses. Alternatively, if Close is called on the object before
		// the timeout elapses, the goroutine exits.
		go func() {
			select {
			case <-time.After(timeout):
				obj.Close()
			case <-obj.closed:
			}
		}()
	}
	return obj
}
