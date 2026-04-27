package main

import (
	"bytes"
	"os"
	"testing"
)

func TestCutField2(t *testing.T) {
	f, err := os.Open("tests/sample.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	if err := cut(f, &buf, 2); err != nil {
		t.Fatal(err)
	}

	// Expected output matches: cut -f2 tests/sample.tsv
	want := "f1\n1\n6\n11\n16\n21\n"
	got := buf.String()
	if got != want {
		t.Errorf("cut -f2 output mismatch:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestCutField9OutOfRange(t *testing.T) {
	f, err := os.Open("tests/sample.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	if err := cut(f, &buf, 9); err != nil {
		t.Fatal(err)
	}

	// Field 9 is out of range for all rows — expect empty output
	if got := buf.String(); got != "" {
		t.Errorf("expected empty output for out-of-range field, got: %q", got)
	}
}
