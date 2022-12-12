// Package token contains constants which are used when lexing a program
// written in the monkey language, as done by the parser.
package token

// Type is a string
type Type string

// Position points to a particular location in an input string.
// It helps track offset since the beginning of the input as well
// as line offsets.
type Position struct {
	Value     rune
	Char      int
	LineStart int
	Line      int
	Column    int
}

func (p Position) LineNumber() int {
	return p.Line + 1
}

func (p Position) ColumnNumber() int {
	return p.Column + 1
}

// Token struct represent the lexer token
type Token struct {
	Type          Type
	Literal       string
	StartPosition Position
	EndPosition   Position
}

// pre-defined Type
const (
	AND             = "&&"
	ASSIGN          = "="
	ASTERISK        = "*"
	ASTERISK_EQUALS = "*="
	BACKTICK        = "`"
	BANG            = "!"
	CASE            = "case"
	COLON           = ":"
	COMMA           = ","
	CONST           = "CONST"
	CONTAINS        = "~="
	DECLARE         = ":="
	DEFAULT         = "DEFAULT"
	FUNC            = "FUNC"
	DOTDOT          = ".."
	ELSE            = "ELSE"
	EOF             = "EOF"
	EQ              = "=="
	FALSE           = "FALSE"
	FLOAT           = "FLOAT"
	FOR             = "FOR"
	GT              = ">"
	GT_EQUALS       = ">="
	IDENT           = "IDENT"
	IF              = "IF"
	ILLEGAL         = "ILLEGAL"
	INT             = "INT"
	LBRACE          = "{"
	LBRACKET        = "["
	LET             = "LET"
	LPAREN          = "("
	LT              = "<"
	LT_EQUALS       = "<="
	MINUS           = "-"
	MINUS_EQUALS    = "-="
	MINUS_MINUS     = "--"
	MOD             = "%"
	NOT_CONTAINS    = "!~"
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
	REGEXP          = "REGEXP"
	RETURN          = "RETURN"
	RPAREN          = ")"
	SEMICOLON       = ";"
	SLASH           = "/"
	SLASH_EQUALS    = "/="
	STRING          = "STRING"
	FSTRING         = "FSTRING"
	SWITCH          = "switch"
	TRUE            = "TRUE"
	NEWLINE         = "EOL"
	IMPORT          = "IMPORT"
	BREAK           = "BREAK"
)

// reserved keywords
var keywords = map[string]Type{
	"case":    CASE,
	"const":   CONST,
	"default": DEFAULT,
	"else":    ELSE,
	"false":   FALSE,
	"for":     FOR,
	"func":    FUNC,
	"if":      IF,
	"let":     LET,
	"nil":     NIL,
	"return":  RETURN,
	"switch":  SWITCH,
	"true":    TRUE,
	"import":  IMPORT,
	"break":   BREAK,
}

// LookupIdentifier used to determinate whether identifier is keyword nor not
func LookupIdentifier(identifier string) Type {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENT
}
