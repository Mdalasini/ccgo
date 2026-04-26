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
