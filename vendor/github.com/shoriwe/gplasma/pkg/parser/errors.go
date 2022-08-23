package parser

import (
	"errors"
	"fmt"
)

var (
	UnknownToken  = errors.New("unknown token")
	BeginRepeated = errors.New("repeated begin statement")
	EndRepeated   = errors.New("repeated end statement")
)

func (parser *Parser) newError(message string) error {
	if parser.lineStack.HasNext() {
		return fmt.Errorf("SyntaxError: %s, at line %d", message, parser.lineStack.Peek())
	}
	return fmt.Errorf("SyntaxError: %s", message)
}

func (parser *Parser) newSyntaxError(nodeType string) error {
	return parser.newError(fmt.Sprintf("invalid %s syntax", nodeType))
}

func (parser *Parser) expectingExpressionError(nodeType string) error {
	return parser.newError(fmt.Sprintf("expecting expression but received %s", nodeType))
}

func (parser *Parser) expectingIdentifier(nodeType string) error {
	return parser.newError(fmt.Sprintf("expecting identifier but received %s", nodeType))
}

func (parser *Parser) statementNeverEndedError(nodeType string) error {
	return parser.newError(fmt.Sprintf("statement %s never ended", nodeType))
}

func (parser *Parser) invalidTokenKind() error {
	return parser.newError(fmt.Sprintf("invalid token kind"))
}

func (parser *Parser) expressionNeverClosedError(nodeType string) error {
	return parser.newError(fmt.Sprintf("expression %s never closed", nodeType))
}

func (parser *Parser) expectingFunctionDefinition(nodeType string) error {
	return parser.newError(fmt.Sprintf("expecting function definition but received %s", nodeType))
}
