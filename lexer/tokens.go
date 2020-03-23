package lexer

type TokenType string

const (
	IDENTIFIER = "IDENTIFIER"
	INT        = "INTEGER"
	LET        = "LET"
	FUNCTION   = "FUNCTION"
	IF         = "IF"
	ELSE       = "ELSE"
	RETURN     = "RETURN"
	TRUE       = "TRUE"
	FALSE      = "FALSE"

	GRT = ">"
	LES = "<"

	BANG     = "!"
	PLUS     = "+"
	MINUS    = "-"
	SLASH    = "/"
	ASTERISK = "*"

	LBRACKET  = "{"
	RBRACKET  = "}"
	LPAREN    = "("
	RPAREN    = ")"
	SEMICOLON = ";"
	COMMA     = ","

	EQ     = "=="
	NOT_EQ = "!="

	ASSIGN = "="

	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

const (
	identifierRegex = `^[a-zA-Z|_][a-zA-Z|\d|_]*\b$`
	integerRegex    = `^\d+$`
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

type Token struct {
	Type  TokenType
	Value string
}
