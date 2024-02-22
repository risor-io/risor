package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestResponseWriter(t *testing.T) {
	svr := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := &ResponseWriter{writer: w}
			writeFnObj, ok := writer.GetAttr("write")
			require.True(t, ok)
			writeFn, ok := writeFnObj.(*object.Builtin)
			require.True(t, ok)
			result := writeFn.Call(context.Background(), object.NewString("FOO"))
			require.Equal(t, object.NewInt(3), result)
		}))
	defer svr.Close()
	r, err := http.NewRequest("GET", svr.URL, nil)
	require.Nil(t, err)
	w := httptest.NewRecorder()
	svr.Config.Handler.ServeHTTP(w, r)
	require.Equal(t, 200, w.Code)
	require.Equal(t, "FOO", w.Body.String())
}
