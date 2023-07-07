package object_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestMarshalFile(t *testing.T) {
	ctx := context.Background()
	f := object.NewFile(ctx, nil, "")
	bytes, err := json.Marshal(f)
	require.Nil(t, bytes)
	require.NotNil(t, err)
	require.Equal(t, "json: error calling MarshalJSON for type *object.File: type error: unable to marshal file", err.Error())
}
