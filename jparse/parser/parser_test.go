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

func TestParseValidEmptyObject(t *testing.T) {
	tokens := tokenize(t, `{}`)
	if err := Parse(tokens); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestParseValidSimpleObject(t *testing.T) {
	tokens := tokenize(t, `{"key": "value"}`)
	if err := Parse(tokens); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestParseValidMultiLineObject(t *testing.T) {
	tokens := tokenize(t, `{
  "key": "value",
  "key2": "value"
}`)
	if err := Parse(tokens); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestParseValidMixedTypes(t *testing.T) {
	tokens := tokenize(t, `{
  "key1": true,
  "key2": false,
  "key3": null,
  "key4": "value",
  "key5": 101
}`)
	if err := Parse(tokens); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestParseValidNested(t *testing.T) {
	tokens := tokenize(t, `{
  "key": "value",
  "key-n": 101,
  "key-o": {},
  "key-l": []
}`)
	if err := Parse(tokens); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestParseValidDeeplyNested(t *testing.T) {
	tokens := tokenize(t, `{
  "key": "value",
  "key-n": 101,
  "key-o": {
    "inner key": "inner value"
  },
  "key-l": ["list value"]
}`)
	if err := Parse(tokens); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestParseInvalidEmpty(t *testing.T) {
	tokens := tokenize(t, ``)
	if err := Parse(tokens); err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestParseInvalidTrailingComma(t *testing.T) {
	tokens := tokenize(t, `{"key": "value",}`)
	if err := Parse(tokens); err == nil {
		t.Fatal("expected error for trailing comma, got nil")
	}
}
