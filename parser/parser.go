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
	p.registerPrefix(token.BACKTICK, p.parseString)
	p.registerPrefix(token.BANG, p.parsePrefixExpr)
	p.registerPrefix(token.EOF, p.illegalToken)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.FLOAT, p.parseFloat)
	p.registerPrefix(token.FOR, p.parseFor)
	p.registerPrefix(token.FSTRING, p.parseString)
	p.registerPrefix(token.FUNC, p.parseFunc)
	p.registerPrefix(token.IDENT, p.parseIdent)
	p.registerPrefix(token.IF, p.parseIf)
	p.registerPrefix(token.ILLEGAL, p.illegalToken)
	p.registerPrefix(token.IMPORT, p.parseImport)
	p.registerPrefix(token.INT, p.parseInt)
	p.registerPrefix(token.LBRACE, p.parseMapOrSet)
	p.registerPrefix(token.LBRACKET, p.parseList)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpr)
	p.registerPrefix(token.MINUS, p.parsePrefixExpr)
	p.registerPrefix(token.NEWLINE, p.parseNewline)
	p.registerPrefix(token.NIL, p.parseNil)
	p.registerPrefix(token.PIPE, p.parsePrefixExpr)
	p.registerPrefix(token.RANGE, p.parseRange)
	p.registerPrefix(token.STRING, p.parseString)
	p.registerPrefix(token.SWITCH, p.parseSwitch)
	p.registerPrefix(token.TRUE, p.parseBoolean)

	// Register infix functions
	p.registerInfix(token.AND, p.parseInfixExpr)
	p.registerInfix(token.ASSIGN, p.parseAssign)
	p.registerInfix(token.ASTERISK_EQUALS, p.parseAssign)
	p.registerInfix(token.ASTERISK, p.parseInfixExpr)
	p.registerInfix(token.EQ, p.parseInfixExpr)
	p.registerInfix(token.GT_EQUALS, p.parseInfixExpr)
	p.registerInfix(token.GT, p.parseInfixExpr)
	p.registerInfix(token.LBRACKET, p.parseIndex)
	p.registerInfix(token.LPAREN, p.parseCall)
	p.registerInfix(token.LT_EQUALS, p.parseInfixExpr)
	p.registerInfix(token.LT, p.parseInfixExpr)
	p.registerInfix(token.MINUS_EQUALS, p.parseAssign)
	p.registerInfix(token.MINUS, p.parseInfixExpr)
	p.registerInfix(token.MOD, p.parseInfixExpr)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpr)
	p.registerInfix(token.OR, p.parseInfixExpr)
	p.registerInfix(token.PERIOD, p.parseGetAttr)
	p.registerInfix(token.PIPE, p.parsePipe)
	p.registerInfix(token.PLUS_EQUALS, p.parseAssign)
	p.registerInfix(token.PLUS, p.parseInfixExpr)
	p.registerInfix(token.POW, p.parseInfixExpr)
	p.registerInfix(token.QUESTION, p.parseTernary)
	p.registerInfix(token.SLASH_EQUALS, p.parseAssign)
	p.registerInfix(token.SLASH, p.parseInfixExpr)
	p.registerInfix(token.IN, p.parseIn)

	// Register postfix functions
	p.registerPostfix(token.MINUS_MINUS, p.parsePostfixExpr)
	p.registerPostfix(token.PLUS_PLUS, p.parsePostfixExpr)
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
	var statements []ast.Node
	for p.curToken.Type != token.EOF {
		// Check for context timeout
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		if err := p.nextTokenWithError(); err != nil {
			return nil, err
		}
	}
	return ast.NewProgram(statements), p.err
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
		Message:       fmt.Sprintf("invalid syntax (unexpected %q)", t.Literal),
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

