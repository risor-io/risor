package doctest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/risor-io/risor/rdoc"
	"github.com/stretchr/testify/require"
)

func hasModuleFunc(dir string) bool {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".go") {
			continue
		}
		textData, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			continue
		}
		if strings.Contains(string(textData), "func Module() *object.Module") {
			return true
		}
	}
	return false
}

func TestModuleDocs(t *testing.T) {
	// All modules that have a Module() function should also have a markdown
	// documentation file in the same directory
	mods, err := os.ReadDir("../modules")
	require.Nil(t, err)
	for _, mod := range mods {
		if mod.IsDir() {
			name := mod.Name()
			modPath := filepath.Join("..", "modules", name)
			if !hasModuleFunc(modPath) {
				continue
			}
			t.Run(name, func(t *testing.T) {
				filename := fmt.Sprintf("%s.md", name)
				docPath := filepath.Join("..", "modules", name, filename)
				md, err := os.ReadFile(docPath)
				require.Nil(t, err, "Expected module markdown doc to exist at %s", docPath)
				doc := rdoc.Parse(string(md))
				require.NotNil(t, doc)
				require.Equal(t, name, doc.Name)
			})
		}
	}
}
