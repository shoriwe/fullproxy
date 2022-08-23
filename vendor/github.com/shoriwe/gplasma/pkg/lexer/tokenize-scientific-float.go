package lexer

import "errors"

var (
	ScientificFloatInvalidSyntax = errors.New("invalid scientific float syntax")
)

func (lexer *Lexer) tokenizeScientificFloat() error {
	if !lexer.reader.HasNext() {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return ScientificFloatInvalidSyntax
	}
	direction := lexer.reader.Char()
	if (direction != '-') && (direction != '+') {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return ScientificFloatInvalidSyntax
	}
	lexer.reader.Next()
	// Ensure next is a number
	if !lexer.reader.HasNext() {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return ScientificFloatInvalidSyntax
	}
	lexer.currentToken.append(direction)
	nextDigit := lexer.reader.Char()
	if !('0' <= nextDigit && nextDigit <= '9') {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return ScientificFloatInvalidSyntax
	}
	lexer.currentToken.append(nextDigit)
	lexer.reader.Next()
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !('0' <= nextDigit && nextDigit <= '9') && nextDigit != '_' {
			break
		}
		lexer.currentToken.append(nextDigit)
	}
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = ScientificFloat
	return nil
}
