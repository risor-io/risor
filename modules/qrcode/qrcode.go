package qrcode

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/risor-io/risor/modules/image"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"github.com/risor-io/risor/os"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

const QRCODE object.Type = "qrcode"

type QRCode struct {
	value *qrcode.QRCode
	width uint8
}

func (q *QRCode) Type() object.Type {
	return QRCODE
}

func (q *QRCode) Inspect() string {
	return fmt.Sprintf("qrcode.qrcode(width=%d)", q.width)
}

func (q *QRCode) Interface() interface{} {
	return q.value
}

func (q *QRCode) Equals(other object.Object) object.Object {
	if q == other {
		return object.True
	}
	return object.False
}

func (q *QRCode) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "save":
		return object.NewBuiltin("save", q.Save), true
	case "dimension":
		return object.NewBuiltin("dimension", q.Dimension), true
	case "bytes":
		return object.NewBuiltin("bytes", q.Bytes), true
	case "base64":
		return object.NewBuiltin("base64", q.Base64), true
	case "width":
		return object.NewInt(int64(q.width)), true
	}
	return nil, false
}

func (q *QRCode) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("type error: cannot set %q on %s object", name, QRCODE)
}

func (q *QRCode) IsTruthy() bool {
	return true
}

func (q *QRCode) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("type error: unsupported operation for qrcode")
}

func (q *QRCode) Cost() int {
	return 0
}

func New(value *qrcode.QRCode, width uint8) *QRCode {
	return &QRCode{value: value, width: width}
}

// convertStyleOptions converts a Risor map of options to standard.ImageOption objects
func convertStyleOptions(opts *object.Map) ([]standard.ImageOption, error) {
	var options []standard.ImageOption

	// Background options
	bgTransparent := opts.Get("bg_transparent")
	if bgTransparent != object.Nil {
		if transparent, err := object.AsBool(bgTransparent); err != nil {
			return nil, err
		} else if transparent {
			options = append(options, standard.WithBgTransparent())
		}
	}

	// Background color options
	bgColorHex := opts.Get("bg_color_hex")
	if bgColorHex != object.Nil {
		if hex, err := object.AsString(bgColorHex); err != nil {
			return nil, err
		} else {
			options = append(options, standard.WithBgColorRGBHex(hex))
		}
	}

	// Foreground color options
	fgColorHex := opts.Get("fg_color_hex")
	if fgColorHex != object.Nil {
		if hex, err := object.AsString(fgColorHex); err != nil {
			return nil, err
		} else {
			options = append(options, standard.WithFgColorRGBHex(hex))
		}
	}

	// Logo image from Risor image object
	logoImage := opts.Get("logo_image")
	if logoImage != object.Nil {
		// Check if it's a Risor image type
		img, ok := logoImage.(*image.Image)
		if !ok {
			return nil, fmt.Errorf("logo_image must be an image object (got %s)", logoImage.Type())
		}
		// Add the image as a logo
		options = append(options, standard.WithLogoImage(img.Value()))
	}

	// Shape options
	shape := opts.Get("shape")
	if shape != object.Nil {
		if shapeStr, err := object.AsString(shape); err != nil {
			return nil, err
		} else if shapeStr == "circle" {
			options = append(options, standard.WithCircleShape())
		} else if shapeStr != "rectangle" {
			return nil, fmt.Errorf("unsupported shape: %s (use 'circle' or 'rectangle')", shapeStr)
		}
	}

	// Border width
	borderWidth := opts.Get("border_width")
	if borderWidth != object.Nil {
		width, err := object.AsInt(borderWidth)
		if err != nil {
			return nil, err
		}
		options = append(options, standard.WithBorderWidth(int(width)))
	}

	// Format options
	format := opts.Get("format")
	if format != object.Nil {
		if formatStr, err := object.AsString(format); err != nil {
			return nil, err
		} else {
			switch formatStr {
			case "png":
				options = append(options, standard.WithBuiltinImageEncoder(standard.PNG_FORMAT))
			case "jpeg", "jpg":
				options = append(options, standard.WithBuiltinImageEncoder(standard.JPEG_FORMAT))
			default:
				return nil, fmt.Errorf("unsupported format: %s (use 'png' or 'jpeg')", formatStr)
			}
		}
	}

	return options, nil
}

