package main

import (
	"bufio"
	"fmt"
	"io"
)

func main() {
	fmt.Println("hello from huff")
}

func FreqMap(r io.Reader) (map[rune]int, error) {
	freq := make(map[rune]int)
	br := bufio.NewReader(r)
	for {
		ch, _, err := br.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		freq[ch]++
	}
	return freq, nil
}
