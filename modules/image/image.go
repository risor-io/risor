package image

import (
	"context"
	"image"
	"io"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/cloudcmds/tamarin/v2/internal/arg"
	"github.com/cloudcmds/tamarin/v2/object"
)

func Open(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("image.open", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.String:
		img, err := imgio.Open(arg.Value())
		if err != nil {
			return object.NewError(err)
		}
		return object.NewImage(img)
	case io.Reader:
		img, _, err := image.Decode(arg)
		if err != nil {
			return object.NewError(err)
		}
		return object.NewImage(img)
	default:
		return object.Errorf("type error: image.open() expected a string or reader (got %s)", args[0].Type())
	}
}

func Save(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("image.save", 2, 3, args); err != nil {
		return err
	}
	img, ok := args[0].(*object.Image)
	if !ok {
		return object.Errorf("type error: image.save() expected an image (got %s)", args[0].Type())
	}
	encoding := "png"
	if len(args) == 3 {
		encObj, ok := args[2].(*object.String)
		if !ok {
			return object.Errorf("type error: image.save() expected a string (got %s)", args[2].Type())
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
		return object.Errorf("type error: image.save() unsupported encoding %s", encoding)
	}
	switch dst := args[1].(type) {
	case *object.String:
		if err := imgio.Save(dst.Value(), img.Value(), encoder); err != nil {
			return object.NewError(err)
		}
	case io.Writer:
		if err := encoder(dst, img.Value()); err != nil {
			return object.NewError(err)
		}
	}
	return object.Nil
}

func Module() *object.Module {
	return object.NewBuiltinsModule("image", map[string]object.Object{
		"open": object.NewBuiltin("open", Open),
		"save": object.NewBuiltin("save", Save),
	})
}
