package limits

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLimitsContext(t *testing.T) {
	ctx := WithLimits(context.Background(), New(WithMaxCost(13)))
	lim, ok := GetLimits(ctx)
	require.True(t, ok)
	require.IsType(t, &StandardLimits{}, lim)

	err := TrackCost(ctx, 12)
	require.Nil(t, err)

	err = TrackCost(ctx, 2)
	require.Error(t, err)
	require.Equal(t, "limit error: reached maximum processing cost (13)", err.Error())

	require.Nil(t, TrackCost(context.Background(), 1))
}

func TestLimitReadAll(t *testing.T) {
	buf := bytes.NewBufferString("1234567890")
	_, err := ReadAll(buf, NoLimit)
	require.Nil(t, err)

	buf = bytes.NewBufferString("1234567890")
	data, err := ReadAll(buf, 20)
	require.Nil(t, err)
	require.Equal(t, "1234567890", string(data))
}
