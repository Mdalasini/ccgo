package lexer

import (
	"bufio"
	"strings"
	"testing"
)

func tokenize(t *testing.T, input string) ([]Token, bool) {
	t.Helper()
	tokens, err := Tokenize(bufio.NewReader(strings.NewReader(input)))
	return tokens, err == nil
}

func TestTokenizeValidEmptyObject(t *testing.T) {
	input := `{}`
	tokens, ok := tokenize(t, input)
	if !ok {
		t.Fatal("expected no error for empty object")
	}
	expected := []Token{
		{Kind: OPEN_BRACE, Value: "{"},
		{Kind: CLOSE_BRACE, Value: "}"},
		{Kind: EOF, Value: ""},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		if tok != expected[i] {
			t.Errorf("token %d: expected %+v, got %+v", i, expected[i], tok)
		}
	}
}

func TestTokenizeValidSimpleObject(t *testing.T) {
	input := `{"key": "value"}`
	_, ok := tokenize(t, input)
	if !ok {
		t.Fatal("expected no error for simple object")
	}
}

func TestTokenizeValidMultiLineObject(t *testing.T) {
	input := `{
  "key": "value",
  "key2": "value"
}`
	_, ok := tokenize(t, input)
	if !ok {
		t.Fatal("expected no error for multi-line object")
	}
}

func TestTokenizeValidMixedTypes(t *testing.T) {
	input := `{
  "key1": true,
  "key2": false,
  "key3": null,
  "key4": "value",
  "key5": 101
}`
	_, ok := tokenize(t, input)
	if !ok {
		t.Fatal("expected no error for mixed types")
	}
}

func TestTokenizeValidNested(t *testing.T) {
	input := `{
  "key": "value",
  "key-n": 101,
  "key-o": {},
  "key-l": []
}`
	_, ok := tokenize(t, input)
	if !ok {
		t.Fatal("expected no error for nested objects")
	}
}

func TestTokenizeValidDeeplyNested(t *testing.T) {
	input := `{
  "key": "value",
  "key-n": 101,
  "key-o": {
    "inner key": "inner value"
  },
  "key-l": ["list value"]
}`
	_, ok := tokenize(t, input)
	if !ok {
		t.Fatal("expected no error for deeply nested objects")
	}
}

func TestTokenizeInvalidBareKey(t *testing.T) {
	input := `{
  "key": "value",
  key2: "value"
}`
	_, ok := tokenize(t, input)
	if ok {
		t.Fatal("expected error for bare key")
	}
}

func TestTokenizeInvalidCapitalizedBoolean(t *testing.T) {
	input := `{
  "key1": true,
  "key2": False,
  "key3": null,
  "key4": "value",
  "key5": 101
}`
	_, ok := tokenize(t, input)
	if ok {
		t.Fatal("expected error for capitalized boolean")
	}
}

func TestTokenizeInvalidSingleQuote(t *testing.T) {
	input := `{
  "key": "value",
  "key-n": 101,
  "key-o": {
    "inner key": "inner value"
  },
  "key-l": ['list value']
}`
	_, ok := tokenize(t, input)
	if ok {
		t.Fatal("expected error for single-quoted string")
	}
}
