package object

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type HttpResponse struct {
	Response *http.Response
}

func (r *HttpResponse) Type() Type {
	return HTTP_RESPONSE
}

func (r *HttpResponse) Inspect() string {
	return "HttpResponse()"
}

func (r *HttpResponse) GetAttr(name string) (Object, bool) {
	switch name {
	case "json":
		return &Builtin{
			Name: "http_response.json",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("json", 0, len(args))
				}
				obj, err := r.JSON()
				if err != nil {
					return &Result{Err: &Error{Message: err.Error()}}
				}
				scriptObj := FromGoType(obj)
				if scriptObj == nil {
					return NewError("type error: unmarshal failed")
				}
				return &Result{Ok: scriptObj}
			},
		}, true
	case "text":
		return &Builtin{
			Name: "http_response.text",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("text", 0, len(args))
				}
				text, err := r.Text()
				if err != nil {
					return &Result{Err: &Error{Message: err.Error()}}
				}
				return &Result{Ok: &String{Value: text}}
			},
		}, true
	}
	return nil, false
}

func (r *HttpResponse) ToInterface() interface{} {
	return r.Response
}

func (r *HttpResponse) JSON() (target interface{}, err error) {
	defer r.Response.Body.Close()
	err = json.NewDecoder(r.Response.Body).Decode(&target)
	return
}

func (r *HttpResponse) Text() (string, error) {
	defer r.Response.Body.Close()
	bytes, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (r *HttpResponse) Equals(other Object) Object {
	if other.Type() != HTTP_RESPONSE {
		return False
	}
	return NewBool(r.Response == other.(*HttpResponse).Response)
}
