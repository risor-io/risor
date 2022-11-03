// Package lexer contains the code to lex input-programs into a stream
// of tokens, such that they may be parsed.
package lexer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/myzie/tamarin/internal/token"
)

// Lexer holds our object-state.
type Lexer struct {
	// The current character position
	position int

	// The next character position
	readPosition int

	// The current character
	ch rune

	// A rune slice of our input string
	characters []rune

	// Previous token
	prevToken token.Token

	// Line number of the current character
	lineNumber int

	// Line position of the current character
	linePosition int

	// Line position of the start of the current token
	tokenStartPosition int
}

// New returns a Lexer instance for a given string input.
func New(input string) *Lexer {
	l := &Lexer{
		characters:   []rune(input),
		linePosition: -1,
	}
	l.readChar()
	return l
}

// read one forward character
func (l *Lexer) readChar() {
	// Track current line and position within the line
	if l.ch == rune('\n') {
		l.lineNumber++
		l.linePosition = 0
	} else {
		l.linePosition++
	}
	// Set the current character and overall read position
	if l.readPosition >= len(l.characters) {
		l.ch = rune(0)
	} else {
		l.ch = l.characters[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	// fmt.Println("readChar done; linePos:", l.linePosition,
	// 	"lineNum:", l.lineNumber, "char:", string(l.ch))
}

func (l *Lexer) newToken(typ token.Type, literal string) token.Token {
	t := token.Token{
		Type:          typ,
		Literal:       literal,
		Line:          l.lineNumber,
		StartPosition: l.tokenStartPosition,
		EndPosition:   l.linePosition,
	}
	// fmt.Printf("%+v\n", t)
	return t
}

// NextToken to read next token, skipping the white space.
func (l *Lexer) NextToken() token.Token {

	var tok token.Token
	l.skipWhitespace()
	l.tokenStartPosition = l.linePosition

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
			// Slash is mostly division, but could be the start of a regex
			// We exclude:
			//   a[b] / c        -> RBRACKET
			//   ( a + b ) / c   -> RPAREN
			//   a / c           -> IDENT
			//   3.2 / c         -> FLOAT
			//   1 / c           -> IDENT
			if l.prevToken.Type == token.RBRACKET ||
				l.prevToken.Type == token.RPAREN ||
				l.prevToken.Type == token.IDENT ||
				l.prevToken.Type == token.INT ||
				l.prevToken.Type == token.FLOAT {
				tok = l.newToken(token.SLASH, string(l.ch))
			} else {
				str, err := l.readRegexp()
				if err == nil {
					tok = l.newToken(token.REGEXP, str)
				} else {
					tok = l.newToken(token.REGEXP, str)
				}
			}
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
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.CONTAINS, string(ch)+string(l.ch))
		}
	case rune('!'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.NOT_EQ, string(ch)+string(l.ch))
		} else {
			if l.peekChar() == rune('~') {
				ch := l.ch
				l.readChar()
				tok = l.newToken(token.NOT_CONTAINS, string(ch)+string(l.ch))
			} else {
				tok = l.newToken(token.BANG, string(l.ch))
			}
		}
	case rune('"'):
		tok = l.newToken(token.STRING, l.readString())
	case rune('`'):
		tok = l.newToken(token.BACKTICK, l.readBacktick())
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
			return tok
		}
		ident := l.readIdentifier()
		tok = l.newToken(token.LookupIdentifier(ident), ident)
		l.readChar()
		l.prevToken = tok
		return tok
	}
	l.readChar()
	l.prevToken = tok
	return tok
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

// Read string
func (l *Lexer) readString() string {
	out := ""
	for {
		l.readChar()
		if l.ch == '"' {
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
			if l.ch == rune('"') {
				l.ch = '"'
			}
			if l.ch == rune('\\') {
				l.ch = '\\'
			}
		}
		out = out + string(l.ch)
	}
	return out
}

// Read a regexp, including flags.
func (l *Lexer) readRegexp() (string, error) {
	out := ""

	for {
		l.readChar()

		if l.ch == rune(0) {
			return "unterminated regular expression", fmt.Errorf("unterminated regular expression")
		}
		if l.ch == '/' {

			// consume the terminating "/".
			l.readChar()

			// prepare to look for flags
			flags := ""

			// two flags are supported:
			//   i -> Ignore-case
			//   m -> Multiline
			//
			for l.ch == rune('i') || l.ch == rune('m') {

				// save the char - unless it is a repeat
				if !strings.Contains(flags, string(l.ch)) {

					// we're going to sort the flags
					tmp := strings.Split(flags, "")
					tmp = append(tmp, string(l.ch))
					flags = strings.Join(tmp, "")

				}

				// read the next
				l.readChar()
			}

			// convert the regexp to go-lang
			if len(flags) > 0 {
				out = "(?" + flags + ")" + out
			}
			break
		}
		out = out + string(l.ch)
	}

	return out, nil
}

func (l *Lexer) readBacktick() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '`' {
			break
		}
	}
	out := string(l.characters[position:l.position])
	return out
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.characters) {
		return rune(0)
	}
	return l.characters[l.readPosition]
}

func isIdentifier(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_'
}

func isWhitespace(ch rune) bool {
	return ch == rune(' ') || ch == rune('\t')
	// ch == rune('\r')
	// ch == rune('\n') ||
}

func isWhitespaceWithNewlines(ch rune) bool {
	return ch == rune(' ') ||
		ch == rune('\t') ||
		ch == rune('\r') ||
		ch == rune('\n')
}

func isDigit(ch rune) bool {
	return rune('0') <= ch && ch <= rune('9')
}
