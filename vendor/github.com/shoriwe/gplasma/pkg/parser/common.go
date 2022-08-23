package parser

import (
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) removeNewLines() error {
	for parser.matchDirectValue(lexer.NewLine) {
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return tokenizingError
		}
	}
	return nil
}

func (parser *Parser) matchDirectValue(directValue lexer.DirectValue) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.DirectValue == directValue
}

func (parser *Parser) matchKind(kind lexer.Kind) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.Kind == kind
}

func (parser *Parser) matchString(value string) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.String() == value
}
