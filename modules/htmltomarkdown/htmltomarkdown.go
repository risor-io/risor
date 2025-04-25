package htmltomarkdown

import (
	"context"

	htmltomd "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/risor-io/risor/object"
)

func Convert(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("convert", 1, len(args))
	}
	html, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	markdown, convErr := htmltomd.ConvertString(html)
	if convErr != nil {
		return object.NewError(convErr)
	}
	return object.NewString(markdown)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("htmltomarkdown", map[string]object.Object{
		"convert": object.NewBuiltin("convert", Convert),
	})
}
