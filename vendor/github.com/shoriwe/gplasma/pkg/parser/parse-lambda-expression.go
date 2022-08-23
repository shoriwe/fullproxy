package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseLambdaExpression() (*ast.LambdaExpression, error) {
	var arguments []*ast.Identifier
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	for parser.hasNext() {
		if parser.matchDirectValue(lexer.Colon) {
			break
		}
		newLinesRemoveError := parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}

		identifier, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := identifier.(*ast.Identifier); !ok {
			return nil, parser.expectingIdentifier(LambdaExpression)
		}
		arguments = append(arguments, identifier.(*ast.Identifier))
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !parser.matchDirectValue(lexer.Colon) {
			return nil, parser.newSyntaxError(LambdaExpression)
		}
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchDirectValue(lexer.Colon) {
		return nil, parser.newSyntaxError(LambdaExpression)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	code, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := code.(ast.Expression); !ok {
		return nil, parser.expectingExpressionError(LambdaExpression)
	}
	return &ast.LambdaExpression{
		Arguments: arguments,
		Code:      code.(ast.Expression),
	}, nil
}
