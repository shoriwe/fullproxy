package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var results []ast.Expression
	for parser.hasNext() {
		if parser.matchKind(lexer.Separator) || parser.matchKind(lexer.EOF) {
			break
		}

		result, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := result.(ast.Expression); !ok {
			return nil, parser.expectingExpressionError(ReturnStatement)
		}
		results = append(results, result.(ast.Expression))
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !(parser.matchKind(lexer.Separator) || parser.matchKind(lexer.EOF)) {
			return nil, parser.newSyntaxError(ReturnStatement)
		}
	}
	return &ast.ReturnStatement{
		Results: results,
	}, nil
}
