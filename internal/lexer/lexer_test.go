package lexer

import (
	"github.com/maxild/monkey/internal/token"
	"testing"
)

func TestNextTokenCanBeCalledMultipleTimesAfterEof(t *testing.T) {
	input := "5;"
	l := New(input)

	tok := l.NextToken()

	if tok.Type != token.INT {
		t.Fatalf("wrong type")
	}

	tok = l.NextToken()
	if tok.Type != token.SEMICOLON {
		t.Fatalf("wrong type")
	}

	tok = l.NextToken()
	if tok.Type != token.EOF || tok.Lexeme != "" {
		t.Fatalf("wrong type")
	}
	tok = l.NextToken()
	if tok.Type != token.EOF || tok.Lexeme != "" {
		t.Fatalf("NextToken cannot be called again and again on EOF")
	}
}

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
`

	// table-based testing (a la Theory in xunit.net)`
	tests := []struct{
		expectedType token.Type
		expectedLexeme string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, test := range tests {
		tok := l.NextToken()

		// TODO: Convert to using github.com/stretchr/testify/assert
		// See https://github.com/stretchr/testify
		if tok.Type != test.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, test.expectedType, tok.Type)
		}

		if tok.Lexeme != test.expectedLexeme {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q",
				i, test.expectedLexeme, tok.Lexeme)
		}
	}
}

func TestBug(t *testing.T) {
	input := "-a*b"
	l := New(input)

	tok := l.NextToken()
	if tok.Type != token.MINUS || tok.Lexeme != "-" {
		t.Fatalf("Wanted (-, '-'), got (%s, '%s')", tok.Type, tok.Lexeme)
	}

	tok = l.NextToken()
	if tok.Type != token.IDENT || tok.Lexeme != "a" {
		t.Fatalf("Wanted (IDENT, 'a'), got (%s, '%s')", tok.Type, tok.Lexeme)
	}

	tok = l.NextToken()
	if tok.Type != token.ASTERISK || tok.Lexeme != "*" {
		t.Fatalf("Wanted (*, '*'), got (%s, '%s')", tok.Type, tok.Lexeme)
	}

	tok = l.NextToken()
	if tok.Type != token.IDENT || tok.Lexeme != "b" {
		t.Fatalf("Wanted (IDENT, 'b'), got (%s, '%s')", tok.Type, tok.Lexeme)
	}

	tok = l.NextToken()
	if tok.Type != token.EOF || tok.Lexeme != "" {
		t.Fatalf("Wanted (EOF, ''), got (%s, '%s')", tok.Type, tok.Lexeme)
	}
}

