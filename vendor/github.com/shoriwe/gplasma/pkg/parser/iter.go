package parser

import (
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) hasNext() bool {
	return !parser.complete
}

func (parser *Parser) next() error {
	token, tokenizingError := parser.lexer.Next()
	if tokenizingError != nil {
		return tokenizingError
	}
	if token.Kind == lexer.EOF {
		parser.complete = true
	}
	parser.currentToken = token
	return nil
}
