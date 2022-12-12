package object

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HttpResponse struct {
	Response *http.Response
}

func (r *HttpResponse) Type() Type {
	return HTTP_RESPONSE
}

func (r *HttpResponse) Inspect() string {
	return "<HTTP_RESPONSE>"
}

func (r *HttpResponse) InvokeMethod(method string, args ...Object) Object {
	if method == "json" {
		obj, err := r.JSON()
		if err != nil {
			return &Result{Err: &Error{Message: err.Error()}}
		}
		scriptObj := FromGoType(obj)
		if scriptObj == nil {
			return NewError("type error: unmarshal failed")
		}
		return &Result{Ok: scriptObj}
	} else if method == "text" {
		text, err := r.Text()
		if err != nil {
			return &Result{Err: &Error{Message: err.Error()}}
		}
		return &Result{Ok: &String{Value: text}}
	}
	return nil
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
	bytes, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
