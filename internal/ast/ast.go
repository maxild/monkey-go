package ast

// TODO: We need sum types in Rust, F# (GO does not support sum types)

import (
	"bytes"
	"github.com/maxild/monkey/internal/token"
)

type Node interface {
	TokenLiteral() string 	// only used for debugging and testing
	String() string			// make every Node a 'Stringer' for test/debugging purposes
}

// Fields and methods of anonymous (embedded) field are called promoted.
// They behave like regular fields but can’t be used in struct literals:

// An interface type may specify methods explicitly through method specifications,
// or it may embed methods of other interfaces through interface type names.
type Statement interface {
	Node // embedded TokenLiteral() method
	statementNode()
}

type Expression interface {
	Node // embedded TokenLiteral() method
	expressionNode()
}

// A series of statements is the root Node in our AST
type Program struct {
	Statements []Statement
}

// Program implements Node this way in GO
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type LetStatement struct {
	Token token.Token 	// The token.LET token (kind)
	Name *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Lexeme }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	// TODO: How can Value be nil?
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token token.Token 	// The token.RETURN token (kind)
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Lexeme }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

// Make expression a statement (it is really not a statement). We only make a statement in order
// to be able to add ot to the top-level sequence of statements in the root AST node. This way a program like
//
// let x = 10;
// x + 67;
//
// can be interpreted as the value 77.
type ExpressionStatement struct {
	Token token.Token	// The first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Lexeme }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type Identifier struct {
	Token token.Token 	// The token.IDENT token (kind)
	Value string		// TODO: Wouldn't Name be a better term?
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Lexeme }
func (i *Identifier) String() string { return i.Value }

type IntegerLiteral struct {
	Token token.Token 	// The token.INT token (kind)
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Lexeme }
func (il *IntegerLiteral) String() string { return il.Token.Lexeme }

// Aka UnaryExpression
type PrefixExpression struct {
	Token token.Token	// The prefix token kind (e.g. ! or -)
	Operator string
	Right Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Lexeme }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	// wrap parenthesis around to see connection
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// Aka BinaryExpression
type InfixExpression struct {
	Token token.Token	// The infix token (kind) (e.g. +, -, *, /, ==, !=, > or <)
	Operator string
	Left Expression
	Right Expression
}

func (be *InfixExpression) expressionNode() {}
func (be *InfixExpression) TokenLiteral() string { return be.Token.Lexeme }
func (be *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(be.Left.String())
	out.WriteString(" " + be.Operator + " ")
	out.WriteString(be.Right.String())
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	Token token.Token // The TRUE or FALSE token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string { return b.Token.Lexeme }
func (b *Boolean) String() string { return b.Token.Lexeme }

// Aka ConditionalExpression (catamorphism, ternary operator)
type IfExpression struct {
	Token token.Token		// The IF token
	Condition Expression
	IfArm *BlockStatement 	// each arm can have many statements
	ElseArm *BlockStatement	// do.
}

func (ie *IfExpression) expressionNode() {}
func (ie* IfExpression) TokenLiteral() string { return ie.Token.Lexeme }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(ie.IfArm.String())
	if ie.ElseArm != nil {
		out.WriteString("else")
		out.WriteString(ie.ElseArm.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token token.Token		// The '{' token
	Statements []Statement
}

func (bs *BlockStatement) expressionNode() {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Lexeme }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	// TODO: What about surrounding { ... }
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
