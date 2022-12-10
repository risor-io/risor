package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cloudcmds/tamarin/exec"
	"github.com/cloudcmds/tamarin/scope"
)

func main() {

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(initialModel(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type errMsg error

type model struct {
	input  textarea.Model
	output textarea.Model
	err    error
	scope  *scope.Scope
}

func initialModel() model {

	ti := textarea.New()
	ti.Placeholder = "x := 42"
	ti.SetHeight(10)
	ti.SetWidth(40)
	ti.Focus()

	to := textarea.New()
	to.Placeholder = "output"
	to.SetHeight(10)
	to.ShowLineNumbers = false

	return model{input: ti, output: to}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *model) Execute() {
	s := scope.New(scope.Opts{Name: "global"})
	if err := exec.AutoImport(s, nil, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	result, err := exec.Execute(context.Background(), exec.Opts{
		Input:             string(m.input.Value()),
		Scope:             s,
		DisableAutoImport: true,
	})
	if err != nil {
		m.output.SetValue(err.Error())
		return
	}
	m.output.SetValue(result.Inspect())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var result, inCmd, outCmd tea.Cmd
	m.input, inCmd = m.input.Update(msg)
	m.output, outCmd = m.output.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.input.SetWidth(msg.Width / 2)
		m.input.SetHeight(msg.Height - 10)
		m.output.SetWidth(msg.Width / 2)
		m.output.SetHeight(msg.Height - 10)
	case tea.KeyMsg:
		switch msg.String() {
		case "alt+enter", "ctrl+r":
			m.Execute()
		}
		switch msg.Type {
		case tea.KeyEsc:
			if m.input.Focused() {
				m.input.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			if !m.input.Focused() {
				m.Execute()
			}
		default:
			if !m.input.Focused() {
				result = m.input.Focus()
			}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}
	return m, tea.Batch(inCmd, outCmd, result)
}

func (m model) View() string {

	main := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.input.View(),
		m.output.View())

	return fmt.Sprintf(
		"Tamarin\n\n%s\n\n%s",
		main,
		"(ctrl+c to quit)",
	) + "\n\n"
}
