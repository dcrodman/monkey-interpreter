package lexer

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER"
	INT        = "INTEGER"
	LET        = "LET"
	FUNCTION   = "FUNCTION"

	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"

	LBRACKET  = "{"
	RBRACKET  = "}"
	LPAREN    = "("
	RPAREN    = ")"
	SEMICOLON = ";"
	COMMA     = ","
)

const (
	identifierRegex = `^[a-zA-Z|_][a-zA-Z|\d|_]*\b$`
	integerRegex    = `^\d+$`
)

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

type Token struct {
	Type  TokenType
	Value string
}
