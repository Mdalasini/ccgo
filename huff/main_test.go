package main

import (
	"os"
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
