// Package lexer contains the code to lex input-programs into a stream
// of tokens, such that they may be parsed.
package lexer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/cloudcmds/tamarin/token"
)

// Opts contains Lexer initialization options
type Opts struct {
	Input string
	File  string
}

// Lexer holds our object-state.
type Lexer struct {
	// The index of the current character
	position int

	// The index of the next character
	nextPosition int

	// The current character
	ch rune

	// A rune slice of our input string
	characters []rune

	// Previous token
	prevToken token.Token

	// The index of the current line, for the current character
	line int

	// The character index where the current line began
	lineStart int

	// The column, within the current line, of the current character
	column int

	// The position of the start of the current token
	tokenStartPosition token.Position

	// Name of the file be read
	file string
}

// New returns a Lexer instance for a given string input.
func New(input string) *Lexer {
	l := &Lexer{
		characters:   []rune(input),
		column:       -1, // -1 = before the first column
		position:     -1, // -1 = before the first character
		nextPosition: 0,  //  0 = read the first character next
	}
	l.readChar()
	return l
}

// NewWithOptions returns a Lexer instance with the given options.
func NewWithOptions(opts Opts) *Lexer {
	l := New(opts.Input)
	l.file = opts.File
	return l
}

// The name of the file being read.
func (l *Lexer) File() string {
	return l.file
}

// Advance one charaacter in the input string. Calling this when
// we are already at the end of the input has no effect.
func (l *Lexer) readChar() {

	// Return if we are already at the end of the input. Note that
	// when position == len(l.characters) the current character is
	// considered to be EOF, so that position is considered valid.
	if l.position > len(l.characters) {
		return
	}

	// The rune we were at before we advanced. This is used to
	// understand whether we advanced to a new line.
	prevCh := l.ch

	// Move forward. In the very first call, this moves us from
	// position -1 to position 0.
	l.position = l.nextPosition
	l.nextPosition++

	// Set the current character based on the current position
	if l.position < len(l.characters) {
		l.ch = l.characters[l.position]
	} else {
		l.ch = rune(0) // EOF
	}

	// Track three values, all zero-indexed:
	//  * line: the current line number
	//  * column: the current column number, within the line
	//  * lineStart: where the current line started in the input
	if prevCh == rune('\n') {
		l.column = 0
		l.line++
		l.lineStart = l.position
	} else {
		l.column++
	}
}

// CurrentPosition returns a Position object for the current read position.
func (l *Lexer) CurrentPosition() token.Position {
	return token.Position{
		Value:     l.ch,
		Char:      l.position,
		LineStart: l.lineStart,
		Line:      l.line,
		Column:    l.column,
	}
}

func (l *Lexer) GetTokenLineText(t token.Token) string {

	if len(l.characters) == 0 {
		return ""
	}

	// This is the position where the token begins
	tokenStart := t.StartPosition

	// Guard against programming errors. Raise panics to catch mistakes early.
	// Note that for the EOF token, the char offset will be equal to the number
	// of chars in the input, hence the ">" and not ">=".
	if tokenStart.Char < 0 || tokenStart.Char > len(l.characters) {
		panic(fmt.Errorf("invalid token start position: %d (input length: %d)",
			tokenStart.Char, len(l.characters)))
	}
	if tokenStart.Line < 0 {
		panic(fmt.Errorf("invalid token start line: %d", tokenStart.Line))
	}

	// Find the start of the line containing the given token
	start := tokenStart.Char
	if t.Type == token.EOF {
		start--
	}
	for start > 0 && l.characters[start-1] != rune('\n') {
		start--
	}
	// Find the end of that line
	end := tokenStart.Char
	for end < len(l.characters) && l.characters[end] != rune('\n') {
		end++
	}
	// Return the line, excluding the newline character
	return string(l.characters[start:end])
}

func (l *Lexer) newToken(typ token.Type, literal string) token.Token {
	t := token.Token{
		Type:          typ,
		Literal:       literal,
		StartPosition: l.tokenStartPosition,
		EndPosition:   l.CurrentPosition(),
	}
	return t
}

