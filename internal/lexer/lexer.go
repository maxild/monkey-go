package lexer

import "github.com/maxild/monkey/internal/token"

type Lexer struct {
	input string
	position int		// current position in input (points to current char)
	readPosition int	// current reading position in input (after current char)
	ch byte				// TODO: Implement Unicode support (byte to rune)
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // position := 0, readPosition := 1, ch := 0 or input[0]
	return l
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// read the next character
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// signal EOF
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	// we need to update the position/readPosition (even when we reach EOF)
	// because other procedures in the lexer uses them to slice out values
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	// Important to scan for keywords, operators, punctuation before calling anything an identifier
	case '=':
		next := l.peekChar()
		if next == '=' {
			prev := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Lexeme: string(prev) + string(next)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		next := l.peekChar()
		if next == '=' {
			prev := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Lexeme: string(prev) + string(next)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		//tok = token.Token{Type: token.EOF, Lexeme: ""}
		tok.Type = token.EOF
		tok.Lexeme = ""
	default:
		if isLetter(l.ch) {
			lexeme := l.readIdentifier()
			tok.Lexeme = lexeme
			tok.Type = token.LookupIdent(lexeme)
			return tok // do not call readChar, readIdentifier has done it already
		} else if isNumber(l.ch) {
			tok.Lexeme = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
			//tok.Type = token.ILLEGAL
			//tok.Lexeme = string(l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(t token.Type, ch byte) token.Token {
	return token.Token{
		Type:   t,
		Lexeme: string(ch),
	}
}

// [a-zA-Z_]*
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// [0-9]+ is a very simplified regex for defining numbers
// We are missing
//  - float
//  - decimal
//  - octal
//  - hex
//  - scientific notation
//  - negative numbers
func (l *Lexer) readNumber() string {
	position := l.position
	for isNumber(l.ch) {
		l.readChar()
	}
	return l.input[position: l.position]
}

// [a-zA-Z_]
func isLetter(ch byte) bool {
	return 'A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' || ch == '_'
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}
