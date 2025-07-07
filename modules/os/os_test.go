package os

import (
	"context"
	"os"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestStdio(t *testing.T) {
	m := Module()

	type testCase struct {
		name     string
		attr     string
		expected *os.File
	}

	tests := []testCase{
		{
			name:     "stdin",
			attr:     "stdin",
			expected: os.Stdin,
		},
		{
			name:     "stdout",
			attr:     "stdout",
			expected: os.Stdout,
		},
		{
			name:     "stderr",
			attr:     "stderr",
			expected: os.Stderr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr, ok := m.GetAttr(tt.attr)
			require.True(t, ok)

			stdin, ok := attr.(*object.DynamicAttr)
			require.True(t, ok)

			result, err := stdin.ResolveAttr(context.Background(), "")
			require.NoError(t, err)

			file, ok := result.(*object.File)
			require.True(t, ok)
			require.Equal(t, tt.expected, file.Value())
		})
	}
}
