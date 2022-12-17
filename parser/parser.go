// Package parser is used to parse an input program from its tokens and produce
// an abstract syntax tree (AST) as output.
//
// A parser is created by calling New() with a lexer as input. The parser should
// then be used only once, by calling parser.Parse() to produce the AST.
package parser

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/lexer"
	"github.com/cloudcmds/tamarin/tmpl"
	"github.com/cloudcmds/tamarin/token"
)

type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func() ast.Expression
)

// Parse is a shortcut that can be used to parse the given Tamarin source code.
// The lexer and parser are created internally and not exposed. ParseWithOpts
// should be used in production in order to pass a context.
func Parse(input string) (*ast.Program, error) {
	return New(lexer.New(input)).Parse(context.Background())
}

// ParseWithOpts is a shortcut that can be used to parse the given Tamarin source code.
// The lexer and parser are created internally and not exposed.
func ParseWithOpts(ctx context.Context, opts Opts) (*ast.Program, error) {
	lexerOpts := lexer.Opts{
		Input: opts.Input,
		File:  opts.File,
	}
	return New(lexer.NewWithOptions(lexerOpts)).Parse(ctx)
}

// Opts contains options for the parser.
type Opts struct {
	// Input is the string being parsed.
	Input string
	// File is the name of the file being parsed (optional).
	File string
}

// Parser object
type Parser struct {
	// l is our lexer
	l *lexer.Lexer

	// prevToken holds the previous token, which we already processed.
	prevToken token.Token

	// curToken holds the current token from the lexer.
	curToken token.Token

	// peekToken holds the next token from the lexer.
	peekToken token.Token

	// the parsing error, if any
	err ParserError

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
		prefixParseFns:  map[token.Type]prefixParseFn{},
		infixParseFns:   map[token.Type]infixParseFn{},
		postfixParseFns: map[token.Type]postfixParseFn{},
	}
	p.nextToken() // makes curToken=<empty>, peekToken=token[0]
	p.nextToken() // makes curToken=token[0], peekToken=token[1]

	// Register prefix-functions
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.FUNC, p.parseFunctionDefinition)
	p.registerPrefix(token.EOF, p.illegalToken)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.FOR, p.parseForLoopExpression)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.ILLEGAL, p.illegalToken)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPrefix(token.LBRACKET, p.parseListLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.NIL, p.parseNil)
	p.registerPrefix(token.REGEXP, p.parseRegexpLiteral)
	p.registerPrefix(token.REGEXP, p.parseRegexpLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BACKTICK, p.parseStringLiteral)
	p.registerPrefix(token.FSTRING, p.parseStringLiteral)
	p.registerPrefix(token.SWITCH, p.parseSwitchStatement)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.PIPE, p.parsePrefixExpression)
	p.registerPrefix(token.NEWLINE, p.parseNewlineLiteral)
	p.registerPrefix(token.IMPORT, p.parseImportStatement)

	// Register infix functions
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK_EQUALS, p.parseAssignExpression)
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
// prevToken, curToken, and peekToken. Any lexer error is captured but
// ignored.
func (p *Parser) nextToken() {
	p.nextTokenWithError()
}

// nextToken moves to the next token from the lexer, updating all of
// prevToken, curToken, and peekToken.
func (p *Parser) nextTokenWithError() error {
	// If we have an error, we can't move forward
	if p.err != nil {
		return p.err
	}
	var err error
	p.prevToken = p.curToken
	p.curToken = p.peekToken
	p.peekToken, err = p.l.NextToken()
	if err == nil {
		return nil // success
	}
	// The lexer encountered an error. We consider all lexer errors
	// "syntax errors" and parsing will now be considered broken.
	p.err = NewSyntaxError(ErrorOpts{
		Cause:         err,
		File:          p.l.File(),
		StartPosition: p.peekToken.StartPosition,
		EndPosition:   p.peekToken.EndPosition,
		SourceCode:    p.l.GetTokenLineText(p.peekToken),
	})
	return p.err
}

