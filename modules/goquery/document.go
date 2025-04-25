package goquery

import (
	"context"
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

var _ object.Object = (*Document)(nil)

const DOCUMENT object.Type = "goquery.document"

type Document struct {
	value *goquery.Document
}

func (d *Document) Value() *goquery.Document {
	return d.value
}

func (d *Document) Type() object.Type {
	return DOCUMENT
}

func (d *Document) Inspect() string {
	return fmt.Sprintf("%s()", DOCUMENT)
}

func (d *Document) IsTruthy() bool {
	return d.value != nil
}

func (d *Document) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("type error: cannot set %q on %s object", name, DOCUMENT)
}

func (d *Document) String() string {
	if d.value == nil {
		return "nil"
	}
	html, err := d.value.Html()
	if err != nil {
		return fmt.Sprintf("%s()", DOCUMENT)
	}
	return html
}

func (d *Document) Interface() interface{} {
	return d.value
}

func (d *Document) Equals(other object.Object) object.Object {
	otherDoc, ok := other.(*Document)
	if !ok {
		return object.False
	}
	return object.NewBool(otherDoc.Value() == d.value)
}

func (d *Document) Cost() int {
	return 0
}

func (d *Document) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for %s: %v", DOCUMENT, opType)
}

func (d *Document) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "find":
		return object.NewBuiltin("find", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.document.find", 1, args); err != nil {
				return err
			}
			selector, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			return NewSelection(d.value.Find(selector))
		}), true
	case "html":
		return object.NewBuiltin("html", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.document.html", 0, args); err != nil {
				return err
			}
			html, err := d.value.Html()
			if err != nil {
				return object.NewError(err)
			}
			return object.NewString(html)
		}), true
	case "text":
		return object.NewBuiltin("text", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.document.text", 0, args); err != nil {
				return err
			}
			return object.NewString(d.value.Text())
		}), true
	default:
		return nil, false
	}
}

func NewDocument(doc *goquery.Document) *Document {
	return &Document{value: doc}
}

func NewDocumentFromReader(r io.Reader) (*Document, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	return NewDocument(doc), nil
}
