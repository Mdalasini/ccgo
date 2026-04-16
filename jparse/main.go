package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/mdalasini/ccgo/jparse/lexer"
)

func usage() {
	fmt.Println("Usage: jparse <file>")
	flag.PrintDefaults()
}

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func newParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) advance() {
	p.pos++
}

func (p *Parser) parseJson() error {
	if err := p.expect(lexer.OPEN_BRACE); err != nil {
		return err
	}
	p.advance()
	if err := p.expect(lexer.CLOSE_BRACE); err != nil {
		return err
	}
	p.advance()
	if err := p.expect(lexer.EOF); err != nil {
		return err
	}
	p.advance()
	return nil
}

func (p *Parser) expect(tokenkind lexer.TokenKind) error {
	peek := p.tokens[p.pos]
	if peek.Kind != tokenkind {
		return fmt.Errorf("jparse: expected %s, got %s", tokenkind.String(), peek.Kind.String())
	}
	return nil
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
		fmt.Printf("jparse: %v", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	tokens, err := lexer.Tokenize(reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	parser := newParser(tokens)
	if err := parser.parseJson(); err != nil {
		fmt.Println(err)
	}
}
