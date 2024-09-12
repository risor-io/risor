package fmt

import (
	"context"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/os"
	"github.com/stretchr/testify/require"
)

func TestPrintln(t *testing.T) {
	ctx := context.Background()
	stdout := os.NewInMemoryFile(nil)
	ros := os.NewVirtualOS(ctx, os.WithStdout(stdout))
	ctx = os.WithOS(ctx, ros)
	result := Println(ctx, object.NewString("hello"), object.NewString("world"))
	require.Equal(t, object.Nil, result)
	stdout.Rewind()
	require.Equal(t, "hello world\n", string(stdout.Bytes()))
}

func TestPrintf(t *testing.T) {
	ctx := context.Background()
	stdout := os.NewInMemoryFile(nil)
	ros := os.NewVirtualOS(ctx, os.WithStdout(stdout))
	ctx = os.WithOS(ctx, ros)
	result := Printf(ctx, object.NewString("hello %s\n"), object.NewString("world"))
	require.Equal(t, object.Nil, result)
	stdout.Rewind()
	require.Equal(t, "hello world\n", string(stdout.Bytes()))
}

func TestErrorf(t *testing.T) {
	ctx := context.Background()
	result := Errorf(ctx, object.NewString("hello %s\n"), object.NewString("world"))
	require.IsType(t, &object.Error{}, result)
	require.Equal(t, "hello world\n", result.(*object.Error).Message().Value())
}
