package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestResponseWriterWrite(t *testing.T) {
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

func TestResponseWriterWriteMap(t *testing.T) {
	svr := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := &ResponseWriter{writer: w}
			writeFnObj, ok := writer.GetAttr("write")
			require.True(t, ok)
			writeFn, ok := writeFnObj.(*object.Builtin)
			require.True(t, ok)
			result := writeFn.Call(context.Background(),
				object.NewMap(map[string]object.Object{
					"foo": object.NewString("FOO"),
				}))
			require.Equal(t, object.NewInt(13), result)
		}))
	defer svr.Close()
	r, err := http.NewRequest("GET", svr.URL, nil)
	require.Nil(t, err)
	w := httptest.NewRecorder()
	svr.Config.Handler.ServeHTTP(w, r)
	require.Equal(t, 200, w.Code)
	require.Equal(t, `{"foo":"FOO"}`, w.Body.String())
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestResponseWriterWriteHeader(t *testing.T) {
	svr := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := &ResponseWriter{writer: w}
			writeHeaderObj, ok := writer.GetAttr("write_header")
			require.True(t, ok)
			writeFn, ok := writeHeaderObj.(*object.Builtin)
			require.True(t, ok)
			result := writeFn.Call(context.Background(), object.NewInt(404))
			require.Equal(t, object.Nil, result)
		}))
	defer svr.Close()
	r, err := http.NewRequest("GET", svr.URL, nil)
	require.Nil(t, err)
	w := httptest.NewRecorder()
	svr.Config.Handler.ServeHTTP(w, r)
	require.Equal(t, 404, w.Code)
}

func TestResponseWriterMisc(t *testing.T) {
	svr := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := &ResponseWriter{writer: w}
			require.Equal(t, object.Type("http.response_writer"), writer.Type())
			_, ok := writer.Interface().(http.ResponseWriter)
			require.True(t, ok)
		}))
	defer svr.Close()
	r, err := http.NewRequest("GET", svr.URL, nil)
	require.Nil(t, err)
	w := httptest.NewRecorder()
	svr.Config.Handler.ServeHTTP(w, r)
	require.Equal(t, 200, w.Code)
}
