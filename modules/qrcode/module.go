package qrcode

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"

	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
	qrcode "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// Create creates a new QR code with the given content
//
// Arguments:
//   - content: the string to encode in the QR code
//   - options: (optional) a map of configuration options:
//   - encoding_mode: "numeric", "alphanumeric", or "byte"
//   - error_correction: "low", "medium", "high", or "highest"
func Create(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("qrcode.create", 1, 2, args); err != nil {
		return err
	}

	content, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	var options []qrcode.EncodeOption

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
	}

	qrc, newErr := qrcode.NewWith(content, options...)
	if newErr != nil {
		return object.NewError(newErr)
	}

	return New(qrc)
}

// Save saves the QR code to a PNG file
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

	width := 256 // default width
	if len(args) > 2 {
		w, err := object.AsInt(args[2])
		if err != nil {
			return err
		}
		width = int(w)
		if width > 255 {
			width = 255 // ensure width fits in uint8
		}
	}

	writer, newErr := standard.New(path, standard.WithQRWidth(uint8(width)))
	if newErr != nil {
		return object.NewError(newErr)
	}

	if saveErr := qr.value.Save(writer); saveErr != nil {
		return object.NewError(saveErr)
	}

	return object.Nil
}

// ToBase64 returns a base64 encoded string of the QR code PNG image
func ToBase64(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("qrcode.to_base64", 1, 2, args); err != nil {
		return err
	}

	qr, ok := args[0].(*QRCode)
	if !ok {
		return object.TypeErrorf("first argument to to_base64 must be a qrcode (got %s)", args[0].Type())
	}

	width := 256 // default width
	if len(args) > 1 {
		w, err := object.AsInt(args[1])
		if err != nil {
			return err
		}
		width = int(w)
		if width > 255 {
			width = 255 // ensure width fits in uint8
		}
	}

	buf := bytes.NewBuffer(nil)
	writerCloser := nopCloser{Writer: buf}
	writer := standard.NewWithWriter(writerCloser, standard.WithQRWidth(uint8(width)))

	err := qr.value.Save(writer)
	if err != nil {
		return object.NewError(err)
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return object.NewString(encoded)
}

// nopCloser is a wrapper around an io.Writer that implements io.WriteCloser
type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func Module() *object.Module {
	return object.NewBuiltinsModule("qrcode", map[string]object.Object{
		"create":    object.NewBuiltin("create", Create),
		"save":      object.NewBuiltin("save", Save),
		"to_base64": object.NewBuiltin("to_base64", ToBase64),
	}, Create)
}
