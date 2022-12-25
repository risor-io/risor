// Package ast contains the definitions of the abstract-syntax tree
// that our parse produces, and our interpreter executes.
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

	// StartToken returns the token where this node begins.
	StartToken() token.Token

	// TokenLiteral returns the literal of the token.
	TokenLiteral() string

	// String returns this object as a string.
	String() string
}

// Statement represents a single statement.
type Statement interface {
	// Node is the node holding the actual statement
	Node

	statementNode()
}

// Expression represents a single expression.
type Expression interface {
	// Node is the node holding the expression.
	Node
	expressionNode()
}

// Program represents a complete program.
type Program struct {
	// Statements is the set of statements which the program is comprised
	// of.
	Statements []Statement
}

func (p *Program) StartToken() token.Token {
	if len(p.Statements) > 0 {
		return p.Statements[0].StartToken()
	}
	return token.Token{}
}

// TokenLiteral returns the literal token of our program.
func (p *Program) TokenLiteral() string {
	return p.StartToken().Literal
}

// String returns this object as a string.
func (p *Program) String() string {
	var out bytes.Buffer
	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
	}
	return out.String()
}

// VarStatement holds a var-statemnt
type VarStatement struct {
	// Token holds the token
	Token token.Token

	// Name is the name of the variable to which we're assigning
	Name *Identifier

	// Value is the thing we're storing in the variable.
	Value Expression

	IsWalrus bool
}

func (s *VarStatement) statementNode() {}

func (s *VarStatement) StartToken() token.Token {
	return s.Token
}

func (s *VarStatement) TokenLiteral() string {
	return s.Token.Literal
}

// String returns this object as a string.
func (s *VarStatement) String() string {
	var out bytes.Buffer
	if s.IsWalrus {
		out.WriteString(s.Name.TokenLiteral() + " := ")
		// out.WriteString(s.TokenLiteral() + " ")
		out.WriteString(s.Value.String())
		return out.String()
	}
	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Name.TokenLiteral())
	out.WriteString(" = ")
	if s.Value != nil {
		out.WriteString(s.Value.String())
	}
	return out.String()
}

// ConstStatement is the same as var-statement, but the value
// can't be changed later.
type ConstStatement struct {
	// Token is the token
	Token token.Token

	// Name is the name of the variable we're setting
	Name *Identifier

	// Value contains the value which is to be set
	Value Expression
}

func (cs *ConstStatement) statementNode() {}

func (cs *ConstStatement) StartToken() token.Token {
	return cs.Token
}

func (cs *ConstStatement) TokenLiteral() string {
	return cs.Token.Literal
}

// String returns this object as a string.
func (cs *ConstStatement) String() string {
	var out bytes.Buffer
	out.WriteString(cs.TokenLiteral() + " ")
	out.WriteString(cs.Name.TokenLiteral())
	out.WriteString(" = ")
	if cs.Value != nil {
		out.WriteString(cs.Value.String())
	}
	return out.String()
}

// Identifier holds a single identifier.
type Identifier struct {
	// Token is the literal token
	Token token.Token

	// Value is the name of the identifier
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) StartToken() token.Token {
	return i.Token
}

// TokenLiteral returns the literal token.
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// String returns this object as a string.
func (i *Identifier) String() string {
	return i.Value
}

// ReturnStatement stores a return-statement
type ReturnStatement struct {
	// Token contains the literal token.
	Token token.Token

	// ReturnValue is the value whichis to be returned.
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) StartToken() token.Token { return rs.Token }

// TokenLiteral returns the literal token.
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String returns this object as a string.
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.TokenLiteral())
	}
	out.WriteString(";")
	return out.String()
}

// BreakStatement stores a break statement
type BreakStatement struct {
	Token token.Token
}

func (s *BreakStatement) statementNode() {}

func (s *BreakStatement) StartToken() token.Token { return s.Token }

func (s *BreakStatement) TokenLiteral() string { return s.Token.Literal }

