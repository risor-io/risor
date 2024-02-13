package risor

import (
	"sort"

	"github.com/risor-io/risor/builtins"
	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/importer"
	modBase64 "github.com/risor-io/risor/modules/base64"
	modBytes "github.com/risor-io/risor/modules/bytes"
	modDns "github.com/risor-io/risor/modules/dns"
	modExec "github.com/risor-io/risor/modules/exec"
	modFilepath "github.com/risor-io/risor/modules/filepath"
	modFmt "github.com/risor-io/risor/modules/fmt"
	modHTTP "github.com/risor-io/risor/modules/http"
	modJSON "github.com/risor-io/risor/modules/json"
	modMath "github.com/risor-io/risor/modules/math"
	modOs "github.com/risor-io/risor/modules/os"
	modRand "github.com/risor-io/risor/modules/rand"
	modRegexp "github.com/risor-io/risor/modules/regexp"
	modStrconv "github.com/risor-io/risor/modules/strconv"
	modStrings "github.com/risor-io/risor/modules/strings"
	modTime "github.com/risor-io/risor/modules/time"
	modYAML "github.com/risor-io/risor/modules/yaml"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/vm"
)

// Config assists in configuring a Risor evaluation.
type Config struct {
	Globals               map[string]any
	DefaultGlobals        map[string]object.Object
	Importer              importer.Importer
	LocalImportPath       string
	WithoutDefaultGlobals bool
	WithConcurrency       bool
}

func NewConfig() *Config {
	cfg := &Config{
		Globals:        map[string]any{},
		DefaultGlobals: map[string]object.Object{},
	}
	cfg.addDefaultGlobals()
	return cfg
}

// CombinedGlobals returns a map of all global variables that should be
// available in a Risor evaluation.
func (cfg *Config) CombinedGlobals() map[string]any {
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
func (cfg *Config) GlobalNames() []string {
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

func (cfg *Config) addDefaultGlobals() {
	addGlobals := func(globals map[string]object.Object) {
		for k, v := range globals {
			cfg.DefaultGlobals[k] = v
		}
	}
	// Add default builtin functions
	builtins := []map[string]object.Object{
		builtins.Builtins(),
		modHTTP.Builtins(),
		modFmt.Builtins(),
		modOs.Builtins(),
		modDns.Builtins(),
	}
	for _, b := range builtins {
		addGlobals(b)
	}
	// Add default modules
	modules := map[string]object.Object{
		"base64":   modBase64.Module(),
		"bytes":    modBytes.Module(),
		"exec":     modExec.Module(),
		"filepath": modFilepath.Module(),
		"fmt":      modFmt.Module(),
		"http":     modHTTP.Module(),
		"json":     modJSON.Module(),
		"math":     modMath.Module(),
		"os":       modOs.Module(),
		"rand":     modRand.Module(),
		"regexp":   modRegexp.Module(),
		"strconv":  modStrconv.Module(),
		"strings":  modStrings.Module(),
		"time":     modTime.Module(),
		"yaml":     modYAML.Module(),
	}
	addGlobals(modules)
}

// CompilerOpts returns compiler options derived from this configuration.
func (cfg *Config) CompilerOpts() []compiler.Option {
	globalNames := cfg.GlobalNames()
	var opts []compiler.Option
	if len(globalNames) > 0 {
		opts = append(opts, compiler.WithGlobalNames(globalNames))
	}
	return opts
}

// VMOpts returns virtual machine options derived from this configuration.
func (cfg *Config) VMOpts() []vm.Option {
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
	if cfg.WithConcurrency {
		opts = append(opts, vm.WithConcurrency())
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
