package ast

import (
	"bytes"
	"strings"

	"github.com/cloudcmds/tamarin/internal/token"
)

// SetLiteral holds a set definition
type SetLiteral struct {
	Token token.Token // the '{' token
	Items []Expression
}

func (sl *SetLiteral) expressionNode() {}

func (sl *SetLiteral) TokenLiteral() string { return sl.Token.Literal }

func (sl *SetLiteral) String() string {
	var out bytes.Buffer
	items := make([]string, 0, len(sl.Items))
	for _, key := range sl.Items {
		items = append(items, key.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("}")
	return out.String()
}
