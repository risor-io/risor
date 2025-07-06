package github

import (
	"context"
	"net/http"

	"github.com/google/go-github/v73/github"
	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
)

// Create creates a new GitHub client with optional authentication
func Create(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("github.client", 0, 1, args); err != nil {
		return err
	}

	var client *github.Client
	if len(args) == 0 {
		// Create unauthenticated client
		client = github.NewClient(nil)
	} else {
		// Create authenticated client with token
		token, err := object.AsString(args[0])
		if err != nil {
			return err
		}
		client = github.NewClient(&http.Client{
			Transport: &github.BasicAuthTransport{
				Username: "token",
				Password: token,
			},
		})
	}

	return New(client)
}

// Module returns the GitHub module definition
func Module() *object.Module {
	return object.NewBuiltinsModule("github", map[string]object.Object{
		"client": object.NewBuiltin("client", Create),
	}, Create)
}
