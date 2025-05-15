package color

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

// Color formatting codes
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Italic = "\033[3m"

	// Foreground colors
	FgRed     = "\033[31m"
	FgGreen   = "\033[32m"
	FgYellow  = "\033[33m"
	FgBlue    = "\033[34m"
	FgMagenta = "\033[35m"
	FgCyan    = "\033[36m"
	FgWhite   = "\033[37m"
	FgHiCyan  = "\033[96m"

	// Background colors
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

// Color control variables
var (
	// NoColor controls whether color output is enabled
	NoColor bool

	// ForceColor forces color output regardless of terminal support
	ForceColor bool
)

func init() {
	// Initialize color settings based on environment and terminal
	NoColor = shouldDisableColor()
	ForceColor = os.Getenv("FORCE_COLOR") != ""
}

// shouldDisableColor determines if colors should be disabled based on
// environment variables and terminal support
func shouldDisableColor() bool {
	// Check environment variables first
	if os.Getenv("NO_COLOR") != "" {
		return true
	}

	// Check if TERM indicates a dumb terminal
	if os.Getenv("TERM") == "dumb" && os.Getenv("FORCE_COLOR") == "" {
		return true
	}

	// Check if we're running in Go test mode
	if isRunningTest() {
		return true
	}

	// Check if output is going to a terminal
	return !isTerminal(os.Stdout)
}

// isRunningTest returns true if the code is running as part of Go tests
func isRunningTest() bool {
	// Go test sets the current executable name to "*.test"
	if strings.HasSuffix(os.Args[0], ".test") {
		return true
	}

	// Also look for the "test.v" flag commonly used in Go tests
	for _, arg := range os.Args {
		if arg == "-test.v" || strings.HasPrefix(arg, "-test.") {
			return true
		}
	}

	return false
}

// isTerminal returns true if the given writer is a terminal
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}

	// Basic heuristic: stdout and stderr are usually terminals
	fd := f.Fd()
	return fd == 1 || fd == 2
}

// EnableColors enables colorized output
func EnableColors() {
	NoColor = false
}

// DisableColors disables colorized output
func DisableColors() {
	NoColor = true
}

// Color represents a colored text formatter
type Color struct {
	codes []string
}

// New creates a new Color with the given color codes
func New(codes ...string) *Color {
	return &Color{codes: codes}
}

// Sprint formats using the default formats and returns the resulting string with color
func (c *Color) Sprint(a ...interface{}) string {
	return c.wrap(fmt.Sprint(a...))
}

// Sprintf formats according to a format specifier and returns the resulting string with color
func (c *Color) Sprintf(format string, a ...interface{}) string {
	return c.wrap(fmt.Sprintf(format, a...))
}

// wrap adds color codes around the given string if colors are enabled
func (c *Color) wrap(s string) string {
	// Skip coloring if:
	// - No color codes provided
	// - Colors are disabled and not forced
	// - Running on Windows without adequate support
	if len(c.codes) == 0 ||
		(NoColor && !ForceColor) ||
		(runtime.GOOS == "windows" && !supportsWindowsAnsi() && !ForceColor) {
		return s
	}

	prefix := strings.Join(c.codes, "")
	return prefix + s + Reset
}

// supportsWindowsAnsi determines if the Windows console supports ANSI colors
func supportsWindowsAnsi() bool {
	// Windows 10+ with Windows Terminal or modern console supports ANSI
	if os.Getenv("WT_SESSION") != "" {
		return true
	}

	// Check for other terminal emulators that support ANSI
	termProg := os.Getenv("TERM_PROGRAM")
	return termProg != "" // Programs like VSCode, ConEmu, cmder set this
}
