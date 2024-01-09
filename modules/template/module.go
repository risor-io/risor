package template

import (
	"context"
	"strings"

	"github.com/risor-io/risor/object"
)

func RenderTemplate(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)

	if numArgs != 2 {
		return object.NewArgsError("template", 2, numArgs)
	}

	data, argsErr := object.AsMap(args[0])
	if argsErr != nil {
		return argsErr
	}

	template, argsErr := object.AsString(args[1])
	if argsErr != nil {
		return argsErr
	}

	buf := new(strings.Builder)

	if err := Render(ctx, buf, template, data.Interface()); err != nil {
		return object.NewError(err)
	}

	return object.NewString(buf.String())
}

func Builtins() map[string]object.Object {
	return map[string]object.Object{
		"render": object.NewBuiltin("render", RenderTemplate),
	}
}

func Module() *object.Module {
	return object.NewBuiltinsModule("template", map[string]object.Object{
		"new": object.NewBuiltin("new", New),
	})
}