func (p *Parser) parseStatement() ast.Node {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVar()
	case token.CONST:
		return p.parseConst()
	case token.RETURN:
		return p.parseReturn()
	case token.BREAK:
		return p.parseBreak()
	case token.CONTINUE:
		return p.parseContinue()
	case token.NEWLINE:
		return nil
	case token.IDENT:
		if p.peekTokenIs(token.DECLARE) || p.peekTokenIs(token.COMMA) {
			return p.parseDeclaration()
		}
		// intentional fallthrough!
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseVar() ast.Node {
	tok := p.curToken
	if !p.expectPeek("var statement", token.IDENT) {
		return nil
	}
	idents := []*ast.Ident{ast.NewIdent(p.curToken)}
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		if !p.expectPeek("var statement", token.IDENT) {
			return nil
		}
		idents = append(idents, ast.NewIdent(p.curToken))
	}
	if !p.expectPeek("var statement", token.ASSIGN) {
		return nil
	}
	p.nextToken()
	value := p.parseAssignmentValue()
	if value == nil {
		return nil
	}
	if len(idents) > 1 {
		return ast.NewMultiVar(tok, idents, value, false)
	}
	return ast.NewVar(tok, idents[0], value)
}

func (p *Parser) parseDeclaration() ast.Node {
	tok := p.curToken
	idents := []*ast.Ident{ast.NewIdent(p.curToken)}
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		if !p.expectPeek("declaration statement", token.IDENT) {
			return nil
		}
		idents = append(idents, ast.NewIdent(p.curToken))
	}
	if !p.expectPeek("declaration statement", token.DECLARE) {
		return nil
	}
	p.nextToken()
	value := p.parseAssignmentValue()
	if value == nil {
		return nil
	}
	if len(idents) > 1 {
		return ast.NewMultiVar(tok, idents, value, true)
	}
	return ast.NewDeclaration(tok, idents[0], value)
}

func (p *Parser) parseConst() *ast.Const {
	tok := p.curToken
	if !p.expectPeek("const statement", token.IDENT) {
		return nil
	}
	ident := ast.NewIdent(p.curToken)
	if !p.expectPeek("const statement", token.ASSIGN) {
		return nil
	}
	p.nextToken()
	value := p.parseAssignmentValue()
	if value == nil {
		return nil
	}
	return ast.NewConst(tok, ident, value)
}

