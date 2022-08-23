package parser

import "github.com/shoriwe/gplasma/pkg/ast"

func (parser *Parser) parseDeferStatement() (*ast.DeferStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	x, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := x.(*ast.MethodInvocationExpression); !ok {
		return nil, parser.expectingExpressionError(DeferStatement)
	}
	return &ast.DeferStatement{
		X: x.(*ast.MethodInvocationExpression),
	}, nil
}
