package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseGeneratorDefinitionStatement() (*ast.GeneratorDefinitionStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, parser.newSyntaxError(GeneratorDefinitionStatement)
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchDirectValue(lexer.OpenParentheses) {
		return nil, parser.newSyntaxError(GeneratorDefinitionStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var arguments []*ast.Identifier
	for parser.hasNext() {
		if parser.matchDirectValue(lexer.CloseParentheses) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.matchKind(lexer.IdentifierKind) {
			return nil, parser.newSyntaxError(GeneratorDefinitionStatement)
		}
		argument := &ast.Identifier{
			Token: parser.currentToken,
		}
		arguments = append(arguments, argument)
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
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
			return nil, parser.newSyntaxError(GeneratorDefinitionStatement)
		}
	}
	if !parser.matchDirectValue(lexer.CloseParentheses) {
		return nil, parser.newSyntaxError(GeneratorDefinitionStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(GeneratorDefinitionStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
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
		return nil, parser.statementNeverEndedError(GeneratorDefinitionStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	body = append(body, &ast.ReturnStatement{
		Results: []ast.Expression{&ast.BasicLiteralExpression{
			Token: &lexer.Token{
				Contents:    []rune(lexer.NoneString),
				DirectValue: lexer.None,
				Kind:        lexer.NoneType,
			},
			Kind:        lexer.NoneType,
			DirectValue: lexer.None,
		}},
	})
	return &ast.GeneratorDefinitionStatement{
		Name:      name,
		Arguments: arguments,
		Body:      body,
	}, nil
}
