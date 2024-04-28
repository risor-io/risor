package bcrypt

import (
	"context"
	"strings"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestModule(t *testing.T) {
	m := Module()
	require.NotNil(t, m)

	hashObj, ok := m.GetAttr("hash")
	require.True(t, ok)

	hash, ok := hashObj.(*object.Builtin)
	require.True(t, ok)

	result := hash.Call(context.Background(), object.NewByteSlice([]byte("secretpw")))
	require.NotNil(t, result)

	s := string(result.Interface().([]byte))
	require.True(t, strings.HasPrefix(s, "$2a$"))
	require.Len(t, s, 60)

	compareObj, ok := m.GetAttr("compare")
	require.True(t, ok)

	compare, ok := compareObj.(*object.Builtin)
	require.True(t, ok)

	result = compare.Call(
		context.Background(),
		object.NewByteSlice([]byte(s)),
		object.NewByteSlice([]byte("secretpw")),
	)
	require.Equal(t, object.True, result)

	result = compare.Call(
		context.Background(),
		object.NewByteSlice([]byte(s)),
		object.NewByteSlice([]byte("wrongpw")),
	)
	resultErr, ok := result.(*object.Error)
	require.True(t, ok)
	require.Equal(t,
		"crypto/bcrypt: hashedPassword is not the hash of the given password",
		resultErr.Value().Error())
}
