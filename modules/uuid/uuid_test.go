package uuid

import (
	"context"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestModule(t *testing.T) {
	m := Module()
	require.NotNil(t, m)

	versions := []string{"v4", "v6", "v7"}

	for _, version := range versions {
		fnObj, ok := m.GetAttr(version)
		require.True(t, ok)
		fn, ok := fnObj.(*object.Builtin)
		require.True(t, ok)
		result := fn.Call(context.Background())
		require.NotNil(t, result)
		uuidObj, ok := result.(*object.String)
		require.True(t, ok)
		require.Len(t, uuidObj.Value(), 36)
	}
}

func TestV5(t *testing.T) {
	m := Module()
	require.NotNil(t, m)

	fnObj, ok := m.GetAttr("v5")
	require.True(t, ok)
	fn, ok := fnObj.(*object.Builtin)
	require.True(t, ok)

	namespace := object.NewString("64bcbbf7-9bb8-4aee-b708-68bae49b3306")
	name := object.NewString("joe")

	result := fn.Call(context.Background(), namespace, name)
	require.NotNil(t, result)
	uuidObj, ok := result.(*object.String)
	require.True(t, ok)
	require.Equal(t, "a9f24ca5-222b-5f69-bb9c-ea34f555a295", uuidObj.Value())
}
