package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseBinaryExpression(precedence lexer.DirectValue) (ast.Node, error) {
	var leftHandSide ast.Node
	var rightHandSide ast.Node
	var parsingError error
	leftHandSide, parsingError = parser.parseUnaryExpression()
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := leftHandSide.(ast.Statement); ok {
		return leftHandSide, nil
	}
	for parser.hasNext() {
		if !parser.matchKind(lexer.Operator) &&
			!parser.matchKind(lexer.Comparator) {
			break
		}
		newLinesRemoveError := parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		operator := parser.currentToken
		operatorPrecedence := parser.currentToken.DirectValue
		if operatorPrecedence < precedence {
			return leftHandSide, nil
		}
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}

		rightHandSide, parsingError = parser.parseBinaryExpression(operatorPrecedence + 1)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := rightHandSide.(ast.Expression); !ok {
			return nil, parser.expectingExpressionError(BinaryExpression)
		}

		leftHandSide = &ast.BinaryExpression{
			LeftHandSide:  leftHandSide.(ast.Expression),
			Operator:      operator,
			RightHandSide: rightHandSide.(ast.Expression),
		}
	}
	return leftHandSide, nil
}
