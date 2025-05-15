package main

import (
	"errors"
	"io"
	"os"

	"github.com/risor-io/risor"
	"github.com/risor-io/risor/modules/aws"
	"github.com/risor-io/risor/modules/bcrypt"
	"github.com/risor-io/risor/modules/cli"
	"github.com/risor-io/risor/modules/color"
	"github.com/risor-io/risor/modules/gha"
	"github.com/risor-io/risor/modules/goquery"
	"github.com/risor-io/risor/modules/htmltomarkdown"
	"github.com/risor-io/risor/modules/image"
	"github.com/risor-io/risor/modules/isatty"
	"github.com/risor-io/risor/modules/jmespath"
	k8s "github.com/risor-io/risor/modules/kubernetes"
	"github.com/risor-io/risor/modules/net"
	"github.com/risor-io/risor/modules/pgx"
	"github.com/risor-io/risor/modules/playwright"
	"github.com/risor-io/risor/modules/qrcode"
	"github.com/risor-io/risor/modules/sched"
	"github.com/risor-io/risor/modules/semver"
	"github.com/risor-io/risor/modules/shlex"
	"github.com/risor-io/risor/modules/slack"
	"github.com/risor-io/risor/modules/sql"
	"github.com/risor-io/risor/modules/tablewriter"
	"github.com/risor-io/risor/modules/template"
	"github.com/risor-io/risor/modules/uuid"
	"github.com/risor-io/risor/modules/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Returns a Risor option for global variable configuration.
func getGlobals() risor.Option {
	if viper.GetBool("no-default-globals") {
		return risor.WithoutDefaultGlobals()
	}

	//************************************************************************//
	// Default modules
	//************************************************************************//

	globals := map[string]any{
		"bcrypt":         bcrypt.Module(),
		"cli":            cli.Module(),
		"color":          color.Module(),
		"gha":            gha.Module(),
		"goquery":        goquery.Module(),
		"htmltomarkdown": htmltomarkdown.Module(),
		"image":          image.Module(),
		"isatty":         isatty.Module(),
		"net":            net.Module(),
		"pgx":            pgx.Module(),
		"playwright":     playwright.Module(),
		"qrcode":         qrcode.Module(),
		"sched":          sched.Module(),
		"sql":            sql.Module(),
		"slack":          slack.Module(),
		"tablewriter":    tablewriter.Module(),
		"template":       template.Module(),
		"uuid":           uuid.Module(),
		"semver":         semver.Module(),
		"shlex":          shlex.Module(),
	}

	//************************************************************************//
	// Modules that contribute top-level built-in functions
	//************************************************************************//

	for k, v := range jmespath.Builtins() {
		globals[k] = v
	}
	for k, v := range template.Builtins() {
		globals[k] = v
	}

	//************************************************************************//
	// Modules which are optionally present (depending on build tags).
	// If the build tag is not set then the returned module is nil.
	//************************************************************************//

	if mod := aws.Module(); mod != nil {
		globals["aws"] = mod
	}
	if mod := k8s.Module(); mod != nil {
		globals["k8s"] = mod
	}
	if mod := vault.Module(); mod != nil {
		globals["vault"] = mod
	}

	return risor.WithGlobals(globals)
}

func getRisorOptions() []risor.Option {
	opts := []risor.Option{
		risor.WithConcurrency(),
		risor.WithListenersAllowed(),
		getGlobals(),
	}
	if modulesDir := viper.GetString("modules"); modulesDir != "" {
		opts = append(opts, risor.WithLocalImporter(modulesDir))
	}
	return opts
}

func shouldRunRepl(cmd *cobra.Command, args []string) bool {
	if viper.GetBool("no-repl") || viper.GetBool("stdin") {
		return false
	}
	if cmd.Flags().Lookup("code").Changed {
		return false
	}
	if len(args) > 0 {
		return false
	}
	return isTerminalIO()
}

func getRisorCode(cmd *cobra.Command, args []string) (string, error) {
	// Determine what code is to be executed. There three possibilities:
	// 1. --code <code>
	// 2. --stdin (read code from stdin)
	// 3. path as args[0]
	var codeFlagSet bool
	if f := cmd.Flags().Lookup("code"); f != nil && f.Changed {
		codeFlagSet = true
	}
	var stdinFlagSet bool
	if f := cmd.Flags().Lookup("stdin"); f != nil && f.Changed {
		stdinFlagSet = true
	}
	pathSupplied := len(args) > 0
	// Error if multiple input sources are specified
	if pathSupplied && (codeFlagSet || stdinFlagSet) {
		return "", errors.New("multiple input sources specified")
	} else if codeFlagSet && stdinFlagSet {
		return "", errors.New("multiple input sources specified")
	}
	if stdinFlagSet {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(data), nil
	} else if pathSupplied {
		bytes, err := os.ReadFile(args[0])
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
	return viper.GetString("code"), nil
}
