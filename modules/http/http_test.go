package http

import (
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestModule(t *testing.T) {
	m := Module()
	require.NotNil(t, m)
	reqObj, ok := m.GetAttr("request")
	require.True(t, ok)
	req, ok := reqObj.(*object.Builtin)
	require.True(t, ok)
	require.Equal(t, "http.request", req.Name())

	_, ok = m.GetAttr("listen_and_serve")
	require.False(t, ok)
}

func TestModuleWithListeners(t *testing.T) {
	m := Module(ModuleOpts{ListenersAllowed: true})
	require.NotNil(t, m)
	_, ok := m.GetAttr("listen_and_serve")
	require.True(t, ok)
}
