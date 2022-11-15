// Package parser is used to parse an input program from its tokens and produce
// an abstract syntax tree (AST) as output.
//
// A parser is created by calling New() with a lexer as input. The parser should
// then be used only once, by calling parser.Parse() to produce the AST.
package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudcmds/tamarin/internal/ast"
	"github.com/cloudcmds/tamarin/internal/lexer"
	"github.com/cloudcmds/tamarin/internal/token"
	"github.com/hashicorp/go-multierror"
)

type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func() ast.Expression
)

// Parse is a shortcut that can be used to parse the given Tamarin source code.
// The lexer and parser are created internally and not exposed.
func Parse(input string) (*ast.Program, error) {
	return New(lexer.New(input)).Parse()
}

// Parser object
type Parser struct {
	// l is our lexer
	l *lexer.Lexer

	// prevToken holds the previous token from our lexer.
	// (used for "++" + "--")
	prevToken token.Token

	// curToken holds the current token from our lexer.
	curToken token.Token

	// peekToken holds the next token which will come from the lexer.
	peekToken token.Token

	// errors holds parsing-errors.
	errors []string

	// prefixParseFns holds a map of parsing methods for
	// prefix-based syntax.
	prefixParseFns map[token.Type]prefixParseFn

	// infixParseFns holds a map of parsing methods for
	// infix-based syntax.
	infixParseFns map[token.Type]infixParseFn

	// postfixParseFns holds a map of parsing methods for
	// postfix-based syntax.
	postfixParseFns map[token.Type]postfixParseFn

	// are we inside a ternary expression?
	//
	// Nested ternary expressions are illegal :)
	tern bool
}

// New returns a Parser for the program provided by the lexer.
func New(l *lexer.Lexer) *Parser {

	// Create the parser and prime the token pump
	p := &Parser{
		l:               l,
		errors:          []string{},
		prefixParseFns:  map[token.Type]prefixParseFn{},
		infixParseFns:   map[token.Type]infixParseFn{},
		postfixParseFns: map[token.Type]postfixParseFn{},
	}
	p.nextToken() // loads peekToken with token0
	p.nextToken() // loads curToken with token0 and peekToken with token1

	// Register prefix-functions
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.FUNC, p.parseFunctionDefinition)
	p.registerPrefix(token.EOF, p.parsingBroken)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.FOR, p.parseForLoopExpression)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.ILLEGAL, p.parsingBroken)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.NULL, p.parseNull)
	p.registerPrefix(token.REGEXP, p.parseRegexpLiteral)
	p.registerPrefix(token.REGEXP, p.parseRegexpLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.SWITCH, p.parseSwitchStatement)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.PIPE, p.parsePrefixExpression)
	p.registerPrefix(token.NEWLINE, p.parseNewlineLiteral)
	p.registerPrefix(token.IMPORT, p.parseImportStatement)

	// Register infix functions
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignExpression)
	p.registerInfix(token.DECLARE, p.parseAssignExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK_EQUALS, p.parseAssignExpression)
	p.registerInfix(token.CONTAINS, p.parseInfixExpression)
	p.registerInfix(token.DOTDOT, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GT_EQUALS, p.parseInfixExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LT_EQUALS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS_EQUALS, p.parseAssignExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.NOT_CONTAINS, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.PERIOD, p.parseMethodCallExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.PLUS_EQUALS, p.parseAssignExpression)
	p.registerInfix(token.POW, p.parseInfixExpression)
	p.registerInfix(token.QUESTION, p.parseTernaryExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.SLASH_EQUALS, p.parseAssignExpression)
	p.registerInfix(token.PIPE, p.parsePipeExpression)

	// Register postfix functions
	p.registerPostfix(token.MINUS_MINUS, p.parsePostfixExpression)
	p.registerPostfix(token.PLUS_PLUS, p.parsePostfixExpression)
	return p
}

// nextToken moves to the next token from the lexer, updating all of
// prevToken, curToken, and peekToken.
func (p *Parser) nextToken() {
	p.prevToken = p.curToken
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Parse the program that is provided via the lexer.
func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{Statements: []ast.Statement{}}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	var combinedErrors error
	for _, err := range p.Errors() {
		combinedErrors = multierror.Append(combinedErrors, errors.New(err))
	}
	return program, combinedErrors
}

