package ast

import "github.com/maxild/monkey/internal/token"

type Node interface {
	TokenLiteral() string // only used for debugging and testing
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

type LetStatement struct {
	Token token.Token 	// The token.LET token (kind)
	Name *Identifier
	Value Expression
}

func (l *LetStatement) statementNode() {}
func (l *LetStatement) TokenLiteral() string { return l.Token.Lexeme }


type Identifier struct {
	Token token.Token // The token.IDENT token (kind)
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Lexeme }
