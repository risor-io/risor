// Package token defines language keywords and tokens used when lexing source code.
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
	AND             Type = "&&"
	AS              Type = "AS"
	ASSIGN          Type = "="
	ASTERISK        Type = "*"
	ASTERISK_EQUALS Type = "*="
	BACKTICK        Type = "`"
	BANG            Type = "!"
	BREAK           Type = "BREAK"
	CASE            Type = "CASE"
	COLON           Type = ":"
	COMMA           Type = ","
	CONST           Type = "CONST"
	CONTINUE        Type = "CONTINUE"
	DECLARE         Type = ":="
	DEFAULT         Type = "DEFAULT"
	DEFER           Type = "DEFER"
	ELLIPSIS        Type = "..."
	ELSE            Type = "ELSE"
	EOF             Type = "EOF"
	EQ              Type = "=="
	FALSE           Type = "FALSE"
	FLOAT           Type = "FLOAT"
	FOR             Type = "FOR"
	FROM            Type = "FROM"
	FSTRING         Type = "'"
	FUNC            Type = "FUNC"
	GO              Type = "GO"
	GT              Type = ">"
	GT_EQUALS       Type = ">="
	GT_GT           Type = ">>"
	IDENT           Type = "IDENT"
	IF              Type = "IF"
	ILLEGAL         Type = "ILLEGAL"
	IMPORT          Type = "IMPORT"
	IN              Type = "IN"
	INT             Type = "INT"
	LBRACE          Type = "{"
	LBRACKET        Type = "["
	LPAREN          Type = "("
	LT              Type = "<"
	LT_EQUALS       Type = "<="
	LT_LT           Type = "<<"
	MINUS           Type = "-"
	MINUS_EQUALS    Type = "-="
	MINUS_MINUS     Type = "--"
	MOD             Type = "%"
	NEWLINE         Type = "EOL"
	NIL             Type = "nil"
	NOT_EQ          Type = "!="
	OR              Type = "||"
	PERIOD          Type = "."
	PIPE            Type = "|"
	PLUS            Type = "+"
	PLUS_EQUALS     Type = "+="
	PLUS_PLUS       Type = "++"
	POW             Type = "**"
	QUESTION        Type = "?"
	RANGE           Type = "RANGE"
	RBRACE          Type = "}"
	RBRACKET        Type = "]"
	RETURN          Type = "RETURN"
	RPAREN          Type = ")"
	SEMICOLON       Type = ";"
	SEND            Type = "<-"
	SLASH           Type = "/"
	SLASH_EQUALS    Type = "/="
	STRING          Type = "STRING"
	STRUCT          Type = "STRUCT"
	SWITCH          Type = "SWITCH"
	TRUE            Type = "TRUE"
	VAR             Type = "VAR"
)

// Reserved keywords
var keywords = map[string]Type{
	"as":       AS,
	"break":    BREAK,
	"case":     CASE,
	"const":    CONST,
	"continue": CONTINUE,
	"default":  DEFAULT,
	"defer":    DEFER,
	"else":     ELSE,
	"false":    FALSE,
	"for":      FOR,
	"from":     FROM,
	"func":     FUNC,
	"go":       GO,
	"if":       IF,
	"import":   IMPORT,
	"in":       IN,
	"nil":      NIL,
	"range":    RANGE,
	"return":   RETURN,
	"struct":   STRUCT,
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
