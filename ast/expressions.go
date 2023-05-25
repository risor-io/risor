// Package ast contains the definitions of the abstract syntax tree that our
// parser produces and that our interpreter executes.
package ast

import (
	"bytes"
	"strings"

	"github.com/cloudcmds/tamarin/token"
)

// Ident is an expression that refers to a variable by name.
type Ident struct {
	token token.Token
	value string
}

func NewIdent(token token.Token) *Ident {
	return &Ident{token: token, value: token.Literal}
}

func (i *Ident) ExpressionNode() {}

func (i *Ident) Token() token.Token { return i.token }

func (i *Ident) Literal() string { return i.value }

func (i *Ident) String() string { return i.value }

// Prefix is an operator expression where the operator precedes the operand.
// Examples include "!false" and "-x".
type Prefix struct {
	token token.Token

	// operator holds the operator being invoked (e.g. "!")
	operator string

	// right holds the thing to be operated upon
	right Expression
}

func NewPrefix(token token.Token, right Expression) *Prefix {
	return &Prefix{token: token, operator: token.Literal, right: right}
}

func (p *Prefix) ExpressionNode() {}

func (p *Prefix) Token() token.Token { return p.token }

func (p *Prefix) Literal() string { return p.token.Literal }

func (p *Prefix) Operator() string { return p.operator }

func (p *Prefix) Right() Expression { return p.right }

func (p *Prefix) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.operator)
	out.WriteString(p.right.String())
	out.WriteString(")")
	return out.String()
}

// Infix is an operator expression where the operator is between the operands.
// Examples include "x + y" and "5 - 1".
type Infix struct {
	token token.Token

	// left side expression
	left Expression

	// operator e.g. "+", "-", "==", etc.
	operator string

	// right side expression
	right Expression
}

func NewInfix(token token.Token, left Expression, operator string, right Expression) *Infix {
	return &Infix{token: token, left: left, operator: operator, right: right}
}

func (i *Infix) ExpressionNode() {}

func (i *Infix) Token() token.Token { return i.token }

func (i *Infix) Literal() string { return i.token.Literal }

func (i *Infix) Left() Expression { return i.left }

func (i *Infix) Operator() string { return i.operator }

func (i *Infix) Right() Expression { return i.right }

func (i *Infix) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.left.String())
	out.WriteString(" " + i.operator + " ")
	out.WriteString(i.right.String())
	out.WriteString(")")
	return out.String()
}

// Postfix defines a postfix expression like "x++".
type Postfix struct {
	token token.Token
	// operator holds the postfix token, e.g. ++
	operator string
}

func NewPostfix(token token.Token, operator string) *Postfix {
	return &Postfix{token: token, operator: operator}
}

func (p *Postfix) ExpressionNode() {}

func (p *Postfix) Token() token.Token { return p.token }

func (p *Postfix) Literal() string { return p.token.Literal }

func (p *Postfix) Operator() string { return p.operator }

func (p *Postfix) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.token.Literal)
	out.WriteString(p.operator)
	out.WriteString(")")
	return out.String()
}

// If holds an if statement.
type If struct {
	token       token.Token // the "if" token
	condition   Expression  // the condition to be evaluated
	consequence *Block      // block to evaluate if the condition is true
	alternative *Block      // block to evaluate if the condition is false
}

func NewIf(token token.Token, condition Expression, consequence *Block, alternative *Block) *If {
	return &If{token: token, condition: condition, consequence: consequence, alternative: alternative}
}

func (i *If) ExpressionNode() {}

func (i *If) Token() token.Token { return i.token }

func (i *If) Literal() string { return i.token.Literal }

func (i *If) Condition() Expression { return i.condition }

func (i *If) Consequence() *Block { return i.consequence }

func (i *If) Alternative() *Block { return i.alternative }

func (i *If) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(i.condition.String())
	out.WriteString(" ")
	out.WriteString(i.consequence.String())
	if i.alternative != nil {
		out.WriteString(" else ")
		out.WriteString(i.alternative.String())
	}
	return out.String()
}

// Ternary holds a ternary expression.
type Ternary struct {
	token token.Token

	// condition is the thing that is evaluated to determine
	// which expression should be returned
	condition Expression

	// ifTrue is the expression to return if the condition is true.
	ifTrue Expression

	// ifFalse is the expression to return if the condition is not true.
	ifFalse Expression
}

