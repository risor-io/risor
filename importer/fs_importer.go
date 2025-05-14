package importer

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"sync"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/object"
)

// FSImporterOptions configure an Importer that can read from a filesystem
// implementing the `fs.FS` interface.
type FSImporterOptions struct {
	// Global names that should be available when the module is compiled.
	GlobalNames []string

	// The filesystem to search for Risor modules.
	SourceFS fs.FS

	// Optional list of file extensions to try when locating a Risor module.
	Extensions []string
}

// FSImporter is an Importer that can read Risor code modules from a filesystem
// implementing the `fs.FS` interface.
type FSImporter struct {
	globalNames []string
	codeCache   map[string]*compiler.Code
	sourceFS    fs.FS
	extensions  []string
	mutex       sync.Mutex
}

// NewFSImporter returns an Importer that can read Risor code modules from a
// filesystem implementing the `fs.FS` interface.
func NewFSImporter(opts FSImporterOptions) *FSImporter {
	if opts.Extensions == nil {
		opts.Extensions = defaultExtensions
	}
	return &FSImporter{
		globalNames: opts.GlobalNames,
		codeCache:   map[string]*compiler.Code{},
		sourceFS:    opts.SourceFS,
		extensions:  opts.Extensions,
	}
}

// Import a module by name.
func (i *FSImporter) Import(ctx context.Context, name string) (*object.Module, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if code, ok := i.codeCache[name]; ok {
		return object.NewModule(name, code), nil
	}

	source, fullPath, found := i.readFileWithExtensions(name, i.extensions)
	if !found {
		return nil, fmt.Errorf("import error: module %q not found", name)
	}

	code, err := parseAndCompile(ctx, source, fullPath, i.globalNames)
	if err != nil {
		return nil, err
	}

	i.codeCache[name] = code

	return object.NewModule(name, code), nil
}

func (i *FSImporter) readFileWithExtensions(name string, extensions []string) (string, string, bool) {
	for _, ext := range extensions {
		fullName := name + ext
		f, err := i.sourceFS.Open(fullName)
		if err != nil {
			continue
		}

		b, err := io.ReadAll(f)
		if err != nil {
			f.Close()
			continue
		}

		f.Close()
		return string(b), fullName, true
	}
	return "", "", false
}
