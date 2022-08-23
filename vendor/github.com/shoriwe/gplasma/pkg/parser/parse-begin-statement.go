package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseBeginStatement() (*ast.BeginStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(BeginStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
	for parser.hasNext() {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirectValue(lexer.End) {
				break
			}
			continue
		}
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, bodyNode)
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(BeginStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.BeginStatement{
		Body: body,
	}, nil
}
