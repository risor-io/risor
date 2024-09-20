package image

import (
	"bufio"
	"bytes"
	"context"
	"image"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/builtins"
	"github.com/risor-io/risor/object"
)

func Decode(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("image.decode", 1, args); err != nil {
		return err
	}
	reader, typErr := object.AsReader(args[0])
	if typErr != nil {
		return typErr
	}
	img, format, err := image.Decode(reader)
	if err != nil {
		return object.NewError(err)
	}
	return NewImage(img, format)
}

func Encode(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("image.encode", 1, 2, args); err != nil {
		return err
	}
	img, ok := args[0].(*Image)
	if !ok {
		return object.Errorf("type error: image.encode() expected an image (got %s)", args[0].Type())
	}
	encoding := "png"
	if len(args) == 2 {
		encObj, ok := args[1].(*object.String)
		if !ok {
			return object.Errorf("type error: image.encode() expected a string (got %s)", args[1].Type())
		}
		encoding = encObj.Value()
	}
	var encoder imgio.Encoder
	switch encoding {
	case "png":
		encoder = imgio.PNGEncoder()
	case "jpg":
		encoder = imgio.JPEGEncoder(100)
	case "bmp":
		encoder = imgio.BMPEncoder()
	default:
		return object.Errorf("type error: image.encode() unsupported encoding %s", encoding)
	}
	buf := &bytes.Buffer{}
	writer := bufio.NewWriter(buf)
	if err := encoder(writer, img.Value()); err != nil {
		return object.NewError(err)
	}
	return object.NewByteSlice(buf.Bytes())
}

func encodePNG(ctx context.Context, obj object.Object) object.Object {
	img, ok := obj.(*Image)
	if !ok {
		return object.Errorf("type error: expected an image (got %s)", obj.Type())
	}
	encoder := imgio.PNGEncoder()
	buf := object.NewBuffer(nil)
	if err := encoder(buf, img.Value()); err != nil {
		return object.NewError(err)
	}
	return buf
}

func encodeJPG(ctx context.Context, obj object.Object) object.Object {
	img, ok := obj.(*Image)
	if !ok {
		return object.Errorf("type error: expected an image (got %s)", obj.Type())
	}
	encoder := imgio.JPEGEncoder(100)
	buf := object.NewBuffer(nil)
	if err := encoder(buf, img.Value()); err != nil {
		return object.NewError(err)
	}
	return buf
}

func encodeBMP(ctx context.Context, obj object.Object) object.Object {
	img, ok := obj.(*Image)
	if !ok {
		return object.Errorf("type error: expected an image (got %s)", obj.Type())
	}
	encoder := imgio.BMPEncoder()
	buf := object.NewBuffer(nil)
	if err := encoder(buf, img.Value()); err != nil {
		return object.NewError(err)
	}
	return buf
}

func decodeAny(ctx context.Context, obj object.Object) object.Object {
	reader, err := object.AsReader(obj)
	if err != nil {
		return err
	}
	img, format, decodeErr := image.Decode(reader)
	if decodeErr != nil {
		return object.NewError(decodeErr)
	}
	return NewImage(img, format)
}

func init() {
	builtins.RegisterCodec("png", &builtins.Codec{Encode: encodePNG, Decode: decodeAny})
	builtins.RegisterCodec("jpg", &builtins.Codec{Encode: encodeJPG, Decode: decodeAny})
	builtins.RegisterCodec("bmp", &builtins.Codec{Encode: encodeBMP, Decode: decodeAny})
}

func Module() *object.Module {
	return object.NewBuiltinsModule("image", map[string]object.Object{
		"encode": object.NewBuiltin("image.encode", Encode),
		"decode": object.NewBuiltin("image.decode", Decode),
	})
}
