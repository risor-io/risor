package op

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetInfo(t *testing.T) {
	info := GetInfo(LoadClosure)
	require.Equal(t, "LOAD_CLOSURE", info.Name)
	require.Equal(t, 2, info.OperandCount)
	require.Equal(t, LoadClosure, info.Code)
}
