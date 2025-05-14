package importer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalImporter_Import(t *testing.T) {

	t.Run("successfully imports existing module", func(t *testing.T) {
		importer := NewLocalImporter(LocalImporterOptions{
			SourceDir: "fixtures",
		})

		module, err := importer.Import(context.Background(), "foo/bar")
		require.NoError(t, err)
		require.NotNil(t, module)
		require.Equal(t, "foo/bar", module.Name().Value())
		require.NotNil(t, module.Code())

		code := module.Code()
		require.Equal(t, []string{"bar"}, code.GlobalNames())
	})

	t.Run("returns error for nonexistent module", func(t *testing.T) {
		importer := NewLocalImporter(LocalImporterOptions{
			SourceDir: "fixtures",
		})

		module, err := importer.Import(context.Background(), "nonexistent")
		require.Error(t, err)
		require.Nil(t, module)
		require.Contains(t, err.Error(), "module \"nonexistent\" not found")
	})

	t.Run("uses custom extensions", func(t *testing.T) {
		importer := NewLocalImporter(LocalImporterOptions{
			SourceDir:  "fixtures/foo",
			Extensions: []string{".rriissoorr"},
		})

		module, err := importer.Import(context.Background(), "example")
		require.NoError(t, err)
		require.NotNil(t, module)
		require.Equal(t, "example", module.Name().Value())
	})
}