func (s *BreakStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

// ExpressionStatement is an expression
type ExpressionStatement struct {
	// Token is the literal token
	Token token.Token

	// Expression holds the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) StartToken() token.Token { return es.Token }

// TokenLiteral returns the literal token.
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// String returns this object as a string.
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral holds an integer
type IntegerLiteral struct {
	// Token is the literal token
	Token token.Token

	// Value holds the integer.
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) StartToken() token.Token { return il.Token }

// TokenLiteral returns the literal token.
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// String returns this object as a string.
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// FloatLiteral holds a floating-point number
type FloatLiteral struct {
	// Token is the literal token
	Token token.Token

	// Value holds the floating-point number.
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}

func (fl *FloatLiteral) StartToken() token.Token { return fl.Token }

// TokenLiteral returns the literal token.
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }

// String returns this object as a string.
func (fl *FloatLiteral) String() string { return fl.Token.Literal }

// PrefixExpression holds a prefix-based expression
type PrefixExpression struct {
	// Token holds the token.  e.g. "!"
	Token token.Token

	// Operator holds the operator being invoked (e.g. "!" ).
	Operator string

	// Right holds the thing to be operated upon
	Right Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) StartToken() token.Token { return pe.Token }

// TokenLiteral returns the literal token.
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// String returns this object as a string.
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression stores an infix expression.
type InfixExpression struct {
	// Token holds the literal expression
	Token token.Token

	// Left holds the left-most argument
	Left Expression

	// Operator holds the operation to be carried out (e.g. "+", "-" )
	Operator string

	// Right holds the right-most argument
	Right Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) StartToken() token.Token { return ie.Token }

// TokenLiteral returns the literal token.
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns this object as a string.
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// PostfixExpression holds a postfix-based expression
type PostfixExpression struct {
	// Token holds the token we're operating upon
	Token token.Token
	// Operator holds the postfix token, e.g. ++
	Operator string
}

func (pe *PostfixExpression) expressionNode() {}

func (pe *PostfixExpression) StartToken() token.Token { return pe.Token }

func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Token.Literal)
	out.WriteString(pe.Operator)
	out.WriteString(")")
	return out.String()
}

// NilLiteral represents a literal nil
type NilLiteral struct {
	// Token holds the actual token
	Token token.Token
}

func (n *NilLiteral) expressionNode() {}

func (n *NilLiteral) StartToken() token.Token { return n.Token }

func (n *NilLiteral) TokenLiteral() string { return n.Token.Literal }

func (n *NilLiteral) String() string { return n.Token.Literal }

// Bool holds a boolean type
type Bool struct {
	// Token holds the actual token
	Token token.Token

	// Value stores the bools' value: true, or false.
	Value bool
}

func (b *Bool) expressionNode() {}

func (b *Bool) StartToken() token.Token { return b.Token }

func (b *Bool) TokenLiteral() string { return b.Token.Literal }

func (b *Bool) String() string { return b.Token.Literal }

// BlockStatement holds a group of statements, which are treated
// as a block.  (For example the body of an `if` expression.)
type BlockStatement struct {
	// Token holds the actual token
	Token token.Token

	// Statements contain the set of statements within the block
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) StartToken() token.Token { return bs.Token }

func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for i, s := range bs.Statements {
		if i > 0 {
			out.WriteString("\n")
		}
		out.WriteString(s.String())
	}
	return out.String()
}

// IfExpression holds an if-statement
type IfExpression struct {
	// Token is the actual token
	Token token.Token

	// Condition is the thing that is evaluated to determine
	// which block should be executed.
	Condition Expression

	// Consequence is the set of statements executed if the
	// condition is true.
	Consequence *BlockStatement

	// Alternative is the set of statements executed if the
	// condition is not true (optional).
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) StartToken() token.Token { return ie.Token }

func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// ForeachStatement holds a foreach-statement.
type ForeachStatement struct {
	// Token is the actual token
	Token token.Token

	// Index is the variable we'll set with the index, for the blocks' scope
	//
	// This is optional.
	Index string

	// Ident is the variable we'll set with each item, for the blocks' scope
	Ident string

	// Value is the thing we'll range over.
	Value Expression

	// Body is the block we'll execute.
	Body *BlockStatement
}

