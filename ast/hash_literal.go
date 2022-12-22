package ast

import (
	"bytes"
	"strings"

	"github.com/cloudcmds/tamarin/token"
)

// MapLiteral holds a hash definition
type MapLiteral struct {
	// Token holds the token
	Token token.Token // the '{' token

	// Pairs stores the name/value sets of the hash-content
	Pairs map[Expression]Expression
}

func (hl *MapLiteral) expressionNode() {}

// TokenLiteral returns the literal token.
func (hl *MapLiteral) TokenLiteral() string { return hl.Token.Literal }

// String returns this object as a string.
func (hl *MapLiteral) String() string {
	var out bytes.Buffer
	pairs := make([]string, 0)
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
