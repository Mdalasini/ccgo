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

func (l *Lexer) consumeNumber(isNegative bool) (string, error) {
	var value strings.Builder
	if isNegative {
		value.WriteByte('-')
	}

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

func (l *Lexer) readUntilDelimFunc(isDelimFunc func(byte) bool) ([]byte, error) {
	var result bytes.Buffer

	for {
		b, err := l.reader.ReadByte()
		if err != nil {
			if errors.Is(io.EOF, err) && result.Len() > 0 {
				return result.Bytes(), nil
			}
			return result.Bytes(), err
		}

		if isDelimFunc(b) {
			l.reader.UnreadByte() // place back delimiter
			return result.Bytes(), nil
		}

		result.WriteByte(b)
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
		case peek[0] == '[':
			tokens = append(tokens, newToken(OPEN_BRACKET, "["))
			l.reader.ReadByte()
		case peek[0] == ']':
			tokens = append(tokens, newToken(CLOSE_BRACKET, "]"))
			l.reader.ReadByte()
		case peek[0] == '-':
			l.reader.ReadByte()
			peek, err := l.reader.Peek(1)
			if err != nil {
				return nil, fmt.Errorf("minus must be followed by a digit")
			}
			if !unicode.IsDigit(rune(peek[0])) {
				return nil, fmt.Errorf("minus must be followed by a digit")
			}
			isNegative := true
			value, err := l.consumeNumber(isNegative)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, newToken(NUMBER, value))
		case unicode.IsDigit(rune(peek[0])):
			isNegative := false
			value, err := l.consumeNumber(isNegative)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, newToken(NUMBER, value))
		default:
			isDelim := func(b byte) bool {
				return unicode.IsSpace(rune(b)) || bytes.Contains([]byte{',', '}', ':'}, []byte{b})
			}
			rem, err := l.readUntilDelimFunc(isDelim)
			if err != nil {
				return nil, err
			}
			switch {
			case string(rem) == "true":
				tokens = append(tokens, newToken(BOOLEAN, "true"))
			case string(rem) == "false":
				tokens = append(tokens, newToken(BOOLEAN, "false"))
			case string(rem) == "null":
				tokens = append(tokens, newToken(NULL, "null"))
			default:
				return nil, fmt.Errorf("string value must be double-quoted: %s", rem)
			}
		}
	}

	tokens = append(tokens, newToken(EOF, ""))
	return tokens, nil
}
