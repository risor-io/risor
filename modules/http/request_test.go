package http

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestRequestPathValueEmpty(t *testing.T) {
	ctx := context.Background()
	u, err := url.Parse("http://example.com/?foo=bar")
	require.Nil(t, err)

	req := NewRequest(&http.Request{Method: "GET", URL: u})
	require.NotNil(t, req)

	fn, ok := req.GetAttr("path_value")
	require.True(t, ok)

	v, ok := req.GetAttr("query")
	require.True(t, ok)
	require.Equal(t, object.NewString("bar"), v.(*object.Map).Get("foo"))

	pathValue, ok := fn.(*object.Builtin)
	require.True(t, ok)

	result := pathValue.Call(ctx, object.NewString("foo"))
	require.Equal(t, object.NewString(""), result)
}
