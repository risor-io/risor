package cfg

import (
	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/importer"
	"github.com/risor-io/risor/object"
)

type RisorConfig struct {
	Compiler        *compiler.Compiler
	Main            *object.Code
	Builtins        map[string]object.Object
	Importer        importer.Importer
	LocalImportPath string
	Offset          int
}
