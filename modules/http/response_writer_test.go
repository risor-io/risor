package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func recordGetResponse(t *testing.T, handler http.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	svr := httptest.NewServer(handler)
	defer svr.Close()
	r, err := http.NewRequest("GET", svr.URL, nil)
	require.Nil(t, err)
	w := httptest.NewRecorder()
	svr.Config.Handler.ServeHTTP(w, r)
	return w
}

func TestResponseWriterWrite(t *testing.T) {
	expectedBody := "FOO"
	resp := recordGetResponse(t, func(w http.ResponseWriter, r *http.Request) {
		writer := &ResponseWriter{writer: w}
		writeFnObj, ok := writer.GetAttr("write")
		require.True(t, ok)
		writeFn, ok := writeFnObj.(*object.Builtin)
		require.True(t, ok)
		result := writeFn.Call(context.Background(), object.NewString("FOO"))
		require.Equal(t, object.NewInt(int64(len(expectedBody))), result)
	})
	require.Equal(t, 200, resp.Code)
	require.Equal(t, expectedBody, resp.Body.String())
}

func TestResponseWriterWriteMap(t *testing.T) {
	expectedBody := `{"foo":"FOO"}`
	resp := recordGetResponse(t, func(w http.ResponseWriter, r *http.Request) {
		writer := &ResponseWriter{writer: w}
		writeFnObj, ok := writer.GetAttr("write")
		require.True(t, ok)
		writeFn, ok := writeFnObj.(*object.Builtin)
		require.True(t, ok)
		result := writeFn.Call(context.Background(),
			object.NewMap(map[string]object.Object{
				"foo": object.NewString("FOO"),
			}))
		require.Equal(t, object.NewInt(int64(len(expectedBody))), result)
	})
	require.Equal(t, 200, resp.Code)
	require.Equal(t, expectedBody, resp.Body.String())
	require.Equal(t, "application/json", resp.Header().Get("Content-Type"))
}

func TestResponseWriterWriteList(t *testing.T) {
	expectedBody := `["FOO"]`
	resp := recordGetResponse(t, func(w http.ResponseWriter, r *http.Request) {
		writer := &ResponseWriter{writer: w}
		writeFnObj, ok := writer.GetAttr("write")
		require.True(t, ok)
		writeFn, ok := writeFnObj.(*object.Builtin)
		require.True(t, ok)
		result := writeFn.Call(context.Background(),
			object.NewList([]object.Object{
				object.NewString("FOO"),
			}))
		require.Equal(t, object.NewInt(int64(len(expectedBody))), result)
	})
	require.Equal(t, 200, resp.Code)
	require.Equal(t, expectedBody, resp.Body.String())
	require.Equal(t, "application/json", resp.Header().Get("Content-Type"))
}

func TestResponseWriterWriteHeader(t *testing.T) {
	resp := recordGetResponse(t, func(w http.ResponseWriter, r *http.Request) {
		writer := &ResponseWriter{writer: w}
		writeHeaderObj, ok := writer.GetAttr("write_header")
		require.True(t, ok)
		writeFn, ok := writeHeaderObj.(*object.Builtin)
		require.True(t, ok)
		result := writeFn.Call(context.Background(), object.NewInt(404))
		require.Equal(t, object.Nil, result)
	})
	require.Equal(t, 404, resp.Code)
	require.Equal(t, "", resp.Body.String())
	require.Equal(t, "", resp.Header().Get("Content-Type"))
}

func TestResponseWriterMisc(t *testing.T) {
	expectedBody := "HEY"
	resp := recordGetResponse(t, func(w http.ResponseWriter, r *http.Request) {
		writer := &ResponseWriter{writer: w}
		require.Equal(t, object.Type("http.response_writer"), writer.Type())
		_, ok := writer.Interface().(http.ResponseWriter)
		require.True(t, ok)
		require.True(t, writer.IsTruthy())
		require.Equal(t, 0, writer.Cost())
		require.Equal(t, object.True, writer.Equals(writer))
		writer.AddHeader("foo", "bar")
		writer.AddHeader("x", "y")
		writer.DelHeader("x")
		writer.Write(object.NewByteSlice([]byte("HEY")))
	})
	require.Equal(t, 200, resp.Code)
	require.Equal(t, expectedBody, resp.Body.String())
	require.Equal(t, "bar", resp.Header().Get("foo"))
	require.Equal(t, "", resp.Header().Get("x"))
}

func TestResponseWriterInvalidType(t *testing.T) {
	resp := recordGetResponse(t, func(w http.ResponseWriter, r *http.Request) {
		writer := &ResponseWriter{writer: w}
		_, err := writer.Write(object.NewFloat(3.14))
		require.NotNil(t, err)
		require.Equal(t, "type error: unsupported type for http.response_writer.write: *object.Float", err.Error())
	})
	require.Equal(t, 500, resp.Code)
	require.Equal(t, "io error: failed to marshal response", resp.Body.String())
}