// Parse the program that is provided via the lexer.
func (p *Parser) Parse(ctx context.Context) (*ast.Program, error) {
	// It's possible for an error to already exist because we read tokens from
	// the lexer in the constructor. Parsing is already broken if so.
	if p.err != nil {
		return nil, p.err
	}
	// Parse the entire input program as a series of statements.
	// Parsing stops on the first occurrence of an error.
	program := &ast.Program{Statements: []ast.Statement{}}
	for p.curToken.Type != token.EOF {
		// Check for context timeout
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		if err := p.nextTokenWithError(); err != nil {
			return nil, err
		}
	}
	return program, p.err
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

func (p *Parser) noPrefixParseFnError(t token.Token) {
	if p.err != nil {
		return
	}
	p.err = NewParserError(ErrorOpts{
		ErrType:       "parse error",
		Message:       fmt.Sprintf("unexpected token %s", t.Literal),
		File:          p.l.File(),
		StartPosition: t.StartPosition,
		EndPosition:   t.EndPosition,
		SourceCode:    p.l.GetTokenLineText(t),
	})
}

// peekError raises an error if the next token is not the expected type.
func (p *Parser) peekError(context string, expected token.Type, got token.Token) {
	if p.err != nil {
		return
	}
	gotDesc := tokenDescription(got)
	expDesc := tokenTypeDescription(expected)
	p.err = NewParserError(ErrorOpts{
		ErrType: "parse error",
		Message: fmt.Sprintf("unexpected %s while parsing %s (expected %s)",
			gotDesc, context, expDesc),
		File:          p.l.File(),
		StartPosition: got.StartPosition,
		EndPosition:   got.EndPosition,
		SourceCode:    p.l.GetTokenLineText(got),
	})
}

func (p *Parser) setError(err ParserError) {
	if p.err != nil {
		return
	}
	p.err = err
}

// parseStatement parses a single statement.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.CONST:
		return p.parseConstStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.BREAK:
		return p.parseBreakStatement()
	case token.NEWLINE:
		return nil
	default:
		if p.peekTokenIs(token.DECLARE) {
			return p.parseDeclarationStatement()
		}
		return p.parseExpressionStatement()
	}
}

// parseVarStatement parses a var statement.
func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}
	if !p.expectPeek("var statement", token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek("var statement", token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseAssignmentValue()
	if stmt.Value == nil {
		return nil
	}
	return stmt
}

// parseDeclarationStatement parses a "i := value" statement as a "var" statement.
func (p *Parser) parseDeclarationStatement() *ast.VarStatement {
	name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek("declaration statement", token.DECLARE) {
		return nil
	}
	stmt := &ast.VarStatement{Token: p.curToken, Name: name}
	p.nextToken()
	stmt.Value = p.parseAssignmentValue()
	if stmt.Value == nil {
		return nil
	}
	return stmt
}

// parseConstStatement parses a constant declaration.
func (p *Parser) parseConstStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{Token: p.curToken}
	if !p.expectPeek("const statement", token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek("const statement", token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseAssignmentValue()
	if stmt.Value == nil {
		return nil
	}
	return stmt
}

// This is used to parse the right hand side (RHS) of the three types of
// assignment statements: var, const, and :=
func (p *Parser) parseAssignmentValue() ast.Expression {
	// Parse the value being assigned (the right hand side)
	result := p.parseExpression(LOWEST)
	if result == nil {
		p.setError(NewParserError(ErrorOpts{
			ErrType:       "parse error",
			Message:       "assignment is missing a value",
			File:          p.l.File(),
			StartPosition: p.prevToken.EndPosition,
			EndPosition:   p.prevToken.EndPosition,
			SourceCode:    p.l.GetTokenLineText(p.prevToken),
		}))
		return nil
	}
	switch p.peekToken.Type {
	// Assignment statements can be followed by a newline, semicolon, or EOF.
	case token.NEWLINE, token.SEMICOLON, token.EOF:
		p.nextToken()
		return result
	default:
		p.setError(NewParserError(ErrorOpts{
			ErrType:       "parse error",
			Message:       fmt.Sprintf("unexpected token %s following assignment value", p.peekToken.Literal),
			File:          p.l.File(),
			StartPosition: p.peekToken.StartPosition,
			EndPosition:   p.peekToken.EndPosition,
			SourceCode:    p.l.GetTokenLineText(p.peekToken),
		}))
		return nil
	}
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
			p.setError(NewParserError(ErrorOpts{
				ErrType:       "parse error",
				Message:       fmt.Sprintf("unexpected token %s following return value", p.peekToken.Literal),
				File:          p.l.File(),
				StartPosition: p.peekToken.StartPosition,
				EndPosition:   p.peekToken.EndPosition,
				SourceCode:    p.l.GetTokenLineText(p.peekToken),
			}))
			return nil
		}
	}
}

