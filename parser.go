package lang

import (
	"fmt"
	"io"
)

// FunctionStatement represents a function statement.
type FunctionStatement struct {
	Template string
	Vars     map[string]string
}

type TemplateStatement struct {
	Template string
	Values   []string
}

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

/*
f('hello {{ .foo }}', {
foo: 'bar',
});
*/
// Parse parses a function.
func (p *Parser) Parse() (*FunctionStatement, error) {
	stmt := &FunctionStatement{
		Vars: make(map[string]string),
	}

	// First token should be a "f" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != FUNCTION {
		return nil, fmt.Errorf("found %q, expected FUNCTION", lit)
	}

	// next token should be the open bracket
	if tok, lit := p.scanIgnoreWhitespace(); tok != OPEN_BRACKET {
		return nil, fmt.Errorf("found %q, expected OPEN_BRACKET", lit)
	}

	// next token should be the template
	templateToken, templateLiteral := p.scanIgnoreWhitespace()
	if templateToken != TEMPLATE {
		return nil, fmt.Errorf("found %q, expected TEMPLATE", templateLiteral)
	}
	// log.Printf("templateLiteral: %s", templateLiteral)
	stmt.Template = templateLiteral

	// next token should be the comma
	if tok, lit := p.scanIgnoreWhitespace(); tok != COMMA {
		return nil, fmt.Errorf("found %q, expected COMMA", lit)
	}

	// next token should be the open brace
	openBraceToken, openBraceLiteral := p.scanIgnoreWhitespace()
	if openBraceToken != OPEN_BRACE {
		return nil, fmt.Errorf("found %q, expected OPEN_BRACE", openBraceLiteral)
	}
	// log.Printf("openBraceLiteral: %q", openBraceLiteral)

	// Next we should loop over all our comma-delimited K/V pairs.
	for {
		next, nextLiteral := p.scanIgnoreWhitespace()
		if next == CLOSE_BRACE {
			p.unscan()
			break
		}
		// log.Printf("nextLiteral: %q", nextLiteral)

		// next must be a key identifier.
		if next != IDENT {
			return nil, fmt.Errorf("found %q, expected key IDENT", nextLiteral)
		}
		identLiteral := nextLiteral // rename for readability

		if tok, lit := p.scanIgnoreWhitespace(); tok != COLON {
			return nil, fmt.Errorf("found %q, expected COLON", lit)
		}

		// Read the value.
		valToken, valLiteral := p.scanIgnoreWhitespace()
		if valToken != TEMPLATE && valToken != NUM_VALUE {
			return nil, fmt.Errorf("found %q, expected LIT_VALUE", valLiteral)
		}
		// log.Printf("valLiteral: %s", valLiteral)

		// handle duplicate keys
		stmt.Vars[identLiteral] = valLiteral
		// log.Printf("stmt.Vars: %v", stmt.Vars)

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			// check if next token is a closing brace
			p.unscan()
			break
		}
	}

	// Next we should see the "}" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != CLOSE_BRACE {
		return nil, fmt.Errorf("found %q, expected CLOSE_BRACE", lit)
	}

	// finally we should see a closing bracket
	closeBracketToken, closeBracketLiteral := p.scanIgnoreWhitespace()
	if closeBracketToken != CLOSE_BRACKET {
		return nil, fmt.Errorf("found %q, expected CLOSE_BRACKET", closeBracketLiteral)
	}

	// Return the successfully parsed statement.
	return stmt, nil
}

func (p *Parser) ParseTemplate() (TemplateStatement, error) {
	stmt := TemplateStatement{}

	// example: 'hello {{ .foo }}'
	// find the template in a string

	return stmt, nil
}

func (p *FunctionStatement) Execute() (string, error) {
	result := ""
	return result, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
