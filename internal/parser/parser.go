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
	PREFIX // -X or !X
	// function application
	CALL // myFunc(X)
)

// Table of precedence per token (kind) is defined for all infix operators
var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	// not defined for prefix operators (-, !)
}

// TODO: What about left and right associative operators?
// All our operators are left associative, but there could be right-associative operators in a language

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

	// null denotations ("nuds")
	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	// left denotations ("leds")
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

func (p *Parser) match(t token.Type) bool {
	if p.currTokenIs(t) {
		p.nextToken()
		return true
	} else {
		msg := fmt.Sprintf("expected next token to be %s, got %s instead.",
			t, p.currToken.Type)
		p.errors = append(p.errors, msg)
		return false
	}
}

func (p *Parser) matchPeek(t token.Type) bool {
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

func (p *Parser) Errors() []string {
	return p.errors
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
	stmt := &ast.LetStatement{
		Token: p.currToken,
	}

	// eat 'let'
	if !p.matchPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Lexeme,
	}

	p.nextToken() // eat ID

	// eat '='
	if !p.match(token.ASSIGN) {
		return nil
	}

	stmt.Value = p.parseExpression(LOWEST)

	p.eatOptionalSemicolon()

	return stmt
}

// <return_stmt> := RETURN <expr> SEMICOLON
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.currToken,
	}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	p.eatOptionalSemicolon()

	return stmt
}

// wrapper/adapter
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.currToken,
	}

	stmt.Expression = p.parseExpression(LOWEST)

	p.eatOptionalSemicolon()

	return stmt
}

func (p *Parser) eatOptionalSemicolon() {
	// TODO: The semicolon is optional? Is this GOOD??!??
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
}

//
// Pratt (Top-Down Operator Precedence) Parser
//

// TODO: We can only parse prefix expressions (see the table)
// <expr> :=
//         | ID
//         | BANG <expr>
//         | MINUS <expr>
//		   | INT
// NOTE: precedence tells the parseExpression function which expression can be parsed by that call
//       If precedence is lower for the following token than is allowed by the precedence argument,
//       the parser will stop parsing and just return what it has so far.
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

	// precedence represents the "right binding power" of the caller. If precedence is very high, then
	// the loop will not be executed, and no other infixParseFn will get a chance to
	// create a binary expression with leftExpr as the left arm. Instead leftExpr will
	// return as a "right" arm to whatever expression was previously being parsed in order
	// for this expression to bind with higher precedence.

	// The BOOK calls precedence "the right binding power" of the prev operator (should it bind as right "arm")
	// The BOOK calls p.peekPrecedence "the left binding power" of the next operator (should it bind as "left arm")
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// If we get here another infix parser function is going to get our leftExpr as a left arm
		// This means the precedence of the left operator is lower than the precedence of
		// the right operator in the current context
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpr
		}

		// token was an infix token, and we have a binary operator (infix) expression
		p.nextToken()

		// call builder: This will "suck in" the leftExpr as the left "arm" of some infix expression
		leftExpr = infix(leftExpr)
	}

	return leftExpr
}

// NOTE: None of the following registered (parser) functions does call nextToken unless
//       they need to parse 2 or more tokens. If they are parsing 2 tokens they will call
//       nextToken once, parsing 3 tokens means calling nextToken twice etc. Also if the
//       parser function calls back into parseExpression (kind of like a recursive call)
//       this means the same rules apply, and nextToken will be called such that peekToken
//		 is pointing at the following token after the production LHS have been recognized.

// non-recursive
//         | ID
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Lexeme,
	}
}

// non-recursive
//		   | INT
func (p *Parser) parseIntegerLiteral() ast.Expression {
	expr := &ast.IntegerLiteral{
		Token: p.currToken,
	}
	value, err := strconv.ParseInt(p.currToken.Lexeme, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Lexeme)
		p.errors = append(p.errors, msg)
		return nil
	}
	expr.Value = value
	return expr
}

//       TRUE | FALSE
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.currToken,
		Value: p.currTokenIs(token.TRUE),
	}
}

//		LPARAN <expr> RPARAN
func (p *Parser) parseGroupedExpression() ast.Expression {
	// eat '('
	p.nextToken()
	expr := p.parseExpression(LOWEST)
	// eat ')'
	if !p.matchPeek(token.RPAREN) {
		// TODO: Is this good error handling?
		return nil
	}
	return expr
}

//		IF LPARAN <expr> RPARAN LBRACE <expr>+ RBRACE (ELSE LBRACE <expr>+ RBRACE)?
func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.currToken}

	// if (...)
	p.nextToken() // eat the IF
	if !p.match(token.LPAREN) {
		return nil
	}
	expr.Condition = p.parseExpression(LOWEST)
	p.nextToken() // eat prev token
	// eat ')'
	if !p.match(token.RPAREN) {
		return nil
	}

	// eat '{'
	if !p.match(token.LBRACE) {
		return nil
	}
	expr.IfArm = p.parseBlockStatement()
	// eat '}'
	if !p.match(token.RBRACE) {
		return nil
	}

	if p.currTokenIs(token.ELSE) {
		p.nextToken() // eat ELSE
		// eat '{'
		if !p.match(token.LBRACE) {
			return nil
		}
		expr.ElseArm = p.parseBlockStatement()
		// eat '}'
		if !p.match(token.RBRACE) {
			return nil
		}
	}

	return expr
}

//    FUNCTION <params> LBRACE <expr>+ RBRACE
func (p *Parser) parseFunctionLiteral() ast.Expression {
	fun := &ast.FunctionLiteral{Token: p.currToken}

	p.nextToken() // eat 'fn' token
	// eat '('
	if !p.match(token.LPAREN) {
		return nil // '(' is mandatory
	}

	fun.Parameters = p.parseFunctionParameters()

	// eat '{'
	if !p.match(token.LBRACE) {
		return nil
	}
	fun.Body = p.parseBlockStatement()
	// parseBlockStatement will match the '}'
	//if !p.matchPeek(token.RBRACE) {
	//	return nil
	//}

	return fun
}
//  <params> := LPARAN (ID (COMMA ID)+)? RPARAN
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var ids []*ast.Identifier

	// empty params
	if p.currTokenIs(token.RPAREN) {
		p.nextToken() // eat ')'
		return ids
	}

	// first param
	id := &ast.Identifier{Token: p.currToken, Value: p.currToken.Lexeme}
	ids = append(ids, id)
	p.nextToken() // eat ID

	// loop while we see COMMA
	for p.currTokenIs(token.COMMA) {
		p.nextToken() // eat COMMA
		id := &ast.Identifier{Token: p.currToken, Value: p.currToken.Lexeme}
		ids = append(ids, id)
		p.nextToken() // eat ID
	}

	// eat ')'
	if !p.match(token.RPAREN) {
		return nil
	}

	return ids
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currToken}
	block.Statements = []ast.Statement{}

	for !p.currTokenIs(token.RBRACE) && !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

//         | BANG <expr>
//         | MINUS <expr>
func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Lexeme,
	}
	p.nextToken()
	expr.Right = p.parseExpression(PREFIX) // recursive call
	return expr
}

//		   | <expr> OP <expr>
// where OP in (+, -, *, /, ==, !=, >, <)
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Lexeme,
		Left:     left,
   	}
	// The precedence of the binary operator
	precedence := p.currPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence) // recursive call
	return expr
}
