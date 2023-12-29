package base64

import (
	"context"
	_ "embed"
	"encoding/base64"

	"github.com/risor-io/risor/object"
)

func Encode(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs < 1 || nArgs > 2 {
		return object.Errorf("type error: base64.encode() takes 1 or 2 arguments (%d given)", nArgs)
	}
	data, err := object.AsBytes(args[0])
	if err != nil {
		return err
	}
	padding := true
	if nArgs == 2 {
		var err *object.Error
		padding, err = object.AsBool(args[1])
		if err != nil {
			return err
		}
	}
	var enc *base64.Encoding
	if padding {
		enc = base64.StdEncoding
	} else {
		enc = base64.RawStdEncoding
	}
	dst := make([]byte, enc.EncodedLen(len(data)))
	enc.Encode(dst, data)
	return object.NewString(string(dst))
}

func URLEncode(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs < 1 || nArgs > 2 {
		return object.Errorf("type error: base64.url_encode() takes 1 or 2 arguments (%d given)", nArgs)
	}
	data, err := object.AsBytes(args[0])
	if err != nil {
		return err
	}
	padding := true
	if nArgs == 2 {
		var err *object.Error
		padding, err = object.AsBool(args[1])
		if err != nil {
			return err
		}
	}
	var enc *base64.Encoding
	if padding {
		enc = base64.URLEncoding
	} else {
		enc = base64.RawURLEncoding
	}
	dst := make([]byte, enc.EncodedLen(len(data)))
	enc.Encode(dst, data)
	return object.NewString(string(dst))
}

func Decode(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs < 1 || nArgs > 2 {
		return object.Errorf("type error: base64.decode() takes 1 or 2 arguments (%d given)", nArgs)
	}
	data, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	padding := true
	if nArgs == 2 {
		var err *object.Error
		padding, err = object.AsBool(args[1])
		if err != nil {
			return err
		}
	}
	var enc *base64.Encoding
	if padding {
		enc = base64.StdEncoding
	} else {
		enc = base64.RawStdEncoding
	}
	dst := make([]byte, enc.DecodedLen(len(data)))
	count, decodeErr := enc.Decode(dst, []byte(data))
	if decodeErr != nil {
		return object.NewError(decodeErr)
	}
	return object.NewByteSlice(dst[:count])
}

func URLDecode(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs < 1 || nArgs > 2 {
		return object.Errorf("type error: base64.url_decode() takes 1 or 2 arguments (%d given)", nArgs)
	}
	data, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	padding := true
	if nArgs == 2 {
		var err *object.Error
		padding, err = object.AsBool(args[1])
		if err != nil {
			return err
		}
	}
	var enc *base64.Encoding
	if padding {
		enc = base64.URLEncoding
	} else {
		enc = base64.RawURLEncoding
	}
	dst := make([]byte, enc.DecodedLen(len(data)))
	count, decodeErr := enc.Decode(dst, []byte(data))
	if decodeErr != nil {
		return object.NewError(decodeErr)
	}
	return object.NewByteSlice(dst[:count])
}

//go:embed base64.md
var docs string

func Module() *object.Module {
	return object.NewBuiltinsModule("base64", map[string]object.Object{
		"decode":     object.NewBuiltin("decode", Decode),
		"encode":     object.NewBuiltin("encode", Encode),
		"url_decode": object.NewBuiltin("url_decode", URLDecode),
		"url_encode": object.NewBuiltin("url_encode", URLEncode),
	}).WithDocstring(docs)
}