// NextToken to read next token, skipping the white space.
func (l *Lexer) NextToken() (token.Token, error) {

	var tok token.Token
	l.skipWhitespace()
	l.tokenStartPosition = l.CurrentPosition()

	// skip single-line comments
	if l.ch == rune('#') ||
		(l.ch == rune('/') && l.peekChar() == rune('/')) {
		l.skipComment()
		return l.NextToken()
	}
	// multi-line comments
	if l.ch == rune('/') && l.peekChar() == rune('*') {
		l.skipMultiLineComment()
	}

	switch l.ch {
	case rune('&'):
		if l.peekChar() == rune('&') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.AND, string(ch)+string(l.ch))
		}
	case rune('|'):
		if l.peekChar() == rune('|') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.OR, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.PIPE, string(l.ch))
		}
	case rune('='):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.EQ, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.ASSIGN, string(l.ch))
		}
	case rune(';'):
		tok = l.newToken(token.SEMICOLON, string(l.ch))
	case rune('?'):
		tok = l.newToken(token.QUESTION, string(l.ch))
	case rune('('):
		tok = l.newToken(token.LPAREN, string(l.ch))
	case rune(')'):
		tok = l.newToken(token.RPAREN, string(l.ch))
	case rune(','):
		tok = l.newToken(token.COMMA, string(l.ch))
	case rune('.'):
		tok = l.newToken(token.PERIOD, string(l.ch))
	case rune('+'):
		if l.peekChar() == rune('+') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.PLUS_PLUS, string(ch)+string(l.ch))
		} else if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.PLUS_EQUALS, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.PLUS, string(l.ch))
		}
	case rune('%'):
		tok = l.newToken(token.MOD, string(l.ch))
	case rune('{'):
		tok = l.newToken(token.LBRACE, string(l.ch))
	case rune('}'):
		tok = l.newToken(token.RBRACE, string(l.ch))
	case rune('-'):
		if l.peekChar() == rune('-') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.MINUS_MINUS, string(ch)+string(l.ch))
		} else if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.MINUS_EQUALS, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.MINUS, string(l.ch))
		}
	case rune('/'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.SLASH_EQUALS, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.SLASH, string(l.ch))
		}
	case rune('*'):
		if l.peekChar() == rune('*') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.POW, string(ch)+string(l.ch))
		} else if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.ASTERISK_EQUALS, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.ASTERISK, string(l.ch))
		}
	case rune('<'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.LT_EQUALS, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.LT, string(l.ch))
		}
	case rune('>'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.GT_EQUALS, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.GT, string(l.ch))
		}
	case rune('~'):
		return token.Token{}, fmt.Errorf("unexpected character: %q", l.ch)
	case rune('!'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.NOT_EQ, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.BANG, string(l.ch))
		}
	case rune('\''):
		s, err := l.readString('\'')
		if err != nil {
			tok = l.newToken(token.FSTRING, s)
			l.readChar()
			l.prevToken = tok
			return tok, err
		}
		tok = l.newToken(token.FSTRING, s)
	case rune('"'):
		s, err := l.readString('"')
		if err != nil {
			tok = l.newToken(token.STRING, s)
			l.readChar()
			l.prevToken = tok
			return tok, err
		}
		tok = l.newToken(token.STRING, s)
	case rune('`'):
		s, err := l.readBacktick()
		if err != nil {
			tok = l.newToken(token.BACKTICK, s)
			l.readChar()
			l.prevToken = tok
			return tok, err
		}
		tok = l.newToken(token.BACKTICK, s)
	case rune('['):
		tok = l.newToken(token.LBRACKET, string(l.ch))
	case rune(']'):
		tok = l.newToken(token.RBRACKET, string(l.ch))
	case rune(':'):
		if l.peekChar() == rune('=') {
			l.readChar()
			tok = l.newToken(token.DECLARE, ":=")
		} else {
			tok = l.newToken(token.COLON, string(l.ch))
		}
	case rune('\r'):
		if l.peekChar() == rune('\n') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.NEWLINE, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.NEWLINE, string(l.ch))
		}
	case rune('\n'):
		tok = l.newToken(token.NEWLINE, string(l.ch))
	case rune(0):
		tok = l.newToken(token.EOF, "")
	default:
		if isDigit(l.ch) {
			tok = l.readDecimal()
			l.readChar()
			l.prevToken = tok
			return tok, nil
		}
		ident := l.readIdentifier()
		tok = l.newToken(token.LookupIdentifier(ident), ident)
		l.readChar()
		l.prevToken = tok
		return tok, nil
	}
	l.readChar()
	l.prevToken = tok
	return tok, nil
}

