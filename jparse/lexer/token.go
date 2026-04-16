package lexer

import "fmt"

type TokenKind int

const (
	EOF TokenKind = iota
	OPEN_BRACE
	CLOSE_BRACE
)

type Token struct {
	Kind  TokenKind
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("%s: %s", t.Kind.String(), t.Value)
}

func newToken(kind TokenKind, value string) Token {
	return Token{Kind: kind, Value: value}
}

func (tk TokenKind) String() string {
	switch tk {
	case EOF:
		return "EOF"
	case OPEN_BRACE:
		return "OPEN_BRACE"
	case CLOSE_BRACE:
		return "CLOSE_BRACE"
	default:
		return "UNKNOWN"
	}
}
