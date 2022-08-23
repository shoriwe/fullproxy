package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseInterfaceStatement() (*ast.InterfaceStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, parser.newSyntaxError(InterfaceStatement)
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
				return nil, parser.newSyntaxError(InterfaceStatement)
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
				return nil, parser.newSyntaxError(InterfaceStatement)
			}
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.matchDirectValue(lexer.CloseParentheses) {
			return nil, parser.newSyntaxError(InterfaceStatement)
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(InterfaceStatement)
	}
	var methods []*ast.FunctionDefinitionStatement
	var node ast.Node
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
		node, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := node.(*ast.FunctionDefinitionStatement); !ok {
			return nil, parser.expectingFunctionDefinition(InterfaceStatement)
		}
		methods = append(methods, node.(*ast.FunctionDefinitionStatement))
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(InterfaceStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.InterfaceStatement{
		Name:              name,
		Bases:             bases,
		MethodDefinitions: methods,
	}, nil
}
