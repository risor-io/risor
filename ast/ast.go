// Package ast contains the definitions of the abstract syntax tree that our
// parser produces and that our interpreter executes.
package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cloudcmds/tamarin/tmpl"
	"github.com/cloudcmds/tamarin/token"
)

// Node reresents a node.
type Node interface {

	// Token returns the token where this Node begins.
	Token() token.Token

	// Literal returns the string from the first token that defines the Node.
	Literal() string

	// String returns a human friendly representation of the Node. This should
	// be similar to the original source code, but not necessarily identical.
	String() string
}

// Statement represents a single element of execution in a Tamarin program.
// Programs are made of a series of statements.
// type Statement interface {
// 	// Node is embedded here to indicate that all statements are AST nodes.
// 	Node
// 	StatementNode()
// }

type Statement Node

// Expression represents a snippet of Tamarin code that evaluates to a value.
// Expressions are used in statements, as well as in other expressions.
type Expression interface {
	// Node is embedded here to indicate that all expressions are AST nodes.
	Node
	ExpressionNode()
}

// Program represents a complete program.
type Program struct {
	// statements is the set of statements which comprise the program
	statements []Node
}

func NewProgram(statements []Node) *Program {
	return &Program{statements: statements}
}

func (p *Program) Token() token.Token {
	if len(p.statements) > 0 {
		return p.statements[0].Token()
	}
	return token.Token{}
}

func (p *Program) Literal() string { return p.Token().Literal }

func (p *Program) Statements() []Node { return p.statements }

func (p *Program) First() Node {
	if len(p.statements) > 0 {
		return p.statements[0]
	}
	return nil
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, stmt := range p.statements {
		out.WriteString(stmt.String())
	}
	return out.String()
}

// Var holds a variable declaration statement.
type Var struct {
	token token.Token

	// name is the name of the variable to which we're assigning
	name *Ident

	// value is the thing we're storing in the variable.
	value Expression

	// isWalrus is true if this is a ":=" statement.
	isWalrus bool
}

func NewVar(token token.Token, name *Ident, value Expression) *Var {
	return &Var{token: token, name: name, value: value}
}

func NewDeclaration(token token.Token, name *Ident, value Expression) *Var {
	return &Var{token: token, name: name, value: value, isWalrus: true}
}

func (s *Var) StatementNode() {}

func (s *Var) Token() token.Token { return s.token }

func (s *Var) Literal() string { return s.token.Literal }

func (s *Var) Value() (string, Expression) { return s.name.value, s.value }

func (s *Var) IsWalrus() bool { return s.isWalrus }

func (s *Var) String() string {
	var out bytes.Buffer
	if s.isWalrus {
		out.WriteString(s.name.Literal() + " := ")
		out.WriteString(s.value.String())
		return out.String()
	}
	out.WriteString(s.Literal() + " ")
	out.WriteString(s.name.Literal())
	out.WriteString(" = ")
	if s.value != nil {
		out.WriteString(s.value.String())
	}
	return out.String()
}

// MultiVar holds a variable assignment statement for >1 variables.
type MultiVar struct {
	token    token.Token
	names    []*Ident   // names being assigned
	value    Expression // value is the thing we're storing in the variable.
	isWalrus bool       // isWalrus is true if this is a ":=" statement.
}

func NewMultiVar(token token.Token, names []*Ident, value Expression, isWalrus bool) *MultiVar {
	return &MultiVar{token: token, names: names, value: value, isWalrus: isWalrus}
}

func (s *MultiVar) StatementNode() {}

func (s *MultiVar) Token() token.Token { return s.token }

func (s *MultiVar) Literal() string { return s.token.Literal }

func (s *MultiVar) Value() ([]string, Expression) {
	names := make([]string, 0, len(s.names))
	for _, name := range s.names {
		names = append(names, name.value)
	}
	return names, s.value
}

func (s *MultiVar) IsWalrus() bool { return s.isWalrus }

func (s *MultiVar) String() string {
	names, expr := s.Value()
	namesStr := strings.Join(names, ", ")
	var out bytes.Buffer
	if s.isWalrus {
		out.WriteString(namesStr + " := ")
		out.WriteString(expr.String())
		return out.String()
	}
	out.WriteString(s.Literal() + " ")
	out.WriteString(namesStr)
	out.WriteString(" = ")
	out.WriteString(expr.String())
	return out.String()
}

