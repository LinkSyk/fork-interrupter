package parser

import (
	"fmt"
	"interrupter/ast"
	"interrupter/lexer"
	"interrupter/token"
	"interrupter/xlog"
	"io"
	"strconv"
)

// 运算符优先级
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQT:     EQUALS,
	token.NOTEQT:  EQUALS,
	token.LT:      LESSGREATER,
	token.GT:      LESSGREATER,
	token.PLUS:    SUM,
	token.SUB:     SUM,
	token.DIV:     PRODUCT,
	token.MULTI:   PRODUCT,
	token.LPARENT: CALL,
}

// parse statement
type (
	Parser struct {
		lexer     *lexer.Lexer
		curToken  token.Token
		peekToken token.Token
		errors    []string

		prefixParseFns map[token.TokenType]prefixParseFn
		infixParseFns  map[token.TokenType]infixParseFn
	}
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:          l,
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}
	// assign value to curToken and peekToken
	p.nextToken()
	p.nextToken()

	// register prefix expression function
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.SUB, p.parsePrefixExpression)
	p.registerPrefix(token.LPARENT, p.parseGroupedExpression)

	// register infix expression function
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.SUB, p.parseInfixExpression)
	p.registerInfix(token.MULTI, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.EQT, p.parseInfixExpression)
	p.registerInfix(token.NOTEQT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LPARENT, p.parseCallExpression)

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) registerPrefix(t token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[t] = fn
}

func (p *Parser) registerInfix(t token.TokenType, fn infixParseFn) {
	p.infixParseFns[t] = fn
}

func (p *Parser) curTokenAs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenAs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type != t {
		xlog.Debugf("expectPeek fail expect: %s, peek: %s\n", t, p.peekToken.Type)
		p.peekError(t)
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) PrintAllToken(out io.Writer) {
	for !p.curTokenAs(token.EOF) {
		fmt.Fprintf(out, "type: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)
		p.nextToken()
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	statements := []ast.Statement{}
	prog := &ast.Program{}
	for !p.curTokenAs(token.EOF) {
		stmt := p.parseStatement()
		if stmt == nil {
			xlog.Debugf("parse program statement is nil\n")
			return prog
		}
		statements = append(statements, stmt)
		p.nextToken()
	}
	prog.Statements = statements
	return prog
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	// let statement
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	xlog.Debugf("enter parseLetStatement, curToken type: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	stmt.Name = ident

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	// 如果是；就读取
	if p.peekTokenAs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenAs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenAs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix, ok := p.prefixParseFns[p.curToken.Type]
	if !ok {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()
	for !p.peekTokenAs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// peekToken的运算符优先级
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// curToken的运算符优先级
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenAs(token.TRUE)}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	il := &ast.IntegerLiteral{Token: p.curToken}
	v, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	il.Value = v
	return il
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	pe := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	pe.Right = p.parseExpression(PREFIX)
	return pe
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	ie := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	ie.Right = p.parseExpression(precedence)
	return ie
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPARENT) {
		return nil
	}
	return exp
}

// func (p *Parser) parseInfixGroupedExpression(left ast.Expression) ast.Expression {
// 	ie := &ast.InfixExpression{
// 		Token:    p.curToken,
// 		Operator: p.curToken.Literal,
// 		Left:     left,
// 	}
// 	precedence := p.curPrecedence()
// 	p.nextToken()
// 	// p.nextToken()
// 	ie.Right = p.parseExpression(precedence)
// 	if !p.expectPeek(token.RPARENT) {
// 		return nil
// 	}
// 	return ie
// }

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	call := &ast.CallExpression{
		Token:    p.curToken,
		Function: left,
	}
	args := p.parseCallArguments()
	if args == nil {
		return nil
	}
	call.Arguments = args
	return call
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var exps []ast.Expression

	// cur: (
	if p.peekTokenAs(token.RPARENT) {
		p.nextToken()
		return exps
	}

	// cur: arg1
	p.nextToken()
	exps = append(exps, p.parseExpression(LOWEST))

	for p.peekTokenAs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		exp := p.parseExpression(LOWEST)
		if exp == nil {
			return nil
		}
		exps = append(exps, exp)
	}

	if !p.expectPeek(token.RPARENT) {
		return nil
	}
	return exps
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}
