package lexer

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
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
		if unicode.IsSpace(rune(peek[0])) {
			l.reader.ReadByte()
			continue
		}
		break
	}
}

func (l *Lexer) consumeNumber() (string, error) {
	var value strings.Builder
	inFloat := false
	for {
		peek, err := l.reader.Peek(1)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return value.String(), nil
			}
			return "", err
		}
		if peek[0] == '.' {
			if inFloat {
				return "", fmt.Errorf("multiple dots in number")
			}
			inFloat = true
			peek, err = l.reader.Peek(2)
			if err != nil {
				return "", fmt.Errorf("expected end of number after dot")
			}
			if !unicode.IsDigit(rune(peek[1])) {
				return "", fmt.Errorf("expected digit after dot in number")
			}
			value.WriteByte('.')
			l.reader.ReadByte()
			continue
		}
		if !unicode.IsDigit(rune(peek[0])) {
			break
		}
		value.WriteByte(peek[0])
		l.reader.ReadByte()
	}
	return value.String(), nil
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

func (l *Lexer) peekUntil(stoppers []byte, includeSpace bool) []byte {
	for size := 1; ; size++ {
		peek, err := l.reader.Peek(size)
		if err != nil {
			// ran out of bytes before finding a stopper
			return peek
		}
		// check if the last byte is a stopper
		last := peek[size-1]
		if includeSpace {
			if unicode.IsSpace(rune(last)) || bytes.Contains(stoppers, []byte{last}) {
				return peek[:size-1] // don't include the stopper in the result
			}
		} else {
			if bytes.Contains(stoppers, []byte{last}) {
				return peek[:size-1] // don't include the stopper in the result
			}
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
		switch {
		case peek[0] == '{':
			tokens = append(tokens, newToken(OPEN_BRACE, "{"))
			l.reader.ReadByte()
		case peek[0] == '}':
			tokens = append(tokens, newToken(CLOSE_BRACE, "}"))
			l.reader.ReadByte()
		case peek[0] == '"':
			value, err := l.consumeString()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, newToken(STRING, value))
		case peek[0] == ':':
			tokens = append(tokens, newToken(COLON, ":"))
			l.reader.ReadByte()
		case peek[0] == ',':
			tokens = append(tokens, newToken(COMMA, ","))
			l.reader.ReadByte()
		case unicode.IsDigit(rune(peek[0])):
			value, err := l.consumeNumber()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, newToken(NUMBER, value))
		default:
			stoppers := []byte{',', '}', ':'}
			rem := l.peekUntil(stoppers, true)
			switch {
			case string(rem) == "true":
				tokens = append(tokens, newToken(BOOLEAN, "true"))
				for range len(rem) {
					l.reader.ReadByte()
				}
			case string(rem) == "false":
				tokens = append(tokens, newToken(BOOLEAN, "false"))
				for range len(rem) {
					l.reader.ReadByte()
				}
			case string(rem) == "null":
				tokens = append(tokens, newToken(NULL, "null"))
				for range len(rem) {
					l.reader.ReadByte()
				}
			default:
				return nil, fmt.Errorf("string value must be double-quoted: %s", rem)
			}
		}
	}

	tokens = append(tokens, newToken(EOF, ""))
	return tokens, nil
}
