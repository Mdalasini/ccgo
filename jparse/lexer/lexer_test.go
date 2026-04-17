package lexer

import (
	"bufio"
	"strings"
	"testing"
)

func TestTokenizeInvalidBareKey(t *testing.T) {
	input := "{foo}"
	reader := bufio.NewReader(strings.NewReader(input))
	tokens, err := Tokenize(reader)
	if err == nil {
		t.Fatalf("expected error, got nil (tokens: %+v)", tokens)
	}
}

func TestTokenizeValidObject(t *testing.T) {
	input := "{\"foo\": \"bar\", \"baz\": \"qux\"}"
	reader := bufio.NewReader(strings.NewReader(input))
	tokens, err := Tokenize(reader)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []Token{
		{Kind: OPEN_BRACE, Value: "{"},
		{Kind: STRING, Value: "foo"},
		{Kind: COLON, Value: ":"},
		{Kind: STRING, Value: "bar"},
		{Kind: COMMA, Value: ","},
		{Kind: STRING, Value: "baz"},
		{Kind: COLON, Value: ":"},
		{Kind: STRING, Value: "qux"},
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

func TestTokenizeBoolean(t *testing.T) {
	input := "{\"foo\": true}"
	reader := bufio.NewReader(strings.NewReader(input))
	tokens, err := Tokenize(reader)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []Token{
		{Kind: OPEN_BRACE, Value: "{"},
		{Kind: STRING, Value: "foo"},
		{Kind: COLON, Value: ":"},
		{Kind: BOOLEAN, Value: "true"},
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

func TestTokenizeNull(t *testing.T) {
	input := "{\"foo\": null}"
	reader := bufio.NewReader(strings.NewReader(input))
	tokens, err := Tokenize(reader)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []Token{
		{Kind: OPEN_BRACE, Value: "{"},
		{Kind: STRING, Value: "foo"},
		{Kind: COLON, Value: ":"},
		{Kind: NULL, Value: "null"},
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

func TestTokenizeNumber(t *testing.T) {
	input := "{\"foo\": 42.0}"
	reader := bufio.NewReader(strings.NewReader(input))
	tokens, err := Tokenize(reader)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []Token{
		{Kind: OPEN_BRACE, Value: "{"},
		{Kind: STRING, Value: "foo"},
		{Kind: COLON, Value: ":"},
		{Kind: NUMBER, Value: "42.0"},
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