// parseBreakStatement parses a loop break statement.
func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := &ast.BreakStatement{Token: p.curToken}
	for p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.NEWLINE) {
		if err := p.nextTokenWithError(); err != nil {
			return nil
		}
	}
	return stmt
}

// parse Expression Statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	for p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.NEWLINE) {
		if err := p.nextTokenWithError(); err != nil {
			return nil
		}
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if p.curToken.Type == token.EOF {
		return nil
	}
	postfix := p.postfixParseFns[p.curToken.Type]
	if postfix != nil {
		return (postfix())
	}
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken)
		return nil
	}
	leftExp := prefix()
	if p.err != nil {
		return nil
	}
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
		if p.err != nil {
			break
		}
	}
	p.eatNewlines()
	return leftExp
}

func (p *Parser) illegalToken() ast.Expression {
	p.setError(NewParserError(ErrorOpts{
		ErrType:       "parse error",
		Message:       fmt.Sprintf("illegal token %s", p.curToken.Literal),
		File:          p.l.File(),
		StartPosition: p.curToken.StartPosition,
		EndPosition:   p.curToken.EndPosition,
		SourceCode:    p.l.GetTokenLineText(p.curToken),
	}))
	return nil
}

func (p *Parser) setTokenError(t token.Token, msg string, args ...interface{}) ast.Expression {
	p.setError(NewParserError(ErrorOpts{
		ErrType:       "parse error",
		Message:       fmt.Sprintf(msg, args...),
		File:          p.l.File(),
		StartPosition: t.StartPosition,
		EndPosition:   t.EndPosition,
		SourceCode:    p.l.GetTokenLineText(t),
	}))
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
		p.setError(NewParserError(ErrorOpts{
			ErrType:       "parse error",
			Message:       fmt.Sprintf("invalid integer: %s", p.curToken.Literal),
			File:          p.l.File(),
			StartPosition: p.curToken.StartPosition,
			EndPosition:   p.curToken.EndPosition,
			SourceCode:    p.l.GetTokenLineText(p.curToken),
		}))
		return nil
	}
	lit.Value = value
	return lit
}

// parseFloatLiteral parses a float-literal
func (p *Parser) parseFloatLiteral() ast.Expression {
	f := &ast.FloatLiteral{Token: p.curToken}
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.setError(NewParserError(ErrorOpts{
			ErrType:       "parse error",
			Message:       fmt.Sprintf("invalid float: %s", p.curToken.Literal),
			File:          p.l.File(),
			StartPosition: p.curToken.StartPosition,
			EndPosition:   p.curToken.EndPosition,
			SourceCode:    p.l.GetTokenLineText(p.curToken),
		}))
		return nil
	}
	f.Value = value
	return f
}

