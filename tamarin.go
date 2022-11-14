// Package main is the entry point for the Tamarin CLI.
// A path to a Tamarin script should be provided as an
// argument to the program.
//
// Example:
//
//	$ cd path/to/tamarin
//	$ go build
//	$ ./tamarin path/to/my/script.tm
//
// Tamarin may also be imported into another Go program
// to be used as a library. View the `exec` package for
// documentation on using Tamarin as a library.
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"

	"github.com/cloudcmds/tamarin/exec"
	"github.com/cloudcmds/tamarin/internal/evaluator"
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

	// Add print statement support
	evaluator.RegisterPrintBuiltins()

	// Read input script from a file or from stdin
	var err error
	var input []byte
	if len(flag.Args()) > 0 {
		input, err = ioutil.ReadFile(flag.Args()[0])
	} else {
		input, err = ioutil.ReadAll(os.Stdin)
	}
	if err != nil {
		fmt.Println("Input error:", err)
		os.Exit(1)
	}

	ctx := context.Background()
	result, err := exec.Execute(ctx, exec.Opts{
		Input:    string(input),
		Importer: &evaluator.SimpleImporter{},
	})
	if err != nil {
		fmt.Println("Execution error:", err)
		os.Exit(1)
	}
	if result != object.NULL {
		fmt.Println(result.Inspect())
	}
}