func NewTernary(token token.Token, condition Expression, ifTrue Expression, ifFalse Expression) *Ternary {
	return &Ternary{token: token, condition: condition, ifTrue: ifTrue, ifFalse: ifFalse}
}

func (t *Ternary) ExpressionNode() {}

func (t *Ternary) Token() token.Token { return t.token }

func (t *Ternary) Literal() string { return t.token.Literal }

func (t *Ternary) Condition() Expression { return t.condition }

func (t *Ternary) IfTrue() Expression { return t.ifTrue }

func (t *Ternary) IfFalse() Expression { return t.ifFalse }

func (t *Ternary) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(t.condition.String())
	out.WriteString(" ? ")
	out.WriteString(t.ifTrue.String())
	out.WriteString(" : ")
	out.WriteString(t.ifFalse.String())
	out.WriteString(")")
	return out.String()
}

// Call holds the invocation of a method.
type Call struct {
	token     token.Token  // the '(' token
	function  Expression   // the function being called
	arguments []Expression // the arguments supplied to the call
}

func NewCall(token token.Token, function Expression, arguments []Expression) *Call {
	return &Call{token: token, function: function, arguments: arguments}
}

func (c *Call) ExpressionNode() {}

func (c *Call) Token() token.Token { return c.token }

func (c *Call) Literal() string { return c.token.Literal }

func (c *Call) Function() Expression { return c.function }

func (c *Call) Arguments() []Expression { return c.arguments }

func (c *Call) String() string {
	var out bytes.Buffer
	args := make([]string, 0)
	for _, a := range c.arguments {
		args = append(args, a.String())
	}
	out.WriteString(c.function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// GetAttr
type GetAttr struct {
	token token.Token

	// object whose attribute is being accessed
	object Expression

	// The attribute itself
	attribute *Ident
}

func NewGetAttr(token token.Token, object Expression, attribute *Ident) *GetAttr {
	return &GetAttr{token: token, object: object, attribute: attribute}
}

func (e *GetAttr) ExpressionNode() {}

func (e *GetAttr) Token() token.Token { return e.token }

func (e *GetAttr) Literal() string { return e.token.Literal }

func (e *GetAttr) Object() Expression { return e.object }

func (e *GetAttr) Name() string { return e.attribute.value }

func (e *GetAttr) String() string {
	var out bytes.Buffer
	out.WriteString(e.object.String())
	out.WriteString(".")
	out.WriteString(e.attribute.value)
	return out.String()
}

// Pipe holds a series of calls
type Pipe struct {
	token token.Token

	// exprs contains the pipe separated expressions
	exprs []Expression
}

func NewPipe(token token.Token, exprs []Expression) *Pipe {
	return &Pipe{token: token, exprs: exprs}
}

func (p *Pipe) ExpressionNode() {}

func (p *Pipe) Token() token.Token { return p.token }

func (p *Pipe) Literal() string { return p.token.Literal }

func (p *Pipe) Expressions() []Expression { return p.exprs }

func (p *Pipe) String() string {
	var out bytes.Buffer
	args := make([]string, 0)
	for _, a := range p.exprs {
		args = append(args, a.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(args, " | "))
	out.WriteString(")")
	return out.String()
}

// ObjectCall is used when calling a method on an object.
type ObjectCall struct {
	token  token.Token
	object Expression
	call   Expression
}

func NewObjectCall(token token.Token, object Expression, call Expression) *ObjectCall {
	return &ObjectCall{token: token, object: object, call: call}
}

func (c *ObjectCall) ExpressionNode() {}

func (c *ObjectCall) Token() token.Token { return c.token }

func (c *ObjectCall) Literal() string { return c.token.Literal }

func (c *ObjectCall) Object() Expression { return c.object }

func (c *ObjectCall) Call() Expression { return c.call }

func (c *ObjectCall) String() string {
	var out bytes.Buffer
	out.WriteString(c.object.String())
	out.WriteString(".")
	out.WriteString(c.call.String())
	return out.String()
}

// Index holds an index expression
type Index struct {
	token token.Token

	// left is the thing being indexed.
	left Expression

	// index is the value we're indexing
	index Expression
}

func NewIndex(token token.Token, left Expression, index Expression) *Index {
	return &Index{token: token, left: left, index: index}
}

func (i *Index) ExpressionNode() {}

func (i *Index) Token() token.Token { return i.token }

func (i *Index) Literal() string { return i.token.Literal }

func (i *Index) Left() Expression { return i.left }

func (i *Index) Index() Expression { return i.index }

func (i *Index) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.left.String())
	out.WriteString("[")
	out.WriteString(i.index.String())
	out.WriteString("])")
	return out.String()
}

// Index holds an index expression
type Slice struct {
	token token.Token

	// left is the thing being indexed.
	left Expression

	// Optional "from" index for [from:to] style expressions
	fromIndex Expression

	// Optional "to" index for [from:to] style expressions
	toIndex Expression
}

func NewSlice(token token.Token, left Expression, fromIndex Expression, toIndex Expression) *Slice {
	return &Slice{token: token, left: left, fromIndex: fromIndex, toIndex: toIndex}
}

func (i *Slice) ExpressionNode() {}

func (i *Slice) Token() token.Token { return i.token }

func (i *Slice) Literal() string { return i.token.Literal }

func (i *Slice) Left() Expression { return i.left }

func (i *Slice) FromIndex() Expression { return i.fromIndex }

func (i *Slice) ToIndex() Expression { return i.toIndex }

func (i *Slice) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.left.String())
	out.WriteString("[")
	if i.fromIndex != nil {
		out.WriteString(i.fromIndex.String())
	}
	if i.toIndex != nil {
		out.WriteString(":")
		out.WriteString(i.toIndex.String())
	}
	out.WriteString("])")
	return out.String()
}

