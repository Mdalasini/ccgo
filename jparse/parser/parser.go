package parser

import (
	"fmt"

	"github.com/mdalasini/ccgo/jparse/lexer"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
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

func (p *Parser) expect(kind lexer.TokenKind) error {
	if p.peek().Kind != kind {
		return fmt.Errorf("expected %s, got %s", kind.String(), p.peek().Kind.String())
	}
	return nil
}

func (p *Parser) parseArray() error {
	if err := p.expect(lexer.OPEN_BRACKET); err != nil {
		return err
	}
	p.advance()

	if p.peek().Kind == lexer.CLOSE_BRACKET {
		p.advance()
		return nil
	}
	if err := p.parseObject(); err != nil {
		return err
	}

	switch p.peek().Kind {
	case lexer.COMMA:
		for p.peek().Kind == lexer.COMMA {
			p.advance()
			if err := p.parseObject(); err != nil {
				return err
			}
		}
		if err := p.expect(lexer.CLOSE_BRACKET); err != nil {
			return err
		}
		p.advance()
		return nil
	case lexer.CLOSE_BRACKET:
		p.advance()
		return nil
	default:
		return fmt.Errorf("unexpected token %s", p.peek().Kind.String())
	}
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

func (p *Parser) parsePair() error {
	if err := p.expect(lexer.STRING); err != nil {
		return err
	}
	p.advance()
	if err := p.expect(lexer.COLON); err != nil {
		return err
	}
	p.advance()
	if err := p.parseObject(); err != nil {
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

func (p *Parser) parseObject() error {
	switch p.peek().Kind {
	case lexer.STRING:
		p.advance()
	case lexer.NUMBER:
		p.advance()
	case lexer.BOOLEAN:
		p.advance()
	case lexer.NULL:
		p.advance()
	case lexer.OPEN_BRACE:
		if err := p.parseJson(); err != nil {
			return err
		}
	case lexer.OPEN_BRACKET:
		if err := p.parseArray(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected token %s", p.peek().Kind.String())
	}
	return nil
}

func newParser(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

func Parse(tokens []lexer.Token) error {
	p := newParser(tokens)
	return p.parseObject()
}
