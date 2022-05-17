package parser

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/compiler/ast"
	"github.com/shoriwe/gplasma/pkg/compiler/lexer"
	"github.com/shoriwe/gplasma/pkg/errors"
)

const (
	ForStatement                = "For Statement"
	UntilStatement              = "Until Statement"
	ModuleStatement             = "Module Statement"
	FunctionDefinitionStatement = "Function Definition Statement"
	ClassStatement              = "Class Statement"
	RaiseStatement              = "Raise Statement"
	TryStatement                = "Try Statement"
	ExceptBlock                 = "Except Block"
	ElseBlock                   = "Else Block"
	FinallyBlock                = "Finally Block"
	BeginStatement              = "Begin Statement"
	EndStatement                = "End Statement"
	InterfaceStatement          = "Interface Statement"
	BinaryExpression            = "Binary Expression"
	PointerExpression           = "Pointer Expression"
	LambdaExpression            = "Lambda Expression"
	ParenthesesExpression       = "Parentheses Expression"
	TupleExpression             = "Tuple Expression"
	ArrayExpression             = "Array Expression"
	HashExpression              = "Hash Expression"
	WhileStatement              = "While Statement"
	DoWhileStatement            = "Do-While Statement"
	IfStatement                 = "If Statement"
	ElifBlock                   = "Elif Block"
	UnlessStatement             = "Unless Statement"
	SwitchStatement             = "Switch Statement"
	CaseBlock                   = "Targets Block"
	DefaultBlock                = "Default Block"
	ReturnStatement             = "Return Statement"
	YieldStatement              = "Yield Statement"
	SuperStatement              = "Super Statement"
	SelectorExpression          = "Selector Expression"
	MethodInvocationExpression  = "Method Invocation Expression"
	IndexExpression             = "Index Expression"
	IfOneLinerExpression        = "If One Liner Expression"
	UnlessOneLinerExpression    = "Unless One Liner Expression"
	OneLineElseBlock            = "One Line Else Block"
	GeneratorExpression         = "Generator Expression"
	AssignStatement             = "Assign Statement"
)

func newSyntaxError(line int, nodeType string) *errors.Error {
	return errors.New(line, fmt.Sprintf("invalid definition of %s", nodeType), errors.SyntaxError)
}

func newNonExpressionReceivedError(line int, nodeType string) *errors.Error {
	return errors.New(line, fmt.Sprintf("received a non expression in %s", nodeType), errors.SyntaxError)
}

func newNonIdentifierReceivedError(line int, nodeType string) *errors.Error {
	return errors.New(line, fmt.Sprintf("received a non identifier in %s", nodeType), errors.SyntaxError)
}

func newStatementNeverEndedError(line int, nodeType string) *errors.Error {
	return errors.New(line, fmt.Sprintf("%s never ended", nodeType), errors.SyntaxError)
}

func newInvalidKindOfTokenError(line int) *errors.Error {
	return errors.New(line, "invalid kind of token", errors.ParsingError)
}

func newExpressionNeverClosedError(line int, nodeType string) *errors.Error {
	return errors.New(line, fmt.Sprintf("%s never closed", nodeType), errors.SyntaxError)
}

func newNonFunctionDefinitionReceived(line int, nodeType string) *errors.Error {
	return errors.New(line, fmt.Sprintf("non function definition received in %s", nodeType), errors.SyntaxError)
}

type Parser struct {
	lexer        *lexer.Lexer
	complete     bool
	currentToken *lexer.Token
}

func (parser *Parser) hasNext() bool {
	return !parser.complete
}

func (parser *Parser) next() *errors.Error {
	token, tokenizingError := parser.lexer.Next()
	if tokenizingError != nil {
		return tokenizingError
	}
	if token.Kind == lexer.EOF {
		parser.complete = true
	}
	parser.currentToken = token
	return nil
}

func (parser *Parser) directValueMatch(directValue uint8) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.DirectValue == directValue
}

func (parser *Parser) kindMatch(kind uint8) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.Kind == kind
}

func (parser *Parser) currentLine() int {
	if parser.currentToken == nil {
		return 0
	}
	return parser.currentToken.Line
}

func (parser *Parser) matchString(value string) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.String == value
}

