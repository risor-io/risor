package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/risor-io/risor/limits"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const HTTP_RESPONSE object.Type = "http_response"

type HttpResponse struct {
	resp        *http.Response
	readerLimit int64
	once        sync.Once
	closed      chan bool
	bodyData    []byte
}

func (r *HttpResponse) IsTruthy() bool {
	return true
}

func (r *HttpResponse) Type() object.Type {
	return HTTP_RESPONSE
}

func (r *HttpResponse) Inspect() string {
	return fmt.Sprintf("http.response(status: %q, content_length: %d)",
		r.resp.Status, r.resp.ContentLength)
}

func (r *HttpResponse) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: %s object has no attribute %q", HTTP_RESPONSE, name)
}

func (r *HttpResponse) GetAttr(name string) (object.Object, bool) {
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
	case "cookies":
		return r.Cookies(), true
	case "response":
		return object.FromGoType(r.resp), true
	case "json":
		return object.NewBuiltin("http.response.json",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewArgsError("json", 0, len(args))
				}
				return r.JSON()
			}), true
	case "text":
		return object.NewBuiltin("http.response.text",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewArgsError("text", 0, len(args))
				}
				return r.Text()
			}), true
	case "close":
		return object.NewBuiltin("http.response.close",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewArgsError("close", 0, len(args))
				}
				r.Close()
				return object.Nil
			}), true
	}
	return nil, false
}

func (r *HttpResponse) AsReader() (io.Reader, *object.Error) {
	return r.resp.Body, nil
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
	if r.bodyData != nil {
		return r.bodyData, nil
	}
	if r.readerLimit > 0 && r.resp.ContentLength > r.readerLimit {
		return nil, limits.NewLimitsError("limit error: content length exceeded limit of %d bytes (got %d)",
			r.readerLimit, r.resp.ContentLength)
	}
	data, err := limits.ReadAll(r.resp.Body, r.readerLimit)
	if err != nil {
		return nil, err
	}
	r.bodyData = data
	return data, nil
}

func (r *HttpResponse) JSON() object.Object {
	body, err := r.readBody()
	if err != nil {
		return object.NewError(err)
	}
	var target interface{}
	if err := json.Unmarshal(body, &target); err != nil {
		return object.NewError(err)
	}
	scriptObj := object.FromGoType(target)
	if scriptObj == nil {
		return object.Errorf("value error: unmarshal failed")
	}
	return scriptObj
}

func (r *HttpResponse) Text() object.Object {
	body, err := r.readBody()
	if err != nil {
		return object.NewError(err)
	}
	return object.NewString(string(body))
}

func (r *HttpResponse) Status() *object.String {
	return object.NewString(r.resp.Status)
}

func (r *HttpResponse) StatusCode() *object.Int {
	return object.NewInt(int64(r.resp.StatusCode))
}

func (r *HttpResponse) Proto() *object.String {
	return object.NewString(r.resp.Proto)
}

func (r *HttpResponse) ContentLength() *object.Int {
	return object.NewInt(r.resp.ContentLength)
}

func (r *HttpResponse) Header() *object.Map {
	hdr := r.resp.Header
	m := make(map[string]object.Object, len(hdr))
	for k, v := range hdr {
		m[k] = object.NewStringList(v)
	}
	return object.NewMap(m)
}

func (r *HttpResponse) Cookies() *object.Map {
	cookies := r.resp.Cookies()
	m := make(map[string]object.Object, len(cookies))
	for _, cookie := range cookies {
		m[cookie.Name] = object.FromGoType(cookie)
	}
	return object.NewMap(m)
}

func (r *HttpResponse) Equals(other object.Object) object.Object {
	if other.Type() != HTTP_RESPONSE {
		return object.False
	}
	return object.NewBool(r.resp == other.(*HttpResponse).resp)
}

func (r *HttpResponse) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("eval error: unsupported operation for http.response: %v", opType))
}

func (r *HttpResponse) Cost() int {
	return 8 + int(r.resp.ContentLength)
}

func (r *HttpResponse) MarshalJSON() ([]byte, error) {
	var text string
	var jsonObj object.Object
	data, err := r.readBody()
	if err != nil {
		return nil, err
	}
	if r.resp.Header.Get("Content-Type") == "application/json" {
		jsonObj = r.JSON()
	} else {
		text = string(data)
	}
	return json.Marshal(struct {
		Status        string         `json:"status"`
		StatusCode    int            `json:"status_code"`
		Proto         string         `json:"proto"`
		ContentLength int64          `json:"content_length"`
		Header        http.Header    `json:"header"`
		Cookies       []*http.Cookie `json:"cookies,omitempty"`
		Text          string         `json:"text,omitempty"`
		JSON          object.Object  `json:"json,omitempty"`
	}{
		Status:        r.resp.Status,
		StatusCode:    r.resp.StatusCode,
		Proto:         r.resp.Proto,
		ContentLength: r.resp.ContentLength,
		Header:        r.resp.Header,
		Cookies:       r.resp.Cookies(),
		Text:          text,
		JSON:          jsonObj,
	})
}

func NewHttpResponse(resp *http.Response, timeout time.Duration, readerLimit int64) *HttpResponse {
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
