package object_test

import (
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {
	b := object.NewBool(true)
	obj, ok := b.GetAttr("foo")
	require.False(t, ok)
	require.Nil(t, obj)

	// err := b.SetAttr("foo", object.NewInt(int64(1)))
	// require.Error(t, err)
}