func (parser *Parser) parseForStatement() (*ast.ForLoopStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var receivers []*ast.Identifier
	for parser.hasNext() {
		if parser.directValueMatch(lexer.In) {
			break
		} else if !parser.kindMatch(lexer.IdentifierKind) {
			return nil, newSyntaxError(parser.currentLine(), ForStatement)
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		receivers = append(receivers, &ast.Identifier{
			Token: parser.currentToken,
		})
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.directValueMatch(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if parser.directValueMatch(lexer.In) {
			break
		} else {
			return nil, newSyntaxError(parser.currentLine(), ForStatement)
		}
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.directValueMatch(lexer.In) {
		return nil, newSyntaxError(parser.currentLine(), ForStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	source, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := source.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(parser.currentLine(), ForStatement)
	}
	if !parser.directValueMatch(lexer.NewLine) {
		return nil, newSyntaxError(parser.currentLine(), ForStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var body []ast.Node
	var bodyNode ast.Node
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.directValueMatch(lexer.End) {
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
	if !parser.directValueMatch(lexer.End) {
		return nil, newStatementNeverEndedError(parser.currentLine(), ForStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.ForLoopStatement{
		Receivers: receivers,
		Source:    source.(ast.IExpression),
		Body:      body,
	}, nil
}

func (parser *Parser) parseUntilStatement() (*ast.UntilLoopStatement, *errors.Error) {
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
	if _, ok := condition.(ast.IExpression); !ok {
		return nil, newSyntaxError(parser.currentLine(), UntilStatement)
	}
	if !parser.directValueMatch(lexer.NewLine) {
		return nil, newSyntaxError(parser.currentLine(), UntilStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var bodyNode ast.Node
	var body []ast.Node
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.directValueMatch(lexer.End) {
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
	if !parser.directValueMatch(lexer.End) {
		return nil, newStatementNeverEndedError(parser.currentLine(), UntilStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.UntilLoopStatement{
		Condition: condition.(ast.IExpression),
		Body:      body,
	}, nil
}

func (parser *Parser) parseModuleStatement() (*ast.ModuleStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.kindMatch(lexer.IdentifierKind) {
		return nil, newSyntaxError(parser.currentLine(), ModuleStatement)
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.directValueMatch(lexer.NewLine) {
		return nil, newSyntaxError(parser.currentLine(), ModuleStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError *errors.Error
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.directValueMatch(lexer.End) {
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
	if !parser.directValueMatch(lexer.End) {
		return nil, newStatementNeverEndedError(parser.currentLine(), ModuleStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.ModuleStatement{
		Name: name,
		Body: body,
	}, nil
}

func (parser *Parser) parseFunctionDefinitionStatement() (*ast.FunctionDefinitionStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.kindMatch(lexer.IdentifierKind) {
		return nil, newSyntaxError(parser.currentLine(), FunctionDefinitionStatement)
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
	if !parser.directValueMatch(lexer.OpenParentheses) {
		return nil, newSyntaxError(parser.currentLine(), FunctionDefinitionStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var arguments []*ast.Identifier
	for parser.hasNext() {
		if parser.directValueMatch(lexer.CloseParentheses) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.kindMatch(lexer.IdentifierKind) {
			return nil, newSyntaxError(parser.currentLine(), FunctionDefinitionStatement)
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
		if parser.directValueMatch(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if parser.directValueMatch(lexer.CloseParentheses) {
			break
		} else {
			return nil, newSyntaxError(parser.currentLine(), FunctionDefinitionStatement)
		}
	}
	if !parser.directValueMatch(lexer.CloseParentheses) {
		return nil, newSyntaxError(parser.currentLine(), FunctionDefinitionStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.directValueMatch(lexer.NewLine) {
		return nil, newSyntaxError(parser.currentLine(), FunctionDefinitionStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError *errors.Error
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.directValueMatch(lexer.End) {
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
	if !parser.directValueMatch(lexer.End) {
		return nil, newStatementNeverEndedError(parser.currentLine(), FunctionDefinitionStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.FunctionDefinitionStatement{
		Name:      name,
		Arguments: arguments,
		Body:      body,
	}, nil
}

func (parser *Parser) parseClassStatement() (*ast.ClassStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.kindMatch(lexer.IdentifierKind) {
		return nil, newSyntaxError(parser.currentLine(), ClassStatement)
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var bases []ast.IExpression
	var base ast.Node
	var parsingError *errors.Error
	if parser.directValueMatch(lexer.OpenParentheses) {
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
			if _, ok := base.(ast.IExpression); !ok {
				return nil, newNonExpressionReceivedError(parser.currentLine(), ClassStatement)
			}
			bases = append(bases, base.(ast.IExpression))
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			if parser.directValueMatch(lexer.Comma) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
			} else if parser.directValueMatch(lexer.CloseParentheses) {
				break
			} else {
				return nil, newSyntaxError(parser.currentLine(), ClassStatement)
			}
		}
		if !parser.directValueMatch(lexer.CloseParentheses) {
			return nil, newSyntaxError(parser.currentLine(), ClassStatement)
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if !parser.directValueMatch(lexer.NewLine) {
		return nil, newSyntaxError(parser.currentLine(), ClassStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.directValueMatch(lexer.End) {
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
	if !parser.directValueMatch(lexer.End) {
		return nil, newStatementNeverEndedError(parser.currentLine(), ClassStatement)
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

func (parser *Parser) parseRaiseStatement() (*ast.RaiseStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	x, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := x.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(parser.currentLine(), RaiseStatement)
	}
	return &ast.RaiseStatement{
		X: x.(ast.IExpression),
	}, nil
}

func (parser *Parser) parseTryStatement() (*ast.TryStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.directValueMatch(lexer.NewLine) {
		return nil, newSyntaxError(parser.currentLine(), TryStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError *errors.Error
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.directValueMatch(lexer.End) ||
				parser.directValueMatch(lexer.Except) ||
				parser.directValueMatch(lexer.Else) ||
				parser.directValueMatch(lexer.Finally) {
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
	var exceptBlocks []*ast.ExceptBlock
	for parser.hasNext() {
		if !parser.directValueMatch(lexer.Except) {
			break
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		var targets []ast.IExpression
		var target ast.Node
		for parser.hasNext() {
			if parser.directValueMatch(lexer.NewLine) ||
				parser.directValueMatch(lexer.As) {
				break
			}
			target, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := target.(ast.IExpression); !ok {
				return nil, newSyntaxError(parser.currentLine(), ExceptBlock)
			}
			targets = append(targets, target.(ast.IExpression))
			if parser.directValueMatch(lexer.Comma) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
			}
		}
		var captureName *ast.Identifier
		if parser.directValueMatch(lexer.As) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			newLinesRemoveError := parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			if !parser.kindMatch(lexer.IdentifierKind) {
				return nil, newSyntaxError(parser.currentLine(), ExceptBlock)
			}
			captureName = &ast.Identifier{
				Token: parser.currentToken,
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
		if !parser.directValueMatch(lexer.NewLine) {
			return nil, newSyntaxError(parser.currentLine(), TryStatement)
		}
		var exceptBody []ast.Node
		var exceptBodyNode ast.Node
		for parser.hasNext() {
			if parser.kindMatch(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.directValueMatch(lexer.End) ||
					parser.directValueMatch(lexer.Except) ||
					parser.directValueMatch(lexer.Else) ||
					parser.directValueMatch(lexer.Finally) {
					break
				}
				continue
			}
			exceptBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			exceptBody = append(exceptBody, exceptBodyNode)
		}
		exceptBlocks = append(exceptBlocks, &ast.ExceptBlock{
			Targets:     targets,
			CaptureName: captureName,
			Body:        exceptBody,
		})
	}
	var elseBody []ast.Node
	var elseBodyNode ast.Node
	if parser.directValueMatch(lexer.Else) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.directValueMatch(lexer.NewLine) {
			return nil, newSyntaxError(parser.currentLine(), ElseBlock)
		}
		for parser.hasNext() {
			if parser.kindMatch(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.directValueMatch(lexer.End) ||
					parser.directValueMatch(lexer.Finally) {
					break
				}
				continue
			}
			elseBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			elseBody = append(elseBody, elseBodyNode)
		}
	}
	var finallyBody []ast.Node
	var finallyBodyNode ast.Node
	if parser.directValueMatch(lexer.Finally) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.directValueMatch(lexer.NewLine) {
			return nil, newSyntaxError(parser.currentLine(), FinallyBlock)
		}
		for parser.hasNext() {
			if parser.kindMatch(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.directValueMatch(lexer.End) {
					break
				}
				continue
			}
			finallyBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			finallyBody = append(finallyBody, finallyBodyNode)
		}
	}
	if !parser.directValueMatch(lexer.End) {
		return nil, newSyntaxError(parser.currentLine(), TryStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.TryStatement{
		Body:         body,
		ExceptBlocks: exceptBlocks,
		Else:         elseBody,
		Finally:      finallyBody,
	}, nil
}

func (parser *Parser) parseBeginStatement() (*ast.BeginStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.directValueMatch(lexer.NewLine) {
		return nil, newSyntaxError(parser.currentLine(), BeginStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError *errors.Error
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.directValueMatch(lexer.End) {
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
	if !parser.directValueMatch(lexer.End) {
		return nil, newStatementNeverEndedError(parser.currentLine(), BeginStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.BeginStatement{
		Body: body,
	}, nil
}

func (parser *Parser) parseEndStatement() (*ast.EndStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.directValueMatch(lexer.NewLine) {
		return nil, newSyntaxError(parser.currentLine(), EndStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError *errors.Error
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.directValueMatch(lexer.End) {
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
	if !parser.directValueMatch(lexer.End) {
		return nil, newStatementNeverEndedError(parser.currentLine(), EndStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.EndStatement{
		Body: body,
	}, nil
}

func (parser *Parser) parseInterfaceStatement() (*ast.InterfaceStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.kindMatch(lexer.IdentifierKind) {
		return nil, newSyntaxError(parser.currentLine(), InterfaceStatement)
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var bases []ast.IExpression
	var base ast.Node
	var parsingError *errors.Error
	if parser.directValueMatch(lexer.OpenParentheses) {
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
			if _, ok := base.(ast.IExpression); !ok {
				return nil, newSyntaxError(parser.currentLine(), InterfaceStatement)
			}
			bases = append(bases, base.(ast.IExpression))
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			if parser.directValueMatch(lexer.Comma) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
			} else if parser.directValueMatch(lexer.CloseParentheses) {
				break
			} else {
				return nil, newSyntaxError(parser.currentLine(), InterfaceStatement)
			}
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.directValueMatch(lexer.CloseParentheses) {
			return nil, newSyntaxError(parser.currentLine(), InterfaceStatement)
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if !parser.directValueMatch(lexer.NewLine) {
		return nil, newSyntaxError(parser.currentLine(), InterfaceStatement)
	}
	var methods []*ast.FunctionDefinitionStatement
	var node ast.Node
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.directValueMatch(lexer.End) {
				break
			}
			continue
		}
		node, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := node.(*ast.FunctionDefinitionStatement); !ok {
			return nil, newNonFunctionDefinitionReceived(parser.currentLine(), InterfaceStatement)
		}
		methods = append(methods, node.(*ast.FunctionDefinitionStatement))
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.directValueMatch(lexer.End) {
		return nil, newStatementNeverEndedError(parser.currentLine(), InterfaceStatement)
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

func (parser *Parser) parseLiteral() (ast.IExpression, *errors.Error) {
	if !parser.kindMatch(lexer.Literal) &&
		!parser.kindMatch(lexer.Boolean) &&
		!parser.kindMatch(lexer.NoneType) {
		return nil, newInvalidKindOfTokenError(parser.currentLine())
	}

	switch parser.currentToken.DirectValue {
	case lexer.SingleQuoteString, lexer.DoubleQuoteString, lexer.ByteString,
		lexer.Integer, lexer.HexadecimalInteger, lexer.BinaryInteger, lexer.OctalInteger,
		lexer.Float, lexer.ScientificFloat,
		lexer.True, lexer.False, lexer.None:
		currentToken := parser.currentToken
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast.BasicLiteralExpression{
			Token:       currentToken,
			Kind:        currentToken.Kind,
			DirectValue: currentToken.DirectValue,
		}, nil
	}
	return nil, newInvalidKindOfTokenError(parser.currentLine())
}

func (parser *Parser) parseBinaryExpression(precedence uint8) (ast.Node, *errors.Error) {
	var leftHandSide ast.Node
	var rightHandSide ast.Node
	var parsingError *errors.Error
	leftHandSide, parsingError = parser.parseUnaryExpression()
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := leftHandSide.(ast.Statement); ok {
		return leftHandSide, nil
	}
	for parser.hasNext() {
		if !parser.kindMatch(lexer.Operator) &&
			!parser.kindMatch(lexer.Comparator) {
			break
		}
		newLinesRemoveError := parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		operator := parser.currentToken
		operatorPrecedence := parser.currentToken.DirectValue
		if operatorPrecedence < precedence {
			return leftHandSide, nil
		}
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		line := parser.currentLine()
		rightHandSide, parsingError = parser.parseBinaryExpression(operatorPrecedence + 1)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := rightHandSide.(ast.IExpression); !ok {
			return nil, newNonExpressionReceivedError(line, BinaryExpression)
		}

		leftHandSide = &ast.BinaryExpression{
			LeftHandSide:  leftHandSide.(ast.IExpression),
			Operator:      operator,
			RightHandSide: rightHandSide.(ast.IExpression),
		}
	}
	return leftHandSide, nil
}

func (parser *Parser) parseUnaryExpression() (ast.Node, *errors.Error) {
	// Do something to parse Unary
	if parser.kindMatch(lexer.Operator) {
		switch parser.currentToken.DirectValue {
		case lexer.Sub, lexer.Add, lexer.NegateBits, lexer.SignNot, lexer.Not:
			operator := parser.currentToken
			tokenizingError := parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			line := parser.currentLine()
			x, parsingError := parser.parseUnaryExpression()
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := x.(ast.IExpression); !ok {
				return nil, newNonExpressionReceivedError(line, PointerExpression)
			}
			return &ast.UnaryExpression{
				Operator: operator,
				X:        x.(ast.IExpression),
			}, nil
		}
	}
	return parser.parsePrimaryExpression()
}

func (parser *Parser) removeNewLines() *errors.Error {
	for parser.directValueMatch(lexer.NewLine) {
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return tokenizingError
		}
	}
	return nil
}

func (parser *Parser) parseLambdaExpression() (*ast.LambdaExpression, *errors.Error) {
	var arguments []*ast.Identifier
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	for parser.hasNext() {
		if parser.directValueMatch(lexer.Colon) {
			break
		}
		newLinesRemoveError := parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		line := parser.currentLine()
		identifier, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := identifier.(*ast.Identifier); !ok {
			return nil, newNonIdentifierReceivedError(line, LambdaExpression)
		}
		arguments = append(arguments, identifier.(*ast.Identifier))
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.directValueMatch(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !parser.directValueMatch(lexer.Colon) {
			return nil, newSyntaxError(line, LambdaExpression)
		}
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.directValueMatch(lexer.Colon) {
		return nil, newSyntaxError(parser.currentLine(), LambdaExpression)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	line := parser.currentLine()
	code, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := code.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(line, LambdaExpression)
	}
	return &ast.LambdaExpression{
		Arguments: arguments,
		Code:      code.(ast.IExpression),
	}, nil
}

func (parser *Parser) parseParentheses() (ast.IExpression, *errors.Error) {
	/*
		This should also parse generators
	*/
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if parser.directValueMatch(lexer.CloseParentheses) {
		return nil, newSyntaxError(parser.currentLine(), ParenthesesExpression)
	}
	line := parser.currentLine()
	firstExpression, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := firstExpression.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(line, ParenthesesExpression)
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if parser.directValueMatch(lexer.For) {
		return parser.parseGeneratorExpression(firstExpression.(ast.IExpression))
	}
	if parser.directValueMatch(lexer.CloseParentheses) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast.ParenthesesExpression{
			X: firstExpression.(ast.IExpression),
		}, nil
	}
	if !parser.directValueMatch(lexer.Comma) {
		return nil, newSyntaxError(parser.currentLine(), ParenthesesExpression)
	}
	var values []ast.IExpression
	values = append(values, firstExpression.(ast.IExpression))
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var nextValue ast.Node
	for parser.hasNext() {
		if parser.directValueMatch(lexer.CloseParentheses) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		line = parser.currentLine()
		nextValue, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := nextValue.(ast.IExpression); !ok {
			return nil, newNonExpressionReceivedError(parser.currentLine(), ParenthesesExpression)
		}
		values = append(values, nextValue.(ast.IExpression))
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.directValueMatch(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !parser.directValueMatch(lexer.CloseParentheses) {
			return nil, newSyntaxError(parser.currentLine(), TupleExpression)
		}
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.directValueMatch(lexer.CloseParentheses) {
		return nil, newExpressionNeverClosedError(parser.currentLine(), TupleExpression)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.TupleExpression{
		Values: values,
	}, nil
}
func (parser *Parser) parseArrayExpression() (*ast.ArrayExpression, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var values []ast.IExpression
	for parser.hasNext() {
		if parser.directValueMatch(lexer.CloseSquareBracket) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		line := parser.currentLine()
		value, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := value.(ast.IExpression); !ok {
			return nil, newNonExpressionReceivedError(line, ArrayExpression)
		}
		values = append(values, value.(ast.IExpression))
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.directValueMatch(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !parser.directValueMatch(lexer.CloseSquareBracket) {
			return nil, newSyntaxError(parser.currentLine(), ArrayExpression)
		}
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.ArrayExpression{
		Values: values,
	}, nil
}

func (parser *Parser) parseHashExpression() (*ast.HashExpression, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var values []*ast.KeyValue
	var leftHandSide ast.Node
	var rightHandSide ast.Node
	var parsingError *errors.Error
	for parser.hasNext() {
		if parser.directValueMatch(lexer.CloseBrace) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		line := parser.currentLine()
		leftHandSide, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := leftHandSide.(ast.IExpression); !ok {
			return nil, newNonExpressionReceivedError(line, HashExpression)
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.directValueMatch(lexer.Colon) {
			return nil, newSyntaxError(parser.currentLine(), HashExpression)
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		line = parser.currentLine()
		rightHandSide, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := rightHandSide.(ast.IExpression); !ok {
			return nil, newNonExpressionReceivedError(parser.currentLine(), HashExpression)
		}
		values = append(values, &ast.KeyValue{
			Key:   leftHandSide.(ast.IExpression),
			Value: rightHandSide.(ast.IExpression),
		})
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.directValueMatch(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.HashExpression{
		Values: values,
	}, nil
}

func (parser *Parser) parseWhileStatement() (*ast.WhileLoopStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.IExpression); !ok {
		return nil, newSyntaxError(parser.currentLine(), WhileStatement)
	}
	if !parser.directValueMatch(lexer.NewLine) {
		return nil, newSyntaxError(parser.currentLine(), WhileStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var whileChild ast.Node
	var body []ast.Node
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.directValueMatch(lexer.End) {
				break
			}
			continue
		}
		whileChild, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, whileChild)
	}
	if !parser.directValueMatch(lexer.End) {
		return nil, newStatementNeverEndedError(parser.currentLine(), WhileStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.WhileLoopStatement{
		Condition: condition.(ast.IExpression),
		Body:      body,
	}, nil
}

func (parser *Parser) parseDoWhileStatement() (*ast.DoWhileStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError *errors.Error
	if !parser.directValueMatch(lexer.NewLine) {
		return nil, newSyntaxError(parser.currentLine(), DoWhileStatement)
	}
	// Parse Body
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.directValueMatch(lexer.While) {
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
	if !parser.directValueMatch(lexer.While) {
		return nil, newSyntaxError(parser.currentLine(), DoWhileStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	line := parser.currentLine()
	var condition ast.Node
	condition, parsingError = parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(line, WhileStatement)
	}
	return &ast.DoWhileStatement{
		Condition: condition.(ast.IExpression),
		Body:      body,
	}, nil
}

func (parser *Parser) parseIfStatement() (*ast.IfStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	line := parser.currentLine()
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(line, IfStatement)
	}
	if !parser.directValueMatch(lexer.NewLine) {
		return nil, newSyntaxError(line, IfStatement)
	}
	// Parse If
	root := &ast.IfStatement{
		Condition: condition.(ast.IExpression),
		Body:      []ast.Node{},
		Else:      []ast.Node{},
	}
	var bodyNode ast.Node
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.directValueMatch(lexer.Elif) ||
				parser.directValueMatch(lexer.Else) ||
				parser.directValueMatch(lexer.End) {
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
	if parser.directValueMatch(lexer.Elif) {
		var elifBody []ast.Node
	parsingElifLoop:
		for parser.hasNext() {
			if parser.kindMatch(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.directValueMatch(lexer.Else) ||
					parser.directValueMatch(lexer.End) {
					break
				}
				continue
			}
			if !parser.directValueMatch(lexer.Elif) {
				return nil, newSyntaxError(line, ElifBlock)
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			var elifCondition ast.Node
			elifCondition, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := elifCondition.(ast.IExpression); !ok {
				return nil, newSyntaxError(line, ElifBlock)
			}
			if !parser.directValueMatch(lexer.NewLine) {
				return nil, newSyntaxError(line, ElifBlock)
			}
			var elifBodyNode ast.Node
			for parser.hasNext() {
				if parser.kindMatch(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					if parser.directValueMatch(lexer.Else) ||
						parser.directValueMatch(lexer.End) ||
						parser.directValueMatch(lexer.Elif) {
						break
					}
					continue
				}
				elifBodyNode, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				elifBody = append(elifBody, elifBodyNode)
			}
			lastCondition.Else = append(
				lastCondition.Else,
				&ast.IfStatement{
					Condition: elifCondition.(ast.IExpression),
					Body:      elifBody,
				},
			)
			lastCondition = lastCondition.Else[0].(*ast.IfStatement)
			if parser.directValueMatch(lexer.Else) ||
				parser.directValueMatch(lexer.End) {
				break parsingElifLoop
			}
		}
	}
	// Parse Default
	if parser.directValueMatch(lexer.Else) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		var elseBodyNode ast.Node
		if !parser.directValueMatch(lexer.NewLine) {
			return nil, newSyntaxError(parser.currentLine(), ElseBlock)
		}
		for parser.hasNext() {
			if parser.kindMatch(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.directValueMatch(lexer.End) {
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
	if !parser.directValueMatch(lexer.End) {
		return nil, newStatementNeverEndedError(parser.currentLine(), IfStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return root, nil
}

func (parser *Parser) parseUnlessStatement() (*ast.UnlessStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	line := parser.currentLine()
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(line, UnlessStatement)
	}
	if !parser.directValueMatch(lexer.NewLine) {
		return nil, newSyntaxError(line, UnlessStatement)
	}
	// Parse Unless
	root := &ast.UnlessStatement{
		Condition: condition.(ast.IExpression),
		Body:      []ast.Node{},
		Else:      []ast.Node{},
	}
	var bodyNode ast.Node
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.directValueMatch(lexer.Elif) ||
				parser.directValueMatch(lexer.Else) ||
				parser.directValueMatch(lexer.End) {
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
	if parser.directValueMatch(lexer.Elif) {
		var elifBody []ast.Node
	parsingElifLoop:
		for parser.hasNext() {
			if parser.kindMatch(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.directValueMatch(lexer.Else) ||
					parser.directValueMatch(lexer.End) {
					break
				}
				continue
			}
			if !parser.directValueMatch(lexer.Elif) {
				return nil, newSyntaxError(line, ElifBlock)
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			var elifCondition ast.Node
			elifCondition, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := elifCondition.(ast.IExpression); !ok {
				return nil, newSyntaxError(line, ElifBlock)
			}
			if !parser.directValueMatch(lexer.NewLine) {
				return nil, newSyntaxError(line, ElifBlock)
			}
			var elifBodyNode ast.Node
			for parser.hasNext() {
				if parser.kindMatch(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					if parser.directValueMatch(lexer.Else) ||
						parser.directValueMatch(lexer.End) ||
						parser.directValueMatch(lexer.Elif) {
						break
					}
					continue
				}
				elifBodyNode, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				elifBody = append(elifBody, elifBodyNode)
			}
			lastCondition.Else = append(
				lastCondition.Else,
				&ast.UnlessStatement{
					Condition: elifCondition.(ast.IExpression),
					Body:      elifBody,
				},
			)
			lastCondition = lastCondition.Else[0].(*ast.UnlessStatement)
			if parser.directValueMatch(lexer.Else) ||
				parser.directValueMatch(lexer.End) {
				break parsingElifLoop
			}
		}
	}
	// Parse Default
	if parser.directValueMatch(lexer.Else) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		var elseBodyNode ast.Node
		if !parser.directValueMatch(lexer.NewLine) {
			return nil, newSyntaxError(parser.currentLine(), ElseBlock)
		}
		for parser.hasNext() {
			if parser.kindMatch(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.directValueMatch(lexer.End) {
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
	if !parser.directValueMatch(lexer.End) {
		return nil, newStatementNeverEndedError(parser.currentLine(), UnlessStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return root, nil
}

func (parser *Parser) parseSwitchStatement() (*ast.SwitchStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	line := parser.currentLine()
	target, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := target.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(line, SwitchStatement)
	}
	if !parser.directValueMatch(lexer.NewLine) {
		return nil, newSyntaxError(line, SwitchStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	// parse Cases
	var caseBlocks []*ast.CaseBlock
	if parser.directValueMatch(lexer.Case) {
		for parser.hasNext() {
			if parser.directValueMatch(lexer.Default) ||
				parser.directValueMatch(lexer.End) {
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
			var cases []ast.IExpression
			var caseTarget ast.Node
			for parser.hasNext() {
				caseTarget, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				if _, ok := caseTarget.(ast.IExpression); !ok {
					return nil, newNonExpressionReceivedError(line, CaseBlock)
				}
				cases = append(cases, caseTarget.(ast.IExpression))
				if parser.directValueMatch(lexer.NewLine) {
					break
				} else if parser.directValueMatch(lexer.Comma) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
				} else {
					return nil, newSyntaxError(line, CaseBlock)
				}
			}
			if !parser.directValueMatch(lexer.NewLine) {
				return nil, newSyntaxError(parser.currentLine(), CaseBlock)
			}
			// Targets Body
			var caseBody []ast.Node
			var caseBodyNode ast.Node
			for parser.hasNext() {
				if parser.kindMatch(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					if parser.directValueMatch(lexer.Case) ||
						parser.directValueMatch(lexer.Default) ||
						parser.directValueMatch(lexer.End) {
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
	if parser.directValueMatch(lexer.Default) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.directValueMatch(lexer.NewLine) {
			return nil, newSyntaxError(parser.currentLine(), DefaultBlock)
		}
		var defaultBodyNode ast.Node
		for parser.hasNext() {
			if parser.kindMatch(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.directValueMatch(lexer.End) {
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
	if !parser.directValueMatch(lexer.End) {
		return nil, newStatementNeverEndedError(parser.currentLine(), SwitchStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.SwitchStatement{
		Target:     target.(ast.IExpression),
		CaseBlocks: caseBlocks,
		Default:    defaultBody,
	}, nil
}

func (parser *Parser) parseReturnStatement() (*ast.ReturnStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var results []ast.IExpression
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) || parser.kindMatch(lexer.EOF) {
			break
		}
		line := parser.currentLine()
		result, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := result.(ast.IExpression); !ok {
			return nil, newNonExpressionReceivedError(line, ReturnStatement)
		}
		results = append(results, result.(ast.IExpression))
		if parser.directValueMatch(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !(parser.kindMatch(lexer.Separator) || parser.kindMatch(lexer.EOF)) {
			return nil, newSyntaxError(parser.currentLine(), ReturnStatement)
		}
	}
	return &ast.ReturnStatement{
		Results: results,
	}, nil
}

func (parser *Parser) parseYieldStatement() (*ast.YieldStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var results []ast.IExpression
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) || parser.kindMatch(lexer.EOF) {
			break
		}
		line := parser.currentLine()
		result, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := result.(ast.IExpression); !ok {
			return nil, newNonExpressionReceivedError(line, YieldStatement)
		}
		results = append(results, result.(ast.IExpression))
		if parser.directValueMatch(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !(parser.kindMatch(lexer.Separator) || parser.kindMatch(lexer.EOF)) {
			return nil, newSyntaxError(parser.currentLine(), YieldStatement)
		}
	}
	return &ast.YieldStatement{
		Results: results,
	}, nil
}

func (parser *Parser) parseContinueStatement() (*ast.ContinueStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.ContinueStatement{}, nil
}

func (parser *Parser) parseBreakStatement() (*ast.BreakStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.BreakStatement{}, nil
}

func (parser *Parser) parseRedoStatement() (*ast.RedoStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.RedoStatement{}, nil
}

func (parser *Parser) parsePassStatement() (*ast.PassStatement, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.PassStatement{}, nil
}

func (parser *Parser) parseOperand() (ast.Node, *errors.Error) {
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
	case lexer.Keyboard:
		switch parser.currentToken.DirectValue {
		case lexer.Lambda:
			return parser.parseLambdaExpression()
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
		case lexer.Interface:
			return parser.parseInterfaceStatement()
		case lexer.Class:
			return parser.parseClassStatement()
		case lexer.Raise:
			return parser.parseRaiseStatement()
		case lexer.Try:
			return parser.parseTryStatement()
		case lexer.Return:
			return parser.parseReturnStatement()
		case lexer.Yield:
			return parser.parseYieldStatement()
		case lexer.Continue:
			return parser.parseContinueStatement()
		case lexer.Break:
			return parser.parseBreakStatement()
		case lexer.Redo:
			return parser.parseRedoStatement()
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
	return nil, errors.New(parser.currentLine(), "Unknown Token", errors.ParsingError)
}

func (parser *Parser) parseSelectorExpression(expression ast.IExpression) (*ast.SelectorExpression, *errors.Error) {
	selector := expression
	for parser.hasNext() {
		if !parser.directValueMatch(lexer.Dot) {
			break
		}
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		identifier := parser.currentToken
		if identifier.Kind != lexer.IdentifierKind {
			return nil, newSyntaxError(parser.currentLine(), SelectorExpression)
		}
		selector = &ast.SelectorExpression{
			X: selector,
			Identifier: &ast.Identifier{
				Token: identifier,
			},
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	return selector.(*ast.SelectorExpression), nil
}

func (parser *Parser) parseMethodInvocationExpression(expression ast.IExpression) (*ast.MethodInvocationExpression, *errors.Error) {
	var arguments []ast.IExpression
	// The first token is open parentheses
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	for parser.hasNext() {
		if parser.directValueMatch(lexer.CloseParentheses) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		line := parser.currentLine()
		argument, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := argument.(ast.IExpression); !ok {
			return nil, newNonExpressionReceivedError(line, MethodInvocationExpression)
		}
		arguments = append(arguments, argument.(ast.IExpression))
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.directValueMatch(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.MethodInvocationExpression{
		Function:  expression,
		Arguments: arguments,
	}, nil
}

func (parser *Parser) parseIndexExpression(expression ast.IExpression) (*ast.IndexExpression, *errors.Error) {
	tokenizationError := parser.next()
	if tokenizationError != nil {
		return nil, tokenizationError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	// var rightIndex ast.Node
	line := parser.currentLine()
	index, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := index.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(line, IndexExpression)
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.directValueMatch(lexer.CloseSquareBracket) {
		return nil, newSyntaxError(parser.currentLine(), IndexExpression)
	}
	tokenizationError = parser.next()
	if tokenizationError != nil {
		return nil, tokenizationError
	}
	return &ast.IndexExpression{
		Source: expression,
		Index:  index.(ast.IExpression),
	}, nil
}

func (parser *Parser) parseIfOneLinerExpression(result ast.IExpression) (*ast.IfOneLinerExpression, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	line := parser.currentLine()
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(line, IfOneLinerExpression)
	}
	if !parser.directValueMatch(lexer.Else) {
		return &ast.IfOneLinerExpression{
			Result:     result,
			Condition:  condition.(ast.IExpression),
			ElseResult: nil,
		}, nil
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var elseResult ast.Node
	elseResult, parsingError = parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := elseResult.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(line, OneLineElseBlock)
	}
	return &ast.IfOneLinerExpression{
		Result:     result,
		Condition:  condition.(ast.IExpression),
		ElseResult: elseResult.(ast.IExpression),
	}, nil
}

func (parser *Parser) parseUnlessOneLinerExpression(result ast.IExpression) (*ast.UnlessOneLinerExpression, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	line := parser.currentLine()
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(line, UnlessOneLinerExpression)
	}
	if !parser.directValueMatch(lexer.Else) {
		return &ast.UnlessOneLinerExpression{
			Result:     result,
			Condition:  condition.(ast.IExpression),
			ElseResult: nil,
		}, nil
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var elseResult ast.Node
	line = parser.currentLine()
	elseResult, parsingError = parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(line, OneLineElseBlock)
	}
	return &ast.UnlessOneLinerExpression{
		Result:     result,
		Condition:  condition.(ast.IExpression),
		ElseResult: elseResult.(ast.IExpression),
	}, nil
}

func (parser *Parser) parseGeneratorExpression(operation ast.IExpression) (*ast.GeneratorExpression, *errors.Error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var variables []*ast.Identifier
	numberOfVariables := 0
	for parser.hasNext() {
		if parser.directValueMatch(lexer.In) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.kindMatch(lexer.IdentifierKind) {
			return nil, newSyntaxError(parser.currentLine(), GeneratorExpression)
		}
		variables = append(variables, &ast.Identifier{
			Token: parser.currentToken,
		})
		numberOfVariables++
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.directValueMatch(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
	}
	if numberOfVariables == 0 {
		return nil, newSyntaxError(parser.currentLine(), GeneratorExpression)
	}
	if !parser.directValueMatch(lexer.In) {
		return nil, newSyntaxError(parser.currentLine(), GeneratorExpression)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	line := parser.currentLine()
	source, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := source.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(line, GeneratorExpression)
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	// Finally detect the closing parentheses
	if !parser.directValueMatch(lexer.CloseParentheses) {
		return nil, newSyntaxError(line, GeneratorExpression)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.GeneratorExpression{
		Operation: operation,
		Receivers: variables,
		Source:    source.(ast.IExpression),
	}, nil
}

func (parser *Parser) parseAssignmentStatement(leftHandSide ast.IExpression) (*ast.AssignStatement, *errors.Error) {
	assignmentToken := parser.currentToken
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	line := parser.currentLine()
	rightHandSide, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := rightHandSide.(ast.IExpression); !ok {
		return nil, newNonExpressionReceivedError(line, AssignStatement)
	}
	return &ast.AssignStatement{
		LeftHandSide:   leftHandSide,
		AssignOperator: assignmentToken,
		RightHandSide:  rightHandSide.(ast.IExpression),
	}, nil
}

func (parser *Parser) parsePrimaryExpression() (ast.Node, *errors.Error) {
	var parsedNode ast.Node
	var parsingError *errors.Error
	parsedNode, parsingError = parser.parseOperand()
	if parsingError != nil {
		return nil, parsingError
	}
expressionPendingLoop:
	for {
		switch parser.currentToken.DirectValue {
		case lexer.Dot: // Is selector
			parsedNode, parsingError = parser.parseSelectorExpression(parsedNode.(ast.IExpression))
		case lexer.OpenParentheses: // Is function Call
			parsedNode, parsingError = parser.parseMethodInvocationExpression(parsedNode.(ast.IExpression))
		case lexer.OpenSquareBracket: // Is indexing
			parsedNode, parsingError = parser.parseIndexExpression(parsedNode.(ast.IExpression))
		case lexer.If: // One line If
			parsedNode, parsingError = parser.parseIfOneLinerExpression(parsedNode.(ast.IExpression))
		case lexer.Unless: // One line Unless
			parsedNode, parsingError = parser.parseUnlessOneLinerExpression(parsedNode.(ast.IExpression))
		default:
			if parser.kindMatch(lexer.Assignment) {
				parsedNode, parsingError = parser.parseAssignmentStatement(parsedNode.(ast.IExpression))
			}
			break expressionPendingLoop
		}
	}
	if parsingError != nil {
		return nil, parsingError
	}
	return parsedNode, nil
}

func (parser *Parser) Parse() (*ast.Program, *errors.Error) {
	result := &ast.Program{
		Begin: nil,
		End:   nil,
		Body:  nil,
	}
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var beginStatement *ast.BeginStatement
	var endStatement *ast.EndStatement
	var parsedExpression ast.Node
	var parsingError *errors.Error
	for parser.hasNext() {
		if parser.kindMatch(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			continue
		}
		if parser.directValueMatch(lexer.BEGIN) {
			beginStatement, parsingError = parser.parseBeginStatement()
			if result.Begin != nil {
				return nil, errors.New(parser.currentLine(), "multiple declarations of BEGIN statement at line", errors.ParsingError)
			}
			result.Begin = beginStatement
		} else if parser.directValueMatch(lexer.END) {
			endStatement, parsingError = parser.parseEndStatement()
			if result.End != nil {
				return nil, errors.New(parser.currentLine(), "multiple declarations of END statement at line", errors.ParsingError)
			}
			result.End = endStatement
		} else {
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
		lexer:        lexer_,
		complete:     false,
		currentToken: nil,
	}
}
