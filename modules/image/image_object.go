package image

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

type Image struct {
	image  image.Image
	format string
	buf    *bytes.Buffer
	dim    *object.Map
}

func (img *Image) Inspect() string {
	img.image.ColorModel()
	b := img.image.Bounds()
	width := b.Max.X - b.Min.X
	height := b.Max.Y - b.Min.Y
	return fmt.Sprintf("image(width=%d, height=%d)", width, height)
}

func (img *Image) Type() object.Type {
	return "image"
}

func (img *Image) Value() image.Image {
	return img.image
}

func (img *Image) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "width":
		width, _ := img.Size()
		return object.NewInt(int64(width)), true
	case "height":
		_, height := img.Size()
		return object.NewInt(int64(height)), true
	case "dimensions":
		return object.NewBuiltin("image.dimensions",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewArgsError("image.dimensions", 0, len(args))
				}
				return img.Dimensions()
			}), true
	case "bounds":
		return object.NewBuiltin("image.bounds",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewArgsError("image.bounds", 0, len(args))
				}
				return img.Bounds()
			}), true
	case "at":
		return object.NewBuiltin("image.at",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 2 {
					return object.Errorf(fmt.Sprintf("type error: image.at() takes exactly 2 arguments (%d given)", len(args)))
				}
				x, err := object.AsInt(args[0])
				if err != nil {
					return object.Errorf("type error: image.at() expects argument 1 to be an integer")
				}
				y, err := object.AsInt(args[1])
				if err != nil {
					return object.Errorf("type error: image.at() expects argument 2 to be an integer")
				}
				return object.NewColor(img.image.At(int(x), int(y)))
			}), true
	}
	return nil, false
}

func (img *Image) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: image object has no attribute %q", name)
}

func (img *Image) Interface() interface{} {
	return img.image
}

func (img *Image) String() string {
	return fmt.Sprintf("image(%s)", img.image)
}

func (img *Image) Compare(other object.Object) (int, error) {
	return 0, errors.New("type error: unable to compare images")
}

func (img *Image) Equals(other object.Object) object.Object {
	if img == other {
		return object.True
	}
	return object.False
}

func (img *Image) IsTruthy() bool {
	width, height := img.Size()
	return width > 0 && height > 0
}

func (img *Image) Size() (int, int) {
	b := img.image.Bounds()
	width := b.Max.X - b.Min.X
	height := b.Max.Y - b.Min.Y
	return width, height
}

func (img *Image) Dimensions() *object.Map {
	if img.dim == nil {
		width, height := img.Size()
		img.dim = object.NewMap(map[string]object.Object{
			"width":  object.NewInt(int64(width)),
			"height": object.NewInt(int64(height)),
		})
	}
	return img.dim
}

func (img *Image) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("eval error: unsupported operation for image: %v ", opType)
}

func (img *Image) Bounds() *object.Map {
	b := img.image.Bounds()
	min := object.NewMap(map[string]object.Object{"x": object.NewInt(int64(b.Min.X)), "y": object.NewInt(int64(b.Min.Y))})
	max := object.NewMap(map[string]object.Object{"x": object.NewInt(int64(b.Max.X)), "y": object.NewInt(int64(b.Max.Y))})
	return object.NewMap(map[string]object.Object{"min": min, "max": max})
}

// Read implements the io.Reader interface, by encoding the image to its
// default or original encoding.
func (img *Image) Read(p []byte) (int, error) {
	if img.buf == nil {
		img.buf = &bytes.Buffer{}
		var encoder imgio.Encoder = imgio.PNGEncoder()
		if err := encoder(bufio.NewWriter(img.buf), img.image); err != nil {
			return 0, err
		}
	}
	return img.buf.Read(p)
}

func (img *Image) Cost() int {
	width, height := img.Size()
	return width * height * 24
}

func NewImage(image image.Image, format string) *Image {
	return &Image{image: image, format: format}
}
