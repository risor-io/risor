package object

import (
	"context"
	"errors"
	"fmt"
	"image"

	"github.com/cloudcmds/tamarin/v2/op"
)

type Image struct {
	image image.Image
}

func (img *Image) Inspect() string {
	img.image.ColorModel()
	b := img.image.Bounds()
	width := b.Max.X - b.Min.X
	height := b.Max.Y - b.Min.Y
	return fmt.Sprintf("image(width=%d, height=%d)", width, height)
}

func (img *Image) Type() Type {
	return IMAGE
}

func (img *Image) Value() image.Image {
	return img.image
}

func (img *Image) GetAttr(name string) (Object, bool) {
	switch name {
	case "bounds":
		return &Builtin{
			name: "image.bounds",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("image.bounds", 0, len(args))
				}
				return img.Bounds()
			},
		}, true
	case "at":
		return &Builtin{
			name: "image.at",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 2 {
					return Errorf(fmt.Sprintf("type error: image.at() takes exactly 2 arguments (%d given)", len(args)))
				}
				x, err := AsInt(args[0])
				if err != nil {
					return Errorf("type error: image.at() expects argument 1 to be an integer")
				}
				y, err := AsInt(args[1])
				if err != nil {
					return Errorf("type error: image.at() expects argument 2 to be an integer")
				}
				return NewColor(img.image.At(int(x), int(y)))
			},
		}, true
	}
	return nil, false
}

func (img *Image) Interface() interface{} {
	return img.image
}

func (img *Image) String() string {
	return fmt.Sprintf("image(%s)", img.image)
}

func (img *Image) Compare(other Object) (int, error) {
	return 0, errors.New("type error: unable to compare images")
}

func (img *Image) Equals(other Object) Object {
	switch other := other.(type) {
	case *Image:
		if img.image == other.image {
			return True
		}
	}
	return False
}

func (img *Image) IsTruthy() bool {
	b := img.image.Bounds()
	width := b.Max.X - b.Min.X
	height := b.Max.Y - b.Min.Y
	return width > 0 && height > 0
}

func (img *Image) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for image: %v ", opType))
}

func (img *Image) Bounds() Object {
	b := img.image.Bounds()
	min := NewMap(map[string]Object{"x": NewInt(int64(b.Min.X)), "y": NewInt(int64(b.Min.Y))})
	max := NewMap(map[string]Object{"x": NewInt(int64(b.Max.X)), "y": NewInt(int64(b.Max.Y))})
	return NewMap(map[string]Object{"min": min, "max": max})
}

func NewImage(image image.Image) *Image {
	return &Image{image: image}
}
