package lexer

import (
	"bufio"
	"fmt"
)

type Lexer struct {
	reader *bufio.Reader
}

func NewLexer(reader *bufio.Reader) *Lexer {
	return &Lexer{reader: reader}
}

func (l *Lexer) skipWhitespace() {
	for {
		peek, err := l.reader.Peek(1)
		if err != nil {
			break
		}
		if peek[0] != ' ' && peek[0] != '\t' && peek[0] != '\n' && peek[0] != '\r' {
			break
		}
		l.reader.ReadByte()
	}
}

func Tokenize(reader *bufio.Reader) ([]Token, error) {
	l := NewLexer(reader)
	tokens := make([]Token, 0)
	for {
		l.skipWhitespace()
		peek, err := l.reader.Peek(1)
		if err != nil {
			break
		}
		switch peek[0] {
		case '{':
			tokens = append(tokens, newToken(OPEN_BRACE, "{"))
			l.reader.ReadByte()
		case '}':
			tokens = append(tokens, newToken(CLOSE_BRACE, "}"))
			l.reader.ReadByte()
		default:
			return nil, fmt.Errorf("unexpected character: %c", peek[0])
		}
	}

	tokens = append(tokens, newToken(EOF, ""))
	return tokens, nil
}