// Parses the right hand side of an assignment statement.
func (p *Parser) parseAssignmentValue() ast.Expression {
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
	// Assignment statements can be followed by a newline, semicolon, EOF, or }
	case token.NEWLINE, token.SEMICOLON, token.EOF:
		p.nextToken()
		return result
	case token.RBRACE, token.LBRACE:
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

func (p *Parser) parseReturn() *ast.Control {
	returnToken := p.curToken
	p.nextToken()
	value := p.parseExpression(LOWEST)
	for {
		switch p.peekToken.Type {
		case token.SEMICOLON, token.NEWLINE, token.EOF:
			p.nextToken()
			return ast.NewControl(returnToken, value)
		case token.RBRACE:
			return ast.NewControl(returnToken, value)
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

func (p *Parser) parseBreak() *ast.Control {
	stmt := ast.NewControl(p.curToken, nil)
	for p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.NEWLINE) {
		if err := p.nextTokenWithError(); err != nil {
			return nil
		}
	}
	return stmt
}

func (p *Parser) parseContinue() *ast.Control {
	stmt := ast.NewControl(p.curToken, nil)
	for p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.NEWLINE) {
		if err := p.nextTokenWithError(); err != nil {
			return nil
		}
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Expression {
	expr := p.parseExpression(LOWEST)
	if expr == nil {
		p.setTokenError(p.curToken, "invalid syntax")
	}
	for p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.NEWLINE) {
		if err := p.nextTokenWithError(); err != nil {
			return nil
		}
	}
	return expr
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if p.curToken.Type == token.EOF || p.err != nil {
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
	if p.err != nil || leftExp == nil {
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

func (p *Parser) parseIdent() ast.Expression {
	return ast.NewIdent(p.curToken)
}

func (p *Parser) parseInt() ast.Expression {
	tok, lit := p.curToken, p.curToken.Literal
	var value int64
	var err error
	if strings.HasPrefix(lit, "0b") {
		value, err = strconv.ParseInt(lit[2:], 2, 64)
	} else if strings.HasPrefix(lit, "0x") {
		value, err = strconv.ParseInt(lit[2:], 16, 64)
	} else {
		value, err = strconv.ParseInt(lit, 10, 64)
	}
	if err != nil {
		p.setError(NewParserError(ErrorOpts{
			ErrType:       "parse error",
			Message:       fmt.Sprintf("invalid integer: %s", lit),
			File:          p.l.File(),
			StartPosition: tok.StartPosition,
			EndPosition:   tok.EndPosition,
			SourceCode:    p.l.GetTokenLineText(tok),
		}))
		return nil
	}
	return ast.NewInt(tok, value)
}

func (p *Parser) parseFloat() ast.Expression {
	tok, lit := p.curToken, p.curToken.Literal
	value, err := strconv.ParseFloat(lit, 64)
	if err != nil {
		p.setError(NewParserError(ErrorOpts{
			ErrType:       "parse error",
			Message:       fmt.Sprintf("invalid float: %s", lit),
			File:          p.l.File(),
			StartPosition: p.curToken.StartPosition,
			EndPosition:   p.curToken.EndPosition,
			SourceCode:    p.l.GetTokenLineText(p.curToken),
		}))
		return nil
	}
	return ast.NewFloat(tok, value)
}

func (p *Parser) parseSwitch() ast.Expression {
	switchToken := p.curToken
	p.nextToken()
	switchValue := p.parseExpression(LOWEST)
	if switchValue == nil {
		return nil
	}
	if !p.expectPeek("switch statement", token.LBRACE) {
		return nil
	}
	p.nextToken()
	p.eatNewlines()
	// Process the switch case statements
	var cases []*ast.Case
	var defaultCaseCount int
	// Each time through this loop we process one case statement
	for !p.curTokenIs(token.RBRACE) {
		if p.curTokenIs(token.EOF) {
			p.setTokenError(p.prevToken, "unterminated switch statement")
			return nil
		}
		if p.curToken.Literal != "case" && p.curToken.Literal != "default" {
			p.setTokenError(p.curToken, "expected 'case' or 'default' (got %s)", p.curToken.Literal)
			return nil
		}
		caseToken := p.curToken
		var isDefaultCase bool
		var caseExprs []ast.Expression
		if p.curTokenIs(token.DEFAULT) {
			isDefaultCase = true
		} else if p.curTokenIs(token.CASE) {
			p.nextToken() // move to the token following "case"
			caseExprs = append(caseExprs, p.parseExpression(LOWEST))
			for p.peekTokenIs(token.COMMA) {
				p.nextToken() // move to the comma
				p.nextToken() // move to the following expression
				caseExprs = append(caseExprs, p.parseExpression(LOWEST))
			}
		} else {
			p.setTokenError(p.curToken, "expected 'case' or 'default' (got %s)", p.curToken.Literal)
			return nil
		}
		if !p.expectPeek("switch statement", token.COLON) {
			return nil
		}
		// Now we are at the block of code to be executed for this case
		p.nextToken()
		p.eatNewlines()
		blockFirstToken := p.curToken
		var blockStatements []ast.Node
		for {
			stmt := p.parseStatement()
			if stmt == nil {
				return nil
			}
			blockStatements = append(blockStatements, stmt)
			p.eatNewlines()
			if p.curTokenIs(token.CASE) || p.curTokenIs(token.DEFAULT) || p.curTokenIs(token.RBRACE) {
				break
			}
		}
		block := ast.NewBlock(blockFirstToken, blockStatements)
		if isDefaultCase {
			defaultCaseCount++
			if defaultCaseCount > 1 {
				p.setTokenError(caseToken, "switch statement has multiple default blocks")
				return nil
			}
			cases = append(cases, ast.NewDefaultCase(caseToken, block))
		} else {
			cases = append(cases, ast.NewCase(caseToken, caseExprs, block))
		}
	}
	return ast.NewSwitch(switchToken, switchValue, cases)
}

func (p *Parser) parseImport() ast.Expression {
	importToken := p.curToken
	if !p.expectPeek("an import statement", token.IDENT) {
		return nil
	}
	return ast.NewImport(importToken, ast.NewIdent(p.curToken))
}

func (p *Parser) parseBoolean() ast.Expression {
	return ast.NewBool(p.curToken, p.curTokenIs(token.TRUE))
}

func (p *Parser) parseNil() ast.Expression {
	return ast.NewNil(p.curToken)
}

func (p *Parser) parsePrefixExpr() ast.Expression {
	operator := p.curToken
	p.nextToken()
	right := p.parseExpression(PREFIX)
	if right == nil {
		p.setTokenError(p.curToken, "invalid prefix expression")
		return nil
	}
	return ast.NewPrefix(operator, right)
}

func (p *Parser) parseNewline() ast.Expression {
	p.nextToken()
	return nil
}

func (p *Parser) parsePostfixExpr() ast.Expression {
	return ast.NewPostfix(p.prevToken, p.curToken.Literal)
}

func (p *Parser) parseInfixExpr(left ast.Expression) ast.Expression {
	firstToken := p.curToken
	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)
	if right == nil {
		p.setTokenError(p.curToken, "invalid expression")
		return nil
	}
	return ast.NewInfix(firstToken, left, firstToken.Literal, right)
}

func (p *Parser) parseTernary(condition ast.Expression) ast.Expression {
	if p.tern {
		p.setTokenError(p.curToken, "nested ternary expression detected")
		return nil
	}
	p.tern = true
	defer func() { p.tern = false }()

	firstToken := p.curToken // the "?"
	p.nextToken()            // move past the '?'
	precedence := p.curPrecedence()
	ifTrue := p.parseExpression(precedence)
	if ifTrue == nil {
		p.setTokenError(p.curToken, "invalid syntax in ternary if true expression")
	}
	if !p.expectPeek("ternary expression", token.COLON) { // moves to the ":"
		return nil
	}
	p.nextToken() // moves after the ":"
	ifFalse := p.parseExpression(precedence)
	if ifFalse == nil {
		p.setTokenError(p.curToken, "invalid syntax in ternary if false expression")
	}
	return ast.NewTernary(firstToken, condition, ifTrue, ifFalse)
}

func (p *Parser) parseGroupedExpr() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek("grouped expression", token.RPAREN) {
		return nil
	}
	return exp
}

// Parses an entire if, else if, else block. Else-ifs are handled recursively.
func (p *Parser) parseIf() ast.Expression {
	ifToken := p.curToken
	p.nextToken() // move past the "if"
	cond := p.parseExpression(LOWEST)
	if cond == nil {
		return nil
	}
	if !p.expectPeek("an if expression", token.LBRACE) { // move to the "{"
		return nil
	}
	consequence := p.parseBlock()
	if consequence == nil {
		return nil
	}
	var alternative *ast.Block
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()                // move to the "else"
		if p.peekTokenIs(token.IF) { // this is an "else if"
			p.nextToken() // move to the "if"
			nestedIfToken := p.curToken
			nestedIf := p.parseIf()
			alternative := ast.NewBlock(nestedIfToken, []ast.Node{nestedIf})
			return ast.NewIf(ifToken, cond, consequence, alternative)
		}
		if !p.expectPeek("an if expression", token.LBRACE) {
			return nil
		}
		alternative = p.parseBlock()
		if alternative == nil {
			return nil
		}
	}
	return ast.NewIf(ifToken, cond, consequence, alternative)
}

func (p *Parser) parseFor() ast.Expression {
	forToken := p.curToken
	// Check for simple form: "for { ... }"
	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		consequence := p.parseBlock()
		if consequence == nil {
			return nil
		}
		return ast.NewSimpleFor(forToken, consequence)
	}
	p.nextToken()
	forExprToken := p.curToken
	firstExpr := p.parseStatement()
	if firstExpr == nil {
		p.setTokenError(forExprToken, "invalid for loop expression")
		p.nextToken()
		return nil
	}
	// v, ok := firstExpr.(*ast.Var)
	// if ok {
	// 	vs, vv := v.Value()
	// 	rng, ok := vv.(*ast.Range)
	// 	if ok {
	// 		fmt.Println("FIRST:", vs, rng, reflect.TypeOf(rng))
	// 	}
	// }
	// Check for while loop form: "for condition { ... }"
	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		consequence := p.parseBlock()
		if consequence == nil {
			return nil
		}
		return ast.NewFor(forToken, firstExpr, consequence, nil, nil)
	}
	if !p.curTokenIs(token.SEMICOLON) {
		p.setTokenError(p.curToken, "expected a semicolon (got %s)", p.curToken.Literal)
		return nil
	}
	p.nextToken() // move past the ";"
	condition := p.parseExpression(LOWEST)
	if !p.expectPeek("for loop", token.SEMICOLON) {
		return nil
	}
	if !p.expectPeek("for loop", token.IDENT) {
		return nil
	}
	var postExpr ast.Expression
	if p.peekTokenIs(token.PLUS_PLUS) || p.peekTokenIs(token.MINUS_MINUS) {
		p.nextToken()
		postExpr = p.parsePostfixExpr()
	} else {
		postExpr = p.parseExpression(LOWEST)
	}
	if postExpr == nil {
		return nil
	}
	if !p.expectPeek("for loop", token.LBRACE) {
		return nil
	}
	consequence := p.parseBlock()
	if consequence == nil {
		return nil
	}
	return ast.NewFor(forToken, condition, consequence, firstExpr, postExpr)
}

func (p *Parser) parseBlock() *ast.Block {
	lbrace := p.curToken
	var statements []ast.Node
	p.nextToken() // move past the "{"
	for !p.curTokenIs(token.RBRACE) {
		if p.curTokenIs(token.EOF) {
			p.setTokenError(lbrace, "unterminated block statement")
			return nil
		}
		if s := p.parseStatement(); s != nil {
			statements = append(statements, s)
		}
		if err := p.nextTokenWithError(); err != nil {
			return nil
		}
	}
	return ast.NewBlock(lbrace, statements)
}

func (p *Parser) parseFunc() ast.Expression {
	funcToken := p.curToken
	var ident *ast.Ident
	if p.peekTokenIs(token.IDENT) { // Read optional function name
		p.nextToken()
		ident = ast.NewIdent(p.curToken)
	}
	if !p.expectPeek("function", token.LPAREN) { // Move to the "("
		return nil
	}
	defaults, params := p.parseFuncParams()
	if !p.expectPeek("function", token.LBRACE) { // move to the "{"
		return nil
	}
	return ast.NewFunc(funcToken, ident, params, defaults, p.parseBlock())
}

func (p *Parser) parseFuncParams() (map[string]ast.Expression, []*ast.Ident) {
	// If the next parameter is ")", then there are no parameters
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return map[string]ast.Expression{}, nil
	}
	defaults := map[string]ast.Expression{}
	params := make([]*ast.Ident, 0)
	p.nextToken()
	for !p.curTokenIs(token.RPAREN) { // Keep going until we find a ")"
		if p.curTokenIs(token.EOF) {
			p.setTokenError(p.prevToken, "unterminated function parameters")
			return nil, nil
		}
		if !p.curTokenIs(token.IDENT) {
			p.setTokenError(p.curToken, "expected an identifier (got %s)", p.curToken.Literal)
			return nil, nil
		}
		ident := ast.NewIdent(p.curToken)
		params = append(params, ident)
		if err := p.nextTokenWithError(); err != nil {
			return nil, nil
		}
		// If there is "=expr" after the name then expr is a default value
		if p.curTokenIs(token.ASSIGN) {
			p.nextToken()
			expr := p.parseExpression(LOWEST)
			if expr == nil {
				return nil, nil
			}
			defaults[ident.String()] = expr
			p.nextToken()
		}
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	return defaults, params
}

func (p *Parser) parseString() ast.Expression {
	strToken := p.curToken
	if strToken.Type == token.BACKTICK || strToken.Type == token.STRING {
		return ast.NewString(strToken)
	}
	if !strings.Contains(strToken.Literal, "{") {
		return ast.NewString(strToken)
	}
	// Template string with interpolation
	tmpl, err := tmpl.Parse(strToken.Literal)
	if err != nil {
		p.setTokenError(strToken, err.Error())
		return nil
	}
	var exprs []ast.Expression
	for _, e := range tmpl.Fragments {
		if !e.IsVariable {
			continue
		}
		tmplAst, err := Parse(e.Value)
		if err != nil {
			p.setTokenError(strToken, err.Error())
			return nil
		}
		statements := tmplAst.Statements()
		if len(statements) == 0 {
			exprs = append(exprs, nil)
		} else if len(statements) > 1 {
			p.setTokenError(strToken, "template contains more than one expression")
			return nil
		} else {
			stmt := statements[0]
			expr, ok := stmt.(ast.Expression)
			if !ok {
				p.setTokenError(strToken, "template contains an unexpected statement type")
				return nil
			}
			exprs = append(exprs, expr)
		}
	}
	return ast.NewTemplatedString(strToken, tmpl, exprs)
}

func (p *Parser) parseList() ast.Expression {
	bracket := p.curToken
	items := p.parseExprList(token.RBRACKET)
	return ast.NewList(bracket, items)
}

func (p *Parser) parseExprList(end token.Type) []ast.Expression {
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
	expr := p.parseExpression(LOWEST)
	if expr == nil {
		p.setTokenError(p.curToken, "invalid syntax in list expression")
		return nil
	}
	list = append(list, expr)
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

func (p *Parser) parseIndex(left ast.Expression) ast.Expression {
	indexToken := p.curToken
	var firstIndex, secondIndex ast.Expression
	if !p.peekTokenIs(token.COLON) {
		p.nextToken() // move to the first index
		firstIndex = p.parseExpression(LOWEST)
		if p.peekTokenIs(token.RBRACKET) {
			p.nextToken() // move to the "]"
			return ast.NewIndex(indexToken, left, firstIndex)
		}
	}
	if p.peekTokenIs(token.COLON) {
		p.nextToken() // move to the ":"
		if p.peekTokenIs(token.RBRACKET) {
			p.nextToken() // move to the "]"
			return ast.NewSlice(indexToken, left, firstIndex, nil)
		}
		p.nextToken() // move to the second index
		secondIndex = p.parseExpression(LOWEST)
	}
	if !p.expectPeek("an index expression", token.RBRACKET) {
		return nil
	}
	return ast.NewSlice(indexToken, left, firstIndex, secondIndex)
}

func (p *Parser) parseAssign(name ast.Expression) ast.Expression {
	operator := p.curToken
	var ident *ast.Ident
	var index *ast.Index
	switch node := name.(type) {
	case *ast.Ident:
		ident = node
	case *ast.Index:
		index = node
	default:
		p.setTokenError(operator, "unexpected token for assignment: %s", name.Literal())
		return nil
	}
	switch operator.Type {
	case token.PLUS_EQUALS, token.MINUS_EQUALS, token.SLASH_EQUALS,
		token.ASTERISK_EQUALS, token.DECLARE, token.ASSIGN:
		// this is a valid operator
	default:
		p.setTokenError(operator, "unsupported operator for assignment: %s", operator.Literal)
		return nil
	}
	p.nextToken() // move to the RHS value
	right := p.parseExpression(LOWEST)
	if right == nil {
		p.setTokenError(p.curToken, "invalid assignment statement value")
		return nil
	}
	if index != nil {
		return ast.NewAssignIndex(operator, index, right)
	}
	return ast.NewAssign(operator, ident, right)
}

func (p *Parser) parseCall(function ast.Expression) ast.Expression {
	callToken := p.curToken
	arguments := p.parseExprList(token.RPAREN)
	if arguments == nil {
		return nil
	}
	return ast.NewCall(callToken, function, arguments)
}

func (p *Parser) parsePipe(first ast.Expression) ast.Expression {
	pipeToken := p.curToken
	exprs := []ast.Expression{first}
	for {
		// Move past the pipe operator itself
		if err := p.nextTokenWithError(); err != nil {
			return nil
		}
		// Advance across any extra newlines
		p.eatNewlines()
		// Parse the next expression and add it to the ast.Pipe Arguments
		expr := p.parseExpression(PIPE)
		if expr == nil {
			p.setTokenError(p.curToken, "invalid pipe expression")
			return nil
		}
		exprs = append(exprs, expr)
		// Another pipe character continues the expression
		if p.peekTokenIs(token.PIPE) {
			p.nextToken() // move to the next "|"
			continue
		} else {
			// Anything else indicates the end of the pipe expression
			break
		}
	}
	return ast.NewPipe(pipeToken, exprs)
}

func (p *Parser) parseIn(left ast.Expression) ast.Expression {
	inToken := p.curToken
	if err := p.nextTokenWithError(); err != nil {
		return nil
	}
	right := p.parseExpression(IN)
	if right == nil {
		p.setTokenError(p.curToken, "invalid in expression")
		return nil
	}
	return ast.NewIn(inToken, left, right)
}

func (p *Parser) parseRange() ast.Expression {
	rangeToken := p.curToken
	if err := p.nextTokenWithError(); err != nil {
		return nil
	}
	container := p.parseExpression(RANGE)
	if container == nil {
		p.setTokenError(p.curToken, "invalid range expression")
		return nil
	}
	return ast.NewRange(rangeToken, container)
}

func (p *Parser) parseMapOrSet() ast.Expression {
	firstToken := p.curToken
	for p.peekTokenIs(token.NEWLINE) {
		if err := p.nextTokenWithError(); err != nil {
			return nil
		}
	}
	// Empty {} turns into an empty map (not a set)
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return ast.NewMap(firstToken, nil)
	}
	p.nextToken() // move to the first key
	firstKey := p.parseExpression(LOWEST)
	if p.peekTokenIs(token.COLON) { // This is a map
		p.nextToken() // move to the ":"
		p.nextToken() // move to the first value
		firstValue := p.parseExpression(LOWEST)
		pairs := map[ast.Expression]ast.Expression{firstKey: firstValue}
		for !p.peekTokenIs(token.RBRACE) {
			if !p.expectPeek("map", token.COMMA) {
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
			key, value := p.parseKeyValue()
			if key == nil || value == nil {
				return nil
			}
			pairs[key] = value
			if !p.peekTokenIs(token.COMMA) {
				break
			}
		}
		for p.peekTokenIs(token.NEWLINE) {
			p.nextToken()
		}
		if !p.expectPeek("map", token.RBRACE) {
			return nil
		}
		return ast.NewMap(firstToken, pairs)
	} else { // This is a set
		items := []ast.Expression{firstKey}
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		} else if p.peekTokenIs(token.RBRACE) {
			p.nextToken()
			return ast.NewSet(firstToken, items)
		} else {
			p.setTokenError(p.peekToken, "expected , or } after set element")
			return nil
		}
		for p.peekTokenIs(token.NEWLINE) {
			if err := p.nextTokenWithError(); err != nil {
				return nil
			}
		}
		for !p.peekTokenIs(token.RBRACE) {
			if err := p.nextTokenWithError(); err != nil {
				return nil
			}
			key := p.parseExpression(LOWEST)
			items = append(items, key)
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
		return ast.NewSet(firstToken, items)
	}
}

func (p *Parser) parseKeyValue() (ast.Expression, ast.Expression) {
	p.nextToken()
	key := p.parseExpression(LOWEST)
	if !p.expectPeek("hash value", token.COLON) {
		return nil, nil
	}
	p.nextToken()
	value := p.parseExpression(LOWEST)
	return key, value
}

func (p *Parser) parseGetAttr(obj ast.Expression) ast.Expression {
	period := p.curToken
	p.nextToken()
	p.eatNewlines()
	if !p.curTokenIs(token.IDENT) {
		p.setTokenError(p.curToken, "expected an identifier after %q", ".")
		return nil
	}
	name := p.parseIdent().(*ast.Ident)
	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		callExpr := p.parseCall(name)
		return ast.NewObjectCall(period, obj, callExpr)
	}
	return ast.NewGetAttr(period, obj, name)
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
