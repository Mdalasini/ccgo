package main

import (
	"bufio"
	"log"
	"os"
	"unicode/utf8"
)

func countCharFrequencies(file *os.File) map[rune]int {
	frequencies := make(map[rune]int)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		r, _ := utf8.DecodeRuneInString(scanner.Text())
		frequencies[r]++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	return frequencies
}

func main() {
	filepath := "tests/test.txt"
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("unable to open %v: %v", filepath, err)
	}
	defer file.Close()
}
