package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/cloudcmds/tamarin/exec"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

const clearLine = "\033[2K\r"

func main() {

	ctx := context.Background()

	fmt.Println("Tamarin")
	fmt.Println("")

	// Global scope
	sc := scope.New(scope.Opts{Name: "global"})

	// Automatically import standard modules
	if err := exec.AutoImport(sc, nil, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	run := func(line string) {
		result, err := exec.Execute(ctx, exec.Opts{
			Input:             string(line),
			Scope:             sc,
			DisableAutoImport: true,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		switch result.(type) {
		case *object.Null:
		default:
			fmt.Println(result.Inspect())
		}
	}

	fmt.Print(">>> ")
	var historyIndex int
	var history []string
	var accumulate string

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.CtrlC:
			os.Exit(0)
			return true, nil
		case keys.Enter:
			fmt.Printf("\n")
			run(accumulate)
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
		default:
			if key.Runes != nil {
				fmt.Println(key)
			}
		}
		return false, nil // Return false to continue listening
	})

	for {
		time.Sleep(time.Second)
	}
}
