package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"

	"github.com/myzie/tamarin/internal/evaluator"
	"github.com/myzie/tamarin/internal/exec"
	"github.com/myzie/tamarin/internal/object"
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
	evaluator.RegisterRestrictedBuiltins()

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
	result, err := exec.Execute(ctx, string(input), &evaluator.SimpleImporter{})
	if err != nil {
		fmt.Println("Execution error:", err)
		os.Exit(1)
	}
	if result != object.NULL {
		fmt.Println(result.Inspect())
	}
}
