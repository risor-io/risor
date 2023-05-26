package ast

import (
	"bytes"
	"strings"

	"github.com/cloudcmds/tamarin/token"
)

// Var is a statement that assigns a value to a variable. It may be a
// declaration or an assignment. If it's a declaration, isWalrus is true.
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

func (s *Var) IsExpression() bool { return false }

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

// MultiVar is a statement that assigns values to more than one variable.
// The right hand side must be a container type, with the same number of
// elements as the number of variables on the left hand side.
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

func (s *MultiVar) IsExpression() bool { return false }

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

func (c *Const) IsExpression() bool { return false }

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

// Control defines a return, break, or continue statement.
type Control struct {
	token token.Token // "return", "break", or "continue"
	value Expression  // optional value, for return statements
}

func NewControl(token token.Token, value Expression) *Control {
	return &Control{token: token, value: value}
}

func (c *Control) StatementNode() {}

func (c *Control) IsExpression() bool { return false }

func (c *Control) Token() token.Token { return c.token }

func (c *Control) Literal() string { return c.token.Literal }

func (c *Control) Value() Expression { return c.value }

func (c *Control) IsReturn() bool {
	return c.token.Type == token.RETURN
}

func (c *Control) String() string {
	var out bytes.Buffer
	out.WriteString(c.Literal() + " ")
	if c.value != nil {
		out.WriteString(c.value.Literal())
	}
	out.WriteString(";")
	return out.String()
}

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

func (b *Block) IsExpression() bool { return false }

func (b *Block) Token() token.Token { return b.token }

func (b *Block) Literal() string { return b.token.Literal }

func (b *Block) Statements() []Node { return b.statements }

func (b *Block) EndsWithReturn() bool {
	count := len(b.statements)
	if count == 0 {
		return false
	}
	last := b.statements[count-1]
	if cntrl, ok := last.(*Control); ok {
		return cntrl.IsReturn()
	}
	return false
}

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
	post Node
}

func NewSimpleFor(token token.Token, consequence *Block) *For {
	return &For{token: token, consequence: consequence}
}

func NewFor(token token.Token, condition Node, consequence *Block, init Node, post Node) *For {
	return &For{token: token, condition: condition, consequence: consequence, init: init, post: post}
}

func (f *For) StatementNode() {}

func (f *For) IsExpression() bool { return false }

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

func (f *For) Post() Node { return f.post }

func (f *For) String() string {
	var out bytes.Buffer
	// Simple for {} loop
	if f.IsSimpleLoop() {
		out.WriteString("for { ")
		out.WriteString(f.consequence.String())
		out.WriteString(" }")
		return out.String()
	}
	if f.init == nil {
		out.WriteString("for ")
		out.WriteString(f.condition.String())
		out.WriteString(" { ")
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

// Assign is generally used for a simple assignment like "x = y". We also
// support other operators like "+=", "-=", "*=", and "/=".
type Assign struct {
	token    token.Token
	name     *Ident // this may be nil, e.g. `[0, 1, 2][0] = 3`
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

func (a *Assign) StatementNode() {}

func (a *Assign) IsExpression() bool { return false }

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

// Import holds an import statement
type Import struct {
	token token.Token // the "import" token
	name  *Ident      // name of the module to import
}

func NewImport(token token.Token, name *Ident) *Import {
	return &Import{token: token, name: name}
}

func (i *Import) StatementNode() {}

func (i *Import) IsExpression() bool { return false }

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

// Postfix defines a postfix expression like "x++".
type Postfix struct {
	token token.Token
	// operator holds the postfix token, e.g. ++
	operator string
}

func NewPostfix(token token.Token, operator string) *Postfix {
	return &Postfix{token: token, operator: operator}
}

func (p *Postfix) StatementNode() {}

func (p *Postfix) IsExpression() bool { return false }

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
