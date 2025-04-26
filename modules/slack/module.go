package slack

import (
	"context"

	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
	"github.com/slack-go/slack"
)

func Create(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("slack.client", 1, 1, args); err != nil {
		return err
	}
	token, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	slack := slack.New(token)
	return New(slack)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("slack", map[string]object.Object{
		"client": object.NewBuiltin("client", Create),
	}, Create)
}
