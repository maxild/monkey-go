package parser

import (
	"github.com/maxild/monkey/internal/ast"
	"github.com/maxild/monkey/internal/lexer"
	"github.com/maxild/monkey/internal/token"
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

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
func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	// TODO: Add more branches (arms)
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

// <let_stmt> -> LET IDENT ASSIGN <expr>
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
		return false
	}
}
