package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestFreqMapLesMiserables(t *testing.T) {
	f, err := os.Open("tests/test.txt")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer f.Close()

	freq, err := FreqMap(f)
	if err != nil {
		t.Fatalf("FreqMap returned error: %v", err)
	}

	if got := freq['X']; got != 333 {
		t.Errorf("X count = %d, want 333", got)
	}

	if got := freq['t']; got != 223000 {
		t.Errorf("t count = %d, want 223000", got)
	}
}

func TestBuildCodeTable(t *testing.T) {
	freq := map[rune]int{
		'C': 32,
		'D': 42,
		'E': 120,
		'K': 7,
		'L': 42,
		'M': 24,
		'U': 37,
		'Z': 2,
	}
	want := map[rune]string{
		'C': "1110",
		'D': "101",
		'E': "0",
		'K': "111101",
		'L': "110",
		'M': "11111",
		'U': "100",
		'Z': "111100",
	}

	root := BuildHuffTree(freq)
	got := BuildCodeTable(root)

	for r, wantCode := range want {
		if got[r] != wantCode {
			t.Errorf("code for %q: got %q, want %q", r, got[r], wantCode)
		}
	}
}

func TestBuildCodeTableSingle(t *testing.T) {
	root := BuildHuffTree(map[rune]int{'x': 10})
	codes := BuildCodeTable(root)
	if codes['x'] != "0" {
		t.Errorf("single char code: got %q, want %q", codes['x'], "0")
	}
}

func TestEncodeDecodeRoundtrip(t *testing.T) {
	tests := []string{
		"hello world",
		"a",
		"ab",
		"banana bandana",
		"this is a test of the emergency broadcast system",
		"aaaaaabbbbbccccdde",
		"Mississippi",
		"go gopher go",
	}

	for _, input := range tests {
		var buf bytes.Buffer
		if err := Encode(strings.NewReader(input), &buf); err != nil {
			t.Fatalf("Encode(%q): %v", input, err)
		}

		var out bytes.Buffer
		if err := Decode(&buf, &out); err != nil {
			t.Fatalf("Decode(%q): %v", input, err)
		}

		if got := out.String(); got != input {
			t.Errorf("roundtrip mismatch:\n  input: %q\n  got:   %q", input, got)
		}
	}
}

func TestEncodeDecodeEmpty(t *testing.T) {
	var buf bytes.Buffer
	if err := Encode(strings.NewReader(""), &buf); err != nil {
		t.Fatalf("Encode empty: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %d bytes", buf.Len())
	}

	var out bytes.Buffer
	if err := Decode(&buf, &out); err != nil {
		t.Fatalf("Decode empty: %v", err)
	}
	if out.Len() != 0 {
		t.Errorf("expected empty output, got %q", out.String())
	}
}

func TestEncodeDecodeLargeFile(t *testing.T) {
	f, err := os.Open("tests/test.txt")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer f.Close()

	original, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	var compressed bytes.Buffer
	if err := Encode(bytes.NewReader(original), &compressed); err != nil {
		t.Fatalf("Encode: %v", err)
	}

	t.Logf("original size: %d, compressed size: %d, ratio: %.2f%%",
		len(original), compressed.Len(),
		float64(compressed.Len())/float64(len(original))*100)

	var decompressed bytes.Buffer
	if err := Decode(&compressed, &decompressed); err != nil {
		t.Fatalf("Decode: %v", err)
	}

	if !bytes.Equal(decompressed.Bytes(), original) {
		t.Errorf("roundtrip mismatch on large file")
	}
}

func TestBuildCodeLengths(t *testing.T) {
	freq := map[rune]int{'a': 5, 'b': 2, 'c': 1, 'd': 1}
	root := BuildHuffTree(freq)
	lengths := BuildCodeLengths(root)

	codes := BuildCanonicalCodes(lengths)

	used := make(map[string]bool)
	for _, code := range codes {
		for i := 1; i <= len(code); i++ {
			prefix := code[:i]
			if used[prefix] {
				t.Errorf("prefix %q is not unique (violates prefix code property)", prefix)
			}
		}
		used[code] = true
	}
}
