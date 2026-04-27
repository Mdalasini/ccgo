package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// cut reads delimiter-separated lines from r and writes the specified fields (1-indexed) from each line to w.
// Multiple fields are joined with the delimiter in the output.
func cut(r io.Reader, w io.Writer, fields []int, delim string) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, delim)
		var selected []string
		for _, f := range fields {
			if f > 0 && f <= len(parts) {
				selected = append(selected, parts[f-1])
			}
		}
		if len(selected) > 0 {
			fmt.Fprintln(w, strings.Join(selected, delim))
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