// parseSwitchStatement handles a switch statement
func (p *Parser) parseSwitchStatement() ast.Expression {
	expression := &ast.SwitchExpression{Token: p.curToken}
	p.nextToken()
	expression.Value = p.parseExpression(LOWEST)
	if expression.Value == nil {
		return nil
	}
	if !p.expectPeek("switch statement", token.LBRACE) {
		return nil
	}
	p.nextToken()
	p.eatNewlines()

	// Process the block which we think will contain various case-statements
	for !p.curTokenIs(token.RBRACE) {
		if p.curTokenIs(token.EOF) {
			p.setTokenError(p.prevToken, "unterminated switch statement")
			return nil
		}
		if p.curToken.Literal != "case" && p.curToken.Literal != "default" {
			p.setTokenError(p.curToken, "expected 'case' or 'default' (got %s)", p.curToken.Literal)
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
			p.setTokenError(p.curToken, "expected 'case' or 'default' (got %s)", p.curToken.Literal)
			return nil
		}
		if !p.expectPeek("switch statement", token.COLON) {
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
	var lastDefault token.Token
	for _, c := range expression.Choices {
		if c.Default {
			count++
		}
		lastDefault = c.Token
	}
	if count > 1 {
		p.setTokenError(lastDefault, "switch statement has multiple default blocks")
		return nil
	}
	return expression
}

func (p *Parser) parseImportStatement() ast.Expression {
	if !p.expectPeek("an import statement", token.IDENT) {
		return nil
	}
	return &ast.ImportStatement{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
}

// parseBoolean parses a boolean token.
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Bool{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// parseNil parses a nil keyword
func (p *Parser) parseNil() ast.Expression {
	return &ast.NilLiteral{Token: p.curToken}
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
		p.setTokenError(p.curToken, "nested ternary expression detected")
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

	if !p.expectPeek("ternary expression", token.COLON) { //skip the ":"
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
	if !p.expectPeek("grouped expression", token.RPAREN) {
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
	if !p.expectPeek("an if expression", token.LBRACE) {
		return nil
	}
	// The consequence
	expression.Consequence = p.parseBlockStatement()
	if expression.Consequence == nil {
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
		if !p.expectPeek("an if expression", token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
		if expression.Alternative == nil {
			return nil
		}
	}
	return expression
}

// parseForLoopExpression parses a for loop.
func (p *Parser) parseForLoopExpression() ast.Expression {
	expression := &ast.ForLoopExpression{Token: p.curToken}
	peekToken := p.peekToken

	var initExpression *ast.VarStatement
	switch peekToken.Type {
	case token.VAR:
		p.nextToken()
		initExpression = p.parseVarStatement()
	case token.IDENT:
		p.nextToken()
		initExpression = p.parseDeclarationStatement()
	case token.LBRACE:
		p.nextToken()
		expression.Consequence = p.parseBlockStatement()
		return expression
	default:
		desc := tokenDescription(p.peekToken)
		p.setTokenError(p.peekToken, "unexpected token in for loop: %s", desc)
		p.nextToken()
		return nil
	}
	expression.InitStatement = initExpression

	if !p.curTokenIs(token.SEMICOLON) {
		p.setTokenError(p.curToken, "expected a semicolon (got %s)", p.curToken.Literal)
		return nil
	}
	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek("for loop", token.SEMICOLON) {
		return nil
	}
	if !p.expectPeek("for loop", token.IDENT) {
		return nil
	}
	if p.peekTokenIs(token.PLUS_PLUS) || p.peekTokenIs(token.MINUS_MINUS) {
		p.nextToken()
		expression.PostStatement = p.parsePostfixExpression()
	} else {
		expression.PostStatement = p.parseExpression(LOWEST)
	}
	if expression.PostStatement == nil {
		return nil
	}
	if !p.expectPeek("for loop", token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()
	return expression
}

// parseBlockStatement parses a block.
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) {
		if p.curTokenIs(token.EOF) {
			p.setTokenError(block.Token, "unterminated block statement")
			return nil
		}
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		if err := p.nextTokenWithError(); err != nil {
			return nil
		}
	}
	return block
}

// parseFunctionLiteral parses a function-literal.
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek("function", token.LPAREN) {
		return nil
	}
	lit.Defaults, lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek("function", token.LBRACE) {
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
	if !p.expectPeek("function", token.LPAREN) {
		return nil
	}
	lit.Defaults, lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek("function", token.LBRACE) {
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
			p.setTokenError(p.prevToken, "unterminated function parameters")
			return nil, nil
		}
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
		if err := p.nextTokenWithError(); err != nil {
			return nil, nil
		}
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
	s := p.curToken.Literal
	if p.curToken.Type == token.BACKTICK || p.curToken.Type == token.STRING {
		return &ast.StringLiteral{Token: p.curToken, Value: s}
	}
	if !strings.Contains(s, "{") {
		return &ast.StringLiteral{Token: p.curToken, Value: s}
	}
	// Formatted/template string here (FSTRING)
	tmpl, err := tmpl.Parse(s)
	if err != nil {
		p.setTokenError(p.curToken, err.Error())
		return nil
	}
	var templateExps []*ast.ExpressionStatement
	for _, e := range tmpl.Fragments {
		if e.IsVariable {
			tmplAst, err := Parse(e.Value)
			if err != nil {
				p.setTokenError(p.curToken, err.Error())
				return nil
			}
			if len(tmplAst.Statements) == 0 {
				templateExps = append(templateExps, nil)
			} else if len(tmplAst.Statements) > 1 {
				p.setTokenError(p.curToken, "template contains more than one expression")
				return nil
			} else {
				stmt := tmplAst.Statements[0]
				exprStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					p.setTokenError(p.curToken, "template contains an unexpected statement type")
					return nil
				}
				templateExps = append(templateExps, exprStmt)
			}
		}
	}
	return &ast.StringLiteral{
		Token:               p.curToken,
		Value:               s,
		Template:            tmpl,
		TemplateExpressions: templateExps,
	}
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

func (p *Parser) parseListLiteral() ast.Expression {
	ll := &ast.ListLiteral{Token: p.curToken}
	ll.Items = p.parseExpressionList(token.RBRACKET)
	return ll
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	list := make([]ast.Expression, 0)
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	for p.peekTokenIs(token.NEWLINE) {
		if err := p.nextTokenWithError(); err != nil {
			return nil
		}
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		// move to the comma
		if err := p.nextTokenWithError(); err != nil {
			return nil
		}
		// advance across any extra newlines
		for p.peekTokenIs(token.NEWLINE) {
			if err := p.nextTokenWithError(); err != nil {
				return nil
			}
		}
		// check if the list has ended after the newlines
		if p.peekTokenIs(end) {
			break
		}
		// move to the next expression
		if err := p.nextTokenWithError(); err != nil {
			return nil
		}
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expectPeek("an expression list", end) {
		return nil
	}
	return list
}

// parseIndexExpression parses an array index expression.
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeek("an index expression", token.RBRACKET) {
		return nil
	}
	return exp
}

// parseAssignExpression parses a bare assignment, without a `var`.
func (p *Parser) parseAssignExpression(name ast.Expression) ast.Expression {
	stmt := &ast.AssignStatement{Token: p.curToken}
	if n, ok := name.(*ast.Identifier); ok {
		stmt.Name = n
	} else {
		p.setTokenError(p.curToken, "unexpected token for assignment: %s", name.TokenLiteral())
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
		// Move past the pipe operator itself
		if err := p.nextTokenWithError(); err != nil {
			return nil
		}
		// Parse the next expression and add it to the ast.PipeExpression Arguments
		expr := p.parseExpression(PIPE)
		if expr == nil {
			p.setTokenError(p.curToken, "invalid pipe expression")
			return nil
		}
		exp.Arguments = append(exp.Arguments, expr)
		// Another pipe character continues the expression
		if p.peekTokenIs(token.PIPE) {
			p.nextToken()
			continue
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
		if err := p.nextTokenWithError(); err != nil {
			return nil
		}
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
			if !p.expectPeek("hash", token.COMMA) {
				return nil
			}
			for p.peekTokenIs(token.NEWLINE) {
				if err := p.nextTokenWithError(); err != nil {
					return nil
				}
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
		if !p.expectPeek("hash", token.RBRACE) {
			return nil
		}
		return hash
	} else { // This is a set
		if !p.expectPeek("set", token.COMMA) {
			return nil
		}
		for p.peekTokenIs(token.NEWLINE) {
			if err := p.nextTokenWithError(); err != nil {
				return nil
			}
		}
		set := &ast.SetLiteral{
			Token: p.curToken,
			Items: []ast.Expression{firstKey},
		}
		for !p.peekTokenIs(token.RBRACE) {
			if err := p.nextTokenWithError(); err != nil {
				return nil
			}
			key := p.parseExpression(LOWEST)
			set.Items = append(set.Items, key)
			if !p.peekTokenIs(token.COMMA) {
				break
			}
			p.nextToken() // move to the comma
			for p.peekTokenIs(token.NEWLINE) {
				if err := p.nextTokenWithError(); err != nil {
					return nil
				}
			}
		}
		if !p.expectPeek("set", token.RBRACE) {
			return nil
		}
		return set
	}
}

func (p *Parser) parseHashKeyValue() (ast.Expression, ast.Expression) {
	p.nextToken()
	key := p.parseExpression(LOWEST)
	if !p.expectPeek("hash value", token.COLON) {
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
func (p *Parser) expectPeek(context string, t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(context, t, p.peekToken)
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
		if err := p.nextTokenWithError(); err != nil {
			return
		}
	}
}
