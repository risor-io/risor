package qrcode

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

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

// generateQRCode generates QR code data into a buffer
func (q *QRCode) generateQRCode(width uint8) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	writerCloser := &nopCloser{Writer: buf}
	writer := standard.NewWithWriter(writerCloser, standard.WithQRWidth(width))

	if err := q.value.Save(writer); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Save implements the Save method of the QRCode type as a Risor method
func (q *QRCode) Save(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=1", len(args)))
	}
	writerObj := args[0]

	// If the writer is provided by the standard module, use it directly
	if stdWriter, ok := writerObj.Interface().(standard.Writer); ok {
		err := q.value.Save(stdWriter)
		if err != nil {
			return object.NewError(err)
		}
		return object.Nil
	}

	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	// Generate QR code data
	qrData, genErr := q.generateQRCode(q.width)
	if genErr != nil {
		return object.NewError(genErr)
	}

	// Use Risor OS to write the buffer to a file
	osObj := os.GetDefaultOS(ctx)
	if writeErr := osObj.WriteFile(path, qrData, 0644); writeErr != nil {
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

func (q *QRCode) Bytes(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=0", len(args)))
	}

	qrData, err := q.generateQRCode(q.width)
	if err != nil {
		return object.NewError(err)
	}

	return object.NewByteSlice(qrData)
}

func (q *QRCode) Base64(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=0", len(args)))
	}

	qrData, err := q.generateQRCode(q.width)
	if err != nil {
		return object.NewError(err)
	}

	return object.NewString(base64.StdEncoding.EncodeToString(qrData))
}