// registerPrefix registers a function for handling a prefix-based statement.
func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix registers a function for handling an infix-based statement.
func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// registerPostfix registers a function for handling a postfix-based statement.
func (p *Parser) registerPostfix(tokenType token.Type, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found around line %d",
		t, p.curToken.Line+1)
	p.errors = append(p.errors, msg)
}

// Errors returns all error messages accumulated during program parsing.
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError raises an error if the next token is not the expected type.
func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead around line %d:%d",
		t, p.peekToken.Type, p.curToken.Line+1, p.curToken.EndPosition+2)
	p.errors = append(p.errors, msg)
}

// parseStatement parses a single statement.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.CONST:
		return p.parseConstStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.NEWLINE:
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement parses a let statement.
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.NEWLINE) {
		if p.curTokenIs(token.EOF) {
			p.errors = append(p.errors, "unterminated let statement")
			return nil
		}
		p.nextToken()
	}
	return stmt
}

// parseConstStatement parses a constant declaration.
func (p *Parser) parseConstStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.NEWLINE) {
		if p.curTokenIs(token.EOF) {
			p.errors = append(p.errors, "unterminated const statement")
			return nil
		}
		p.nextToken()
	}
	return stmt
}

// parseReturnStatement parses a function return statement.
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)
	for {
		switch p.peekToken.Type {
		case token.SEMICOLON, token.NEWLINE, token.EOF:
			p.nextToken()
			return stmt
		default:
			p.errors = append(p.errors, fmt.Sprintf("unexpected token in return statement: %s", p.peekToken.Literal))
			return nil
		}
	}
}

// parse Expression Statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	for p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.NEWLINE) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	postfix := p.postfixParseFns[p.curToken.Type]
	if postfix != nil {
		return (postfix())
	}
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	p.eatNewlines()
	return leftExp
}

// parsingBroken is hit if we see an EOF in our input-stream
// this means we're screwed
func (p *Parser) parsingBroken() ast.Expression {
	return nil
}

// parseIdentifier parses an identifier.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral parses an integer literal.
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	var value int64
	var err error

	if strings.HasPrefix(p.curToken.Literal, "0b") {
		value, err = strconv.ParseInt(p.curToken.Literal[2:], 2, 64)
	} else if strings.HasPrefix(p.curToken.Literal, "0x") {
		value, err = strconv.ParseInt(p.curToken.Literal[2:], 16, 64)
	} else {
		value, err = strconv.ParseInt(p.curToken.Literal, 10, 64)
	}

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer around line %d",
			p.curToken.Literal, p.curToken.Line+1)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

// parseFloatLiteral parses a float-literal
func (p *Parser) parseFloatLiteral() ast.Expression {
	flo := &ast.FloatLiteral{Token: p.curToken}
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float around line %d",
			p.curToken.Literal, p.curToken.Line+1)
		p.errors = append(p.errors, msg)
		return nil
	}
	flo.Value = value
	return flo
}