func (fes *ForeachStatement) expressionNode() {}

func (fes *ForeachStatement) StartToken() token.Token { return fes.Token }

func (fes *ForeachStatement) TokenLiteral() string { return fes.Token.Literal }

func (fes *ForeachStatement) String() string {
	var out bytes.Buffer
	out.WriteString("foreach ")
	out.WriteString(fes.Ident)
	out.WriteString(" ")
	out.WriteString(fes.Value.String())
	out.WriteString(fes.Body.String())
	return out.String()
}

// TernaryExpression holds a ternary-expression.
type TernaryExpression struct {
	// Token is the actual token.
	Token token.Token

	// Condition is the thing that is evaluated to determine
	// which expression should be returned
	Condition Expression

	// IfTrue is the expression to return if the condition is true.
	IfTrue Expression

	// IFFalse is the expression to return if the condition is not true.
	IfFalse Expression
}

func (te *TernaryExpression) expressionNode() {}

func (te *TernaryExpression) StartToken() token.Token { return te.Token }

func (te *TernaryExpression) TokenLiteral() string { return te.Token.Literal }

func (te *TernaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(te.Condition.String())
	out.WriteString(" ? ")
	out.WriteString(te.IfTrue.String())
	out.WriteString(" : ")
	out.WriteString(te.IfFalse.String())
	out.WriteString(")")
	return out.String()
}

// ForLoopExpression holds a for-loop
type ForLoopExpression struct {
	// Token is the actual token
	Token token.Token

	// Condition is the expression used to determine if the loop
	// is still running.
	Condition Expression

	// Consequence is the set of statements to be executed for the
	// loop body.
	Consequence *BlockStatement

	// Initialization statement which is executed once before evaluating the
	// condition for the first iteration
	InitStatement *VarStatement

	// Statement which is executed after each execution of the block
	// (and only if the block was executed)
	PostStatement Expression
}

func (fle *ForLoopExpression) expressionNode() {}

func (fle *ForLoopExpression) StartToken() token.Token { return fle.Token }

func (fle *ForLoopExpression) TokenLiteral() string { return fle.Token.Literal }

func (fle *ForLoopExpression) IsSimpleLoop() bool {
	return fle.Consequence != nil && fle.InitStatement == nil
}

func (fle *ForLoopExpression) String() string {
	var out bytes.Buffer
	// Simple for {} loop
	if fle.IsSimpleLoop() {
		out.WriteString("for { ")
		out.WriteString(fle.Consequence.String())
		out.WriteString(" }")
		return out.String()
	}
	// Full style for loop
	out.WriteString("for ")
	out.WriteString(fle.InitStatement.String() + "; ")
	out.WriteString(fle.Condition.String() + "; ")
	out.WriteString(fle.PostStatement.String())
	out.WriteString(" { ")
	out.WriteString(fle.Consequence.String())
	out.WriteString(" }")
	return out.String()
}

// FunctionLiteral holds a function-definition
//
// See-also FunctionDefineLiteral.
type FunctionLiteral struct {
	// Token is the actual token
	Token token.Token

	// Parameters is the list of parameters the function receives.
	Parameters []*Identifier

	// Defaults holds any default values for arguments which aren't
	// specified
	Defaults map[string]Expression

	// Body contains the set of statements within the function.
	Body *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) StartToken() token.Token { return fl.Token }

func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := make([]string, 0)
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()

}

// FunctionDefineLiteral holds a function-definition.
//
// See-also FunctionLiteral.
type FunctionDefineLiteral struct {
	// Token holds the token
	Token token.Token

	// Paremeters holds the function parameters.
	Parameters []*Identifier

	// Defaults holds any default-arguments.
	Defaults map[string]Expression

	// Body holds the set of statements in the functions' body.
	Body *BlockStatement
}

func (fl *FunctionDefineLiteral) expressionNode() {}

func (fl *FunctionDefineLiteral) StartToken() token.Token { return fl.Token }

func (fl *FunctionDefineLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionDefineLiteral) String() string {
	var out bytes.Buffer
	params := make([]string, 0)
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()

}

