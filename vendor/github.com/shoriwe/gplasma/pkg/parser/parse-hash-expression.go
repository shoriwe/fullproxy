package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseHashExpression() (*ast.HashExpression, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var values []*ast.KeyValue
	var leftHandSide ast.Node
	var rightHandSide ast.Node
	var parsingError error
	for parser.hasNext() {
		if parser.matchDirectValue(lexer.CloseBrace) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}

		leftHandSide, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := leftHandSide.(ast.Expression); !ok {
			return nil, parser.expectingExpressionError(HashExpression)
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.matchDirectValue(lexer.Colon) {
			return nil, parser.newSyntaxError(HashExpression)
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}

		rightHandSide, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := rightHandSide.(ast.Expression); !ok {
			return nil, parser.expectingExpressionError(HashExpression)
		}
		values = append(values, &ast.KeyValue{
			Key:   leftHandSide.(ast.Expression),
			Value: rightHandSide.(ast.Expression),
		})
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.HashExpression{
		Values: values,
	}, nil
}
