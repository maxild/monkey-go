package parser

import (
	"fmt"
	"github.com/maxild/monkey/internal/ast"
	"github.com/maxild/monkey/internal/lexer"
	"github.com/maxild/monkey/internal/token"
	"strconv"
)

// Precedence: Highest binds the most/first
const (
	_ int = iota
	LOWEST
	// binary operators
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	// Unary operators
	PREFIX      // -X or !X
	// function application
	CALL        // myFunc(X)
)

var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	// TODO: Add line, column number to errors and lexer/tokens
	errors []string

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)

	// read two tokens so currToken and peekToken are both defined
	// (even though this seems a little weird, l.NextToken can be called multiple times after EOF)
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) currTokenIs(t token.Type) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead.",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if precedence, ok := precedences[p.currToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

// TODO: Write out the productions (grammar, CFG)
// TODO: Calculate First/Follow and assert grammar is LL(1)

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// program := statement+
	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// <stmt> -> <let_stmt>
//         | <return_stmt>
//         | <expression_stmt>
func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// Pratt parser methods (never advance the currToken passed the last token in the expression)

// <let_stmt> -> LET IDENT ASSIGN <expr> SEMICOLON
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Lexeme}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	p.eatOptionalSemicolon()

	return stmt
}

// <return_stmt> := RETURN <expr> SEMICOLON
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	p.eatOptionalSemicolon()

	return stmt
}

// wrapper/adapter
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	p.eatOptionalSemicolon()

	return stmt
}

func (p *Parser) eatOptionalSemicolon() {
	// TODO: The semicolon optional? Is this GOOD??!??
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
}

// TODO: We can only parse prefix expressions (see the table)
// <expr> :=
//         | ID
//         | BANG <expr>
//         | MINUS <expr>
//		   | INT
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// table-driven parser functions
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		msg := fmt.Sprintf("No prefix parse function for %s found.", p.currToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	leftExpr := prefix()

	// Find the deepest possible expression to parse (evaluate)
	// while the lookahead token is a registered infix operator token
	// with higher precedence we call the associated parser function
	// That is we call any associated parser function until we encounter
	// a token/operator that has higher precedence
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpr
		}

		// token was an infix token, and we have a binary operator (infix) expression
		p.nextToken()

		// call builder
		leftExpr = infix(leftExpr)
	}

	return leftExpr
}

//         | ID
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Lexeme}
}

//		   | INT
func (p *Parser) parseIntegerLiteral() ast.Expression {
	expr := &ast.IntegerLiteral{Token: p.currToken}
	value, err := strconv.ParseInt(p.currToken.Lexeme, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Lexeme)
		p.errors = append(p.errors, msg)
		return nil
	}
	expr.Value = value
	return expr
}

//         | BANG <expr>
//         | MINUS <expr>
func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{Token: p.currToken, Operator: p.currToken.Lexeme}
	p.nextToken()
	expr.Right = p.parseExpression(PREFIX) // recursive call
	return expr
}

//		   | <expr> OP <expr>
// where OP in (+, -, *, /, ==, !=, >, <)
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{Token: p.currToken, Operator: p.currToken.Lexeme, Left: left}
	// The precedence of the binary operator
	precedence := p.currPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence) // recursive call
	return expr
}
