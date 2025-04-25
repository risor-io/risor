package goquery

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
)

func Parse(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("goquery.parse", 1, args); err != nil {
		return err
	}

	var reader io.Reader

	switch arg := args[0].(type) {
	case *object.String:
		reader = strings.NewReader(arg.Value())
	case *object.ByteSlice:
		reader = bytes.NewReader(arg.Value())
	case *object.File:
		reader = arg.Value()
	case io.Reader:
		reader = arg
	default:
		return object.TypeErrorf("type error: expected reader (got %s)", args[0].Type())
	}

	doc, err := NewDocumentFromReader(reader)
	if err != nil {
		return object.NewError(err)
	}
	return doc
}

func Module() *object.Module {
	return object.NewBuiltinsModule("goquery", map[string]object.Object{
		"parse": object.NewBuiltin("parse", Parse),
	}, Parse)
}
