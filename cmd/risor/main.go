package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fatih/color"
	"github.com/risor-io/risor"
	"github.com/risor-io/risor/modules/all"
	"github.com/risor-io/risor/object"
	tos "github.com/risor-io/risor/os"
	"github.com/risor-io/risor/os/s3fs"
	"github.com/risor-io/risor/parser"
	"github.com/risor-io/risor/repl"
)

func main() {
	var noColor, showTiming, virtualOS bool
	var profilerOutputPath, code, breakpoints string
	flag.BoolVar(&noColor, "no-color", false, "Disable color output")
	flag.BoolVar(&showTiming, "timing", false, "Show timing information")
	flag.StringVar(&code, "c", "", "Code to execute")
	flag.StringVar(&profilerOutputPath, "profile", "", "Enable profiling")
	flag.StringVar(&breakpoints, "breakpoints", "", "Comma-separated list of breakpoints")
	flag.BoolVar(&virtualOS, "virtual-os", false, "Enable virtual OS")
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

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithSharedConfigProfile(os.Getenv("RISOR_TEST_AWS_PROFILE")),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", red(err.Error()))
		os.Exit(1)
	}

	s3Client := s3.NewFromConfig(cfg)

	fs, err := s3fs.New(ctx,
		s3fs.WithBase("/stat"),
		s3fs.WithClient(s3Client),
		s3fs.WithBucket(os.Getenv("RISOR_TEST_S3_BUCKET")))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", red(err.Error()))
		os.Exit(1)
	}

	if virtualOS {
		ctx = tos.WithOS(ctx, tos.NewVirtualOS(ctx,
			tos.WithMounts(map[string]*tos.Mount{
				"/": {
					Source: fs,
					Target: "/",
				},
			})))
	}

	// Input can only come from one source
	nArgs := len(flag.Args())
	if nArgs > 0 && len(code) > 0 {
		fmt.Fprintf(os.Stderr, "%s\n", red("error: cannot provide both a script file and -c input\n"))
		os.Exit(1)
	} else if nArgs == 0 && len(code) == 0 {
		// Run REPL
		if err := repl.Run(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", red(err.Error()))
			os.Exit(1)
		}
		return
	}

	// Otherwise, use input from either -c or the first argument
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

	start := time.Now()

	result, err := risor.Eval(ctx,
		string(input),
		risor.WithBuiltins(all.Builtins()),
		risor.WithDefaultBuiltins(),
		risor.WithDefaultModules())
	if err != nil {
		parserErr, ok := err.(parser.ParserError)
		if ok {
			fmt.Fprintf(os.Stderr, "%s\n", red(parserErr.FriendlyErrorMessage()))
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", red(err.Error()))
		}
		os.Exit(1)
	}
	if result != object.Nil {
		fmt.Println(result.Inspect())
	}
	if showTiming {
		fmt.Printf("%.03f\n", time.Since(start).Seconds())
	}
}
