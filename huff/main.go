package main

import (
	"bufio"
	"log"
	"os"
)

func countCharFrequencies(file *os.File) map[rune]int {
	frequencies := make(map[rune]int)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		char := scanner.Text()
		frequencies[rune(char[0])]++
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

	countCharFrequencies(file)
}
