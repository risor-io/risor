package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
	"github.com/mitchellh/go-homedir"
	"github.com/risor-io/risor"
	"github.com/risor-io/risor/cmd/risor/repl"
	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/modules/aws"
	"github.com/risor-io/risor/modules/image"
	"github.com/risor-io/risor/modules/pgx"
	"github.com/risor-io/risor/modules/uuid"
	"github.com/risor-io/risor/object"
	ros "github.com/risor-io/risor/os"
	"github.com/risor-io/risor/os/s3fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	red     = color.New(color.FgRed).SprintfFunc()
)

func init() {

	cobra.OnInitialize(initConfig)
	viper.SetEnvPrefix("risor")

	// Global flags

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.risor.yaml)")
	rootCmd.PersistentFlags().StringP("code", "c", "", "Code to evaluate")
	rootCmd.PersistentFlags().Bool("stdin", false, "Read code from stdin")
	rootCmd.PersistentFlags().String("cpu-profile", "", "Capture a CPU profile")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().Bool("virtual-os", false, "Enable a virtual operating system")
	rootCmd.PersistentFlags().StringArrayP("mount", "m", []string{}, "Mount a filesystem")
	rootCmd.PersistentFlags().Bool("no-default-globals", false, "Disable the default globals")
	rootCmd.PersistentFlags().String("modules", ".", "Path to library modules")
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Help for Risor")

	viper.BindPFlag("code", rootCmd.PersistentFlags().Lookup("code"))
	viper.BindPFlag("stdin", rootCmd.PersistentFlags().Lookup("stdin"))
	viper.BindPFlag("cpu-profile", rootCmd.PersistentFlags().Lookup("cpu-profile"))
	viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
	viper.BindPFlag("virtual-os", rootCmd.PersistentFlags().Lookup("virtual-os"))
	viper.BindPFlag("mount", rootCmd.PersistentFlags().Lookup("mount"))
	viper.BindPFlag("no-default-globals", rootCmd.PersistentFlags().Lookup("no-default-globals"))
	viper.BindPFlag("modules", rootCmd.PersistentFlags().Lookup("modules"))
	viper.BindPFlag("help", rootCmd.PersistentFlags().Lookup("help"))

	// Root command flags

	rootCmd.Flags().Bool("timing", false, "Show timing information")
	rootCmd.Flags().StringP("output", "o", "", "Set the output format")
	rootCmd.Flags().SetInterspersed(false)
	viper.BindPFlag("timing", rootCmd.Flags().Lookup("timing"))
	viper.BindPFlag("output", rootCmd.Flags().Lookup("output"))

	viper.AutomaticEnv()
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Search config in home directory with name ".risor"
		viper.AddConfigPath(home)
		viper.SetConfigName(".risor")
	}
	viper.ReadInConfig()
}

