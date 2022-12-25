package stack

import (
	"fmt"
	"strings"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/scope"
)

type Frame struct {
	name      string
	statement ast.Statement
	scope     *scope.Scope
}

type FrameOpts struct {
	Name      string
	Statement ast.Statement
	Scope     *scope.Scope
}

func NewFrame(opts FrameOpts) *Frame {
	return &Frame{
		name:      opts.Name,
		statement: opts.Statement,
		scope:     opts.Scope,
	}
}

func (f *Frame) Statement() ast.Statement {
	return f.statement
}

func (f *Frame) Scope() *scope.Scope {
	return f.scope
}

func (f *Frame) Name() string {
	return f.name
}

// Stack represents the call stack of a Tamarin program. Push and Pop are called
// to add and remove frames from the stack, respectively.
type Stack struct {
	frames []*Frame
}

// Push adds a new frame to the stack.
func (s *Stack) Push(f *Frame) {
	s.frames = append(s.frames, f)
}

// Pop removes the top frame from the stack and returns it.
func (s *Stack) Pop() *Frame {
	size := len(s.frames)
	if size == 0 {
		return nil
	}
	top := s.frames[size-1]
	s.frames = s.frames[:size-1]
	return top
}

// Top returns the top frame on the stack without removing it.
func (s *Stack) Top() *Frame {
	size := len(s.frames)
	if size == 0 {
		return nil
	}
	return s.frames[size-1]
}

// Size returns the number of frames on the stack.
func (s *Stack) Size() int {
	return len(s.frames)
}

// TrackStatement is used to mark the statement being executed in the current frame
func (s *Stack) TrackStatement(statement ast.Statement, sc *scope.Scope) *Frame {
	if s.Size() == 0 {
		f := NewFrame(FrameOpts{
			Name:      "main",
			Statement: statement,
			Scope:     sc,
		})
		s.Push(f)
		return f
	}
	f := s.Top()
	f.statement = statement
	return f
}

func (s *Stack) String() string {
	var frames []string
	for i, frame := range s.frames {
		var s string
		if frame.statement != nil {
			tok := frame.statement.StartToken()
			loc := fmt.Sprintf("%s:%d", tok.StartPosition.File, tok.StartPosition.LineNumber())
			s = fmt.Sprintf("%s - in %s | %s",
				loc, frame.name, frame.statement.String())
		} else {
			s = fmt.Sprintf("in %s", frame.name)
		}
		var pad string
		if i > 0 {
			pad = strings.Repeat("  ", i)
		}
		frames = append(frames, fmt.Sprintf("%s%s", pad, s))
	}
	return strings.Join(frames, "\n")
}

// New returns a new Stack.
func New() *Stack {
	return &Stack{}
}
