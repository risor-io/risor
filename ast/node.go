package ast

import "github.com/cloudcmds/tamarin/token"

// Node reresents a portion of the syntax tree. All nodes have a token, which is
// the token that begins the node. A Node may be an Expression, in which case
// it evaluates to a value.
type Node interface {

	// Token returns the token where this Node begins.
	Token() token.Token

	// Literal returns the string from the first token that defines the Node.
	Literal() string

	// String returns a human friendly representation of the Node. This should
	// be similar to the original source code, but not necessarily identical.
	String() string
}

// Expression represents a snippet of Tamarin code that evaluates to a value.
// Expressions may be embedded within other expressions.
type Expression interface {
	// Node is embedded here to indicate that all expressions are AST nodes.
	Node

	// ExpressionNode signals that this Node is an expression.
	ExpressionNode()
}
