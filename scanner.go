package lang

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	// If we see a digit then consume as a number.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		next := s.peek()
		if isOpenBracket(next) {
			return FUNCTION, string(ch)
		}
		if isColon(next) {
			return IDENT, string(ch)
		}
		s.unread()
		tok, lit := s.scanIdent()
		return tok, fmt.Sprintf("%s%s", string(ch), lit)
	} else if isDigit(ch) {
		s.unread()
		return s.scanNumber()
	} else if isQuote(ch) {
		s.unread()
		return s.scanTemplate()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case '(':
		return OPEN_BRACKET, string(ch)
	case ')':
		return CLOSE_BRACKET, string(ch)
	case '{':
		return OPEN_BRACE, string(ch)
	case '}':
		return CLOSE_BRACE, string(ch)
	case '\'':
		return QUOTE, string(ch)
	case ':':
		return COLON, string(ch)
	case ',':
		return COMMA, string(ch)
	case ';':
		return SEMICOLON, string(ch)
	}

	return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		ch := s.read()
		if ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}

// scanNumber consumes the current rune and all contiguous number runes.
func (s *Scanner) scanNumber() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent number character into the buffer.
	// Non-number characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isDigit(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// Otherwise return as a regular number.
	return NUM_VALUE, buf.String()
}

// scanTemplate read a template that starts and ends with a single quote
func (s *Scanner) scanTemplate() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	quote := s.read()
	if quote != '\'' {
		return ILLEGAL, string(quote)
	}
	buf.WriteRune(quote)

	// Read every subsequent character into the buffer.
	// single quote and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if isQuote(ch) {
			_, _ = buf.WriteRune(ch)
			// check that next character is a comma
			next := s.peek()
			if next != ',' {
				continue
			}
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// return the template and the literal string
	return TEMPLATE, buf.String()
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s Scanner) peek() rune {
	ch, err := s.r.Peek(1)
	if err != nil {
		return eof
	}
	return rune(ch[0])
}

// might not be the best thing to do
// func (s Scanner) peek() rune {
// 	ch := s.read()
// 	s.unread()
// 	return ch
// }

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

func isOpenBracket(ch rune) bool  { return ch == '(' }
func isCloseBracket(ch rune) bool { return ch == ')' }
func isOpenBrace(ch rune) bool    { return ch == '{' }
func isCloseBrace(ch rune) bool   { return ch == '}' }
func isQuote(ch rune) bool        { return ch == '\'' } // maybe support double quotes
func isColon(ch rune) bool        { return ch == ':' }

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