func fatal(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// this works around issues with cobra not passing args after -- to the script
// see: https://github.com/spf13/cobra/issues/1877
// passedargs is everything after the '--'
// If dropped is true, cobra dropped the double dash in its 'args' slice.
// argsout returns the cobra args but without the '--' and the items behind it
// this also supports multiple '--' in the args	list
func getpassthruargs(args []string) (argsout []string, passedargs []string, dropped bool) {
	//		lenArgs := len(args)
	argsout = args
	ddashcnt := 0
	for n, arg := range os.Args {
		if arg == "--" {
			ddashcnt++
			if len(passedargs) == 0 {
				if len(os.Args) > n {
					passedargs = os.Args[n+1:]
				}
			}
		}
	}
	// drop arg0 from count
	noddash := true
	for n2, argz := range args {
		// don't go past the first '--' - this allows one to pass '--' to risor also
		if argz == "--" {
			noddash = false
			argsout = args[:n2]
			break
		}
	}
	// cobra seems to drop the '--' in its args if everything before it was a flag
	// if thats the case then args is just empty b/c evrything else is the '--' and the passed args
	if (noddash || ddashcnt > 1) && len(passedargs) > 0 {
		dropped = true
		argsout = []string{}
	}
	return
}

var rootCmd = &cobra.Command{
	Use:   "risor",
	Short: "Fast and flexible scripting for Go developers and DevOps",
	Long:  `https://risor.io`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()
		var passedargs []string

		args, passedargs, _ = getpassthruargs(args)
		// pass the 'passthru' args to risor's os package
		ros.SetScriptArgs(passedargs)
		// Optionally enable a virtual operating system and add it to
		// the context so that it's made available to Risor VM.
		if viper.GetBool("virtual-os") {
			mounts := map[string]*ros.Mount{}
			m := viper.GetStringSlice("mount")
			for _, v := range m {
				fs, dst, err := mountFromSpec(ctx, v)
				if err != nil {
					fatal(err.Error())
				}
				mounts[dst] = &ros.Mount{
					Source: fs,
					Target: dst,
				}
			}
			vos := ros.NewVirtualOS(ctx, ros.WithMounts(mounts), ros.WithArgs(passedargs))
			ctx = ros.WithOS(ctx, vos)
		}

		// Disable colored output if no-color is specified
		if viper.GetBool("no-color") {
			color.NoColor = true
		}

		// Optionally capture a CPU profile to the given path
		if path := viper.GetString("cpu-profile"); path != "" {
			f, err := os.Create(path)
			if err != nil {
				fatal(red(err.Error()))
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}

		// Build up a list of options to pass to the VM
		var opts []risor.Option
		if viper.GetBool("no-default-globals") {
			opts = append(opts, risor.WithoutDefaultGlobals())
		} else {
			opts = append(opts, risor.WithGlobals(map[string]any{
				"image": image.Module(),
				"pgx":   pgx.Module(),
				"uuid":  uuid.Module(),
			}))
			// AWS support may or may not be compiled in based on build tags
			if aws := aws.Module(); aws != nil {
				opts = append(opts, risor.WithGlobal("aws", aws))
			}
		}
		if modulesDir := viper.GetString("modules"); modulesDir != "" {
			opts = append(opts, risor.WithLocalImporter(modulesDir))
		}

		// Determine what code is to be executed. The code may be supplied
		// via the --code option, a path supplied as an arg, or stdin.
		codeWasSupplied := cmd.Flags().Lookup("code").Changed
		code := viper.GetString("code")
		if len(args) > 0 && codeWasSupplied {
			fatal(red("cannot specify both code and a filepath"))
		}
		if len(args) == 0 && !codeWasSupplied && !viper.GetBool("stdin") {
			if err := repl.Run(ctx, opts); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", red(err.Error()))
				os.Exit(1)
			}
			return
		}
		if viper.GetBool("stdin") {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				fatal(red(err.Error()))
			}
			if len(data) == 0 {
				fatal(red("no code supplied"))
			}
			code = string(data)
		} else if len(args) > 0 {
			bytes, err := os.ReadFile(args[0])
			if err != nil {
				fatal(red(err.Error()))
			}
			code = string(bytes)
		}

		start := time.Now()

		// Execute the code
		result, err := risor.Eval(ctx, code, opts...)
		if err != nil {
			if friendlyErr, ok := err.(errz.FriendlyError); ok {
				fmt.Fprintf(os.Stderr, "%s\n", red(friendlyErr.FriendlyErrorMessage()))
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", red(err.Error()))
			}
			os.Exit(1)
		}

		dt := time.Since(start)

		// Print the result
		output, err := getOutput(result, viper.GetString("output"))
		if err != nil {
			fatal(red(err.Error()))
		} else if output != "" {
			fmt.Println(output)
		}

		// Optionally print the execution time
		if viper.GetBool("timing") {
			fmt.Printf("%v\n", dt)
		}
	},
}

func getOutput(result object.Object, format string) (string, error) {
	switch strings.ToLower(format) {
	case "":
		// With an unspecified format, we'll try to do the most helpful thing:
		//  1. If the result is nil, we want to print nothing
		//  2. If the result marshals to JSON, we'll print that
		//  3. Otherwise, we'll print the result's string representation
		if result == object.Nil {
			return "", nil
		}
		output, err := getOutputJSON(result)
		if err != nil {
			return fmt.Sprintf("%v", result), nil
		}
		return string(output), nil
	case "json":
		output, err := getOutputJSON(result)
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "text":
		return fmt.Sprintf("%v", result), nil
	default:
		return "", fmt.Errorf("unknown output format: %s", format)
	}
}

func getOutputJSON(result object.Object) ([]byte, error) {
	if viper.GetBool("no-color") {
		return json.MarshalIndent(result, "", "  ")
	} else {
		return prettyjson.Marshal(result)
	}
}

func mountFromSpec(ctx context.Context, spec string) (ros.FS, string, error) {
	parts := strings.Split(spec, ",")
	items := map[string]string{}
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return nil, "", fmt.Errorf("invalid mount spec: %s (expected k=v format)", spec)
		}
		items[kv[0]] = kv[1]
	}
	typ, ok := items["type"]
	if !ok || typ == "" {
		return nil, "", fmt.Errorf("invalid mount spec: %q (missing type)", spec)
	}
	src, ok := items["src"]
	if !ok || src == "" {
		return nil, "", fmt.Errorf("invalid mount spec: %q (missing src)", spec)
	}
	dst, ok := items["dst"]
	if !ok || dst == "" {
		return nil, "", fmt.Errorf("invalid mount spec: %q (missing dst)", spec)
	}
	switch typ {
	case "s3":
		var awsOpts []func(*config.LoadOptions) error
		if r, ok := items["region"]; ok {
			awsOpts = append(awsOpts, config.WithRegion(r))
		}
		if p, ok := items["profile"]; ok {
			awsOpts = append(awsOpts, config.WithSharedConfigProfile(p))
		}
		cfg, err := config.LoadDefaultConfig(ctx, awsOpts...)
		if err != nil {
			return nil, "", err
		}
		s3Opts := []s3fs.Option{
			s3fs.WithBucket(src),
			s3fs.WithClient(s3.NewFromConfig(cfg)),
		}
		if p, ok := items["prefix"]; ok && p != "" {
			s3Opts = append(s3Opts, s3fs.WithBase(p))
		}
		fs, err := s3fs.New(ctx, s3Opts...)
		if err != nil {
			return nil, "", err
		}
		return fs, dst, nil
	default:
		return nil, "", fmt.Errorf("unsupported source: %s", src)
	}
}
