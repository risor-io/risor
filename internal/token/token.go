// Package token contains constants which are used when lexing a program
// written in the monkey language, as done by the parser.
package token

// Type is a string
type Type string

// Token struct represent the lexer token
type Token struct {
	Type          Type
	Literal       string
	Line          int
	StartPosition int
	EndPosition   int
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
	NULL            = "null"
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
	SWITCH          = "switch"
	TRUE            = "TRUE"
	NEWLINE         = "EOL"
	IMPORT          = "IMPORT"
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
	"null":    NULL,
	"return":  RETURN,
	"switch":  SWITCH,
	"true":    TRUE,
	"import":  IMPORT,
}

// LookupIdentifier used to determinate whether identifier is keyword nor not
func LookupIdentifier(identifier string) Type {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENT
}