// Read a single identifier
func (l *Lexer) readIdentifier() string {
	var runes []rune
	if isIdentifier(l.ch) {
		runes = append(runes, l.ch)
	}
	for isIdentifier(l.peekChar()) {
		l.readChar()
		runes = append(runes, l.ch)
	}
	return string(runes)
}

// Skip white space
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

// Skip a comment until the end of the line
func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != rune(0) {
		l.readChar()
	}
	l.skipWhitespace()
}

// Consume all tokens until we've had the close of a multi-line comment
func (l *Lexer) skipMultiLineComment() {
	found := false
	for !found {
		// break at the end of our input.
		if l.ch == rune(0) {
			found = true
		}
		// otherwise keep going until we find "*/"
		if l.ch == '*' && l.peekChar() == '/' {
			found = true
			// Our current position is "*", so skip forward to consume the "/"
			l.readChar()
		}
		l.readChar()
	}
	l.skipWhitespace()
}

// Read number. This handles 0x1234 and 0b101010101 too.
func (l *Lexer) readNumber() string {
	str := string(l.ch)
	// We usually just accept digits
	accept := "0123456789"
	// But if we have `0x` as a prefix we accept hexadecimal instead
	if l.ch == '0' && l.peekChar() == 'x' {
		accept = "0x123456789abcdefABCDEF"
	}
	// If we have `0b` as a prefix we accept binary digits only
	if l.ch == '0' && l.peekChar() == 'b' {
		accept = "b01"
	}
	for strings.Contains(accept, string(l.peekChar())) {
		l.readChar()
		str += string(l.ch)
	}
	return str
}

// Read an integer or floating point number
func (l *Lexer) readDecimal() token.Token {
	// Read an integer
	integer := l.readNumber()
	// Check for a period which indicates a floating point
	if l.peekChar() == rune('.') {
		l.readChar()
		if isDigit(l.peekChar()) {
			l.readChar()
			fraction := l.readNumber()
			return l.newToken(token.FLOAT, integer+"."+fraction)
		}
		// This point is reached if the code looks like a method call on an
		// integer, e.g. `42.foo`. TODO: figure out how to handle this. For now,
		// just fall through to create an integer.
	}
	return l.newToken(token.INT, integer)
}

func (l *Lexer) readString(end rune) (string, error) {
	var err error
	var out []string
	for {
		peekChar := l.peekChar()
		if peekChar == rune(0) || peekChar == rune('\n') {
			err = fmt.Errorf("unterminated string literal")
			break
		}
		l.readChar()
		if l.ch == end {
			break
		}
		// Handle \n, \r, \t, \", etc
		if l.ch == '\\' {
			l.readChar()
			if l.ch == rune('n') {
				l.ch = '\n'
			}
			if l.ch == rune('r') {
				l.ch = '\r'
			}
			if l.ch == rune('t') {
				l.ch = '\t'
			}
			if l.ch == rune('\\') {
				l.ch = '\\'
			}
		}
		out = append(out, string(l.ch))
	}
	return strings.Join(out, ""), err
}

func (l *Lexer) readBacktick() (string, error) {
	var err error
	position := l.position + 1
	for {
		peekChar := l.peekChar()
		if peekChar == rune(0) {
			err = fmt.Errorf("unterminated string literal")
			break
		}
		l.readChar()
		if l.ch == '`' {
			break
		}
	}
	return string(l.characters[position:l.position]), err
}

func (l *Lexer) peekChar() rune {
	if l.nextPosition >= len(l.characters) {
		return rune(0)
	}
	return l.characters[l.nextPosition]
}

func isIdentifier(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_'
}

func isWhitespace(ch rune) bool {
	return ch == rune(' ') || ch == rune('\t')
	// ch == rune('\r')
	// ch == rune('\n') ||
}

func isDigit(ch rune) bool {
	return rune('0') <= ch && ch <= rune('9')
}
