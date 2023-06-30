package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/risor-io/risor/limits"
	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	expected := "1234567890"
	gotHeaders := make(http.Header)
	gotMethod := ""
	var gotBody []byte
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			gotHeaders[k] = v
		}
		gotMethod = r.Method
		gotBody, _ = io.ReadAll(r.Body)
		fmt.Fprint(w, expected)
	}))
	defer svr.Close()

	ctx := context.Background()
	ctx = limits.WithLimits(ctx, limits.New())

	result := Fetch(ctx, object.NewString(svr.URL), object.NewMap(map[string]object.Object{
		"method":  object.NewString("PATCH"),
		"timeout": object.NewInt(1000),
		"body":    object.NewString("dummy body"),
		"headers": object.NewMap(map[string]object.Object{
			"Content-Type": object.NewString("application/json"),
			"Foo":          object.NewString("bar"),
		}),
	}))
	if errObj, ok := result.(*object.Error); ok {
		require.Nil(t, errObj, errObj)
	}
	resp, ok := result.(*object.HttpResponse)
	require.True(t, ok)

	require.Equal(t, int64(200), resp.StatusCode().Value())
	require.Equal(t, "200 OK", resp.Status().Value())
	require.Equal(t, int64(len(expected)), resp.ContentLength().Value())
	require.Equal(t, "dummy body", string(gotBody))
	require.Equal(t, "PATCH", gotMethod)
	require.Equal(t, http.Header{
		"Content-Length":  []string{"10"},
		"Content-Type":    []string{"application/json"},
		"Accept-Encoding": []string{"gzip"},
		"User-Agent":      []string{"Go-http-client/1.1"},
		"Foo":             []string{"bar"},
	}, gotHeaders)
}

func TestBasicFetch(t *testing.T) {
	gotHeaders := make(http.Header)
	gotMethod := ""
	var gotBody []byte
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			gotHeaders[k] = v
		}
		gotMethod = r.Method
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer svr.Close()

	ctx := context.Background()
	ctx = limits.WithLimits(ctx, limits.New())

	result := Fetch(ctx, object.NewString(svr.URL))
	if errObj, ok := result.(*object.Error); ok {
		require.Nil(t, errObj, errObj)
	}
	resp, ok := result.(*object.HttpResponse)
	require.True(t, ok)

	require.Equal(t, int64(204), resp.StatusCode().Value())
	require.Equal(t, "204 No Content", resp.Status().Value())
	require.Equal(t, int64(0), resp.ContentLength().Value())
	require.Equal(t, "", string(gotBody))
	require.Equal(t, "GET", gotMethod)
	require.Equal(t, http.Header{
		"Accept-Encoding": []string{"gzip"},
		"User-Agent":      []string{"Go-http-client/1.1"},
	}, gotHeaders)
}
