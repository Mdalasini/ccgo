package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// cut reads tab-separated lines from r and writes the specified field (1-indexed) from each line to w.
func cut(r io.Reader, w io.Writer, field int) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if field <= len(parts) {
			fmt.Fprintln(w, parts[field-1])
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
