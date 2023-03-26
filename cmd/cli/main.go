package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/cloudcmds/tamarin/exec"
	"github.com/cloudcmds/tamarin/internal/vm"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/parser"
	"github.com/cloudcmds/tamarin/repl"
	"github.com/cloudcmds/tamarin/scope"
	"github.com/fatih/color"
)

func main() {
	var noColor bool
	var profilerOutputPath, code, breakpoints string
	flag.BoolVar(&noColor, "no-color", false, "Disable color output")
	flag.StringVar(&code, "c", "", "Code to execute")
	flag.StringVar(&profilerOutputPath, "profile", "", "Enable profiling")
	flag.StringVar(&breakpoints, "breakpoints", "", "Comma-separated list of breakpoints")
	flag.Parse()

	if noColor {
		color.NoColor = true
	}
	red := color.New(color.FgRed).SprintfFunc()

	if profilerOutputPath != "" {
		f, err := os.Create(profilerOutputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", red(err.Error()))
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	ctx := context.Background()
	globalScope := scope.New(scope.Opts{Name: "global"})
	if err := exec.AutoImport(globalScope, nil, nil); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", red(err.Error()))
		os.Exit(1)
	}

	// Input can only come from one source
	nArgs := len(flag.Args())
	if nArgs > 0 && len(code) > 0 {
		fmt.Fprintf(os.Stderr, "%s\n", red("error: cannot provide both a script file and -c input\n"))
		os.Exit(1)
	} else if nArgs == 0 && len(code) == 0 {
		// Run REPL
		if err := repl.Run(ctx, globalScope); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", red(err.Error()))
			os.Exit(1)
		}
		return
	}

	// Otherwise, use input from either -c or the first argument
	var err error
	var input string
	var filename string
	if nArgs == 0 {
		input = code
	} else {
		filename = flag.Args()[0]
		bytes, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", red(err.Error()))
			os.Exit(1)
		}
		input = string(bytes)
	}

	result, err := vm.Run(string(input))
	if err != nil {
		parserErr, ok := err.(parser.ParserError)
		if ok {
			fmt.Fprintf(os.Stderr, "%s\n", red(parserErr.FriendlyMessage()))
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", red(err.Error()))
		}
		os.Exit(1)
	}

	// Print the result
	if result != object.Nil {
		fmt.Println(result.Inspect())
	}
}
