// Package repl implements a read-eval-print-loop for Risor.
package repl

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/fatih/color"
	"github.com/risor-io/risor"
	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/parser"
	"github.com/risor-io/risor/vm"
)

const (
	clearLine   = "\033[2K\r"
	moveBack    = "\033[%dD"
	moveForward = "\033[%dC"
)

func Run(ctx context.Context, options []risor.Option) error {
	color.New(color.Bold).Println("Risor")
	fmt.Println("")
	fmt.Printf(">>> ")

	var column int
	var historyIndex int
	var history []string
	var accumulate string

	// Read execution history just like Python's REPL.
	var historyPath string
	homeDir, err := os.UserHomeDir()
	if err == nil {
		historyPath = path.Join(homeDir, ".risor_history")
		historyData, err := os.ReadFile(historyPath)
		if err == nil {
			history = strings.Split(string(historyData), "\n")
			historyIndex = len(history) - 1
		}
	}

	appendToHistory := func(line string) {
		if historyPath != "" {
			f, err := os.OpenFile(historyPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				return
			}
			defer f.Close()
			f.WriteString(line + "\n")
		}
	}

	getLineText := func() string {
		return clearLine + ">>> " + accumulate
	}

	r := risor.NewConfig()
	for _, opt := range options {
		opt(r)
	}

	evalFunc := getEvaluator(r)

	// This could certainly use a refactor! But it works for now.
	return keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Enter:
			fmt.Printf("\n")
			evalFunc(ctx, accumulate)
			appendToHistory(accumulate)
			history = append(history, accumulate)
			historyIndex = len(history)
			accumulate = ""
			fmt.Print(getLineText())
			column = 0
		case keys.RuneKey, keys.Space, keys.Tab:
			if column < len(accumulate) {
				rest := accumulate[column:]
				restLen := len(rest)
				accumulate = accumulate[:column] + string(key.Runes) + rest
				fmt.Print(getLineText() + fmt.Sprintf(moveBack, restLen))
			} else {
				accumulate += string(key.Runes)
				fmt.Print(getLineText())
			}
			column += len(key.Runes)
		case keys.Backspace:
			if len(accumulate) > 0 {
				if column < len(accumulate) {
					rest := accumulate[column:]
					restLen := len(rest)
					if column > 0 {
						accumulate = accumulate[:column-1] + rest
					}
					fmt.Print(getLineText() + fmt.Sprintf(moveBack, restLen))
				} else {
					accumulate = accumulate[:len(accumulate)-1]
					fmt.Print(getLineText())
				}
				if column > 0 {
					column--
				}
			}
		case keys.Delete:
			if len(accumulate) > 0 {
				if column < len(accumulate) {
					rest := accumulate[column+1:]
					restLen := len(rest)
					if restLen > 0 {
						accumulate = accumulate[:column] + rest
						fmt.Print(getLineText() + fmt.Sprintf(moveBack, restLen))
					} else {
						accumulate = accumulate[:column]
						fmt.Print(getLineText())
					}
				}
			}
		case keys.Up:
			if historyIndex > 0 {
				historyIndex--
			}
			if historyIndex < len(history) {
				accumulate = history[historyIndex]
				column = len(accumulate)
				fmt.Print(getLineText())
			}
		case keys.Down:
			if historyIndex < len(history)-1 {
				historyIndex++
			}
			if historyIndex < len(history) {
				accumulate = history[historyIndex]
				column = len(accumulate)
				fmt.Print(getLineText())
			} else {
				column = 0
				accumulate = ""
				fmt.Print(getLineText())
			}
		case keys.Left:
			if column > 0 {
				fmt.Printf(moveBack, 1)
				column--
			}
		case keys.Right:
			if column < len(accumulate) {
				fmt.Printf(moveForward, 1)
				column++
			}
		case keys.CtrlA:
			fmt.Print(getLineText() + strings.Repeat("\b", len(accumulate)))
			column = 0
		case keys.CtrlE:
			fmt.Printf(moveForward, len(accumulate)-column)
			column = len(accumulate)
		case keys.CtrlC, keys.CtrlD:
			fmt.Println()
			return true, nil
		}
		return false, nil
	})
}

func getEvaluator(cfg *risor.Config) func(ctx context.Context, source string) (object.Object, error) {
	var c *compiler.Compiler
	var v *vm.VirtualMachine

	return func(ctx context.Context, source string) (object.Object, error) {
		if c == nil {
			var err error
			c, err = compiler.New(cfg.CompilerOpts()...)
			if err != nil {
				return nil, err
			}
		}

		ast, err := parser.Parse(ctx, source)
		if err != nil {
			color.Red(err.Error())
			return nil, err
		}

		code, err := c.Compile(ast)
		if err != nil {
			color.Red(err.Error())
			return nil, err
		}

		if v == nil {
			v = vm.New(code, cfg.VMOpts()...)
		}
		if err := v.Run(ctx); err != nil {
			// Update the IP to be after the last instruction, so that next
			// time around we start in the right location.
			v.SetIP(code.InstructionCount())
			color.Red(err.Error())
			return nil, err
		}

		result, ok := v.TOS()
		if !ok || result == nil {
			return object.Nil, nil
		}

		switch result := result.(type) {
		case *object.Error:
			color.Red(result.Value().Error())
		case *object.NilType:
		default:
			fmt.Println(result.Inspect())
		}
		return result, nil
	}
}
