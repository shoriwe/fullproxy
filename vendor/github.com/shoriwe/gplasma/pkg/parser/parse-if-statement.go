package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseIfStatement() (*ast.IfStatement, error) {
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
		return nil, parser.expectingExpressionError(IfStatement)
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(IfStatement)
	}
	// Parse If
	root := &ast.IfStatement{
		Condition: condition.(ast.Expression),
		Body:      []ast.Node{},
		Else:      []ast.Node{},
	}
	var bodyNode ast.Node
	for parser.hasNext() {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirectValue(lexer.Elif) ||
				parser.matchDirectValue(lexer.Else) ||
				parser.matchDirectValue(lexer.End) {
				break
			}
			continue
		}
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		root.Body = append(root.Body, bodyNode)
	}
	// Parse Elifs
	lastCondition := root
	if parser.matchDirectValue(lexer.Elif) {
	elifBlocksParsingLoop:
		for parser.hasNext() {
			block := ast.ElifBlock{
				Condition: nil,
				Body:      nil,
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			condition, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := condition.(ast.Expression); !ok {
				return nil, parser.expectingExpressionError(ElifBlock)
			}
			block.Condition = condition.(ast.Expression)
			if !parser.matchDirectValue(lexer.NewLine) {
				return nil, parser.newSyntaxError(IfStatement)
			}
			for parser.hasNext() {
				if parser.matchKind(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					if parser.matchDirectValue(lexer.Else) ||
						parser.matchDirectValue(lexer.End) {
						root.ElifBlocks = append(root.ElifBlocks, block)
						break elifBlocksParsingLoop
					} else if parser.matchDirectValue(lexer.Elif) {
						break
					}
					continue
				}
				bodyNode, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				block.Body = append(block.Body, bodyNode)
			}
			if !parser.matchDirectValue(lexer.Elif) {
				return nil, parser.newSyntaxError(ElifBlock)
			}
			root.ElifBlocks = append(root.ElifBlocks, block)
		}
	}
	// Parse Default
	if parser.matchDirectValue(lexer.Else) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		var elseBodyNode ast.Node
		if !parser.matchDirectValue(lexer.NewLine) {
			return nil, parser.newSyntaxError(ElseBlock)
		}
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
			elseBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			lastCondition.Else = append(lastCondition.Else, elseBodyNode)
		}
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(IfStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return root, nil
}
