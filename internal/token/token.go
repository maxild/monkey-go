package token

type Type string

type Token struct {
	Type   Type
	Lexeme string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	// Identifiers + literals
	IDENT = "IDENT"
	INT = "INT"

	// Operators
	ASSIGN = "="
	PLUS = "+"

	// Punctuation
	COMMA = ","
	SEMICOLON = ";"

	LPARAN = "("
	RPARAN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET = "LET"
)

var keywords = map[string]Type {
	"fn": FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) Type {
	if t, ok := keywords[ident]; ok {
		return t
	}
	return IDENT
}
