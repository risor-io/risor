package parser

import (
	"fmt"

	"github.com/cloudcmds/tamarin/internal/token"
)

// ErrorOpts is a struct that holds a variety of error data.
// All fields are optional, although one of `Cause` or `Message`
// are recommended. If `Cause` is set, `Message` will be ignored.
type ErrorOpts struct {
	ErrType       string
	Message       string
	Cause         error
	File          string
	StartPosition token.Position
	EndPosition   token.Position
	SourceCode    string
}

// NewBaseParserError returns a new BaseParserError populated with
// the given error data.
func NewParserError(opts ErrorOpts) *BaseParserError {
	return &BaseParserError{
		errType:       opts.ErrType,
		message:       opts.Message,
		cause:         opts.Cause,
		file:          opts.File,
		startPosition: opts.StartPosition,
		endPosition:   opts.EndPosition,
		sourceCode:    opts.SourceCode,
	}
}

// ParserError is an interface that all parser errors implement.
type ParserError interface {
	Type() string
	Message() string
	Cause() error
	File() string
	StartPosition() token.Position
	EndPosition() token.Position
	SourceCode() string
	Error() string
}

// BaseParserError is the simplest implementation of ParserError.
type BaseParserError struct {
	// Type of the error, e.g. "syntax error"
	errType string
	// The error message
	message string
	// The wrapped error
	cause error
	// File where the error occurred
	file string
	// Start position of the error in the input string
	startPosition token.Position
	// End position of the error in the input string
	endPosition token.Position
	// Relevant line of source code text
	sourceCode string
}

func (e *BaseParserError) Error() string {
	var msg string
	if e.cause != nil {
		msg = e.cause.Error()
	} else if e.message != "" {
		msg = e.message
	}
	if e.errType != "" {
		msg = fmt.Sprintf("%s: %s", e.errType, msg)
	}
	return msg
}

func (e *BaseParserError) Cause() error {
	return e.cause
}

func (e *BaseParserError) Message() string {
	return e.message
}

func (e *BaseParserError) Line() int {
	return e.startPosition.Line
}

func (e *BaseParserError) StartPosition() token.Position {
	return e.startPosition
}

func (e *BaseParserError) EndPosition() token.Position {
	return e.endPosition
}

func (e *BaseParserError) File() string {
	return e.file
}

func (e *BaseParserError) SourceCode() string {
	return e.sourceCode
}

func (e *BaseParserError) Unwrap() error {
	return e.cause
}

func (e *BaseParserError) Type() string {
	return e.errType
}

// NewSyntaxError returns a new SyntaxError populated with the given error data
func NewSyntaxError(opts ErrorOpts) *SyntaxError {
	opts.ErrType = "syntax error"
	return &SyntaxError{BaseParserError: NewParserError(opts)}
}

type SyntaxError struct {
	*BaseParserError
}
