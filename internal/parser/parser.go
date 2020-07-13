package parser

import (
	"fmt"
	"github.com/maxild/monkey/internal/ast"
	"github.com/maxild/monkey/internal/lexer"
	"github.com/maxild/monkey/internal/token"
	"strconv"
)

// Precedence
const (
	_ int = iota
	LOWEST
	EQUALS		// ==
	LESSGREATER // > or <
	SUM			// +
	PRODUCT		// *
	PREFIX		// -X or !X
	CALL		// myFunc(X)
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.ExpressionStatement) ast.Expression
	)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	// TODO: Add line, column number to errors and lexer/tokens
	errors []string

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns map[token.Type]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	// read two tokens so currToken and peekToken are both defined
	// (even though this seems a little weird, l.NextToken can be called multiple times after EOF)
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

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

// TODO: Write out the productions (grammar, CFG)
// TODO: Calculate First/Follow and assert grammar is LL(1)

// <stmt> -> <let_stmt>
//         | <return_stmt>
func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	// TODO: Add more branches (arms)
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

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

	// TODO: We are skipping the expression until we encounter a semicolon
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// <return_stmt> := RETURN <expr> SEMICOLON
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	// TODO: We are skipping the expression until we encounter a semicolon
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// TODO: The semicolon optional? Is this GOOD??!??
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		return nil
	}
	leftExpr := prefix()
	return leftExpr
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
		msg := fmt.Sprintf("expected next token to be %s, got %s instead.",
			t, p.peekToken.Type)
		p.errors = append(p.errors, msg)
		return false
	}
}

// Pratt parser methods (never advance the currToken passed the last token in the expression)

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Lexeme}
}

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
