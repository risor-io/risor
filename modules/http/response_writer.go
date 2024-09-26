package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const RESPONSE_WRITER object.Type = "http.response_writer"

type ResponseWriter struct {
	writer http.ResponseWriter
}

func (w *ResponseWriter) IsTruthy() bool {
	return true
}

func (w *ResponseWriter) Type() object.Type {
	return RESPONSE_WRITER
}

func (w *ResponseWriter) Inspect() string {
	return "http.response_writer()"
}

func (w *ResponseWriter) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: %s object has no attribute %q", RESPONSE_WRITER, name)
}

func (w *ResponseWriter) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "add_header":
		return object.NewBuiltin("http.response_writer.add_header",
			func(ctx context.Context, args ...object.Object) object.Object {
				if err := arg.Require("http.response_writer.add_header", 2, args); err != nil {
					return err
				}
				key, errObj := object.AsString(args[0])
				if errObj != nil {
					return object.TypeErrorf("type error: %s", errObj)
				}
				value, errObj := object.AsString(args[1])
				if errObj != nil {
					return object.TypeErrorf("type error: %s", errObj)
				}
				w.AddHeader(key, value)
				return object.Nil
			}), true
	case "del_header":
		return object.NewBuiltin("http.response_writer.del_header",
			func(ctx context.Context, args ...object.Object) object.Object {
				if err := arg.Require("http.response_writer.del_header", 1, args); err != nil {
					return err
				}
				key, errObj := object.AsString(args[0])
				if errObj != nil {
					return object.TypeErrorf("type error: %s", errObj)
				}
				w.DelHeader(key)
				return object.Nil
			}), true
	case "write":
		return object.NewBuiltin("http.response_writer.write",
			func(ctx context.Context, args ...object.Object) object.Object {
				if err := arg.Require("http.response_writer.write", 1, args); err != nil {
					return err
				}
				count, err := w.Write(args[0])
				if err != nil {
					return object.Errorf("io error: %s", err)
				}
				return object.NewInt(int64(count))
			}), true
	case "write_header":
		return object.NewBuiltin("http.response_writer.write_header",
			func(ctx context.Context, args ...object.Object) object.Object {
				if err := arg.Require("http.response_writer.write_header", 1, args); err != nil {
					return err
				}
				statusCode, errObj := object.AsInt(args[0])
				if errObj != nil {
					return object.TypeErrorf("type error: %s", errObj)
				}
				w.WriteHeader(int(statusCode))
				return object.Nil
			}), true
	}
	return nil, false
}

func (w *ResponseWriter) Interface() interface{} {
	return w.writer
}

func (w *ResponseWriter) Equals(other object.Object) object.Object {
	return object.NewBool(w == other)
}

func (w *ResponseWriter) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.EvalErrorf("eval error: unsupported operation for http.response: %v", opType)
}

func (w *ResponseWriter) Cost() int {
	return 0
}

func (w *ResponseWriter) Write(obj object.Object) (int, error) {
	switch obj := obj.(type) {
	case *object.ByteSlice:
		return w.writer.Write(obj.Value())
	case *object.String:
		return w.writer.Write([]byte(obj.Value()))
	case *object.Map:
		w.AddHeader("Content-Type", "application/json")
		data, err := obj.MarshalJSON()
		if err != nil {
			w.writeMarshalError()
			return 0, err
		}
		return w.writer.Write(data)
	case *object.List:
		w.AddHeader("Content-Type", "application/json")
		data, err := obj.MarshalJSON()
		if err != nil {
			w.writeMarshalError()
			return 0, err
		}
		return w.writer.Write(data)
	default:
		w.writeMarshalError()
		return 0, errz.TypeErrorf("type error: unsupported type for http.response_writer.write: %T", obj)
	}
}

func (w *ResponseWriter) writeMarshalError() {
	w.writer.WriteHeader(http.StatusInternalServerError)
	w.writer.Write([]byte("io error: failed to marshal response"))
}

func (w *ResponseWriter) AddHeader(key, value string) {
	w.writer.Header().Add(key, value)
}

func (w *ResponseWriter) DelHeader(key string) {
	w.writer.Header().Del(key)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.writer.WriteHeader(statusCode)
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{writer: w}
}
