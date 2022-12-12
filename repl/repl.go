package repl

import (
	"context"
	"fmt"
	"os"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/cloudcmds/tamarin/exec"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
	"github.com/fatih/color"
)

const clearLine = "\033[2K\r"

func Run(ctx context.Context, sc *scope.Scope) error {

	color.New(color.Bold).Println("Tamarin")
	fmt.Println("")
	fmt.Printf(">>> ")

	var historyIndex int
	var history []string
	var accumulate string

	return keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Enter:
			fmt.Printf("\n")
			execute(ctx, accumulate, sc)
			history = append(history, accumulate)
			historyIndex = len(history)
			accumulate = ""
			fmt.Print(">>> ")
		case keys.RuneKey, keys.Space, keys.Tab:
			accumulate += string(key.Runes)
			fmt.Print(string(key.Runes))
		case keys.Backspace:
			if len(accumulate) > 0 {
				accumulate = accumulate[:len(accumulate)-1]
				fmt.Printf("\b \b")
			}
		case keys.Up:
			if historyIndex > 0 {
				historyIndex--
			}
			if historyIndex < len(history) {
				accumulate = history[historyIndex]
				fmt.Printf(clearLine + ">>> " + accumulate)
			}
		case keys.Down:
			if historyIndex < len(history) {
				historyIndex++
			}
			if historyIndex < len(history) {
				accumulate = history[historyIndex]
				fmt.Printf(clearLine + ">>> " + accumulate)
			} else {
				fmt.Printf(clearLine + ">>> ")
				accumulate = ""
			}
		case keys.Left, keys.Right:
			// Ignore
		case keys.CtrlA:
			fmt.Printf(clearLine + ">>> " + accumulate + strings.Repeat("\b", len(accumulate)))
		case keys.CtrlC, keys.CtrlD:
			fmt.Println()
			os.Exit(0)
		}
		return false, nil // Return false to continue listening
	})
}

func execute(ctx context.Context, code string, sc *scope.Scope) (object.Object, error) {
	result, err := exec.Execute(ctx, exec.Opts{
		Input:             code,
		Scope:             sc,
		DisableAutoImport: true,
	})
	if err != nil {
		color.Red(err.Error())
		return nil, err
	}
	switch result.(type) {
	case *object.NilType:
	default:
		fmt.Println(result.Inspect())
	}
	return result, nil
}
