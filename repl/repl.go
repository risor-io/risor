package repl

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/cloudcmds/tamarin/exec"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
	"github.com/fatih/color"
)

const (
	clearLine   = "\033[2K\r"
	moveBack    = "\033[%dD"
	moveForward = "\033[%dC"
)

func Run(ctx context.Context, sc *scope.Scope) error {

	color.New(color.Bold).Println("Tamarin")
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
		historyPath = path.Join(homeDir, ".tamarin_history")
		historyData, err := os.ReadFile(historyPath)
		if err == nil {
			history = strings.Split(string(historyData), "\n")
			historyIndex = len(history) - 1
		}
	}

	appendToHistory := func(line string) {
		if historyPath != "" {
			f, err := os.OpenFile(historyPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return
			}
			defer f.Close()
			f.WriteString(line + "\n")
		}
	}

	updateLine := func() {
		// var parts []string
		// l := lexer.New(accumulate)
		// for {
		// 	tok, err := l.NextToken()
		// 	if err != nil {
		// 		break
		// 	}
		// 	isKeyword := token.IsKeyword(tok.Literal)
		// 	if isKeyword {
		// 		parts = append(parts, color.MagentaString(tok.Literal))
		// 	} else {
		// 		parts = append(parts, tok.Literal)
		// 	}
		// }
		// s := strings.Join(parts, "")
		fmt.Printf(clearLine + ">>> " + accumulate)
	}

	// This could certainly use a refactor! But it works for now.
	return keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Enter:
			fmt.Printf("\n")
			execute(ctx, accumulate, sc)
			appendToHistory(accumulate)
			history = append(history, accumulate)
			historyIndex = len(history)
			accumulate = ""
			fmt.Print(">>> ")
			column = 0
		case keys.RuneKey, keys.Space, keys.Tab:
			if column < len(accumulate) {
				rest := accumulate[column:]
				restLen := len(rest)
				accumulate = accumulate[:column] + string(key.Runes) + rest
				updateLine()
				fmt.Printf(moveBack, restLen)
				// fmt.Printf(clearLine + ">>> " + accumulate + fmt.Sprintf(moveBack, restLen))
			} else {
				accumulate += string(key.Runes)
				updateLine()
				// fmt.Print(string(key.Runes))
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
					updateLine()
					fmt.Printf(moveBack, restLen)
					// fmt.Printf(clearLine + ">>> " + accumulate + fmt.Sprintf(moveBack, restLen))
				} else {
					accumulate = accumulate[:len(accumulate)-1]
					updateLine()
					fmt.Printf("\b \b")
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
						updateLine()
						fmt.Printf(moveBack, restLen)
						// fmt.Printf(clearLine + ">>> " + accumulate + fmt.Sprintf(moveBack, restLen))
					} else {
						accumulate = accumulate[:column]
						updateLine()
						// fmt.Printf(clearLine + ">>> " + accumulate)
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
				// fmt.Printf(clearLine + ">>> " + accumulate)
				updateLine()
			}
		case keys.Down:
			if historyIndex < len(history)-1 {
				historyIndex++
			}
			if historyIndex < len(history) {
				accumulate = history[historyIndex]
				column = len(accumulate)
				// fmt.Printf(clearLine + ">>> " + accumulate)
				updateLine()
			} else {
				// fmt.Printf(clearLine + ">>> ")
				column = 0
				accumulate = ""
				updateLine()
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
			updateLine()
			fmt.Print(strings.Repeat("\b", len(accumulate)))
			// fmt.Printf(clearLine + ">>> " + accumulate + strings.Repeat("\b", len(accumulate)))
			column = 0
		case keys.CtrlE:
			fmt.Printf(moveForward, len(accumulate)-column)
			column = len(accumulate)
		case keys.CtrlC, keys.CtrlD:
			fmt.Println()
			os.Exit(0)
		}
		return false, nil
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
