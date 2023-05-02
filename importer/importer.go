package importer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudcmds/tamarin/internal/compiler"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/parser"
)

// Importer is an interface used to import Tamarin code modules
type Importer interface {

	// Import a module by name
	Import(ctx context.Context, name string) (*object.Module, error)
}

type LocalImporter struct {
	builtins   map[string]object.Object
	modules    map[string]*object.Module
	sourceDir  string
	extensions []string
}

type LocalImporterOptions struct {
	Builtins   map[string]object.Object
	SourceDir  string
	Extensions []string
}

func NewLocalImporter(opts LocalImporterOptions) *LocalImporter {
	if opts.Builtins == nil {
		opts.Builtins = map[string]object.Object{}
	}
	if opts.Extensions == nil {
		opts.Extensions = []string{".tm"}
	}
	return &LocalImporter{
		builtins:   opts.Builtins,
		modules:    map[string]*object.Module{},
		sourceDir:  opts.SourceDir,
		extensions: opts.Extensions,
	}
}

func (i *LocalImporter) Import(ctx context.Context, name string) (*object.Module, error) {
	if m, ok := i.modules[name]; ok {
		return m, nil
	}
	mid := fmt.Sprintf("module: %s", name)
	source, found := readFileWithExtensions(i.sourceDir, name, i.extensions)
	if !found {
		return nil, fmt.Errorf("module not found: %s", name)
	}
	cmp := compiler.New(compiler.Options{
		Builtins: i.builtins,
		Name:     mid,
	})
	ast, err := parser.Parse(source)
	if err != nil {
		return nil, err
	}
	code, err := cmp.Compile(ast)
	if err != nil {
		return nil, err
	}
	return object.NewModule(name, code), nil
}

func readFileWithExtensions(dir, name string, extensions []string) (string, bool) {
	for _, ext := range extensions {
		fullPath := filepath.Join(dir, name+ext)
		bytes, err := os.ReadFile(fullPath)
		if err == nil {
			return string(bytes), true
		}
	}
	return "", false
}
