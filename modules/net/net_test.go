package net

import (
	"context"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestModule(t *testing.T) {
	m := Module()
	require.NotNil(t, m)

	fnObj, ok := m.GetAttr("parse_cidr")
	require.True(t, ok)
	fn, ok := fnObj.(*object.Builtin)
	require.True(t, ok)

	result := fn.Call(context.Background(), object.NewString("10.2.11.12/16"))
	require.NotNil(t, result)
	net, ok := result.(*IPNet)
	require.True(t, ok)
	require.Equal(t, "net.ipnet(10.2.0.0/16)", net.Inspect())
}
