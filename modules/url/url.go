package url

import (
	"context"
	"net/url"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

func QueryEscape(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("query_escape", 1, args); err != nil {
		return err
	}
	q, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	s := url.QueryEscape(q)
	return object.NewString(s)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("url", map[string]object.Object{
		"query_escape": object.NewBuiltin("v7", QueryEscape),
	}, QueryEscape)
}
