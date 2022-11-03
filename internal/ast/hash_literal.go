package ast

import (
	"bytes"
	"strings"

	"github.com/myzie/tamarin/internal/token"
)

// HashLiteral holds a hash definition
type HashLiteral struct {
	// Token holds the token
	Token token.Token // the '{' token

	// Pairs stores the name/value sets of the hash-content
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}

// TokenLiteral returns the literal token.
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }

// String returns this object as a string.
func (hl *HashLiteral) String() string {
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
