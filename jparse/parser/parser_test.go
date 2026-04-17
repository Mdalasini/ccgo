package parser

import (
	"testing"

	"github.com/mdalasini/ccgo/jparse/lexer"
)

func TestParseValidObject(t *testing.T) {
	tokens := []lexer.Token{
		{Kind: lexer.OPEN_BRACE, Value: "{"},
		{Kind: lexer.STRING, Value: "key"},
		{Kind: lexer.COLON, Value: ":"},
		{Kind: lexer.BOOLEAN, Value: "true"},
		{Kind: lexer.COMMA, Value: ","},
		{Kind: lexer.STRING, Value: "key2"},
		{Kind: lexer.COLON, Value: ":"},
		{Kind: lexer.NULL, Value: "null"},
		{Kind: lexer.COMMA, Value: ","},
		{Kind: lexer.STRING, Value: "key3"},
		{Kind: lexer.COLON, Value: ":"},
		{Kind: lexer.NUMBER, Value: "42.0"},
		{Kind: lexer.CLOSE_BRACE, Value: "}"},
		{Kind: lexer.EOF},
	}

	if err := Parse(tokens); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestParseTrailingComma(t *testing.T) {
	tokens := []lexer.Token{
		{Kind: lexer.OPEN_BRACE, Value: "{"},
		{Kind: lexer.STRING, Value: "key"},
		{Kind: lexer.COLON, Value: ":"},
		{Kind: lexer.STRING, Value: "value"},
		{Kind: lexer.COMMA, Value: ","},
		{Kind: lexer.CLOSE_BRACE, Value: "}"},
		{Kind: lexer.EOF},
	}

	if err := Parse(tokens); err == nil {
		t.Fatal("expected error for trailing comma, got nil")
	}
}

func TestParseMissingColon(t *testing.T) {
	tokens := []lexer.Token{
		{Kind: lexer.OPEN_BRACE, Value: "{"},
		{Kind: lexer.STRING, Value: "key"},
		{Kind: lexer.STRING, Value: "value"},
		{Kind: lexer.CLOSE_BRACE, Value: "}"},
		{Kind: lexer.EOF},
	}

	if err := Parse(tokens); err == nil {
		t.Fatal("expected error for missing colon, got nil")
	}
}

func TestParseInvalidKey(t *testing.T) {
	tokens := []lexer.Token{
		{Kind: lexer.OPEN_BRACE, Value: "{"},
		{Kind: lexer.NUMBER, Value: "42.0"},
		{Kind: lexer.COLON, Value: ":"},
		{Kind: lexer.STRING, Value: "value"},
		{Kind: lexer.CLOSE_BRACE, Value: "}"},
		{Kind: lexer.EOF},
	}

	if err := Parse(tokens); err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
}
