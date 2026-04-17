package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/mdalasini/ccgo/jparse/lexer"
	"github.com/mdalasini/ccgo/jparse/parser"
)

func usage() {
	fmt.Println("Usage: jparse <file>")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		return
	}

	file, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "jparse: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	tokens, err := lexer.Tokenize(reader)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := parser.Parse(tokens); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
