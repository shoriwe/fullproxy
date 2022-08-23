package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseSwitchStatement() (*ast.SwitchStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	target, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := target.(ast.Expression); !ok {
		return nil, parser.expectingExpressionError(SwitchStatement)
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(SwitchStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	// parse Cases
	var caseBlocks []*ast.CaseBlock
	if parser.matchDirectValue(lexer.Case) {
		for parser.hasNext() {
			if parser.matchDirectValue(lexer.Default) ||
				parser.matchDirectValue(lexer.End) {
				break
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			var cases []ast.Expression
			var caseTarget ast.Node
			for parser.hasNext() {
				caseTarget, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				if _, ok := caseTarget.(ast.Expression); !ok {
					return nil, parser.expectingExpressionError(CaseBlock)
				}
				cases = append(cases, caseTarget.(ast.Expression))
				if parser.matchDirectValue(lexer.NewLine) {
					break
				} else if parser.matchDirectValue(lexer.Comma) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
				} else {
					return nil, parser.newSyntaxError(CaseBlock)
				}
			}
			if !parser.matchDirectValue(lexer.NewLine) {
				return nil, parser.newSyntaxError(CaseBlock)
			}
			// Targets Body
			var caseBody []ast.Node
			var caseBodyNode ast.Node
			for parser.hasNext() {
				if parser.matchKind(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					if parser.matchDirectValue(lexer.Case) ||
						parser.matchDirectValue(lexer.Default) ||
						parser.matchDirectValue(lexer.End) {
						break
					}
					continue
				}
				caseBodyNode, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				caseBody = append(caseBody, caseBodyNode)
			}
			// Targets block
			caseBlocks = append(caseBlocks, &ast.CaseBlock{
				Cases: cases,
				Body:  caseBody,
			})
		}
	}
	// Parse Default
	var defaultBody []ast.Node
	if parser.matchDirectValue(lexer.Default) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.matchDirectValue(lexer.NewLine) {
			return nil, parser.newSyntaxError(DefaultBlock)
		}
		var defaultBodyNode ast.Node
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
			defaultBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			defaultBody = append(defaultBody, defaultBodyNode)
		}
	}
	// Finally detect valid end
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(SwitchStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.SwitchStatement{
		Target:     target.(ast.Expression),
		CaseBlocks: caseBlocks,
		Default:    defaultBody,
	}, nil
}
