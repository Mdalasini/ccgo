package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
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

func (l *Lexer) consumeString() (string, error) {
	l.reader.ReadByte()
	var value strings.Builder
	for {
		peek, err := l.reader.Peek(1)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return "", fmt.Errorf("expected closing quote for string")
			}
			return "", err
		}
		switch peek[0] {
		case '"':
			l.reader.ReadByte()
			return value.String(), nil
		case '\\':
			l.reader.ReadByte()
			peek, err := l.reader.Peek(1)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return "", fmt.Errorf("escaped character at end of string")
				}
				return "", err
			}
			switch peek[0] {
			case 'u':
				l.reader.ReadByte()
				peek, err := l.reader.Peek(4)
				if err != nil {
					if errors.Is(err, io.EOF) {
						return "", fmt.Errorf("expected 4 hex digits after \\u")
					}
					return "", err
				}

				// Parse the 4 hex digits as a Unicode code point
				n, err := strconv.ParseUint(string(peek), 16, 32)
				if err != nil {
					return "", fmt.Errorf("invalid hex digit: %s", string(peek))
				}
				value.WriteRune(rune(n))

				for range 4 {
					l.reader.ReadByte()
				}
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				value.WriteString(string(peek[0]))
				l.reader.ReadByte()
			default:
				return "", fmt.Errorf("invalid escape character: %c", peek[0])
			}
		default:
			value.WriteString(string(peek[0]))
			l.reader.ReadByte()
		}
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
		case '"':
			value, err := l.consumeString()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, newToken(STRING, value))
		case ':':
			tokens = append(tokens, newToken(COLON, ":"))
			l.reader.ReadByte()
		case ',':
			tokens = append(tokens, newToken(COMMA, ","))
			l.reader.ReadByte()
		default:
			return nil, fmt.Errorf("unexpected character: %c", peek[0])
		}
	}

	tokens = append(tokens, newToken(EOF, ""))
	return tokens, nil
}

func parseUnicodeEscape()
