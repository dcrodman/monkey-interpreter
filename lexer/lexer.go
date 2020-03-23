package lexer

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
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

func (l *Lexer) NextToken() Token {
	nextRune := l.getNextRune()
	literal := string(nextRune)
	var tokenType TokenType

	switch nextRune {
	case '>':
		tokenType = GRT
	case '<':
		tokenType = LES
	case '=':
		if l.peekNextRune() == '=' {
			l.moveToNextPosition()
			literal, tokenType = EQ, EQ
		} else {
			tokenType = ASSIGN
		}
	case '!':
		if l.peekNextRune() == '=' {
			l.moveToNextPosition()
			literal, tokenType = NOT_EQ, NOT_EQ
		} else {
			tokenType = BANG
		}
	case '+':
		tokenType = PLUS
	case '-':
		tokenType = MINUS
	case '/':
		tokenType = SLASH
	case '*':
		tokenType = ASTERISK
	case '{':
		tokenType = LBRACKET
	case '}':
		tokenType = RBRACKET
	case '(':
		tokenType = LPAREN
	case ')':
		tokenType = RPAREN
	case ';':
		tokenType = SEMICOLON
	case ',':
		tokenType = COMMA
	case 0:
		literal, tokenType = EOF, EOF
	}

	// Assume the token is a literal of some kind of it doesn't match any of
	// the special characters. readLiteral() will handle advancing the parser.
	if tokenType == "" {
		literal, tokenType = l.readLiteral()
	} else {
		l.moveToNextPosition()
	}

	// TODO: Something useful with the ILLEGAL token type.

	return Token{tokenType, literal}
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
func (l *Lexer) readLiteral() (string, TokenType) {
	var literal string

	for r := l.peekCurrentRune(); unicode.IsDigit(r) || unicode.IsLetter(r); r = l.peekCurrentRune() {
		literal += string(r)
		l.moveToNextPosition()
	}

	if unicode.IsDigit(rune(literal[0])) {
		if ok, _ := regexp.Match(integerRegex, []byte(literal)); !ok {
			return literal, ILLEGAL
		}
		return literal, INT
	}

	// If it's not a number then we assume the token is a keyword or identifier.
	if ok, _ := regexp.Match(identifierRegex, []byte(literal)); !ok {
		return literal, ILLEGAL
	}

	if keywordType, ok := keywords[literal]; ok {
		return literal, keywordType
	}
	return literal, IDENTIFIER
}
