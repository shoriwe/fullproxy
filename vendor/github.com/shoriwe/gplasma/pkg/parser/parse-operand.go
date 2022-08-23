package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseOperand() (ast.Node, error) {
	switch parser.currentToken.Kind {
	case lexer.Literal, lexer.Boolean, lexer.NoneType:
		return parser.parseLiteral()
	case lexer.IdentifierKind:
		identifier := parser.currentToken
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast.Identifier{
			Token: identifier,
		}, nil
	case lexer.Keyword:
		switch parser.currentToken.DirectValue {
		case lexer.Lambda:
			return parser.parseLambdaExpression()
		case lexer.Super:
			return parser.parseSuperExpression()
		case lexer.Delete:
			return parser.parseDeleteStatement()
		case lexer.Defer:
			return parser.parseDeferStatement()
		case lexer.While:
			return parser.parseWhileStatement()
		case lexer.For:
			return parser.parseForStatement()
		case lexer.Until:
			return parser.parseUntilStatement()
		case lexer.If:
			return parser.parseIfStatement()
		case lexer.Unless:
			return parser.parseUnlessStatement()
		case lexer.Switch:
			return parser.parseSwitchStatement()
		case lexer.Module:
			return parser.parseModuleStatement()
		case lexer.Def:
			return parser.parseFunctionDefinitionStatement()
		case lexer.Generator:
			return parser.parseGeneratorDefinitionStatement()
		case lexer.Interface:
			return parser.parseInterfaceStatement()
		case lexer.Class:
			return parser.parseClassStatement()
		case lexer.Return:
			return parser.parseReturnStatement()
		case lexer.Yield:
			return parser.parseYieldStatement()
		case lexer.Continue:
			return parser.parseContinueStatement()
		case lexer.Break:
			return parser.parseBreakStatement()
		case lexer.Pass:
			return parser.parsePassStatement()
		case lexer.Do:
			return parser.parseDoWhileStatement()
		}
	case lexer.Punctuation:
		switch parser.currentToken.DirectValue {
		case lexer.OpenParentheses:
			return parser.parseParentheses()
		case lexer.OpenSquareBracket: // Parse Arrays
			return parser.parseArrayExpression()
		case lexer.OpenBrace: // Parse Dictionaries
			return parser.parseHashExpression()
		}
	}
	return nil, UnknownToken
}