// Case handles the case within a switch statement
type Case struct {
	token token.Token

	// Default branch?
	isDefault bool

	// The thing we match
	expr []Expression

	// The code to execute if there is a match
	block *Block
}

func NewCase(token token.Token, expressions []Expression, block *Block) *Case {
	return &Case{token: token, expr: expressions, block: block}
}

func NewDefaultCase(token token.Token, block *Block) *Case {
	return &Case{token: token, isDefault: true, block: block}
}

func (c *Case) ExpressionNode() {}

func (c *Case) Token() token.Token { return c.token }

func (c *Case) Literal() string { return c.token.Literal }

func (c *Case) IsDefault() bool { return c.isDefault }

func (c *Case) Expressions() []Expression { return c.expr }

func (c *Case) Block() *Block { return c.block }

func (c *Case) String() string {
	var out bytes.Buffer
	if c.isDefault {
		out.WriteString("default")
	} else {
		out.WriteString("case ")
		tmp := []string{}
		for _, exp := range c.expr {
			tmp = append(tmp, exp.String())
		}
		out.WriteString(strings.Join(tmp, ","))
	}
	out.WriteString(":\n")
	for i, exp := range c.block.statements {
		if i > 0 {
			out.WriteString("\n")
		}
		out.WriteString("\t" + exp.String())
	}
	out.WriteString("\n")
	return out.String()
}

// Switch represents a switch statement and its cases
type Switch struct {
	token   token.Token // token containing "switch"
	value   Expression  // the expression to switch on
	choices []*Case     // switch cases
}

func NewSwitch(token token.Token, value Expression, choices []*Case) *Switch {
	return &Switch{token: token, value: value, choices: choices}
}

func (s *Switch) ExpressionNode() {}

func (s *Switch) Token() token.Token { return s.token }

func (s *Switch) Literal() string { return s.token.Literal }

func (s *Switch) Value() Expression { return s.value }

func (s *Switch) Choices() []*Case { return s.choices }

func (s *Switch) String() string {
	var out bytes.Buffer
	out.WriteString("\nswitch ")
	out.WriteString(s.value.String())
	out.WriteString(" {\n")
	for _, tmp := range s.choices {
		if tmp != nil {
			out.WriteString(tmp.String())
		}
	}
	out.WriteString("}\n")
	return out.String()
}

// In is an expression that evalutes to a boolean and checks if the left
// expression is contained within the right expression.
type In struct {
	token token.Token
	left  Expression
	right Expression
}

func NewIn(token token.Token, left Expression, right Expression) *In {
	return &In{token: token, left: left, right: right}
}

func (i *In) ExpressionNode() {}

func (i *In) Token() token.Token { return i.token }

func (i *In) Literal() string { return i.token.Literal }

func (i *In) Left() Expression { return i.left }

func (i *In) Right() Expression { return i.right }

func (i *In) String() string {
	var out bytes.Buffer
	out.WriteString(i.left.String())
	out.WriteString(" in ")
	out.WriteString(i.right.String())
	return out.String()
}
