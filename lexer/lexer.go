package lexer

import (
	"fmt"
	"io/ioutil"
	"monkey-interpreter/token"
	"os"
	"unicode"
)

type Lexer struct {
	codeInput []rune

	currentPosition int
	currentLine     int
}

func New(input string) *Lexer {
	return &Lexer{codeInput: []rune(input), currentLine: 1}
}

func NewForFile(filename string) (*Lexer, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Lexer: %s", err)
	}
	defer file.Close()

	sourceFileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to parse source file: %s", err)
	}

	return &Lexer{
		codeInput:   []rune(string(sourceFileContents)),
		currentLine: 1,
	}, nil
}

func (l *Lexer) NextToken() token.Token {
	nextRune := l.getNextRune()
	literal := string(nextRune)
	var tokenType token.TokenType

	switch nextRune {
	case '>':
		tokenType = token.GRT
	case '<':
		tokenType = token.LES
	case '=':
		if l.peekNextRune() == '=' {
			l.moveToNextPosition()
			literal, tokenType = token.EQ, token.EQ
		} else {
			tokenType = token.ASSIGN
		}
	case '!':
		if l.peekNextRune() == '=' {
			l.moveToNextPosition()
			literal, tokenType = token.NOT_EQ, token.NOT_EQ
		} else {
			tokenType = token.BANG
		}
	case '+':
		tokenType = token.PLUS
	case '-':
		tokenType = token.MINUS
	case '/':
		tokenType = token.SLASH
	case '*':
		tokenType = token.ASTERISK
	case '{':
		tokenType = token.LBRACKET
	case '}':
		tokenType = token.RBRACKET
	case '(':
		tokenType = token.LPAREN
	case ')':
		tokenType = token.RPAREN
	case ';':
		tokenType = token.SEMICOLON
	case ',':
		tokenType = token.COMMA
	case '"':
		tokenType = token.STRING
		literal = l.readString()
	case 0:
		literal, tokenType = token.EOF, token.EOF
	}

	// Assume the token is a literal of some kind if it doesn't match any of
	// the special characters. readLiteral() will handle advancing the parser.
	if tokenType == "" {
		literal, tokenType = l.readLiteral()
	} else {
		l.moveToNextPosition()
	}

	// TODO: Something useful with the ILLEGAL token type.

	return token.Token{Type: tokenType, Value: literal}
}

func (l *Lexer) getNextRune() rune {
	nextRune := l.peekCurrentRune()

	// Consume whitespace until we find the next significant character.
	for unicode.IsSpace(nextRune) {
		l.moveToNextPosition()
		nextRune = l.peekCurrentRune()
	}

	return nextRune
}

func (l *Lexer) peekCurrentRune() rune {
	return l.readRune(l.currentPosition)
}

func (l *Lexer) peekNextRune() rune {
	return l.readRune(l.currentPosition + 1)
}

func (l *Lexer) readRune(position int) rune {
	if position >= len(l.codeInput) {
		return 0
	}
	return l.codeInput[position]
}

func (l *Lexer) moveToNextPosition() {
	l.currentPosition += 1

	if string(l.peekCurrentRune()) == "\n" {
		l.currentLine += 1
	}
}

// Note: This method will advance the parser's current position in the file up until
// it runs off the end of the literal (for instance a space or semicolon).
func (l *Lexer) readLiteral() (string, token.TokenType) {
	var literal string

	for r := l.peekCurrentRune(); unicode.IsDigit(r) || unicode.IsLetter(r); r = l.peekCurrentRune() {
		literal += string(r)
		l.moveToNextPosition()
	}

	if unicode.IsDigit(rune(literal[0])) {
		if !token.IsValidInteger(literal) {
			return literal, token.ILLEGAL
		}
		return literal, token.INT
	}

	// If it's not a number then we assume the token is a keyword or identifier.
	if !token.IsValidIdentifier(literal) {
		return literal, token.ILLEGAL
	}

	if ok, keywordType := token.GetKeywordType(literal); ok {
		return literal, keywordType
	}
	return literal, token.IDENTIFIER
}

func (l *Lexer) readString() string {
	position := l.currentPosition + 1

	for {
		l.moveToNextPosition()
		r := l.getNextRune()
		if r == '"' || r == 0 {
			break
		}
	}

	return string(l.codeInput[position:l.currentPosition])
}
