package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
)

func (parser *Parser) parseAssignmentStatement(leftHandSide ast.Expression) (*ast.AssignStatement, error) {
	assignmentToken := parser.currentToken
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	rightHandSide, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := rightHandSide.(ast.Expression); !ok {
		return nil, parser.expectingExpressionError(AssignStatement)
	}
	return &ast.AssignStatement{
		LeftHandSide:   leftHandSide,
		AssignOperator: assignmentToken,
		RightHandSide:  rightHandSide.(ast.Expression),
	}, nil
}
