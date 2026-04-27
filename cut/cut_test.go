package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestCutField2(t *testing.T) {
	f, err := os.Open("tests/sample.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	if err := cut(f, &buf, []int{2}, "\t"); err != nil {
		t.Fatal(err)
	}

	// Expected output matches: cut -f2 tests/sample.tsv
	want := "f1\n1\n6\n11\n16\n21\n"
	got := buf.String()
	if got != want {
		t.Errorf("cut -f2 output mismatch:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestCutFields1And2(t *testing.T) {
	f, err := os.Open("tests/sample.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	if err := cut(f, &buf, []int{1, 2}, "\t"); err != nil {
		t.Fatal(err)
	}

	// Expected output matches: cut -f1,2 tests/sample.tsv
	want := "f0\tf1\n" +
		"0\t1\n" +
		"5\t6\n" +
		"10\t11\n" +
		"15\t16\n" +
		"20\t21\n"
	got := buf.String()
	if got != want {
		t.Errorf("cut -f1,2 output mismatch:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestCutField9OutOfRange(t *testing.T) {
	f, err := os.Open("tests/sample.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	if err := cut(f, &buf, []int{9}, "\t"); err != nil {
		t.Fatal(err)
	}

	// Field 9 is out of range for all rows — expect empty output
	if got := buf.String(); got != "" {
		t.Errorf("expected empty output for out-of-range field, got: %q", got)
	}
}

func TestCutCommaDelim(t *testing.T) {
	f, err := os.Open("tests/fourchords.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	if err := cut(f, &buf, []int{1}, ","); err != nil {
		t.Fatal(err)
	}

	// Expected output matches: cut -f1 -d, tests/fourchords.csv | head -n5
	// The file begins with a UTF-8 BOM (\ufeff), which the system cut tool preserves.
	want := "\ufeffSong title\n" +
		"\"10000 Reasons (Bless the Lord)\"\n" +
		"\"20 Good Reasons\"\n" +
		"\"Adore You\"\n" +
		"\"Africa\"\n"
	got := buf.String()
	if !strings.HasPrefix(got, want) {
		// Show only first N chars so the diff is readable
		gotPrefix := got
		if len(gotPrefix) > len(want) {
			gotPrefix = gotPrefix[:len(want)]
		}
		t.Errorf("cut -f1 -d, output mismatch:\ngot:  %q\nwant: %q", gotPrefix, want)
	}
}
