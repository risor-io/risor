package cfg

import (
	"sort"

	"github.com/risor-io/risor/builtins"
	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/importer"
	modBase64 "github.com/risor-io/risor/modules/base64"
	modBytes "github.com/risor-io/risor/modules/bytes"
	modFetch "github.com/risor-io/risor/modules/fetch"
	modFmt "github.com/risor-io/risor/modules/fmt"
	modHash "github.com/risor-io/risor/modules/hash"
	modJson "github.com/risor-io/risor/modules/json"
	modMath "github.com/risor-io/risor/modules/math"
	modOs "github.com/risor-io/risor/modules/os"
	modRand "github.com/risor-io/risor/modules/rand"
	modStrconv "github.com/risor-io/risor/modules/strconv"
	modStrings "github.com/risor-io/risor/modules/strings"
	modTime "github.com/risor-io/risor/modules/time"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/vm"
)

// RisorConfig assists in configuring a Risor evaluation.
type RisorConfig struct {
	Globals               map[string]any
	DefaultGlobals        map[string]object.Object
	Importer              importer.Importer
	LocalImportPath       string
	WithoutDefaultGlobals bool
}

func NewRisorConfig() *RisorConfig {
	cfg := &RisorConfig{
		Globals:        map[string]any{},
		DefaultGlobals: map[string]object.Object{},
	}
	cfg.addDefaultGlobals()
	return cfg
}

// CombinedGlobals returns a map of all global variables that should be
// available in a Risor evaluation.
func (cfg *RisorConfig) CombinedGlobals() map[string]any {
	combined := map[string]any{}
	for k, v := range cfg.DefaultGlobals {
		combined[k] = v
	}
	for k, v := range cfg.Globals {
		combined[k] = v
	}
	return combined
}

// GlobalNames returns a list of all global variables names that should be
// available in a Risor evaluation.
func (cfg *RisorConfig) GlobalNames() []string {
	nameMap := map[string]bool{}
	for k := range cfg.DefaultGlobals {
		nameMap[k] = true
	}
	for k := range cfg.Globals {
		nameMap[k] = true
	}
	var names []string
	for name := range nameMap {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (cfg *RisorConfig) addDefaultGlobals() {
	addGlobals := func(globals map[string]object.Object) {
		for k, v := range globals {
			cfg.DefaultGlobals[k] = v
		}
	}
	// Add default builtin functions
	builtins := []map[string]object.Object{
		builtins.Builtins(),
		modFetch.Builtins(),
		modFmt.Builtins(),
		modHash.Builtins(),
		modOs.Builtins(),
	}
	for _, b := range builtins {
		addGlobals(b)
	}
	// Add default modules
	modules := map[string]object.Object{
		"math":    modMath.Module(),
		"json":    modJson.Module(),
		"strings": modStrings.Module(),
		"time":    modTime.Module(),
		"rand":    modRand.Module(),
		"strconv": modStrconv.Module(),
		"os":      modOs.Module(),
		"bytes":   modBytes.Module(),
		"base64":  modBase64.Module(),
		"fmt":     modFmt.Module(),
	}
	addGlobals(modules)
}

// CompilerOpts returns compiler options derived from this configuration.
func (cfg *RisorConfig) CompilerOpts() []compiler.Option {
	globalNames := cfg.GlobalNames()
	var opts []compiler.Option
	if len(globalNames) > 0 {
		opts = append(opts, compiler.WithGlobalNames(globalNames))
	}
	return opts
}

// VMOpts returns virtual machine options derived from this configuration.
func (cfg *RisorConfig) VMOpts() []vm.Option {
	var opts []vm.Option
	combinedGlobals := cfg.CombinedGlobals()
	if len(combinedGlobals) > 0 {
		opts = append(opts, vm.WithGlobals(combinedGlobals))
	}
	importer := cfg.Importer
	if importer == nil && cfg.LocalImportPath != "" {
		var names []string
		for name := range combinedGlobals {
			names = append(names, name)
		}
		importer = newLocalImporter(names, cfg.LocalImportPath)
	}
	if importer != nil {
		opts = append(opts, vm.WithImporter(importer))
	}
	return opts
}

func newLocalImporter(globalNames []string, sourceDir string) importer.Importer {
	return importer.NewLocalImporter(importer.LocalImporterOptions{
		GlobalNames: globalNames,
		SourceDir:   sourceDir,
		Extensions:  []string{".risor", ".rsr"},
	})
}

// func WithCompiler(c *compiler.Compiler) Option {
// 	return func(r *cfg.RisorConfig) {
// 		r.Compiler = c
// 	}
// }

// func WithCode(c *compiler.Code) Option {
// 	return func(r *cfg.RisorConfig) {
// 		r.Code = c
// 	}
// }

// func WithInstructionOffset(offset int) Option {
// 	return func(r *cfg.RisorConfig) {
// 		r.InstructionOffset = offset
// 	}
// }
