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

func (p *Parser) peek() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Kind: lexer.EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) parsePair() error {
	if err := p.expect(lexer.STRING); err != nil {
		return err
	}
	p.advance()
	if err := p.expect(lexer.COLON); err != nil {
		return err
	}
	p.advance()
	if err := p.parseValue(); err != nil {
		return err
	}

	switch p.peek().Kind {
	case lexer.COMMA:
		p.advance()
		return p.parsePair()
	case lexer.CLOSE_BRACE:
		p.advance()
		return nil
	default:
		return fmt.Errorf("unexpected token %s, expected COMMA or CLOSE_BRACE", p.peek().Kind.String())
	}
}

func (p *Parser) parseValue() error {
	switch p.peek().Kind {
	case lexer.STRING:
		p.advance()
	default:
		return fmt.Errorf("unexpected token %s", p.peek().Kind.String())
	}
	return nil
}

func (p *Parser) parseJson() error {
	if err := p.expect(lexer.OPEN_BRACE); err != nil {
		return err
	}
	p.advance()

	switch p.peek().Kind {
	case lexer.CLOSE_BRACE:
		p.advance()
		return nil
	case lexer.STRING:
		return p.parsePair()
	default:
		return fmt.Errorf("unexpected token %s", p.peek().Kind.String())
	}
}

func (p *Parser) expect(tokenkind lexer.TokenKind) error {
	peek := p.peek()
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
		return
	}
	if err := parser.expect(lexer.EOF); err != nil {
		fmt.Println(err)
	}
}
