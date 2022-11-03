// Monkey is a scripting language implemented in golang, based upon
// the book "Write an Interpreter in Go", written by Thorsten Ball.
//
// This implementation adds a number of tweaks, improvements, and new
// features.  For example we support file-based I/O, regular expressions,
// the ternary operator, and more.
//
// For full details please consult the project homepage https://github.com/skx/monkey/
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/skx/monkey/evaluator"
	"github.com/skx/monkey/lexer"
	"github.com/skx/monkey/object"
	"github.com/skx/monkey/parser"
)

// This version-string will be updated via CI system for generated binaries.
var version = "master/unreleased"

//go:embed data/stdlib.mon
var stdlib string

//
// Implemention of "version()" function.
//
func versionFun(args ...object.Object) object.Object {
	return &object.String{Value: version}
}

//
// Implemention of "args()" function.
//
func argsFun(args ...object.Object) object.Object {
	l := len(os.Args[1:])
	result := make([]object.Object, l)
	for i, txt := range os.Args[1:] {
		result[i] = &object.String{Value: txt}
	}
	return &object.Array{Elements: result}
}

//
// Execute the supplied string as a program.
//
func Execute(input string) int {

	env := object.NewEnvironment()
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Printf("\t%s\n", msg)
		}
		os.Exit(1)
	}

	// Register a function called version()
	// that the script can call.
	evaluator.RegisterBuiltin("version",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (versionFun(args...))
		})

	// Access to the command-line arguments
	evaluator.RegisterBuiltin("args",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (argsFun(args...))
		})

	//
	//  Parse and evaluate our standard-library.
	//
	initL := lexer.New(stdlib)
	initP := parser.New(initL)
	initProg := initP.ParseProgram()
	evaluator.Eval(initProg, env)

	//
	//  Now evaluate the code the user wanted to load.
	//
	//  Note that here our environment will still contain
	// the code we just loaded from our data-resource
	//
	//  (i.e. Our monkey-based standard library.)
	//
	evaluator.Eval(program, env)
	return 0
}

func main() {

	//
	// Setup some flags.
	//
	eval := flag.String("eval", "", "Code to execute.")
	vers := flag.Bool("version", false, "Show our version and exit.")

	//
	// Parse the flags
	//
	flag.Parse()

	//
	// Showing the version?
	//
	if *vers {
		fmt.Printf("monkey %s\n", version)
		os.Exit(1)
	}

	//
	// Executing code?
	//
	if *eval != "" {
		Execute(*eval)
		os.Exit(1)
	}

	//
	// Otherwise we're either reading from STDIN, or the
	// named file containing source-code.
	//
	var input []byte
	var err error

	if len(flag.Args()) > 0 {
		input, err = ioutil.ReadFile(os.Args[1])
	} else {
		input, err = ioutil.ReadAll(os.Stdin)
	}

	if err != nil {
		fmt.Printf("Error reading: %s\n", err.Error())
	}

	Execute(string(input))
}
