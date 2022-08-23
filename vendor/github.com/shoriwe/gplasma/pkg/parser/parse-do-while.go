package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseDoWhileStatement() (*ast.DoWhileStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(DoWhileStatement)
	}
	// Parse Body
	for parser.hasNext() {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirectValue(lexer.While) {
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
	// Parse Condition
	if !parser.matchDirectValue(lexer.While) {
		return nil, parser.newSyntaxError(DoWhileStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	var condition ast.Node
	condition, parsingError = parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.Expression); !ok {
		return nil, parser.expectingExpressionError(WhileStatement)
	}
	return &ast.DoWhileStatement{
		Condition: condition.(ast.Expression),
		Body:      body,
	}, nil
}
