package parser

import (
	"bufio"
	"strings"
	"testing"

	"github.com/mdalasini/ccgo/jparse/lexer"
)

func tokenize(t *testing.T, input string) []lexer.Token {
	t.Helper()
	tokens, err := lexer.Tokenize(bufio.NewReader(strings.NewReader(input)))
	if err != nil {
		t.Fatalf("lexer error for %q: %v", input, err)
	}
	return tokens
}

func TestParseValidObject(t *testing.T) {
	tokens := tokenize(t, `{"key": true, "key2": null, "key3": 42.0}`)
	if err := Parse(tokens); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestParseTrailingComma(t *testing.T) {
	tokens := tokenize(t, `{"key": "value",}`)
	if err := Parse(tokens); err == nil {
		t.Fatal("expected error for trailing comma, got nil")
	}
}

func TestParseMissingColon(t *testing.T) {
	tokens := tokenize(t, `{"key" "value"}`)
	if err := Parse(tokens); err == nil {
		t.Fatal("expected error for missing colon, got nil")
	}
}

func TestParseInvalidKey(t *testing.T) {
	tokens := tokenize(t, `{42.0: "value"}`)
	if err := Parse(tokens); err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
}

func TestParseArray(t *testing.T) {
	tokens := tokenize(t, `{"foo": [1, 2, 3]}`)
	if err := Parse(tokens); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
