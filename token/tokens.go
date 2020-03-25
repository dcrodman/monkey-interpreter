package token

import (
	"regexp"
)

// TokenType is an enum type to represent the possible categories of tokens supported
// by Monkey during the lexing and parsing phases.
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

// Keywords supported by the
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

// Determines whether or not `literal` is a syntactically valid identifier.
func IsValidIdentifier(literal string) bool {
	ok, _ := regexp.Match(identifierRegex, []byte(literal))
	return ok
}

// Determines whether or not `literal` is a syntactically valid integer.
func IsValidInteger(literal string) bool {
	ok, _ := regexp.Match(integerRegex, []byte(literal))
	return ok
}

// Returns a bool indicating whether or not `literal` is a valid keyword and, if so,
// also returns the TokenType corresponding to that keyword.
func GetKeywordType(literal string) (bool, TokenType) {
	if keywordType, ok := keywords[literal]; ok {
		return true, keywordType
	}
	return false, ""
}
