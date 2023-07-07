// Package token defines the tokens that are produced when lexing Risor code.
package token

// Type describes the type of a token as a string.
type Type string

// Position points to a particular location in an input string.
type Position struct {
	Value     rune
	Char      int
	LineStart int
	Line      int
	Column    int
	File      string
}

// LineNumber returns the 1-indexed line number for this position in the input.
func (p Position) LineNumber() int {
	return p.Line + 1
}

// ColumnNumber returns the 1-indexed column number for this position in the input.
func (p Position) ColumnNumber() int {
	return p.Column + 1
}

// Token represents one token lexed from the input source code.
type Token struct {
	Type          Type
	Literal       string
	StartPosition Position
	EndPosition   Position
}

// Token types
const (
	AND             = "&&"
	ASSIGN          = "="
	ASTERISK        = "*"
	ASTERISK_EQUALS = "*="
	BACKTICK        = "`"
	FSTRING         = "'"
	BANG            = "!"
	CASE            = "case"
	COLON           = ":"
	COMMA           = ","
	CONST           = "CONST"
	DECLARE         = ":="
	DEFAULT         = "DEFAULT"
	FUNC            = "FUNC"
	ELSE            = "ELSE"
	EOF             = "EOF"
	EQ              = "=="
	FALSE           = "FALSE"
	FLOAT           = "FLOAT"
	FOR             = "FOR"
	GT              = ">"
	GT_GT           = ">>"
	GT_EQUALS       = ">="
	IDENT           = "IDENT"
	IF              = "IF"
	ILLEGAL         = "ILLEGAL"
	INT             = "INT"
	LBRACE          = "{"
	LBRACKET        = "["
	LPAREN          = "("
	LT              = "<"
	LT_LT           = "<<"
	LT_EQUALS       = "<="
	MINUS           = "-"
	MINUS_EQUALS    = "-="
	MINUS_MINUS     = "--"
	MOD             = "%"
	NOT_EQ          = "!="
	NIL             = "nil"
	PIPE            = "|"
	OR              = "||"
	PERIOD          = "."
	PLUS            = "+"
	PLUS_EQUALS     = "+="
	PLUS_PLUS       = "++"
	POW             = "**"
	QUESTION        = "?"
	RBRACE          = "}"
	RBRACKET        = "]"
	RETURN          = "RETURN"
	RPAREN          = ")"
	SEMICOLON       = ";"
	SLASH           = "/"
	SLASH_EQUALS    = "/="
	STRING          = "STRING"
	SWITCH          = "switch"
	TRUE            = "TRUE"
	NEWLINE         = "EOL"
	IMPORT          = "IMPORT"
	BREAK           = "BREAK"
	CONTINUE        = "CONTINUE"
	VAR             = "VAR"
	IN              = "IN"
	RANGE           = "RANGE"
)

// Reserved keywords
var keywords = map[string]Type{
	"break":    BREAK,
	"case":     CASE,
	"const":    CONST,
	"continue": CONTINUE,
	"default":  DEFAULT,
	"else":     ELSE,
	"false":    FALSE,
	"for":      FOR,
	"func":     FUNC,
	"if":       IF,
	"import":   IMPORT,
	"in":       IN,
	"nil":      NIL,
	"range":    RANGE,
	"return":   RETURN,
	"switch":   SWITCH,
	"true":     TRUE,
	"var":      VAR,
}

// LookupIdentifier used to determinate whether identifier is keyword nor not
func LookupIdentifier(identifier string) Type {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENT
}
