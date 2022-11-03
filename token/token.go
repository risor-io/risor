// Package token contains constants which are used when lexing a program
// written in the monkey language, as done by the parser.
package token

// Type is a string
type Type string

// Token struct represent the lexer token
type Token struct {
	Type    Type
	Literal string
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
	DEFAULT         = "DEFAULT"
	DEFINE_FUNCTION = "DEFINE_FUNCTION"
	DOTDOT          = ".."
	ELSE            = "ELSE"
	EOF             = "EOF"
	EQ              = "=="
	FALSE           = "FALSE"
	FLOAT           = "FLOAT"
	FOR             = "FOR"
	FOREACH         = "FOREACH"
	FUNCTION        = "FUNCTION"
	GT              = ">"
	GT_EQUALS       = ">="
	IDENT           = "IDENT"
	IF              = "IF"
	ILLEGAL         = "ILLEGAL"
	IN              = "IN"
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
)

// reserved keywords
var keywords = map[string]Type{
	"case":     CASE,
	"const":    CONST,
	"default":  DEFAULT,
	"else":     ELSE,
	"false":    FALSE,
	"fn":       FUNCTION,
	"for":      FOR,
	"foreach":  FOREACH,
	"function": DEFINE_FUNCTION,
	"if":       IF,
	"in":       IN,
	"let":      LET,
	"null":     NULL,
	"return":   RETURN,
	"switch":   SWITCH,
	"true":     TRUE,
}

// LookupIdentifier used to determinate whether identifier is keyword nor not
func LookupIdentifier(identifier string) Type {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENT
}
