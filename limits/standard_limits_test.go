package limits

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLimits(t *testing.T) {
	l := New(
		WithMaxCost(10),
		WithIOTimeout(time.Second),
		WithMaxBufferSize(20),
		WithMaxHttpRequestCount(1),
	)
	require.Equal(t, time.Second, l.IOTimeout())
	require.Equal(t, int64(20), l.MaxBufferSize())

	require.Nil(t, l.TrackCost(5))
	require.Nil(t, l.TrackCost(5))

	err := l.TrackCost(5)
	require.Error(t, err)
	require.Equal(t, err.Error(), "limit error: reached maximum processing cost (10)")

	_, err = l.ReadAll(bytes.NewBufferString("123456789012345678901"))
	require.Error(t, err)
	require.Equal(t, "limit error: data size exceeded limit of 20 bytes", err.Error())

	req := &http.Request{}
	require.Nil(t, l.TrackHTTPRequest(req))
	err = l.TrackHTTPRequest(req)
	require.Error(t, err)
	require.Equal(t, "limit error: reached maximum number of http requests (1)", err.Error())
}
