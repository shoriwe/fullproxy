package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseClassStatement() (*ast.ClassStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, parser.newSyntaxError(ClassStatement)
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var bases []ast.Expression
	var base ast.Node
	var parsingError error
	if parser.matchDirectValue(lexer.OpenParentheses) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		for parser.hasNext() {
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			base, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := base.(ast.Expression); !ok {
				return nil, parser.expectingExpressionError(ClassStatement)
			}
			bases = append(bases, base.(ast.Expression))
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			if parser.matchDirectValue(lexer.Comma) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
			} else if parser.matchDirectValue(lexer.CloseParentheses) {
				break
			} else {
				return nil, parser.newSyntaxError(ClassStatement)
			}
		}
		if !parser.matchDirectValue(lexer.CloseParentheses) {
			return nil, parser.newSyntaxError(ClassStatement)
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(ClassStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
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
		return nil, parser.statementNeverEndedError(ClassStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.ClassStatement{
		Name:  name,
		Bases: bases,
		Body:  body,
	}, nil
}
