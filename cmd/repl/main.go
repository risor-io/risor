package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudcmds/tamarin/exec"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/mattn/go-runewidth"
	term "github.com/nsf/termbox-go"
)

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func displayLines(screen tcell.Screen, lines []string) {
	// w, h := screen.Size()
	screen.Clear()
	style := tcell.StyleDefault.Foreground(tcell.ColorCadetBlue).Background(tcell.ColorWhite)
	for y, line := range lines {
		var x int
		for _, c := range line {
			var comb []rune
			w := runewidth.RuneWidth(c)
			if w == 0 {
				comb = []rune{c}
				c = ' '
				w = 1
			}
			screen.SetContent(x, y, c, comb, style)
			x += w
		}
	}
	screen.Show()
}

func main() {

	ctx := context.Background()

	encoding.Register()

	if err := term.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer term.Close()

	// Global scope
	sc := scope.New(scope.Opts{Name: "global"})

	// Automatically import standard modules
	if err := exec.AutoImport(sc, nil, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf(">>> ")

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
		if _, ok := result.(*object.Null); !ok {
			fmt.Println(result.Inspect())
		}
	}

	var currentLine string
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc:
				term.Sync()
				fmt.Println("ESC pressed")
			case term.KeyF1:
				term.Sync()
				fmt.Println("F1 pressed")
			case term.KeyInsert:
				term.Sync()
				fmt.Println("Insert pressed")
			case term.KeyDelete:
				term.Sync()
				fmt.Println("Delete pressed")
			case term.KeyHome:
				term.Sync()
				fmt.Println("Home pressed")
			case term.KeyEnd:
				term.Sync()
				fmt.Println("End pressed")
			case term.KeyPgup:
				term.Sync()
			case term.KeyArrowRight:
				term.Sync()
				fmt.Println("Arrow Right pressed")
			case term.KeySpace:
				term.Sync()
				fmt.Println("Space pressed")
			case term.KeyBackspace:
				term.Sync()
				fmt.Println("Backspace pressed")
			case term.KeyEnter:
				term.Sync()
				fmt.Println("Enter pressed")
				run(currentLine)
			case term.KeyTab:
				term.Sync()
				fmt.Println("Tab pressed")
			default:
				term.Sync()
				fmt.Println("ASCII : ", ev.Ch)
				currentLine += string(ev.Ch)
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}