// CallExpression holds the invokation of a method-call.
type CallExpression struct {
	// Token stores the literal token
	Token token.Token

	// Function is the function to be invoked.
	Function Expression

	// Arguments are the arguments to be applied
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) StartToken() token.Token { return ce.Token }

func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := make([]string, 0)
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// GetAttributeExpression
type GetAttributeExpression struct {
	// Token stores the literal token
	Token token.Token
	// Object whose attribute is being accessed
	Object Expression
	// The attribute itself
	Attribute *Identifier
}

func (e *GetAttributeExpression) expressionNode() {}

func (e *GetAttributeExpression) StartToken() token.Token { return e.Token }

func (e *GetAttributeExpression) TokenLiteral() string { return e.Token.Literal }

func (e *GetAttributeExpression) String() string {
	var out bytes.Buffer
	out.WriteString(e.Object.String())
	out.WriteString(".")
	out.WriteString(e.Attribute.String())
	return out.String()
}

// PipeExpression holds a series of calls
type PipeExpression struct {
	// Token stores the literal token
	Token token.Token

	// Arguments are the arguments to be applied
	Arguments []Expression
}

func (pe *PipeExpression) expressionNode() {}

func (pe *PipeExpression) StartToken() token.Token { return pe.Token }

func (pe *PipeExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PipeExpression) String() string {
	var out bytes.Buffer
	args := make([]string, 0)
	for _, a := range pe.Arguments {
		args = append(args, a.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(args, " | "))
	out.WriteString(")")
	return out.String()
}

// ObjectCallExpression is used when calling a method on an object.
type ObjectCallExpression struct {
	// Token is the literal token
	Token token.Token

	// Object is the object against which the call is invoked.
	Object Expression

	// Call is the method-name.
	Call Expression
}

func (oce *ObjectCallExpression) expressionNode() {}

func (oce *ObjectCallExpression) StartToken() token.Token { return oce.Token }

func (oce *ObjectCallExpression) TokenLiteral() string {
	return oce.Token.Literal
}

func (oce *ObjectCallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(oce.Object.String())
	out.WriteString(".")
	out.WriteString(oce.Call.String())
	return out.String()
}

// StringLiteral holds a string
type StringLiteral struct {
	// Token is the token
	Token token.Token

	// Value is the value of the string.
	Value string

	// Template is the templatized version of the string, if any
	Template *tmpl.Template

	TemplateExpressions []*ExpressionStatement
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) StartToken() token.Token { return sl.Token }

func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

func (sl *StringLiteral) String() string { return fmt.Sprintf("%q", sl.Token.Literal) }

// ListLiteral holds an inline list
type ListLiteral struct {
	// Token is the token
	Token token.Token

	// Items holds the members of the list.
	Items []Expression
}

func (ll *ListLiteral) expressionNode() {}

func (ll *ListLiteral) StartToken() token.Token { return ll.Token }

func (ll *ListLiteral) TokenLiteral() string { return ll.Token.Literal }

func (ll *ListLiteral) String() string {
	var out bytes.Buffer
	elements := make([]string, 0)
	for _, el := range ll.Items {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// IndexExpression holds an index-expression
type IndexExpression struct {
	// Token is the actual token
	Token token.Token

	// Left is the thing being indexed.
	Left Expression

	// Index is the value we're indexing
	Index Expression

	// Optional "from" index for [from:to] style expressions
	FromIndex Expression

	// Optional "to" index for [from:to] style expressions
	ToIndex Expression
}

func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) StartToken() token.Token { return ie.Token }

func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	if ie.Index != nil {
		out.WriteString(ie.Index.String())
		out.WriteString("])")
		return out.String()
	}
	if ie.FromIndex != nil {
		out.WriteString(ie.FromIndex.String())
	}
	if ie.ToIndex != nil {
		out.WriteString(":")
		out.WriteString(ie.ToIndex.String())
	}
	out.WriteString("])")
	return out.String()
}

// AssignStatement is generally used for a (var-less) assignment,
// such as "x = y", however we allow an operator to be stored ("=" in that
// example), such that we can do self-operations.
//
// Specifically "x += y" is defined as an assignment-statement with
// the operator set to "+=".  The same applies for "+=", "-=", "*=", and
// "/=".
type AssignStatement struct {
	Token    token.Token
	Name     *Identifier
	Index    *IndexExpression
	Operator string
	Value    Expression
}

func (as *AssignStatement) expressionNode() {}

func (as *AssignStatement) StartToken() token.Token { return as.Token }

func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }

