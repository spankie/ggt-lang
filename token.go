package lang

type Token int

const (
	// special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// variables
	IDENT
	KEY
	// LIT_VALUE
	NUM_VALUE

	// function
	FUNCTION
	TEMPLATE

	// special characters
	PERIOD
	OPEN_BRACE
	CLOSE_BRACE
	OPEN_BRACKET
	CLOSE_BRACKET
	QUOTE
	PIPE
	COLON
	COMMA
	SEMICOLON
)
