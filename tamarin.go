// Package main is the entry point for the Tamarin CLI.
// A path to a Tamarin script should be provided as an
// argument to the program.
//
// Example:
//
//	$ cd path/to/tamarin
//	$ go build
//	$ ./tamarin ./examples/math.tm
//
// Tamarin may also be imported into another Go program
// to be used as a library. View the exec package for
// documentation on using Tamarin as a library.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"

	"github.com/cloudcmds/tamarin/evaluator"
	"github.com/cloudcmds/tamarin/exec"
	"github.com/cloudcmds/tamarin/object"
)

func main() {
	var profilerOutputPath string
	flag.StringVar(&profilerOutputPath, "profile", "", "Enable profiling")
	flag.Parse()

	if profilerOutputPath != "" {
		f, err := os.Create(profilerOutputPath)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Determine if input is being provided via stdin
	var isStdinInput bool
	if fi, err := os.Stdin.Stat(); err == nil {
		if fi.Mode()&os.ModeNamedPipe != 0 {
			isStdinInput = true
		}
	}

	// Input can only come from one source
	nArgs := len(flag.Args())
	if nArgs > 0 && isStdinInput {
		fmt.Fprintf(os.Stderr, "error: cannot provide both a script file and stdin input\n")
		os.Exit(1)
	} else if nArgs == 0 && !isStdinInput {
		fmt.Fprintf(os.Stderr, "error: expected one argument, a path to a tamarin script\n\n")
		fmt.Fprintf(os.Stderr, "example:\n ./tamarin ./examples/hello.tm\n\n")
		os.Exit(1)
	}

	// Read input
	var err error
	var input []byte
	if isStdinInput {
		input, err = io.ReadAll(os.Stdin)
	} else {
		input, err = os.ReadFile(flag.Args()[0])
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %s\n", err)
		os.Exit(1)
	}

	// Execute the script
	ctx := context.Background()
	result, err := exec.Execute(ctx, exec.Opts{
		Input:    string(input),
		Importer: &evaluator.SimpleImporter{},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// Print the result
	if result != object.NULL {
		fmt.Println(result.Inspect())
	}
}
