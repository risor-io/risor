package jmespath

import (
	"context"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/api"
	"github.com/jmespath-community/go-jmespath/pkg/parsing"
	"github.com/risor-io/risor/object"
)

func Jmespath(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)

	if numArgs != 2 {
		return object.NewArgsError("jmespath", 2, numArgs)
	}

	data, argsErr := object.AsMap(args[0])
	if argsErr != nil {
		return argsErr
	}

	expression, argsErr := object.AsString(args[1])
	if argsErr != nil {
		return argsErr
	}

	parser := parsing.NewParser()
	if _, err := parser.Parse(expression); err != nil {
		if syntaxError, ok := err.(parsing.SyntaxError); ok {
			return object.NewError(fmt.Errorf("%s\n%s", syntaxError, syntaxError.HighlightLocation()))
		}
		return object.NewError(err)
	}
	result, err := api.Search(expression, data.Interface())
	if argsErr != nil {
		return object.NewError(err)
	}

	return object.FromGoType(result)
}

func Builtins() map[string]object.Object {
	return map[string]object.Object{
		"jmespath": object.NewBuiltin("jmespath", Jmespath),
	}
}
