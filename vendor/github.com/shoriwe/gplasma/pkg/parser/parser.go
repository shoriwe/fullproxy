package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/common"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

type Parser struct {
	lineStack    common.ListStack[int]
	lexer        *lexer.Lexer
	complete     bool
	currentToken *lexer.Token
}

func (parser *Parser) Parse() (*ast.Program, error) {
	result := &ast.Program{
		Begin: nil,
		End:   nil,
		Body:  nil,
	}
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var (
		beginStatement   *ast.BeginStatement
		endStatement     *ast.EndStatement
		parsedExpression ast.Node
		parsingError     error
	)
	for parser.hasNext() {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			continue
		}
		switch {
		case parser.matchDirectValue(lexer.BEGIN):
			if result.Begin != nil {
				return nil, BeginRepeated
			}
			beginStatement, parsingError = parser.parseBeginStatement()
			if parsingError != nil {
				return nil, parsingError
			}
			result.Begin = beginStatement
		case parser.matchDirectValue(lexer.END):
			if result.End != nil {
				return nil, EndRepeated
			}
			endStatement, parsingError = parser.parseEndStatement()
			if parsingError != nil {
				return nil, parsingError
			}
			result.End = endStatement
		default:
			parsedExpression, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			result.Body = append(result.Body, parsedExpression)
		}
	}
	parser.complete = true
	return result, nil
}

func NewParser(lexer_ *lexer.Lexer) *Parser {
	return &Parser{
		lineStack:    common.ListStack[int]{},
		lexer:        lexer_,
		complete:     false,
		currentToken: nil,
	}
}
