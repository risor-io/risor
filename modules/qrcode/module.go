package qrcode

import (
	"context"
	"io"

	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/os"
	qrcode "github.com/yeqown/go-qrcode/v2"
)

// Create creates a new QR code with the given content
//
// Arguments:
//   - content: the string to encode in the QR code
//   - options: (optional) a map of configuration options:
//   - encoding_mode: "numeric", "alphanumeric", or "byte"
//   - error_correction: "low", "medium", "high", or "highest"
//   - width: integer QR code width in pixels (default: 40)
func Create(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("qrcode.create", 1, 2, args); err != nil {
		return err
	}

	content, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	var options []qrcode.EncodeOption
	width := uint8(40) // default width

	// Check for options map
	if len(args) > 1 && args[1] != object.Nil {
		opts, errObj := object.AsMap(args[1])
		if errObj != nil {
			return errObj
		}

		// Handle encoding mode option
		encModeObj := opts.Get("encoding_mode")
		if encModeObj != object.Nil {
			encMode, err := object.AsString(encModeObj)
			if err != nil {
				return err
			}

			switch encMode {
			case "numeric":
				options = append(options, qrcode.WithEncodingMode(qrcode.EncModeNumeric))
			case "alphanumeric":
				options = append(options, qrcode.WithEncodingMode(qrcode.EncModeAlphanumeric))
			case "byte":
				options = append(options, qrcode.WithEncodingMode(qrcode.EncModeByte))
			default:
				return object.Errorf("invalid encoding mode: must be 'numeric', 'alphanumeric', or 'byte'")
			}
		}

		// Handle error correction level option
		errLevelObj := opts.Get("error_correction")
		if errLevelObj != object.Nil {
			errLevel, err := object.AsString(errLevelObj)
			if err != nil {
				return err
			}

			switch errLevel {
			case "low":
				options = append(options, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionLow))
			case "medium":
				options = append(options, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionLow+1))
			case "high":
				options = append(options, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionLow+2))
			case "highest":
				options = append(options, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionLow+3))
			default:
				return object.Errorf("invalid error correction level: must be 'low', 'medium', 'high', or 'highest'")
			}
		}

		// Handle width option
		widthObj := opts.Get("width")
		if widthObj != object.Nil {
			w, err := object.AsInt(widthObj)
			if err != nil {
				return err
			}
			if w < 1 || w > 255 {
				return object.Errorf("invalid width: must be between 1 and 255")
			}
			width = uint8(w)
		}
	}

	qrc, newErr := qrcode.NewWith(content, options...)
	if newErr != nil {
		return object.NewError(newErr)
	}

	return New(qrc, width)
}

// GetOS returns the OS from the context
func GetOS(ctx context.Context) os.OS {
	return os.GetDefaultOS(ctx)
}

// Save saves the QR code to a PNG file using the Risor OS
func Save(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("qrcode.save", 2, 3, args); err != nil {
		return err
	}

	qr, ok := args[0].(*QRCode)
	if !ok {
		return object.TypeErrorf("first argument to save_png must be a qrcode (got %s)", args[0].Type())
	}

	path, err := object.AsString(args[1])
	if err != nil {
		return err
	}

	width := qr.width
	if len(args) > 2 {
		w, err := object.AsInt(args[2])
		if err != nil {
			return err
		}
		if w > 255 {
			w = 255 // ensure width fits in uint8
		}
		width = uint8(w)
	}

	// Generate QR code data
	qrData, genErr := qr.generateQRCode(width)
	if genErr != nil {
		return object.NewError(genErr)
	}

	// Use Risor OS to write the buffer to a file
	osObj := GetOS(ctx)
	if writeErr := osObj.WriteFile(path, qrData, 0644); writeErr != nil {
		return object.NewError(writeErr)
	}

	return object.Nil
}

// nopCloser is a wrapper around an io.Writer that implements io.WriteCloser
type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func Module() *object.Module {
	return object.NewBuiltinsModule("qrcode", map[string]object.Object{
		"create": object.NewBuiltin("create", Create),
		"save":   object.NewBuiltin("save", Save),
	}, Create)
}