// generateQRCode generates QR code data into a buffer
func (q *QRCode) generateQRCode(width uint8, opts ...standard.ImageOption) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	writerCloser := &nopCloser{Writer: buf}

	// Create a writer with the provided options
	options := []standard.ImageOption{standard.WithQRWidth(width)}
	options = append(options, opts...)
	writer := standard.NewWithWriter(writerCloser, options...)

	if err := q.value.Save(writer); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Save implements the Save method of the QRCode type as a Risor method
// The options map can include:
// - bg_transparent: (bool) make the background transparent
// - bg_color_hex: (string) set background color using hex color code (e.g. "#FFFFFF")
// - fg_color_hex: (string) set foreground color using hex color code (e.g. "#000000")
// - logo_image: (image) a Risor image object to use as a logo in the center
// - shape: (string) "circle" or "rectangle" (default: "rectangle")
// - border_width: (int) width of the border around the QR code
// - format: (string) "png" or "jpeg" (default: "png")
func (q *QRCode) Save(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 || len(args) > 2 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=1 or 2", len(args)))
	}

	// First argument should be the path
	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	// Check for style options
	var styleOptions []standard.ImageOption
	if len(args) > 1 && args[1] != object.Nil {
		optsMap, err := object.AsMap(args[1])
		if err != nil {
			return err
		}

		options, convErr := convertStyleOptions(optsMap)
		if convErr != nil {
			return object.NewError(convErr)
		}
		styleOptions = options
	}

	// Generate QR code data with style options
	qrData, genErr := q.generateQRCode(q.width, styleOptions...)
	if genErr != nil {
		return object.NewError(genErr)
	}

	// Use Risor OS to write the buffer to a file
	osObj := os.GetDefaultOS(ctx)
	if writeErr := osObj.WriteFile(path, qrData, 0o644); writeErr != nil {
		return object.NewError(writeErr)
	}

	return object.Nil
}

// Dimension implements the Dimension method of the QRCode type as a Risor method
func (q *QRCode) Dimension(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=0", len(args)))
	}

	// Call the underlying Dimension method
	dimension := q.value.Dimension()

	return object.NewInt(int64(dimension))
}

// Bytes returns the QR code as a byte slice
// The options map can include:
// - bg_transparent: (bool) make the background transparent
// - bg_color_hex: (string) set background color using hex color code (e.g. "#FFFFFF")
// - fg_color_hex: (string) set foreground color using hex color code (e.g. "#000000")
// - logo_image: (image) a Risor image object to use as a logo in the center
// - shape: (string) "circle" or "rectangle" (default: "rectangle")
// - border_width: (int) width of the border around the QR code
// - format: (string) "png" or "jpeg" (default: "png")
func (q *QRCode) Bytes(ctx context.Context, args ...object.Object) object.Object {
	if len(args) > 1 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=0 or 1", len(args)))
	}

	// Check for style options
	var styleOptions []standard.ImageOption
	if len(args) == 1 && args[0] != object.Nil {
		optsMap, err := object.AsMap(args[0])
		if err != nil {
			return err
		}

		options, convErr := convertStyleOptions(optsMap)
		if convErr != nil {
			return object.NewError(convErr)
		}
		styleOptions = options
	}

	qrData, err := q.generateQRCode(q.width, styleOptions...)
	if err != nil {
		return object.NewError(err)
	}

	return object.NewByteSlice(qrData)
}

// Base64 returns the QR code as a base64 encoded string
// The options map can include:
// - bg_transparent: (bool) make the background transparent
// - bg_color_hex: (string) set background color using hex color code (e.g. "#FFFFFF")
// - fg_color_hex: (string) set foreground color using hex color code (e.g. "#000000")
// - logo_image: (image) a Risor image object to use as a logo in the center
// - shape: (string) "circle" or "rectangle" (default: "rectangle")
// - border_width: (int) width of the border around the QR code
// - format: (string) "png" or "jpeg" (default: "png")
func (q *QRCode) Base64(ctx context.Context, args ...object.Object) object.Object {
	if len(args) > 1 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=0 or 1", len(args)))
	}

	// Check for style options
	var styleOptions []standard.ImageOption
	if len(args) == 1 && args[0] != object.Nil {
		optsMap, err := object.AsMap(args[0])
		if err != nil {
			return err
		}

		options, convErr := convertStyleOptions(optsMap)
		if convErr != nil {
			return object.NewError(convErr)
		}
		styleOptions = options
	}

	qrData, err := q.generateQRCode(q.width, styleOptions...)
	if err != nil {
		return object.NewError(err)
	}

	return object.NewString(base64.StdEncoding.EncodeToString(qrData))
}
