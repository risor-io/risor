package evaluator

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseBreakpoints parses a comma-separated list of breakpoints, where each
// breakpoint is formatted as "file:line:flags". The flags are optional and
// can be any combination of "n" (no-stop), "t" (trace), and "d" (disabled).
func ParseBreakpoints(desc string) ([]Breakpoint, error) {
	var breakpoints []Breakpoint
	for _, s := range strings.Split(desc, ",") {
		parts := strings.Split(s, ":")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid breakpoint %q", s)
		}
		line, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid breakpoint %q: %v", s, err)
		}
		var flags string
		if len(parts) > 2 {
			flags = parts[2]
		}
		breakpoints = append(breakpoints, Breakpoint{
			File:     parts[0],
			Line:     line,
			Stop:     !strings.Contains(flags, "n"),
			Trace:    strings.Contains(flags, "t"),
			Disabled: strings.Contains(flags, "d"),
		})
	}
	return breakpoints, nil
}