// Const defines a named constant containing a constant value.
type Const struct {
	token token.Token // the "const" token
	name  *Ident      // name of the constant
	value Expression  // value of the constant
}

func NewConst(token token.Token, name *Ident, value Expression) *Const {
	return &Const{token: token, name: name, value: value}
}

func (c *Const) StatementNode() {}

func (c *Const) Token() token.Token { return c.token }

func (c *Const) Literal() string { return c.token.Literal }

func (c *Const) Value() (string, Expression) { return c.name.value, c.value }

func (c *Const) String() string {
	var out bytes.Buffer
	out.WriteString(c.Literal() + " ")
	out.WriteString(c.name.Literal())
	out.WriteString(" = ")
	if c.value != nil {
		out.WriteString(c.value.String())
	}
	return out.String()
}

// Ident holds a single identifier.
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

// Control defines a return, break, or continue statement.
type Control struct {
	token token.Token // "return", "break", or "continue"
	value Expression  // optional value, for return statements
}

func NewControl(token token.Token, value Expression) *Control {
	return &Control{token: token, value: value}
}

func (c *Control) StatementNode() {}

func (c *Control) Token() token.Token { return c.token }

func (c *Control) Literal() string { return c.token.Literal }

func (c *Control) Value() Expression { return c.value }

func (c *Control) String() string {
	var out bytes.Buffer
	out.WriteString(c.Literal() + " ")
	if c.value != nil {
		out.WriteString(c.value.Literal())
	}
	out.WriteString(";")
	return out.String()
}

// Int holds an integer number
type Int struct {
	token token.Token // the token containing the number
	value int64       // the value of the int
}

func NewInt(token token.Token, value int64) *Int {
	return &Int{token: token, value: value}
}

func (i *Int) ExpressionNode() {}

func (i *Int) Token() token.Token { return i.token }

func (i *Int) Literal() string { return i.token.Literal }

func (i *Int) Value() int64 { return i.value }

func (i *Int) String() string { return i.token.Literal }

// Float holds a floating point number
type Float struct {
	token token.Token // the token containing the number
	value float64     // the value of the float
}

func NewFloat(token token.Token, value float64) *Float {
	return &Float{token: token, value: value}
}

func (f *Float) ExpressionNode() {}

func (f *Float) Token() token.Token { return f.token }

func (f *Float) Literal() string { return f.token.Literal }

func (f *Float) Value() float64 { return f.value }

func (f *Float) String() string { return f.token.Literal }

