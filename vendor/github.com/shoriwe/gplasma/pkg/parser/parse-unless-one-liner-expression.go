package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseUnlessOneLinerExpression(result ast.Expression) (*ast.UnlessOneLinerExpression, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.Expression); !ok {
		return nil, parser.expectingExpressionError(UnlessOneLinerExpression)
	}
	if !parser.matchDirectValue(lexer.Else) {
		return &ast.UnlessOneLinerExpression{
			Result:     result,
			Condition:  condition.(ast.Expression),
			ElseResult: nil,
		}, nil
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var elseResult ast.Node

	elseResult, parsingError = parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.Expression); !ok {
		return nil, parser.expectingExpressionError(OneLineElseBlock)
	}
	return &ast.UnlessOneLinerExpression{
		Result:     result,
		Condition:  condition.(ast.Expression),
		ElseResult: elseResult.(ast.Expression),
	}, nil
}