// parseSwitchStatement handles a switch statement
func (p *Parser) parseSwitchStatement() ast.Expression {
	expression := &ast.SwitchExpression{Token: p.curToken}
	p.nextToken()
	expression.Value = p.parseExpression(LOWEST)
	if expression.Value == nil {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	p.nextToken()
	p.eatNewlines()

	// Process the block which we think will contain various case-statements
	for !p.curTokenIs(token.RBRACE) {
		if p.curTokenIs(token.EOF) {
			p.errors = append(p.errors, "unterminated switch statement")
			return nil
		}
		if p.curToken.Literal != "case" && p.curToken.Literal != "default" {
			p.errors = append(p.errors, fmt.Sprintf("unexpected token %s; expected case or default",
				p.curToken.Literal))
			return nil
		}
		tmp := &ast.CaseExpression{Token: p.curToken}
		// Default will be handled specially
		if p.curTokenIs(token.DEFAULT) {
			tmp.Default = true
		} else if p.curTokenIs(token.CASE) {
			// skip "case"
			p.nextToken()
			// parse the match-expression
			tmp.Expr = append(tmp.Expr, p.parseExpression(LOWEST))
			for p.peekTokenIs(token.COMMA) {
				// skip the comma
				p.nextToken()
				// setup the expression
				p.nextToken()
				tmp.Expr = append(tmp.Expr, p.parseExpression(LOWEST))
			}
		} else {
			// error - unexpected token
			p.errors = append(p.errors, fmt.Sprintf("expected case|default, got %s", p.curToken.Type))
			return nil
		}
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		p.eatNewlines()
		// parse the block and save the choice
		block := &ast.BlockStatement{Token: p.curToken}
		for {
			stmt := p.parseStatement()
			if stmt == nil {
				return nil
			}
			block.Statements = append(block.Statements, stmt)
			p.eatNewlines()
			if p.curTokenIs(token.CASE) || p.curTokenIs(token.DEFAULT) || p.curTokenIs(token.RBRACE) {
				break
			}
		}
		tmp.Block = block
		expression.Choices = append(expression.Choices, tmp)
	}
	// More than one default is a bug
	count := 0
	for _, c := range expression.Choices {
		if c.Default {
			count++
		}
	}
	if count > 1 {
		msg := "A switch-statement should only have one default block"
		p.errors = append(p.errors, msg)
		return nil
	}
	return expression
}

func (p *Parser) parseImportStatement() ast.Expression {
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	return &ast.ImportStatement{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
}

// parseBoolean parses a boolean token.
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// parseNull parses a null keyword
func (p *Parser) parseNull() ast.Expression {
	return &ast.NullLiteral{Token: p.curToken}
}

// parsePrefixExpression parses a prefix-based expression.
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

// parseNewlineLiteral parses a prefix-based expression.
func (p *Parser) parseNewlineLiteral() ast.Expression {
	p.nextToken()
	return nil
}

// parsePostfixExpression parses a postfix-based expression.
func (p *Parser) parsePostfixExpression() ast.Expression {
	expression := &ast.PostfixExpression{
		Token:    p.prevToken,
		Operator: p.curToken.Literal,
	}
	return expression
}

// parseInfixExpression parses an infix-based expression.
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

// parseTernaryExpression parses a ternary expression
func (p *Parser) parseTernaryExpression(condition ast.Expression) ast.Expression {
	if p.tern {
		msg := fmt.Sprintf("nested ternary expressions are illegal, around line %d", p.curToken.Line+1)
		p.errors = append(p.errors, msg)
		return nil
	}
	p.tern = true
	defer func() { p.tern = false }()

	expression := &ast.TernaryExpression{
		Token:     p.curToken,
		Condition: condition,
	}
	p.nextToken() //skip the '?'
	precedence := p.curPrecedence()
	expression.IfTrue = p.parseExpression(precedence)

	if !p.expectPeek(token.COLON) { //skip the ":"
		return nil
	}

	// Get to next token, then parse the else part
	p.nextToken()
	expression.IfFalse = p.parseExpression(precedence)

	p.tern = false
	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	p.nextToken()
	// Look for the condition expression
	expression.Condition = p.parseExpression(LOWEST)
	if expression.Condition == nil {
		return nil
	}
	// Now "{"
	if !p.expectPeek(token.LBRACE) {
		msg := fmt.Sprintf("expected '{' but got %s", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	// The consequence
	expression.Consequence = p.parseBlockStatement()
	if expression.Consequence == nil {
		p.errors = append(p.errors, "unexpected nil expression")
		return nil
	}
	// Optional else block
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		// else if
		if p.peekTokenIs(token.IF) {
			p.nextToken()
			expression.Alternative = &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: p.parseIfExpression(),
					},
				},
			}
			return expression
		}
		// else { block }
		if !p.expectPeek(token.LBRACE) {
			msg := fmt.Sprintf("expected '{' but got %s", p.curToken.Literal)
			p.errors = append(p.errors, msg)
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
		if expression.Alternative == nil {
			p.errors = append(p.errors, "unexpected nil expression")
			return nil
		}
	}
	return expression
}

// parseForLoopExpression parses a for-loop.
func (p *Parser) parseForLoopExpression() ast.Expression {
	expression := &ast.ForLoopExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()

	if p.curToken.Type == token.LET {
		expression.InitStatement = p.parseLetStatement()
		if !p.curTokenIs(token.SEMICOLON) {
			p.errors = append(p.errors, "expected a semicolon")
			return nil
		}
		p.nextToken()
	}

	expression.Condition = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		p.nextToken()
		if !p.curTokenIs(token.IDENT) {
			p.errors = append(p.errors, "expected an identifier")
			return nil
		}
		expr := p.parseIdentifier()
		p.nextToken()
		if fn, ok := p.postfixParseFns[p.curToken.Type]; ok {
			expr = fn()
		}
		expression.PostStatement = expr
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()
	return expression
}

// parseBlockStatement parsea a block.
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) {
		if p.curTokenIs(token.EOF) {
			p.errors = append(p.errors, "unterminated block statement")
			return nil
		}
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// parseFunctionLiteral parses a function-literal.
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Defaults, lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

// parseFunctionDefinition parses the definition of a function.
func (p *Parser) parseFunctionDefinition() ast.Expression {
	if p.peekTokenIs(token.LPAREN) {
		return p.parseFunctionLiteral()
	}
	p.nextToken()
	lit := &ast.FunctionDefineLiteral{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Defaults, lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

// parseFunctionParameters parses the parameters used for a function.
func (p *Parser) parseFunctionParameters() (map[string]ast.Expression, []*ast.Identifier) {
	// Any default parameters
	m := make(map[string]ast.Expression)
	// The argument-definitions
	identifiers := make([]*ast.Identifier, 0)
	// Is the next parameter ")" ?  If so we're done. No args.
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return m, identifiers
	}
	p.nextToken()
	// Keep going until we find a ")"
	for !p.curTokenIs(token.RPAREN) {
		if p.curTokenIs(token.EOF) {
			p.errors = append(p.errors, "unterminated function parameters")
			return nil, nil
		}
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
		p.nextToken()
		// If there is "=x" after the name then that is a default parameter value
		if p.curTokenIs(token.ASSIGN) {
			p.nextToken()
			m[ident.Value] = p.parseExpressionStatement().Expression
			p.nextToken()
		}
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	return m, identifiers
}

// parseStringLiteral parses a string-literal.
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseRegexpLiteral parses a regular-expression.
func (p *Parser) parseRegexpLiteral() ast.Expression {
	flags := ""
	val := p.curToken.Literal
	if strings.HasPrefix(val, "(?") {
		val = strings.TrimPrefix(val, "(?")
		i := 0
		for i < len(val) {
			if val[i] == ')' {
				val = val[i+1:]
				break
			} else {
				flags += string(val[i])
			}
			i++
		}
	}
	return &ast.RegexpLiteral{Token: p.curToken, Value: val, Flags: flags}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	list := make([]ast.Expression, 0)
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	for p.peekTokenIs(token.NEWLINE) {
		p.nextToken()
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		// move to the comma
		p.nextToken()
		// advance across any extra newlines
		for p.peekTokenIs(token.NEWLINE) {
			p.nextToken()
		}
		// check if the list has ended after the newlines
		if p.peekTokenIs(end) {
			break
		}
		// move to the next expression
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return list
}

// parseInfixExpression parsea an array index expression.
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

// parseAssignExpression parses a bare assignment, without a `let`.
func (p *Parser) parseAssignExpression(name ast.Expression) ast.Expression {
	stmt := &ast.AssignStatement{Token: p.curToken}
	if n, ok := name.(*ast.Identifier); ok {
		stmt.Name = n
	} else {
		msg := fmt.Sprintf("expected assign token to be IDENT, got %s instead around line %d",
			name.TokenLiteral(), p.curToken.Line+1)
		p.errors = append(p.errors, msg)
		return nil
	}
	oper := p.curToken
	p.nextToken()
	switch oper.Type {
	case token.PLUS_EQUALS:
		stmt.Operator = "+="
	case token.MINUS_EQUALS:
		stmt.Operator = "-="
	case token.SLASH_EQUALS:
		stmt.Operator = "/="
	case token.ASTERISK_EQUALS:
		stmt.Operator = "*="
	case token.DECLARE:
		stmt.Operator = ":="
	default:
		stmt.Operator = "="
	}
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

// parseCallExpression parses a function-call expression.
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parsePipeExpression(first ast.Expression) ast.Expression {
	exp := &ast.PipeExpression{Token: p.curToken, Arguments: []ast.Expression{first}}
	for {
		if p.curTokenIs(token.EOF) {
			p.errors = append(p.errors, "unterminated pipe expression")
			return nil
		}
		// Move past the pipe operator itself
		p.nextToken()
		// Parse the next expression and add it to the ast.PipeExpression Arguments
		expr := p.parseExpression(PIPE)
		if expr == nil {
			p.errors = append(p.errors, "unable to parse pipe expression")
			return nil
		}
		exp.Arguments = append(exp.Arguments, expr)
		// Another pipe character continues the expression
		if p.peekTokenIs(token.PIPE) {
			p.nextToken()
			continue
		} else if p.peekTokenIs(token.EOF) {
			p.errors = append(p.errors, "unterminated pipe expression")
			return nil
		} else {
			// Anything else indicates the end of the pipe expression
			break
		}
	}
	return exp
}

// parseHashLiteral parses a hash literal
func (p *Parser) parseHashLiteral() ast.Expression {
	for p.peekTokenIs(token.NEWLINE) {
		p.nextToken()
	}
	// Empty {} turns into an empty hash (not a set)
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return &ast.HashLiteral{Token: p.curToken}
	}
	p.nextToken()
	firstKey := p.parseExpression(LOWEST)
	if p.peekTokenIs(token.COLON) { // This is a hash
		p.nextToken() // advance to colon
		p.nextToken() // advance to first key
		firstValue := p.parseExpression(LOWEST)
		hash := &ast.HashLiteral{
			Token: p.curToken,
			Pairs: map[ast.Expression]ast.Expression{
				firstKey: firstValue,
			},
		}
		for !p.peekTokenIs(token.RBRACE) {
			if !p.expectPeek(token.COMMA) {
				return nil
			}
			for p.peekTokenIs(token.NEWLINE) {
				p.nextToken()
			}
			if p.peekTokenIs(token.RBRACE) {
				break
			}
			key, value := p.parseHashKeyValue()
			if key == nil || value == nil {
				return nil
			}
			hash.Pairs[key] = value
			if !p.peekTokenIs(token.COMMA) {
				break
			}
		}
		for p.peekTokenIs(token.NEWLINE) {
			p.nextToken()
		}
		if !p.expectPeek(token.RBRACE) {
			return nil
		}
		return hash
	} else { // This is a set
		if !p.expectPeek(token.COMMA) {
			return nil
		}
		for p.peekTokenIs(token.NEWLINE) {
			p.nextToken()
		}
		set := &ast.SetLiteral{
			Token: p.curToken,
			Items: []ast.Expression{firstKey},
		}
		for !p.peekTokenIs(token.RBRACE) {
			p.nextToken()
			key := p.parseExpression(LOWEST)
			set.Items = append(set.Items, key)
			if !p.peekTokenIs(token.COMMA) {
				break
			}
			p.nextToken() // move to the comma
			for p.peekTokenIs(token.NEWLINE) {
				p.nextToken()
			}
		}
		if !p.expectPeek(token.RBRACE) {
			return nil
		}
		return set
	}
}

func (p *Parser) parseHashKeyValue() (ast.Expression, ast.Expression) {
	p.nextToken()
	key := p.parseExpression(LOWEST)
	if !p.expectPeek(token.COLON) {
		return nil, nil
	}
	p.nextToken()
	value := p.parseExpression(LOWEST)
	return key, value
}

// parseMethodCallExpression parses an object-based method-call.
func (p *Parser) parseMethodCallExpression(obj ast.Expression) ast.Expression {
	p.nextToken() // move to the attribute/method identifier
	name := p.parseIdentifier().(*ast.Identifier)
	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		methodCall := &ast.ObjectCallExpression{Token: p.curToken, Object: obj}
		methodCall.Call = p.parseCallExpression(name)
		return methodCall
	}
	return &ast.GetAttributeExpression{Object: obj, Attribute: name}
}

// curTokenIs tests if the current token has the given type.
func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

// peekTokenIs tests if the next token has the given type.
func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

// expectPeek validates the next token is of the given type,
// and advances if so.  If it is not an error is stored.
func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// peekPrecedence looks up the next token precedence.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// curPrecedence looks up the current token precedence.
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) eatNewlines() {
	for p.curTokenIs(token.NEWLINE) {
		p.nextToken()
	}
}