// Prefix defines a prefix expression like "!false" or "-x".
type Prefix struct {
	token token.Token

	// operator holds the operator being invoked (e.g. "!" ).
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

// Infix defines an infix expression like "x + y" or "5 - 1".
type Infix struct {
	// Token holds the literal expression
	token token.Token

	// left holds the left-most argument
	left Expression

	// operator holds the operation to be carried out (e.g. "+", "-" )
	operator string

	// right holds the right-most argument
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

// Nil represents a literal nil
type Nil struct {
	token token.Token // token containing "nil"
}

func NewNil(token token.Token) *Nil {
	return &Nil{token: token}
}

func (n *Nil) ExpressionNode() {}

func (n *Nil) Token() token.Token { return n.token }

func (n *Nil) Literal() string { return n.token.Literal }

func (n *Nil) String() string { return n.token.Literal }

// Bool holds the boolean "true" or "false"
type Bool struct {
	// Token holds the actual token
	token token.Token

	// value stores the bools' value: true, or false.
	value bool
}

func NewBool(token token.Token, value bool) *Bool {
	return &Bool{token: token, value: value}
}

func (b *Bool) ExpressionNode() {}

func (b *Bool) Token() token.Token { return b.token }

func (b *Bool) Literal() string { return b.token.Literal }

func (b *Bool) Value() bool { return b.value }

func (b *Bool) String() string { return b.token.Literal }

// Block holds a sequence of statements, which are treated as a group.
// This may represent the body of a function, loop, or a conditional block.
type Block struct {
	token      token.Token // the opening "{" token
	statements []Node      // the statements in the block
}

func NewBlock(token token.Token, statements []Node) *Block {
	return &Block{token: token, statements: statements}
}

func (b *Block) StatementNode() {}

func (b *Block) Token() token.Token { return b.token }

func (b *Block) Literal() string { return b.token.Literal }

func (b *Block) Statements() []Node { return b.statements }

func (b *Block) String() string {
	var out bytes.Buffer
	for i, s := range b.statements {
		if i > 0 {
			out.WriteString("\n")
		}
		out.WriteString(s.String())
	}
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

// For defines a for loop.
type For struct {
	token token.Token

	// condition determines whether the loop should continue running.
	condition Node

	// consequence contains the statements that make up the loop body.
	consequence *Block

	// Initialization statement which is executed once before evaluating the
	// condition for the first iteration.
	init Node

	// Statement which is executed after each execution of the block
	// (and only if the block was executed).
	post Expression
}

func NewSimpleFor(token token.Token, consequence *Block) *For {
	return &For{token: token, consequence: consequence}
}

func NewFor(token token.Token, condition Node, consequence *Block, init Node, post Expression) *For {
	return &For{token: token, condition: condition, consequence: consequence, init: init, post: post}
}

// Should be StatementNode only?
func (f *For) ExpressionNode() {}

func (f *For) Token() token.Token { return f.token }

func (f *For) Literal() string { return f.token.Literal }

func (f *For) IsSimpleLoop() bool {
	return f.consequence != nil && f.init == nil && f.condition == nil
}

func (f *For) IsIteratorLoop() bool {
	if f.condition == nil {
		return false
	}
	switch f.condition.(type) {
	case *Var, *MultiVar:
		// The only case where var and multi-var assignments are supported are
		// when an iterator is being used to define the loop. The assignment AST
		// node currently comes through as the loop "condition" in the AST. For
		// loops that use iterators can look like these:
		//   for i := range x {}
		//   for i, j := range x {}
		//   for i, j := myiterator {}
		return true
	}
	return false
}

func (f *For) Condition() Node { return f.condition }

func (f *For) Consequence() *Block { return f.consequence }

func (f *For) Init() Node { return f.init }

func (f *For) Post() Expression { return f.post }

func (f *For) String() string {
	var out bytes.Buffer
	// Simple for {} loop
	if f.IsSimpleLoop() {
		out.WriteString("for { ")
		out.WriteString(f.consequence.String())
		out.WriteString(" }")
		return out.String()
	}
	// Full style for loop
	out.WriteString("for ")
	out.WriteString(f.init.String() + "; ")
	out.WriteString(f.condition.String() + "; ")
	out.WriteString(f.post.String())
	out.WriteString(" { ")
	out.WriteString(f.consequence.String())
	out.WriteString(" }")
	return out.String()
}

// Func holds a function definition.
type Func struct {
	token token.Token

	name *Ident

	// parameters is the list of parameters the function receives.
	parameters []*Ident

	// defaults holds any default values for arguments which aren't specified.
	defaults map[string]Expression

	// body contains the set of statements within the function.
	body *Block
}

func NewFunc(token token.Token, name *Ident, parameters []*Ident, defaults map[string]Expression, body *Block) *Func {
	return &Func{
		token:      token,
		name:       name,
		parameters: parameters,
		defaults:   defaults,
		body:       body,
	}
}

func (f *Func) ExpressionNode() {}

func (f *Func) Token() token.Token { return f.token }

func (f *Func) Literal() string { return f.token.Literal }

func (f *Func) Name() *Ident { return f.name }

func (f *Func) Parameters() []*Ident { return f.parameters }

func (f *Func) Defaults() map[string]Expression { return f.defaults }

func (f *Func) Body() *Block { return f.body }

func (f *Func) String() string {
	var out bytes.Buffer
	params := make([]string, 0)
	for _, p := range f.parameters {
		params = append(params, p.value)
	}
	out.WriteString(f.Literal())
	if f.name != nil {
		out.WriteString(" " + f.name.value)
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.body.String())
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

// String holds a string
type String struct {
	// Token is the token
	token token.Token

	// value is the value of the string.
	value string

	// template is the templatized version of the string, if any
	template *tmpl.Template

	exprs []Expression
}

func NewString(tok token.Token) *String {
	return &String{token: tok, value: tok.Literal}
}

func NewTemplatedString(tok token.Token, template *tmpl.Template, exprs []Expression) *String {
	return &String{token: tok, value: tok.Literal, template: template, exprs: exprs}
}

func (s *String) ExpressionNode() {}

func (s *String) Token() token.Token { return s.token }

func (s *String) Literal() string { return s.token.Literal }

func (s *String) Value() string { return s.value }

func (s *String) Template() *tmpl.Template { return s.template }

func (s *String) TemplateExpressions() []Expression { return s.exprs }

func (s *String) String() string { return fmt.Sprintf("%q", s.token.Literal) }

// List holds an inline list
type List struct {
	// Token is the token
	token token.Token

	// items holds the members of the list.
	items []Expression
}

func NewList(tok token.Token, items []Expression) *List {
	return &List{token: tok, items: items}
}

func (l *List) ExpressionNode() {}

func (l *List) Token() token.Token { return l.token }

func (l *List) Literal() string { return l.token.Literal }

func (l *List) Items() []Expression { return l.items }

func (l *List) String() string {
	var out bytes.Buffer
	elements := make([]string, 0)
	for _, el := range l.items {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
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

// Assign is generally used for a simple assignment like "x = y". We also
// support other operators like "+=", "-=", "*=", and "/=".
type Assign struct {
	token    token.Token
	name     *Ident
	index    *Index
	operator string
	value    Expression
}

func NewAssign(operator token.Token, name *Ident, value Expression) *Assign {
	return &Assign{token: operator, name: name, operator: operator.Literal, value: value}
}

func NewAssignIndex(operator token.Token, index *Index, value Expression) *Assign {
	return &Assign{token: operator, index: index, operator: operator.Literal, value: value}
}

func (a *Assign) ExpressionNode() {}

func (a *Assign) Token() token.Token { return a.token }

func (a *Assign) Literal() string { return a.token.Literal }

func (a *Assign) Name() string { return a.name.value }

func (a *Assign) Index() *Index { return a.index }

func (a *Assign) Operator() string { return a.operator }

func (a *Assign) Value() Expression { return a.value }

func (a *Assign) String() string {
	var out bytes.Buffer
	if a.index != nil {
		out.WriteString(a.index.String())
	} else {
		out.WriteString(a.name.value)
	}
	out.WriteString(" " + a.operator + " ")
	out.WriteString(a.value.String())
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

// Import holds an import statement
type Import struct {
	token token.Token // the "import" token
	name  *Ident      // name of the module to import
}

func NewImport(token token.Token, name *Ident) *Import {
	return &Import{token: token, name: name}
}

func (i *Import) ExpressionNode() {}

func (i *Import) Token() token.Token { return i.token }

func (i *Import) Literal() string { return i.token.Literal }

func (i *Import) Module() *Ident { return i.name }

func (i *Import) String() string {
	var out bytes.Buffer
	out.WriteString(i.Literal() + " ")
	out.WriteString(i.name.Literal())
	out.WriteString(";")
	return out.String()
}

// Map holds a map
type Map struct {
	token token.Token               // the '{' token
	items map[Expression]Expression // items in the map
}

func NewMap(token token.Token, items map[Expression]Expression) *Map {
	return &Map{token: token, items: items}
}

func (m *Map) ExpressionNode() {}

func (m *Map) Token() token.Token { return m.token }

func (m *Map) Literal() string { return m.token.Literal }

func (m *Map) Items() map[Expression]Expression { return m.items }

func (m *Map) String() string {
	var out bytes.Buffer
	pairs := make([]string, 0)
	for key, value := range m.items {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// Set holds a set definition
type Set struct {
	token token.Token  // the '{' token
	items []Expression // items in the set
}

func NewSet(token token.Token, items []Expression) *Set {
	return &Set{token: token, items: items}
}

func (s *Set) ExpressionNode() {}

func (s *Set) Token() token.Token { return s.token }

func (s *Set) Literal() string { return s.token.Literal }

func (s *Set) Items() []Expression { return s.items }

func (s *Set) String() string {
	var out bytes.Buffer
	items := make([]string, 0, len(s.items))
	for _, key := range s.items {
		items = append(items, key.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("}")
	return out.String()
}

// In holds an "in" expression
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

// Range is used to iterator over a container
type Range struct {
	token     token.Token // the "range" token
	container Expression  // the container to iterate over
}

func NewRange(token token.Token, container Expression) *Range {
	return &Range{token: token, container: container}
}

func (r *Range) ExpressionNode() {}

func (r *Range) Token() token.Token { return r.token }

func (r *Range) Literal() string { return r.token.Literal }

func (r *Range) Container() Expression { return r.container }

func (r *Range) String() string {
	var out bytes.Buffer
	out.WriteString(r.Literal())
	out.WriteString(" ")
	out.WriteString(r.container.String())
	return out.String()
}