func (as *AssignStatement) String() string {
	var out bytes.Buffer
	if as.Index != nil {
		out.WriteString(as.Index.String())
	} else {
		out.WriteString(as.Name.String())
	}
	out.WriteString(" " + as.Operator + " ")
	out.WriteString(as.Value.String())
	return out.String()
}

// CaseExpression handles the case within a switch statement
type CaseExpression struct {
	// Token is the actual token
	Token token.Token

	// Default branch?
	Default bool

	// The thing we match
	Expr []Expression

	// The code to execute if there is a match
	Block *BlockStatement
}

func (ce *CaseExpression) expressionNode() {}

func (ce *CaseExpression) StartToken() token.Token { return ce.Token }

func (ce *CaseExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *CaseExpression) String() string {
	var out bytes.Buffer
	if ce.Default {
		out.WriteString("default")
	} else {
		out.WriteString("case ")
		tmp := []string{}
		for _, exp := range ce.Expr {
			tmp = append(tmp, exp.String())
		}
		out.WriteString(strings.Join(tmp, ","))
	}
	out.WriteString(":\n")
	for i, exp := range ce.Block.Statements {
		if i > 0 {
			out.WriteString("\n")
		}
		out.WriteString("\t" + exp.String())
	}
	out.WriteString("\n")
	return out.String()
}

// SwitchExpression handles a switch statement
type SwitchExpression struct {
	// Token is the actual token
	Token token.Token

	// Value is the thing that is evaluated to determine
	// which block should be executed.
	Value Expression

	// The branches we handle
	Choices []*CaseExpression
}

func (se *SwitchExpression) expressionNode() {}

func (se *SwitchExpression) StartToken() token.Token { return se.Token }

func (se *SwitchExpression) TokenLiteral() string { return se.Token.Literal }

func (se *SwitchExpression) String() string {
	var out bytes.Buffer
	out.WriteString("\nswitch ")
	out.WriteString(se.Value.String())
	out.WriteString(" {\n")
	for _, tmp := range se.Choices {
		if tmp != nil {
			out.WriteString(tmp.String())
		}
	}
	out.WriteString("}\n")
	return out.String()
}

// ImportStatement holds an import statement
type ImportStatement struct {
	// Token holds the token
	Token token.Token

	// Name of the module to import
	Name *Identifier
}

func (i *ImportStatement) expressionNode() {}

func (i *ImportStatement) StartToken() token.Token { return i.Token }

func (i *ImportStatement) TokenLiteral() string { return i.Token.Literal }

func (i *ImportStatement) String() string {
	var out bytes.Buffer
	out.WriteString(i.TokenLiteral() + " ")
	out.WriteString(i.Name.TokenLiteral())
	out.WriteString(";")
	return out.String()
}

// MapLiteral holds a map
type MapLiteral struct {
	// Token holds the token
	Token token.Token // the '{' token
	// Pairs stores the name/value sets of the hash-content
	Pairs map[Expression]Expression
}

func (ml *MapLiteral) expressionNode() {}

func (ml *MapLiteral) StartToken() token.Token { return ml.Token }

func (ml *MapLiteral) TokenLiteral() string { return ml.Token.Literal }

func (ml *MapLiteral) String() string {
	var out bytes.Buffer
	pairs := make([]string, 0)
	for key, value := range ml.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// SetLiteral holds a set definition
type SetLiteral struct {
	Token token.Token // the '{' token
	Items []Expression
}

func (sl *SetLiteral) expressionNode() {}

func (sl *SetLiteral) StartToken() token.Token { return sl.Token }

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
